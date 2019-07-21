package cmd

import (
	"net"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusgiger/iot-abe/pkg/acc"
	"github.com/mariusgiger/iot-abe/pkg/cctv"
	"github.com/mariusgiger/iot-abe/pkg/cctv/image"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/spf13/cobra"
)

// serverCmd represents the iot server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts an iot device server",
	Long:  `Starts an iot device server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}

		//TODO profile cpu
		// f, err := os.Create(*cpuprofile)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// pprof.StartCPUProfile(f)
		// defer pprof.StopCPUProfile()

		accManager, err := acc.NewManager(log, wm, client)
		if err != nil {
			return err
		}
		contract := common.HexToAddress(serverOptions.contract)
		device := common.HexToAddress(serverOptions.device)
		cameraService := image.NewCamera(log)

		server := cctv.NewServer(
			net.JoinHostPort(serverOptions.serverInterface, strconv.Itoa(serverOptions.serverPort)),
			log,
			cfg,
			accManager,
			contract,
			device,
			cameraService,
		)

		return server.Run()
	},
}

var (
	serverOptions struct {
		serverInterface string
		serverPort      int
		contract        string
		device          string
	}
)

func init() {
	serverCmd.Flags().StringVarP(&serverOptions.serverInterface, "bind", "", "0.0.0.0", "interface to which the server will bind")
	serverCmd.Flags().IntVarP(&serverOptions.serverPort, "port", "p", 8080, "port on which the server will listen")
	serverCmd.Flags().StringVar(&serverOptions.contract, "contract", "0xC695C023d4A2FfB1C98e0d609A7Ff02e858AF09e", "contract address")
	serverCmd.Flags().StringVar(&serverOptions.device, "device", "0xE1097bAAA914277A8E2AefE464f8E29557e5f046", "device address")

	RootCmd.AddCommand(serverCmd)
}
