package cmd

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/spf13/cobra"
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manages eth wallets",
	Long:  `Manages eth wallets.`,
}

// addWalletCmd represents the add wallet command
var addWalletCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new eth account",
	Long:  `Adds a new eth account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		password, err := wallet.PromptProvidePrivatePass(reader)
		if err != nil {
			return err
		}
		log.Info("keep your password safe, it is uniquely related to this ETH account")

		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		acc, err := wm.NewAccount(string(password))
		if err != nil {
			return err
		}

		log.Infof("new account created: %v", acc.Address.Hex())
		return nil
	},
}

// listWalletCmd represents the list wallet command
var listWalletCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all eth accounts",
	Long:  `Lists all eth accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}

		accounts := wm.Accounts()
		if len(accounts) == 0 {
			log.Infof("no ethereum accounts found under %v", cfg.EthKeystoreDir)
		}

		for _, acc := range accounts {
			balance, err := client.BalanceAt(acc.Address)
			if err != nil {
				return err
			}

			log.Infof("address: %v | %v ETH", acc.Address.String(), float64(balance.Uint64())/float64(core.Ether))
		}

		return nil
	},
}

// importWalletCmd represents the import wallet command
var importWalletCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports an eth account",
	Long:  `Imports an eth account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)

		prompt := "Enter the private key of your existing wallet for importing it"
		reader := bufio.NewReader(os.Stdin)
		priv, err := wallet.PromptProvideSecret(reader, prompt)
		if err != nil {
			return err
		}

		password, err := wallet.PromptProvidePrivatePass(reader)
		if err != nil {
			return err
		}
		fmt.Println("keep your password safe, it is uniquely related to this ETH account")

		acc, err := wm.ImportPrivKey(priv, string(password))
		if err != nil {
			return err
		}
		log.Infof("new account imported: %v", acc.Address.Hex())

		return nil
	},
}

var (
	transferOptions struct {
		from   string
		to     string
		amount float64
	}
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfers ETH between accounts",
	Long:  `Transfers ETH between accounts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wm := wallet.NewManager(log, cfg.EthKeystoreDir)
		client, err := rpc.NewRPCClient(log, cfg)
		if err != nil {
			return err
		}

		from := common.HexToAddress(transferOptions.from)
		to := common.HexToAddress(transferOptions.to)
		value := big.NewInt(int64(transferOptions.amount * core.Ether))

		_, err = wm.AccountByAddress(from)
		if err != nil {
			return err
		}

		log.Infof("sending %v ETH from %v to %v", transferOptions.amount, from.Hex(), to.Hex())

		nonce, err := client.PendingNonceAt(from)
		if err != nil {
			return err
		}

		gasPrice, err := client.SuggestGasPrice()
		if err != nil {
			return err
		}
		gasLimit := uint64(21000)

		tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, nil)
		chainID, err := client.NetworkID()
		if err != nil {
			return err
		}

		signedTx, err := wm.SignTx(types.NewEIP155Signer(chainID), &from, tx, chainID)
		if err != nil {
			return err
		}

		err = client.SendTransaction(signedTx)
		if err != nil {
			return err
		}

		hash := signedTx.Hash()
		log.WithField("gasPrice", gasPrice.Int64()).Infof("sent transaction: %v", toEtherscanURL(hash.Hex()))

		return nil
	},
}

func init() {
	transferCmd.Flags().StringVar(&transferOptions.from, "from", "", "from address (private key must be present in local keystore)")
	transferCmd.Flags().StringVar(&transferOptions.to, "to", "", "to address")
	transferCmd.Flags().Float64Var(&transferOptions.amount, "amount", 0, "amount")

	walletCmd.AddCommand(addWalletCmd, importWalletCmd, listWalletCmd, transferCmd)

	RootCmd.AddCommand(walletCmd)
}

func toEtherscanURL(hash string) string {
	if cfg.UseTestnet {
		return fmt.Sprintf("https://ropsten.etherscan.io/tx/%v", hash)
	}

	return fmt.Sprintf("https://etherscan.io/tx/%v", hash)
}
