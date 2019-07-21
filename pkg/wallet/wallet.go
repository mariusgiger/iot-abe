package wallet

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Manager is the interface for the eth wallet functionality
type Manager interface {
	UsePW(pw string)
	NewAccount(password string) (*accounts.Account, error)
	Accounts() []accounts.Account
	AccountByAddress(address common.Address) (*accounts.Account, error)
	SignTx(signer types.Signer, address *common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	ImportPrivKey(priv []byte, passphrase string) (*accounts.Account, error)
	DecryptMessage(address common.Address, cipherText []byte) ([]byte, error)
}

// ethManager implements an eth wallet manager using go-ethereum's keystore
// also refer to: https://goethereumbook.org/keystore/
type ethManager struct {
	keyStoreDir string
	log         *logrus.Logger
	ks          *keystore.KeyStore
	pw          string
}

// NewManager creates a new wallet manager
func NewManager(log *logrus.Logger, keyStoreDir string) Manager {
	ks := keystore.NewKeyStore(keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	return &ethManager{keyStoreDir: keyStoreDir, log: log, ks: ks}
}

// UsePW sets the password (used for debugging)
func (wm *ethManager) UsePW(pw string) {
	wm.pw = pw
}

// NewAccount creates a new ethereum wallet and account
func (wm *ethManager) NewAccount(password string) (*accounts.Account, error) {
	account, err := wm.ks.NewAccount(password)
	if err != nil {
		return nil, err
	}
	wm.log.WithField("account", account.Address.String()).Info("created a new account")

	return &account, nil
}

// Accounts return a list of all eth accounts
func (wm *ethManager) Accounts() []accounts.Account {
	ks := keystore.NewKeyStore(wm.keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks.Accounts()
}

// AccountByAddress return a an account by address
func (wm *ethManager) AccountByAddress(address common.Address) (*accounts.Account, error) {
	ks := keystore.NewKeyStore(wm.keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	for _, acc := range ks.Accounts() {
		if acc.Address == address {
			return &acc, nil
		}
	}

	return nil, fmt.Errorf("account %v not found in local wallet", address.Hex())
}

// ImportPrivKey imports a private key
func (wm *ethManager) ImportPrivKey(priv []byte, passphrase string) (*accounts.Account, error) {
	ks := keystore.NewKeyStore(wm.keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	privECDSA, err := crypto.HexToECDSA(string(priv))
	if err != nil {
		return nil, err
	}

	acc, err := ks.ImportECDSA(privECDSA, passphrase)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

// SignTx signs an transaction
func (wm *ethManager) SignTx(signer types.Signer, address *common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	ks := keystore.NewKeyStore(wm.keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	acc, err := wm.AccountByAddress(*address)
	if err != nil {
		return nil, err
	}

	var password string
	if wm.pw == "" {
		reader := bufio.NewReader(os.Stdin)
		pwBytes, err := PromptProvideSecret(reader, "Enter password for unlocking wallet and signing transaction")
		if err != nil {
			return nil, err
		}

		password = string(pwBytes)
	} else {
		password = wm.pw
	}

	err = ks.Unlock(*acc, password)
	if err != nil {
		return nil, err
	}

	defer func() {
		utils.IgnoreError(ks.Lock(acc.Address))
	}()

	return ks.SignTx(*acc, tx, chainID)
}

//DecryptMessage implements decrypting a message based on ECIES
func (wm *ethManager) DecryptMessage(address common.Address, cipherText []byte) ([]byte, error) {
	ks := keystore.NewKeyStore(wm.keyStoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	acc, err := wm.AccountByAddress(address)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get account by address (%v", address.Hex())
	}

	reader := bufio.NewReader(os.Stdin)
	password, err := PromptProvideSecret(reader, "Enter password for unlocking wallet and decrypting the message")
	if err != nil {
		return nil, errors.Wrap(err, "could read secret from command line")
	}

	fmt.Printf("passphrase %v\n", string(password))
	password = []byte("1234") //TODO somehow cmd input is buggy?
	tempPass := "pwd"         //should this be generated on the fly?
	encKey, err := ks.Export(*acc, string(password), tempPass)
	if err != nil {
		return nil, errors.Wrap(err, "could not export private key")
	}

	key, err := keystore.DecryptKey(encKey, tempPass)
	if err != nil {
		return nil, errors.Wrap(err, "could decrypt key with temporary passphrase")
	}

	message, err := utils.DecryptMessage(key.PrivateKey, cipherText)
	if err != nil {
		return nil, errors.Wrap(err, "could not decrypt ABE priv key")
	}

	return message, nil
}
