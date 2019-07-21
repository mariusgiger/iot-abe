package acc

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/mariusgiger/iot-abe/pkg/core"

	"github.com/mariusgiger/iot-abe/pkg/utils"

	"github.com/ethereum/go-ethereum/event"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mariusgiger/iot-abe/pkg/contract"
	"github.com/mariusgiger/iot-abe/pkg/crypto"
	"github.com/mariusgiger/iot-abe/pkg/rpc"
	"github.com/mariusgiger/iot-abe/pkg/wallet"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	chainID = big.NewInt(3) //Ropsten
)

// Manager manages access and lifetime of the smart contract
type Manager struct {
	log *logrus.Logger
	wm  wallet.Manager
	rpc rpc.Client
}

// NewManager returns a new Manager
func NewManager(log *logrus.Logger, wm wallet.Manager, rpc rpc.Client) (*Manager, error) {
	return &Manager{
		log: log,
		wm:  wm,
		rpc: rpc,
	}, nil
}

//DeployInfo holds deployment information
type DeployInfo struct {
	ContractAddress common.Address
	ContractTxHash  common.Hash
	PublicKey       string
	MasterKey       string
}

// Deploy deploys the access control contract
func (m *Manager) Deploy(from common.Address) (*DeployInfo, error) {
	pubKeyBytes, masterKeyBytes, err := crypto.Setup()
	if err != nil {
		return nil, err
	}

	pubKey := hexutil.Encode(pubKeyBytes)
	tOps, err := m.transactOps(from, common.Address{}, "", pubKey)
	if err != nil {
		return nil, err
	}

	// deploy contract
	address, tx, _, err := contract.DeployAccessControl(tOps, m.rpc.GetRawClient(), pubKey)
	if err != nil {
		return nil, err
	}

	return &DeployInfo{
		ContractAddress: address,
		ContractTxHash:  tx.Hash(),
		PublicKey:       pubKey,
		MasterKey:       hexutil.Encode(masterKeyBytes),
	}, nil
}

//PubKey retrieves the public key from an access control contract
func (m *Manager) PubKey(contractAddr common.Address) (string, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return "", errors.Wrap(err, "could not bind to existing contract")
	}

	pubKey, err := instance.PubKey(nil)
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve public key")
	}

	return pubKey, nil
}

//Owner retrieves the owner from an access control contract
func (m *Manager) Owner(contractAddr common.Address) (common.Address, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return common.Address{}, errors.Wrap(err, "could not bind to existing contract")
	}

	owner, err := instance.Owner(nil)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "could not retrieve owner")
	}

	return owner, nil
}

//ACLEntry represents an entry in the acl list
type ACLEntry struct {
	Pending      bool
	EncryptedKey string
}

// GetAccessGrant returns the decrypted secret key for a given address
func (m *Manager) GetAccessGrant(subject, contractAddr common.Address) (string, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return "", errors.Wrap(err, "could not bind to existing contract")
	}

	request, err := instance.Acl(nil, subject)
	if err != nil {
		return "", errors.Wrap(err, "could not retrieve acl entry")
	}

	if request.Key == "" {
		return "", errors.Wrapf(err, "no acl entry present for: %v", subject.Hex())
	}

	encryptedKeyBytes, err := hexutil.Decode(request.Key)
	if err != nil {
		return "", errors.Wrap(err, "could not decode encrypte key")
	}

	keyBytes, err := m.wm.DecryptMessage(subject, encryptedKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "could not dump wallet priv key")
	}

	return hexutil.Encode(keyBytes), nil
}

// ACLByAddress returns the acl entry for a given address
func (m *Manager) ACLByAddress(subject, contractAddr common.Address) (*ACLEntry, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	request, err := instance.Acl(nil, subject)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve acl entry")
	}

	//TODO why is attr list not present
	return &ACLEntry{
		Pending:      request.Pending,
		EncryptedKey: request.Key,
	}, nil
}

//DevicePolicyEntry represents an entry in the device policy list
type DevicePolicyEntry struct {
	Policy string
}

// DevicePolicyByAddress returns the device policy entry for a given device address
func (m *Manager) DevicePolicyByAddress(deviceAddr, contractAddr common.Address) (*DevicePolicyEntry, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	request, err := instance.Devices(nil, deviceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve device policy entry")
	}

	return &DevicePolicyEntry{
		Policy: request,
	}, nil
}

// DevicePolicies returns all created device policies for a smart contract
func (m *Manager) DevicePolicies(contractAddr common.Address) ([]*DevicePolicyTx, error) {
	txs, err := m.rpc.TransactionsByAddress(contractAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve txs for contract")
	}

	abi, err := abi.JSON(strings.NewReader(contract.AccessControlABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read contract ABI")
	}

	var policies []*DevicePolicyTx
	for _, tx := range txs {
		if tx.To() == nil {
			continue //contract creation
		}

		txData := tx.Data()
		method, err := abi.MethodById(txData[:4]) // first 4 bytes contain the id
		if err != nil {
			return nil, errors.Wrap(err, "failed to extract contract method from tx")
		}

		if method.Name != "setDevicePolicy" {
			continue
		}

		var policy struct {
			Device common.Address
			Policy string
		}
		err = method.Inputs.Unpack(&policy, txData[4:])
		if err != nil {
			return nil, errors.Wrap(err, "failed to unpack device address")
		}

		from, err := utils.SenderFromTx(tx, chainID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to extract sender from tx")
		}

		policies = append(policies, &DevicePolicyTx{
			From:   from,
			TxHash: tx.Hash(),
			Tx:     tx,
			Policy: policy.Policy,
			Device: policy.Device,
		})
	}

	return policies, nil
}

//DevicePolicyTx hold information about device policy updates
type DevicePolicyTx struct {
	From   common.Address
	TxHash common.Hash
	Tx     *types.Transaction
	Policy string
	Device common.Address
}

//AccessRequest hold information about access requests
type AccessRequest struct {
	From   common.Address
	TxHash common.Hash
	Tx     *types.Transaction
}

// AccessRequests retrieves all access requests from a smart contract
func (m *Manager) AccessRequests(contractAddr common.Address) ([]*AccessRequest, error) {
	txs, err := m.rpc.TransactionsByAddress(contractAddr)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve txs for contract")
	}

	abi, err := abi.JSON(strings.NewReader(contract.AccessControlABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read contract ABI")
	}

	var requests []*AccessRequest
	for _, tx := range txs {
		if tx.To() == nil {
			continue //contract creation
		}

		txData := tx.Data()
		method, err := abi.MethodById(txData[:4]) // first 4 bytes contain the id
		if err != nil {
			return nil, errors.Wrap(err, "failed to extract contract method from tx")
		}

		if method.Name != "requestAccess" {
			continue
		}

		from, err := utils.SenderFromTx(tx, chainID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to extract sender from tx")
		}

		requests = append(requests, &AccessRequest{
			From:   from,
			TxHash: tx.Hash(),
			Tx:     tx,
		})
	}

	return requests, nil
}

//RequestAccess executes the request access request method
func (m *Manager) RequestAccess(requester, contractAddr common.Address) (common.Hash, error) {
	tOps, err := m.transactOps(requester, contractAddr, "requestAccess")
	if err != nil {
		return common.Hash{}, err
	}

	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "could not bind to existing contract")
	}

	tx, err := instance.AccessControlTransactor.RequestAccess(tOps)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "executing request access method failed")
	}

	return tx.Hash(), nil
}

// WatchAccessRequested allows watching access requested events
func (m *Manager) WatchAccessRequested(contractAddr common.Address, events chan<- *contract.AccessControlAccessRequested) (event.Subscription, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	subscr, err := instance.WatchAccessRequested(nil, events, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not subscribe to events")
	}

	return subscr, nil
}

// WatchAccessGranted allows watching access granted events
func (m *Manager) WatchAccessGranted(contractAddr common.Address, events chan<- *contract.AccessControlAccessGranted) (event.Subscription, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	subscr, err := instance.WatchAccessGranted(nil, events, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not subscribe to events")
	}

	return subscr, nil
}

//GrantInfo holds grant information
type GrantInfo struct {
	TxHash             common.Hash
	SecretKey          []byte
	EncryptedSecretKey []byte
	Attributes         []string
}

//GrantAccess grants access for the provided subject by creating a new private key with the given set of attributes
func (m *Manager) GrantAccess(contractAddr, owner, subject common.Address, attributes []string, pubKey string, masterKey string, subjectPubKey *ecdsa.PublicKey) (*GrantInfo, error) {
	var compressedAttrs [][32]byte
	for _, attr := range attributes {
		attrBytes := []byte(attr)
		if len(attrBytes) > 32 {
			return nil, errors.New("attrs must be smaller than 32 byte")
		}

		var compressedAttr [32]byte
		copy(compressedAttr[:], attrBytes)
		compressedAttrs = append(compressedAttrs, compressedAttr)
	}

	pubKeyBytes, err := hexutil.Decode(pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode public key")
	}

	masterKeyBytes, err := hexutil.Decode(masterKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode master key")
	}

	prvKey, err := crypto.GenerateKey(pubKeyBytes, masterKeyBytes, attributes)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate private key")
	}

	encryptedPrvKeyBytes, err := utils.EncryptMessage(subjectPubKey, prvKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not encrypt private key")
	}
	encryptedPrvKey := hexutil.Encode(encryptedPrvKeyBytes)

	tOps, err := m.transactOps(owner, contractAddr, "grantAccess")
	if err != nil {
		return nil, err
	}

	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	tx, err := instance.AccessControlTransactor.GrantAccess(tOps, subject, encryptedPrvKey, compressedAttrs)
	if err != nil {
		return nil, errors.Wrap(err, "executing grant access method failed")
	}

	return &GrantInfo{
		TxHash:             tx.Hash(),
		SecretKey:          prvKey,
		EncryptedSecretKey: encryptedPrvKeyBytes,
		Attributes:         attributes,
	}, nil
}

// SetDevicePolicy sets the policy of a device
func (m *Manager) SetDevicePolicy(name string, deviceAddr, contractAddr, ownerAddr common.Address, policy string) (*types.Transaction, error) {
	tOps, err := m.transactOps(ownerAddr, contractAddr, "setDevicePolicy")
	if err != nil {
		return nil, err
	}

	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	tx, err := instance.AccessControlTransactor.SetDevicePolicy(tOps, deviceAddr, policy)
	if err != nil {
		return nil, errors.Wrap(err, "executing setDevicePolicy method failed")
	}

	return tx, nil
}

// WatchDevicePolicyUpdated allows watching device policy updated events
func (m *Manager) WatchDevicePolicyUpdated(contractAddr common.Address, events chan<- *contract.AccessControlDevicePolicyUpdated) (event.Subscription, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	subscr, err := instance.WatchDevicePolicyUpdated(nil, events, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not subscribe to events")
	}

	return subscr, nil
}

//RemoveDevicePolicy removes a device policy from the contract
func (m *Manager) RemoveDevicePolicy(deviceAddr, contractAddr, ownerAddr common.Address) (*types.Transaction, error) {
	tOps, err := m.transactOps(ownerAddr, contractAddr, "removeDevicePolicy")
	if err != nil {
		return nil, err
	}

	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	tx, err := instance.AccessControlTransactor.RemoveDevicePolicy(tOps, deviceAddr)
	if err != nil {
		return nil, errors.Wrap(err, "executing removeDevicePolicy method failed")
	}

	return tx, nil
}

// WatchDevicePolicyRemoved allows watching device policy deleted events
func (m *Manager) WatchDevicePolicyRemoved(contractAddr common.Address, events chan<- *contract.AccessControlDevicePolicyDeleted) (event.Subscription, error) {
	instance, err := contract.NewAccessControl(contractAddr, m.rpc.GetRawClient())
	if err != nil {
		return nil, errors.Wrap(err, "could not bind to existing contract")
	}

	subscr, err := instance.WatchDevicePolicyDeleted(nil, events, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not subscribe to events")
	}

	return subscr, nil
}

func (m *Manager) transactOps(from, contractAddr common.Address, method string, args ...interface{}) (*bind.TransactOpts, error) {
	context := context.Background()
	signer := func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return m.wm.SignTx(signer, &address, tx, chainID)
	}

	//TODO this does not work properly
	// gasLimit, err := m.rpc.EstimateGas(method, from, contractAddr, args)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not get gas limit")
	// }
	// fmt.Println(gasLimit)

	//TODO hardcoded for reproducability
	// gasPrice, err := m.rpc.SuggestGasPrice()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not get gas price")
	// }
	gasPrice := big.NewInt(2 * core.GWei)

	nonce, err := m.rpc.PendingNonceAt(from)
	if err != nil {
		return nil, errors.Wrap(err, "could not get nonce")
	}

	return &bind.TransactOpts{
		From:     from,
		Nonce:    big.NewInt(0).SetUint64(nonce),
		Signer:   signer,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
		//GasLimit: gasLimit,
		Context: context,
	}, nil
}
