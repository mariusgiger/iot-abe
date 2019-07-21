package acc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MockClient struct {
	rpc             *backends.SimulatedBackend
	contractAddress common.Address
	contractTx      *types.Transaction
	confirmations   int64
}

func (c *MockClient) PendingNonceAt(from common.Address) (uint64, error) {
	return c.rpc.PendingNonceAt(context.Background(), from)
}

func (c *MockClient) SuggestGasPrice() (*big.Int, error) {
	return c.rpc.SuggestGasPrice(context.Background())
}

func (c *MockClient) EstimateGas(method string, from, contractAddr common.Address, args ...interface{}) (uint64, error) {
	return c.rpc.EstimateGas(context.Background(), ethereum.CallMsg{})
}

func (c *MockClient) BalanceAt(address common.Address) (*big.Int, error) {
	return c.rpc.BalanceAt(context.Background(), address, nil)
}

func (c *MockClient) ContractAddress(sendAddress common.Address) (common.Address, *types.Transaction, int64, error) {
	return c.contractAddress, c.contractTx, c.confirmations, nil
}

func (c *MockClient) GetRawClient() bind.ContractBackend {
	return c.rpc
}

func (c *MockClient) TransactionsByAddress(address common.Address) ([]*types.Transaction, error) {
	return nil, nil
}

func (c *MockClient) TransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (c *MockClient) TransactionByHash(txHash common.Hash) (*types.Transaction, bool, error) {
	return nil, false, nil
}

func (c *MockClient) SendTransaction(signedTx *types.Transaction) error {
	return nil
}

func (c *MockClient) NetworkID() (*big.Int, error) {
	return nil, nil
}

func (c *MockClient) BlockByNumber(number *big.Int) (*types.Block, error) {
	return nil, nil
}
