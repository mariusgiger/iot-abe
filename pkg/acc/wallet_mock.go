package acc

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type MockWalletManager struct {
	accounts map[common.Address]*ecdsa.PrivateKey
}

func (wm *MockWalletManager) UsePW(pw string) {
	panic("not implemented")
}

func (wm *MockWalletManager) NewAccount(password string) (*accounts.Account, error) {
	panic("not implemented")
}

func (wm *MockWalletManager) Accounts() []accounts.Account {
	panic("not implemented")
}

func (wm *MockWalletManager) AccountByAddress(address common.Address) (*accounts.Account, error) {
	panic("not implemented")
}

func (wm *MockWalletManager) SignTx(signer types.Signer, address *common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	key, ok := wm.accounts[*address]
	if !ok {
		return nil, errors.New("no account found")
	}

	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, signature)
}

func (wm *MockWalletManager) ImportPrivKey(priv []byte, passphrase string) (*accounts.Account, error) {
	panic("not implemented")
}

func (wm *MockWalletManager) DecryptMessage(address common.Address, cipherText []byte) ([]byte, error) {
	panic("not implemented")
}
