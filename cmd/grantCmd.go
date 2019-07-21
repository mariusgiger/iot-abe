package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/contract"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// grantCmd represents the access grant command
var grantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Manages access rights",
	Long:  `Manages access rights.`,
}

var (
	initOptions struct {
		from string
	}
	watchOptions struct {
		contract string
	}
	grantOptions struct {
		contract   string
		subject    string
		owner      string
		attributes string
	}
)

// Keys stores keys for admin
type Keys struct {
	MasterKey       string
	PublicKey       string
	ContractTxHash  string
	ContractAddress string
	Created         time.Time
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Deploys a new smart contract for managing ABE",
	Long:  `Deploys a new smart contract for managing ABE.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}
		accManager, err := acc.NewManager(log, wm, client)
		if err != nil {
			return err
		}

		from := common.HexToAddress(initOptions.from)
		_, err = wm.AccountByAddress(from)
		if err != nil {
			return err
		}

		info, err := accManager.Deploy(from)
		if err != nil {
			return err
		}
		log.Infof("deployed new access control smart contract at %v", info.ContractAddress.Hex())
		log.Infof("see %v", toEtherscanURL(info.ContractTxHash.Hex()))

		keys := &Keys{
			PublicKey:       info.PublicKey,
			MasterKey:       info.MasterKey,
			ContractTxHash:  info.ContractTxHash.Hex(),
			ContractAddress: info.ContractAddress.Hex(),
			Created:         time.Now(),
		}

		keyPath := path.Join(keysPath, fmt.Sprintf("%v.json", info.ContractAddress.Hex()))
		keysBytes, err := json.Marshal(keys)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(keyPath, keysBytes, 0644)
		if err != nil {
			return err
		}

		log.Infof("wrote keys file to %v", keyPath)
		return nil
	},
}

var watchAccessRequestsCmd = &cobra.Command{
	Use:   "watch-requests",
	Short: "Watches access requests for a smart contract",
	Long:  `Watches access requests for a smart contract.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewWSSRPCClient(log, cfg)
		if err != nil {
			return err
		}
		accManager, err := acc.NewManager(log, wm, client)
		if err != nil {
			return err
		}

		if watchOptions.contract == "" {
			log.Error("contract is empty (use --contract to set it)")
			return nil
		}

		contractAddr := common.HexToAddress(watchOptions.contract)
		events := make(chan *contract.AccessControlAccessRequested)
		subsc, err := accManager.WatchAccessRequested(contractAddr, events)
		if err != nil {
			return err
		}
		defer subsc.Unsubscribe()
		log.Infof("watching AccessRequested events on %v", watchOptions.contract)

		for {
			select {
			case msg := <-events:
				log.Infof("access requested: %v", msg.From.Hex())
			case err := <-subsc.Err():
				log.Errorf("error received: %v", err.Error())
			}
		}
	},
}

var getAccessRequestsCmd = &cobra.Command{
	Use:   "get-requests",
	Short: "Gets all access requests for a smart contract",
	Long:  `Gets all access requests for a smart contract.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}
		accManager, err := acc.NewManager(log, wm, client)
		if err != nil {
			return err
		}

		if watchOptions.contract == "" {
			log.Error("contract is empty (use --contract to set it)")
			return nil
		}

		contract := common.HexToAddress(watchOptions.contract)
		requests, err := accManager.AccessRequests(contract)
		if err != nil {
			return err
		}

		for _, request := range requests {
			log.Infof("access request received: from: %v, tx: %v", request.From.Hex(), request.TxHash.Hex())
		}

		return nil
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists",
	Long:  `Lists managed access control smart contracts.`,
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			keyPath := path.Join(keysPath, fmt.Sprintf("%v.json", args[0]))
			keyBytes, err := ioutil.ReadFile(keyPath)
			if err != nil {
				return err
			}

			key := &Keys{}
			err = json.Unmarshal(keyBytes, key)
			if err != nil {
				return err
			}

			return utils.PrintFormatted(key)
		}

		files, err := ioutil.ReadDir(keysPath)
		if err != nil {
			return err
		}

		var keys []*Keys
		for _, f := range files {
			keyPath := path.Join(keysPath, f.Name())
			keyBytes, err := ioutil.ReadFile(keyPath)
			if err != nil {
				return err
			}

			key := &Keys{}
			err = json.Unmarshal(keyBytes, key)
			if err != nil {
				return err
			}

			keys = append(keys, key)
		}

		data := [][]string{}
		for _, key := range keys {
			data = append(data, []string{
				key.ContractAddress,
				key.Created.Format("2006-01-02 15:04:05"),
				toEtherscanURL(key.ContractTxHash),
			})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Contract address", "Created", "Contract tx"})
		table.SetBorder(true)
		table.AppendBulk(data)
		table.Render()

		return nil
	},
}

// grantAccessRightsCmd represents the grant access command
var grantAccessRightsCmd = &cobra.Command{
	Use:   "access",
	Short: "Grants access for a user",
	Long:  `Grants access for a user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}
		accManager, err := acc.NewManager(log, wm, client)
		if err != nil {
			return err
		}

		attrList := strings.Split(grantOptions.attributes, ",")
		contract := common.HexToAddress(grantOptions.contract)
		subject := common.HexToAddress(grantOptions.subject)
		owner := common.HexToAddress(grantOptions.owner)
		_, err = wm.AccountByAddress(owner)
		if err != nil {
			return err
		}

		log.WithFields(logrus.Fields{
			"contract": contract.Hex(),
		}).Infof("granting access to: %v", subject.Hex())

		keyPath := path.Join(keysPath, fmt.Sprintf("%v.json", grantOptions.contract))
		keyBytes, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return err
		}

		key := &Keys{}
		err = json.Unmarshal(keyBytes, key)
		if err != nil {
			return err
		}

		requestTxs, err := accManager.AccessRequests(contract)
		if err != nil {
			return err
		}

		var requestTx *types.Transaction
		for _, request := range requestTxs {
			if request.From == subject { //subject is expected to only send 1 tx
				requestTx = request.Tx
			}
		}

		if requestTx == nil {
			log.Errorf("could not find request tx for: %v", subject)
			return nil
		}

		subjectPubKey, err := utils.RecoverPubKey(requestTx, big.NewInt(3)) //TODO move to constant
		if err != nil {
			return err
		}

		grantInfo, err := accManager.GrantAccess(contract, owner, subject, attrList, key.PublicKey, key.MasterKey, subjectPubKey)
		if err != nil {
			return err
		}

		log.Infof("granted access rights to user, see %v", toEtherscanURL(grantInfo.TxHash.Hex()))
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initOptions.from, "from", "", "from address (private key must be present in local keystore)")
	watchAccessRequestsCmd.Flags().StringVar(&watchOptions.contract, "contract", "", "contract address to watch")
	getAccessRequestsCmd.Flags().StringVar(&watchOptions.contract, "contract", "", "contract address to watch")

	grantAccessRightsCmd.Flags().StringVar(&grantOptions.contract, "contract", "", "contract address")
	grantAccessRightsCmd.Flags().StringVar(&grantOptions.attributes, "attributes", "", "comma-separated list of attributes")
	grantAccessRightsCmd.Flags().StringVar(&grantOptions.owner, "owner", "", "contract owner address (private key must be present in local keystore)")
	grantAccessRightsCmd.Flags().StringVar(&grantOptions.subject, "for", "", "subject to add access rights")

	grantCmd.AddCommand(initCmd, listCmd, watchAccessRequestsCmd, getAccessRequestsCmd, grantAccessRightsCmd)

	RootCmd.AddCommand(grantCmd)
}
