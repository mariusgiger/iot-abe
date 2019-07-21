package cmd

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/contract"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/spf13/cobra"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manages IoT devices",
	Long:  `Manages IoT devices.`,
}

var (
	addOptions struct {
		name     string
		address  string
		owner    string
		contract string
		policy   string
	}
	getOptions struct {
		contract string
		address  string
	}

	removeOptions struct {
		address  string
		owner    string
		contract string
	}
)

// addDeviceCmd represents the add device command
var addDeviceCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds or sets a new IoT device policy",
	Long:  `Adds or sets a new IoT device policy.`,
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

		//TODO rename to set cmd
		//TODO refine device permission scheme
		//TODO validate policy and input
		address := common.HexToAddress(addOptions.address)
		owner := common.HexToAddress(addOptions.owner)
		contract := common.HexToAddress(addOptions.contract)
		tx, err := accManager.SetDevicePolicy(addOptions.name, address, contract, owner, addOptions.policy)
		if err != nil {
			return err
		}

		log.Infof("new device added: see %v", toEtherscanURL(tx.Hash().Hex()))
		return nil
	},
}

// getDeviceCmd represents the get device command
var getDeviceCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an IoT device policy",
	Long:  `Gets an IoT device policy.`,
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

		if getOptions.address == "" {
			return errors.New("device address must not be empty")
		}

		address := common.HexToAddress(getOptions.address)
		contract := common.HexToAddress(getOptions.contract)
		policy, err := accManager.DevicePolicyByAddress(address, contract)
		if err != nil {
			return err
		}
		log.Infof("device policy for %v: %v", getOptions.address, policy)

		return nil
	},
}

var getAllDevicesCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Gets all IoT device policies",
	Long:  `Gets all IoT device policies.`,
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

		contract := common.HexToAddress(getOptions.contract)
		policies, err := accManager.DevicePolicies(contract)
		if err != nil {
			return err
		}

		for _, policy := range policies {
			log.Infof("device policy for %v, %v, see tx: %v", policy.Device.Hex(), policy.Policy, toEtherscanURL(policy.TxHash.Hex()))
		}

		return nil
	},
}

// removeDeviceCmd represents the remove device command
var removeDeviceCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes an IoT device policy",
	Long:  `Removes an IoT device policy.`,
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

		address := common.HexToAddress(removeOptions.address)
		contract := common.HexToAddress(removeOptions.contract)
		owner := common.HexToAddress(removeOptions.owner)
		tx, err := accManager.RemoveDevicePolicy(address, contract, owner)
		if err != nil {
			return err
		}
		log.Infof("device policy removed: see %v", toEtherscanURL(tx.Hash().Hex()))

		return nil
	},
}

// watchDevicePolicyUpdatedCmd represents the command to watch device policy updated events
var watchDevicePolicyUpdatedCmd = &cobra.Command{
	Use:   "watch-policy-updated",
	Short: "Watches device policy updated events",
	Long:  `Watches device policy updated events.`,
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

		contractAddr := common.HexToAddress(watchOptions.contract)
		events := make(chan *contract.AccessControlDevicePolicyUpdated)
		subsc, err := accManager.WatchDevicePolicyUpdated(contractAddr, events)
		if err != nil {
			return err
		}
		defer subsc.Unsubscribe()
		log.Infof("watching DevicePolicyUpdated events on %v", watchOptions.contract)

		for {
			select {
			case msg := <-events:
				log.Infof("device policy updated: device %v, policy: %v", msg.Device.Hex(), msg.Policy)
			case err := <-subsc.Err():
				log.Errorf("error received: %v", err.Error())
			}
		}
	},
}

// watchDevicePolicyRemovedCmd represents the command to watch device policy deleted events
var watchDevicePolicyRemovedCmd = &cobra.Command{
	Use:   "watch-policy-removed",
	Short: "Watches device policy deleted events",
	Long:  `Watches device policy deleted events.`,
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

		contractAddr := common.HexToAddress(watchOptions.contract)
		events := make(chan *contract.AccessControlDevicePolicyDeleted)
		subsc, err := accManager.WatchDevicePolicyRemoved(contractAddr, events)
		if err != nil {
			return err
		}
		defer subsc.Unsubscribe()
		log.Infof("watching DevicePolicyDeleted events on %v", watchOptions.contract)

		for {
			select {
			case msg := <-events:
				log.Infof("device policy removed: device %v", msg.Device.Hex())
			case err := <-subsc.Err():
				log.Errorf("error received: %v", err.Error())
			}
		}
	},
}

func init() {
	addDeviceCmd.Flags().StringVar(&addOptions.owner, "owner", "", "contract owner address (private key must be present in local keystore)")
	addDeviceCmd.Flags().StringVar(&addOptions.address, "device", "", "device for which a policy should be added")
	addDeviceCmd.Flags().StringVar(&addOptions.contract, "contract", "", "contract address")
	addDeviceCmd.Flags().StringVarP(&addOptions.name, "name", "n", "", "name of the device")
	addDeviceCmd.Flags().StringVarP(&addOptions.policy, "policy", "p", "", "policy to add")

	getDeviceCmd.Flags().StringVar(&getOptions.address, "device", "", "device for which the policy should be retrieved")
	getDeviceCmd.Flags().StringVar(&getOptions.contract, "contract", "", "contract address")
	getAllDevicesCmd.Flags().StringVar(&getOptions.contract, "contract", "", "contract address")

	watchDevicePolicyUpdatedCmd.Flags().StringVar(&watchOptions.contract, "contract", "", "contract address to watch")
	watchDevicePolicyRemovedCmd.Flags().StringVar(&watchOptions.contract, "contract", "", "contract address to watch")

	removeDeviceCmd.Flags().StringVar(&removeOptions.owner, "owner", "", "contract owner address (private key must be present in local keystore)")
	removeDeviceCmd.Flags().StringVar(&removeOptions.address, "device", "", "device for which the policy should be removed")
	removeDeviceCmd.Flags().StringVar(&removeOptions.contract, "contract", "", "contract address")

	devicesCmd.AddCommand(addDeviceCmd, getDeviceCmd, getAllDevicesCmd, watchDevicePolicyUpdatedCmd, removeDeviceCmd, watchDevicePolicyRemovedCmd)
	RootCmd.AddCommand(devicesCmd)
}
