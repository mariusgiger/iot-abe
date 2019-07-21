package cmd

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/contract"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/spf13/cobra"
)

// requestCmd represents the request access command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Manages access right requests",
	Long:  `Manages access right requests.`,
}

var (
	requestOptions struct {
		requester string
		contract  string
	}
)

var requestAccessCmd = &cobra.Command{
	Use:   "access",
	Short: "Requests access for a smart contract",
	Long:  `Requests access for a smart contract.`,
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

		requester := common.HexToAddress(requestOptions.requester)
		contract := common.HexToAddress(requestOptions.contract)
		_, err = wm.AccountByAddress(requester)
		if err != nil {
			return err
		}

		info, err := accManager.RequestAccess(requester, contract)
		if err != nil {
			return err
		}
		log.Infof("requested access, see tx: %v", toEtherscanURL(info.Hex()))

		return nil
	},
}

var watchAccessGrantedCmd = &cobra.Command{
	Use:   "watch-grants",
	Short: "Watches access grants for a smart contract",
	Long:  `Watches access grants for a smart contract.`,
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
		events := make(chan *contract.AccessControlAccessGranted)
		subsc, err := accManager.WatchAccessGranted(contractAddr, events)
		if err != nil {
			return err
		}
		defer subsc.Unsubscribe()
		log.Infof("watching AccessGranted events on %v", watchOptions.contract)

		for {
			select {
			case msg := <-events:
				log.Infof("access granted: %v", msg.Requester.Hex())
			case err := <-subsc.Err():
				log.Errorf("error granted: %v", err.Error())
			}
		}
	},
}

var getAccessGrantCmd = &cobra.Command{
	Use:   "get-grant",
	Short: "Gets access an grant for a smart contract",
	Long:  `Gets access an grant for a smart contract.`,
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

		if requestOptions.contract == "" {
			log.Error("contract is empty (use --contract to set it)")
			return nil
		}

		contractAddr := common.HexToAddress(requestOptions.contract)
		subject := common.HexToAddress(requestOptions.requester)
		decryptedKey, err := accManager.GetAccessGrant(subject, contractAddr)
		if err != nil {
			return err
		}

		log.Infof("access grant received, key: %v", decryptedKey)
		return nil
	},
}

func init() {
	requestAccessCmd.Flags().StringVar(&requestOptions.requester, "for", "", "address for which access is requested (private key must be present in local keystore)")
	requestAccessCmd.Flags().StringVar(&requestOptions.contract, "contract", "", "contract address for which access is requested")
	watchAccessGrantedCmd.Flags().StringVar(&watchOptions.contract, "contract", "", "contract address to watch")

	getAccessGrantCmd.Flags().StringVar(&requestOptions.contract, "contract", "", "contract address")
	getAccessGrantCmd.Flags().StringVar(&requestOptions.requester, "for", "", "address of the subject")

	requestCmd.AddCommand(requestAccessCmd, watchAccessGrantedCmd, getAccessGrantCmd)

	RootCmd.AddCommand(requestCmd)
}
