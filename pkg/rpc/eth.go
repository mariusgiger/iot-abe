package rpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mariusgiger/iot-abe/pkg/contract"
	"github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/pkg/errors"
)

// Client interface for an ETH rpc client
type Client interface {
	PendingNonceAt(from common.Address) (uint64, error)
	SuggestGasPrice() (*big.Int, error)
	BalanceAt(address common.Address) (*big.Int, error)
	EstimateGas(method string, from, contractAddr common.Address, args ...interface{}) (uint64, error)
	GetRawClient() bind.ContractBackend
	TransactionsByAddress(address common.Address) ([]*types.Transaction, error)
	TransactionByHash(txHash common.Hash) (*types.Transaction, bool, error)
	TransactionReceipt(txHash common.Hash) (*types.Receipt, error)
	SendTransaction(signedTx *types.Transaction) error
	NetworkID() (*big.Int, error)
	BlockByNumber(number *big.Int) (*types.Block, error)
}

// ethClient wraps eth rpc calls and implements rpc.ETHClient
type ethClient struct {
	rpcClient       *rpc.Client
	ethClient       *ethclient.Client
	httpClient      *http.Client
	etherscanAPIKey string
	etherscanURL    string
	log             *logrus.Logger
}

// NewRPCClient creates a new ETH rpc client
func NewRPCClient(log *logrus.Logger, cfg *core.Config) (Client, error) {
	httpClient, err := GetHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	log.Infof("dialing RPC client via HTTP at: %v", cfg.ETHNodeURL)
	rpcClient, err := rpc.DialHTTPWithClient(cfg.ETHNodeURL, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial RPC client via HTTP")
	}

	return &ethClient{
		rpcClient:       rpcClient,
		ethClient:       ethclient.NewClient(rpcClient),
		httpClient:      httpClient,
		etherscanAPIKey: cfg.EtherscanAPIKey,
		etherscanURL:    cfg.EtherscanURL,
	}, nil
}

// NewWSSRPCClient creates a new ETH rpc client
func NewWSSRPCClient(log *logrus.Logger, cfg *core.Config) (Client, error) {
	log.Infof("dialing RPC client via WSS at: %v", cfg.ETHWSSNodeURL)
	client, err := ethclient.Dial(cfg.ETHWSSNodeURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial RPC client via WSS")
	}

	return &ethClient{
		ethClient:       client,
		etherscanAPIKey: cfg.EtherscanAPIKey,
		etherscanURL:    cfg.EtherscanURL,
	}, nil
}

// GetRawClient returns the raw rpc client
func (c *ethClient) GetRawClient() bind.ContractBackend {
	return c.ethClient
}

// PendingNonceAt returns the account nonce of the given account in the pending state. This is the nonce that should be used for the next transaction.
func (c *ethClient) PendingNonceAt(from common.Address) (uint64, error) {
	return c.ethClient.PendingNonceAt(context.Background(), from)
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely execution of a transaction.
func (c *ethClient) SuggestGasPrice() (*big.Int, error) {
	return c.ethClient.SuggestGasPrice(context.Background())
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (c *ethClient) EstimateGas(method string, from, contractAddr common.Address, args ...interface{}) (uint64, error) {
	if method == "" {
		return 1000000, nil // contract creation
	}

	// createContractMethodInvocation creates a CMI used for the ETH ABI as a Data attribute
	// https://theethereum.wiki/w/index.php/ERC20_Token_Standard
	abi, err := abi.JSON(strings.NewReader(contract.AccessControlABI))
	if err != nil {
		return 0, errors.Wrap(err, "could not read abi")
	}

	cmi, err := abi.Pack(method, args)
	if err != nil {
		return 0, errors.Wrap(err, "could not pack abi")
	}

	msg := ethereum.CallMsg{
		From:     from,
		To:       &contractAddr,
		Gas:      0,
		GasPrice: big.NewInt(0),
		Value:    big.NewInt(0),
		Data:     cmi,
	}

	return c.ethClient.EstimateGas(context.Background(), msg)
}

// BalanceAt returns the wei balance of the given account. The block number can be nil, in which case the balance is taken from the latest known block.
func (c *ethClient) BalanceAt(address common.Address) (*big.Int, error) {
	return c.ethClient.BalanceAt(context.Background(), address, nil)
}

//SendTransaction injects a signed transaction into the pending pool for execution.
func (c *ethClient) SendTransaction(signedTx *types.Transaction) error {
	return c.ethClient.SendTransaction(context.Background(), signedTx)
}

//NetworkID returns the network ID (also known as the chain ID) for this chain.
func (c *ethClient) NetworkID() (*big.Int, error) {
	return c.ethClient.NetworkID(context.Background())
}

//TransactionByHash returns the transaction with the given hash.
func (c *ethClient) TransactionByHash(txHash common.Hash) (*types.Transaction, bool, error) {
	return c.ethClient.TransactionByHash(context.Background(), txHash)
}

//TransactionReceipt returns the transaction receipt with the given hash.
func (c *ethClient) TransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return c.ethClient.TransactionReceipt(context.Background(), txHash)
}

//BlockByNumber returns the block by number, if number is nil the latest block is returned
func (c *ethClient) BlockByNumber(number *big.Int) (*types.Block, error) {
	return c.ethClient.BlockByNumber(context.Background(), number)
}

// Balance returns the balance for an eth account
func (c *ethClient) TransactionsByAddress(address common.Address) ([]*types.Transaction, error) {
	startblock := 0
	endblock := 99999999

	url := utils.MustParseURL(c.etherscanURL).
		AddPath("/api").
		SetQuery("module", "account").
		SetQuery("action", "txlist").
		SetQuery("address", address.Hex()).
		SetQuery("startblock", strconv.Itoa(startblock)).
		SetQuery("endblock", strconv.Itoa(endblock)).
		SetQuery("sort", "asc").
		SetQuery("apikey", c.etherscanAPIKey).
		String()

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txResp TransactionResponse
	err = json.NewDecoder(resp.Body).Decode(&txResp)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode etherscan result")
	}

	var txs []*types.Transaction
	for _, scanTx := range txResp.Result {
		hash := common.HexToHash(scanTx.Hash)
		tx, _, err := c.ethClient.TransactionByHash(context.Background(), hash)
		if err != nil {
			return nil, errors.Wrap(err, "could not get transaction by hash")
		}

		txs = append(txs, tx)
	}

	return txs, nil
}

// TransactionResponse is the response from etherscan for normal transactions
type TransactionResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Result  []TxResult `json:"result"`
}

// TxResult is the result struct in TransactionResponse
type TxResult struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

// GetHTTPClient creates new HTTP client for network
func GetHTTPClient(cfg *core.Config) (*http.Client, error) {
	// see http.DefaultTransport for reference
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if cfg.ClientCert != "" {
		cert, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientCertKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load X509 certificate")
		}
		clientConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		clientConfig.BuildNameToCertificate()
		transport.TLSClientConfig = clientConfig
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

// compile-time assertion to make sure ethClient implements the interface
var _ Client = (*ethClient)(nil)
