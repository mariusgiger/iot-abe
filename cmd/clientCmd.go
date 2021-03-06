package cmd

import (
	"net"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/cctv"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client for an iot device server",
	Long:  `Client for an iot device server.`,
}

// getDataCmd represents the get data command
var getDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Retrieves encrypted data from an iot device server and decrypts it (if possible) - used for testing",
	Long:  `Retrieves encrypted data from an iot device server and decrypts it (if possible) - used for testing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		rpcClient, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}

		accManager, err := acc.NewManager(log, wm, rpcClient)
		if err != nil {
			return err
		}
		contract := common.HexToAddress(getDataOptions.contract)
		user := common.HexToAddress(getDataOptions.user)

		client := cctv.NewClient(getDataOptions.serverURL, contract, user, accManager, wm, log)
		message, err := client.GetMsg()
		if err != nil {
			return err
		}

		log.Infof("got message: %v", message)
		return nil
	},
}

// serveClientCmd represents the serve decrypted image command
var serveClientCmd = &cobra.Command{
	Use:   "serve",
	Short: "Exposes an endpoint to view decrypted images and video stream.",
	Long:  `Exposes an endpoint to view decrypted images and video stream.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		rpcClient, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}

		accManager, err := acc.NewManager(log, wm, rpcClient)
		if err != nil {
			return err
		}

		contract := common.HexToAddress(serveCaptureOptions.contract)
		user := common.HexToAddress(serveCaptureOptions.user)
		client := cctv.NewClient(serveCaptureOptions.serverURL, contract, user, accManager, wm, log)

		return client.Serve(net.JoinHostPort(serveCaptureOptions.clientInterface, strconv.Itoa(serveCaptureOptions.clientPort)))
	},
}

var (
	getDataOptions struct {
		serverURL string
		contract  string
		user      string
	}
	getImageOptions struct {
		serverURL string
	}
	serveCaptureOptions struct {
		clientInterface string
		clientPort      int
		serverURL       string
		contract        string
		user            string
	}
)

func init() {
	getDataCmd.Flags().StringVarP(&getDataOptions.serverURL, "server", "", "http://localhost:8080", "url of the IoT device server")
	getDataCmd.Flags().StringVar(&getDataOptions.contract, "contract", "0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e", "contract address")
	getDataCmd.Flags().StringVar(&getDataOptions.user, "user", "0x1e52b030261C4890A6aCe85Ed48CaE5f459525A0", "user address (private key must be in the keystore)")

	serveClientCmd.Flags().StringVarP(&serveCaptureOptions.clientInterface, "bind", "", "0.0.0.0", "interface to which the client will bind")
	serveClientCmd.Flags().IntVarP(&serveCaptureOptions.clientPort, "port", "p", 8081, "port on which the client will listen")
	serveClientCmd.Flags().StringVarP(&serveCaptureOptions.serverURL, "server", "", "http://localhost:8080", "url of the IoT device server")
	serveClientCmd.Flags().StringVar(&serveCaptureOptions.contract, "contract", "0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e", "contract address")
	serveClientCmd.Flags().StringVar(&serveCaptureOptions.user, "user", "0x1e52b030261C4890A6aCe85Ed48CaE5f459525A0", "user address (private key must be in the keystore)")

	clientCmd.AddCommand(getDataCmd, serveClientCmd)
	RootCmd.AddCommand(clientCmd)
}
