package acc

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/mariusgiger/iot-abe/pkg/contract"
	units "github.com/mariusgiger/iot-abe/pkg/core"
	"github.com/mariusgiger/iot-abe/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type AccTestSuite struct {
	suite.Suite
	wm         *MockWalletManager
	blockchain *backends.SimulatedBackend
	rpcClient  *MockClient
	manager    *Manager
	alloc      core.GenesisAlloc

	admin     common.Address
	requester common.Address
}

// SetupTest setups a mocked blockchain
func (ts *AccTestSuite) SetupTest() {
	alloc := make(core.GenesisAlloc)
	wm := &MockWalletManager{
		accounts: make(map[common.Address]*ecdsa.PrivateKey),
	}
	ts.wm = wm

	adminPrvKey, adminAddr, err := utils.GenerateKeyPair()
	ts.Require().NoError(err)
	alloc[adminAddr] = core.GenesisAccount{Balance: big.NewInt(1 * units.Ether)}
	wm.accounts[adminAddr] = adminPrvKey
	ts.admin = adminAddr

	requesterPrvKey, requesterAddr, err := utils.GenerateKeyPair()
	ts.Require().NoError(err)
	alloc[requesterAddr] = core.GenesisAccount{Balance: big.NewInt(0.5 * units.Ether)}
	wm.accounts[requesterAddr] = requesterPrvKey
	ts.requester = requesterAddr

	gasLimit := uint64(100 * units.GWei)
	blockchain := backends.NewSimulatedBackend(alloc, gasLimit)
	rpcClient := &MockClient{rpc: blockchain}

	ts.blockchain = blockchain
	manager, err := NewManager(logrus.New(), wm, rpcClient)
	ts.Require().NoError(err)
	ts.manager = manager
	ts.rpcClient = rpcClient
	ts.alloc = alloc
}

func (ts *AccTestSuite) TestDeploy() {
	ts.deploy()
}

func (ts *AccTestSuite) deploy() *DeployInfo {
	//arrange & act
	info, err := ts.manager.Deploy(ts.admin)
	ts.blockchain.Commit()

	//assert
	ts.Require().NoError(err)
	ts.NotEqual(info.ContractAddress, common.Address{})
	ts.NotEqual(info.ContractTxHash, common.Hash{})

	receipt, err := ts.blockchain.TransactionReceipt(context.Background(), info.ContractTxHash)
	ts.Require().NoError(err)
	ts.Require().NotNil(receipt)
	ts.Equal(uint64(1), receipt.Status)

	pubKey, err := ts.manager.PubKey(info.ContractAddress)
	ts.Require().NoError(err)
	ts.Equal(info.PublicKey, pubKey)

	owner, err := ts.manager.Owner(info.ContractAddress)
	ts.Require().NoError(err)
	ts.Equal(ts.admin, owner)

	return info
}

func (ts *AccTestSuite) TestRequestAccess() {
	deployInfo := ts.deploy()
	ts.requestAccess(deployInfo)
}

func (ts *AccTestSuite) requestAccess(deployInfo *DeployInfo) common.Hash {
	//arrange
	events := make(chan *contract.AccessControlAccessRequested)
	subsc, err := ts.manager.WatchAccessRequested(deployInfo.ContractAddress, events)
	ts.Require().NoError(err)
	defer subsc.Unsubscribe()

	//act
	txHash, err := ts.manager.RequestAccess(ts.requester, deployInfo.ContractAddress)
	ts.blockchain.Commit()

	//assert
	ts.Require().NoError(err)
	receipt, err := ts.blockchain.TransactionReceipt(context.Background(), txHash)
	ts.Require().NoError(err)
	ts.Require().NotNil(receipt)
	ts.Equal(uint64(1), receipt.Status)

	acl, err := ts.manager.ACLByAddress(ts.requester, deployInfo.ContractAddress)
	ts.Require().NoError(err)
	ts.True(acl.Pending)
	ts.Empty(acl.EncryptedKey)

	for {
		select {
		case msg := <-events:
			ts.Require().NotNil(msg)
			ts.Equal(ts.requester, msg.From)
			return txHash
		case err := <-subsc.Err():
			ts.Require().NoError(err)
		}
	}
}

func (ts *AccTestSuite) TestGrantAccess() {
	//arrange
	deployInfo := ts.deploy()

	events := make(chan *contract.AccessControlAccessGranted)
	subsc, err := ts.manager.WatchAccessGranted(deployInfo.ContractAddress, events)
	ts.Require().NoError(err)
	defer subsc.Unsubscribe()

	ts.requestAccess(deployInfo)
	attrs := []string{"ceo", "admin"}
	requesterKey := ts.wm.accounts[ts.requester]

	//act
	grantInfo, err := ts.manager.GrantAccess(deployInfo.ContractAddress, ts.admin, ts.requester, attrs, deployInfo.PublicKey, deployInfo.MasterKey, &requesterKey.PublicKey)
	ts.blockchain.Commit()

	//assert
	ts.Require().NoError(err)
	receipt, err := ts.blockchain.TransactionReceipt(context.Background(), grantInfo.TxHash)
	ts.Require().NoError(err)
	ts.Require().NotNil(receipt)
	ts.Equal(uint64(1), receipt.Status)

	acl, err := ts.manager.ACLByAddress(ts.requester, deployInfo.ContractAddress)
	ts.Require().NoError(err)
	ts.False(acl.Pending)

	encryptedKey, err := hexutil.Decode(acl.EncryptedKey)
	ts.Require().NoError(err)
	key, err := utils.DecryptMessage(requesterKey, encryptedKey)
	ts.Require().NoError(err)
	ts.NotNil(key)

	for {
		select {
		case msg := <-events:
			ts.Require().NotNil(msg)
			ts.NotEmpty(msg.Key)
			ts.Equal(ts.requester, msg.Requester)

			var stringAttrs []string
			for _, attr := range msg.Attrs {
				unpaddedAttr := bytes.TrimRight(attr[:], "\x00")
				attrStr := string(unpaddedAttr)
				stringAttrs = append(stringAttrs, attrStr)
			}

			ts.ElementsMatch(attrs, stringAttrs)
			return
		case err := <-subsc.Err():
			ts.Require().NoError(err)
		}
	}
}

func (ts *AccTestSuite) TestAddDevice() {
	deployInfo := ts.deploy()
	ts.addDevice(deployInfo)
}

func (ts *AccTestSuite) addDevice(deployInfo *DeployInfo) common.Address {
	//arrange
	name := "Camera 123"
	policy := "(ceo & admin)"
	_, deviceAddr, err := utils.GenerateKeyPair()
	ts.Require().NoError(err)

	events := make(chan *contract.AccessControlDevicePolicyUpdated)
	subsc, err := ts.manager.WatchDevicePolicyUpdated(deployInfo.ContractAddress, events)
	ts.Require().NoError(err)
	defer subsc.Unsubscribe()

	//act
	tx, err := ts.manager.SetDevicePolicy(name, deviceAddr, deployInfo.ContractAddress, ts.admin, policy)
	ts.blockchain.Commit()

	//assert
	ts.Require().NoError(err)
	receipt, err := ts.blockchain.TransactionReceipt(context.Background(), tx.Hash())
	ts.Require().NoError(err)
	ts.Require().NotNil(receipt)
	ts.Equal(uint64(1), receipt.Status)

	policyEntry, err := ts.manager.DevicePolicyByAddress(deviceAddr, deployInfo.ContractAddress)
	ts.Require().NoError(err)
	ts.Equal(policy, policyEntry.Policy)

	for {
		select {
		case msg := <-events:
			ts.Require().NotNil(msg)
			ts.Equal(deviceAddr, msg.Device)
			ts.Equal(policy, msg.Policy)
			return deviceAddr
		case err := <-subsc.Err():
			ts.Require().NoError(err)
		}
	}
}

func (ts *AccTestSuite) TestRemoveDevice() {
	//arrange
	deployInfo := ts.deploy()
	deviceAddr := ts.addDevice(deployInfo)

	events := make(chan *contract.AccessControlDevicePolicyDeleted)
	subsc, err := ts.manager.WatchDevicePolicyRemoved(deployInfo.ContractAddress, events)
	ts.Require().NoError(err)
	defer subsc.Unsubscribe()

	//act
	tx, err := ts.manager.RemoveDevicePolicy(deviceAddr, deployInfo.ContractAddress, ts.admin)
	ts.blockchain.Commit()

	//assert
	ts.Require().NoError(err)
	receipt, err := ts.blockchain.TransactionReceipt(context.Background(), tx.Hash())
	ts.Require().NoError(err)
	ts.Require().NotNil(receipt)
	ts.Equal(uint64(1), receipt.Status)

	policyEntry, err := ts.manager.DevicePolicyByAddress(deviceAddr, deployInfo.ContractAddress)
	ts.Require().NoError(err)
	ts.Equal("", policyEntry.Policy)

	for {
		select {
		case msg := <-events:
			ts.Require().NotNil(msg)
			ts.Equal(deviceAddr, msg.Device)
			return
		case err := <-subsc.Err():
			ts.Require().NoError(err)
		}
	}
}

func TestAccTestSuite(t *testing.T) {
	suite.Run(t, &AccTestSuite{})
}
