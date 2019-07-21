// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AccessControlABI is the input ABI used to generate the binding from.
const AccessControlABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"acl\",\"outputs\":[{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"pending\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_device\",\"type\":\"address\"}],\"name\":\"removeDevicePolicy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_key\",\"type\":\"string\"},{\"name\":\"_attrs\",\"type\":\"bytes32[]\"}],\"name\":\"grantAccess\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_device\",\"type\":\"address\"},{\"name\":\"_policy\",\"type\":\"string\"}],\"name\":\"setDevicePolicy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pubKey\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"devices\",\"outputs\":[{\"name\":\"policy\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"requestAccess\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_pubKey\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_time\",\"type\":\"uint256\"}],\"name\":\"AccessRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_requester\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_key\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"attrs\",\"type\":\"bytes32[]\"}],\"name\":\"AccessGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_device\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"policy\",\"type\":\"string\"}],\"name\":\"DevicePolicyUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_device\",\"type\":\"address\"}],\"name\":\"DevicePolicyDeleted\",\"type\":\"event\"}]"

// AccessControlBin is the compiled bytecode used for deploying new contracts.
const AccessControlBin = `60806040523480156200001157600080fd5b506040516200119138038062001191833981018060405260208110156200003757600080fd5b8101908080516401000000008111156200005057600080fd5b828101905060208101848111156200006757600080fd5b81518560018202830111640100000000821117156200008557600080fd5b5050929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060019080519060200190620000e5929190620000ed565b50506200019c565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200013057805160ff191683800117855562000161565b8280016001018555821562000161579182015b828111156200016057825182559160200191906001019062000143565b5b50905062000170919062000174565b5090565b6200019991905b80821115620001955760008160009055506001016200017b565b5090565b90565b610fe580620001ac6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80638da5cb5b1161005b5780638da5cb5b146103e3578063ac2a5dfd1461042d578063e7b4cac6146104b0578063eb2f48171461056d57610088565b80632010c0341461008d57806324781dae146101555780634ce65778146101995780635901200f14610308575b600080fd5b6100cf600480360360208110156100a357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610577565b604051808060200183151515158152602001828103825284818151815260200191508051906020019080838360005b838110156101195780820151818401526020810190506100fe565b50505050905090810190601f1680156101465780820380516001836020036101000a031916815260200191505b50935050505060405180910390f35b6101976004803603602081101561016b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610640565b005b610306600480360360608110156101af57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156101ec57600080fd5b8201836020820111156101fe57600080fd5b8035906020019184600183028401116401000000008311171561022057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561028357600080fd5b82018360208201111561029557600080fd5b803590602001918460208302840111640100000000831117156102b757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f82011690508083019250505050505050919291929050505061079a565b005b6103e16004803603604081101561031e57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561035b57600080fd5b82018360208201111561036d57600080fd5b8035906020019184600183028401116401000000008311171561038f57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610a66565b005b6103eb610c36565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610435610c5b565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561047557808201518184015260208101905061045a565b50505050905090810190601f1680156104a25780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6104f2600480360360208110156104c657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610cf9565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610532578082015181840152602081019050610517565b50505050905090810190601f16801561055f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610575610daf565b005b6002602052806000526040600020600091509050806000018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156106235780601f106105f857610100808354040283529160200191610623565b820191906000526020600020905b81548152906001019060200180831161060657829003601f168201915b5050505050908060010160009054906101000a900460ff16905082565b3373ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610702576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b600360008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600080820160006107529190610e5a565b50508073ffffffffffffffffffffffffffffffffffffffff167f3a48ef8a388b1e27a6e19b5a85532fb13d679d686486eecc9eed656a2222be7a60405160405180910390a250565b3373ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461085c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b81600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000190805190602001906108b2929190610ea2565b5080600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206002019080519060200190610909929190610f22565b506000600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548160ff0219169083151502179055508273ffffffffffffffffffffffffffffffffffffffff167fd682788c6c3dccd0663083b0dd0a64e2d58a5dad264fc64bbfccc2425dcbed538383604051808060200180602001838103835285818151815260200191508051906020019080838360005b838110156109e35780820151818401526020810190506109c8565b50505050905090810190601f168015610a105780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019060200280838360005b83811015610a4c578082015181840152602081019050610a31565b5050505090500194505050505060405180910390a2505050565b3373ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610b28576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b80600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000019080519060200190610b7e929190610ea2565b508173ffffffffffffffffffffffffffffffffffffffff167f81d8162d314cded5fb1930e8dd0a90f29252cdb994add44a133f7a64485d980f826040518080602001828103825283818151815260200191508051906020019080838360005b83811015610bf8578082015181840152602081019050610bdd565b50505050905090810190601f168015610c255780820380516001836020036101000a031916815260200191505b509250505060405180910390a25050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610cf15780601f10610cc657610100808354040283529160200191610cf1565b820191906000526020600020905b815481529060010190602001808311610cd457829003601f168201915b505050505081565b6003602052806000526040600020600091509050806000018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610da55780601f10610d7a57610100808354040283529160200191610da5565b820191906000526020600020905b815481529060010190602001808311610d8857829003601f168201915b5050505050905081565b6001600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160006101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167f383cc34b0f756994f1690edba5293974b3590b2f802f0eec957436ab88218a0f426040518082815260200191505060405180910390a2565b50805460018160011615610100020316600290046000825580601f10610e805750610e9f565b601f016020900490600052602060002090810190610e9e9190610f6f565b5b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610ee357805160ff1916838001178555610f11565b82800160010185558215610f11579182015b82811115610f10578251825591602001919060010190610ef5565b5b509050610f1e9190610f6f565b5090565b828054828255906000526020600020908101928215610f5e579160200282015b82811115610f5d578251825591602001919060010190610f42565b5b509050610f6b9190610f94565b5090565b610f9191905b80821115610f8d576000816000905550600101610f75565b5090565b90565b610fb691905b80821115610fb2576000816000905550600101610f9a565b5090565b9056fea165627a7a7230582068cc28d97ec533fe8c2f1e3bef7b2cae0fee04ee146c2d4dc16eb6cee756d3f10029`

// DeployAccessControl deploys a new Ethereum contract, binding an instance of AccessControl to it.
func DeployAccessControl(auth *bind.TransactOpts, backend bind.ContractBackend, _pubKey string) (common.Address, *types.Transaction, *AccessControl, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccessControlBin), backend, _pubKey)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessControl{AccessControlCaller: AccessControlCaller{contract: contract}, AccessControlTransactor: AccessControlTransactor{contract: contract}, AccessControlFilterer: AccessControlFilterer{contract: contract}}, nil
}

// AccessControl is an auto generated Go binding around an Ethereum contract.
type AccessControl struct {
	AccessControlCaller     // Read-only binding to the contract
	AccessControlTransactor // Write-only binding to the contract
	AccessControlFilterer   // Log filterer for contract events
}

// AccessControlCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessControlCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessControlTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessControlFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessControlSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessControlSession struct {
	Contract     *AccessControl    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccessControlCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessControlCallerSession struct {
	Contract *AccessControlCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AccessControlTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessControlTransactorSession struct {
	Contract     *AccessControlTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AccessControlRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessControlRaw struct {
	Contract *AccessControl // Generic contract binding to access the raw methods on
}

// AccessControlCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessControlCallerRaw struct {
	Contract *AccessControlCaller // Generic read-only contract binding to access the raw methods on
}

// AccessControlTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessControlTransactorRaw struct {
	Contract *AccessControlTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessControl creates a new instance of AccessControl, bound to a specific deployed contract.
func NewAccessControl(address common.Address, backend bind.ContractBackend) (*AccessControl, error) {
	contract, err := bindAccessControl(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessControl{AccessControlCaller: AccessControlCaller{contract: contract}, AccessControlTransactor: AccessControlTransactor{contract: contract}, AccessControlFilterer: AccessControlFilterer{contract: contract}}, nil
}

// NewAccessControlCaller creates a new read-only instance of AccessControl, bound to a specific deployed contract.
func NewAccessControlCaller(address common.Address, caller bind.ContractCaller) (*AccessControlCaller, error) {
	contract, err := bindAccessControl(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlCaller{contract: contract}, nil
}

// NewAccessControlTransactor creates a new write-only instance of AccessControl, bound to a specific deployed contract.
func NewAccessControlTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessControlTransactor, error) {
	contract, err := bindAccessControl(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessControlTransactor{contract: contract}, nil
}

// NewAccessControlFilterer creates a new log filterer instance of AccessControl, bound to a specific deployed contract.
func NewAccessControlFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessControlFilterer, error) {
	contract, err := bindAccessControl(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessControlFilterer{contract: contract}, nil
}

// bindAccessControl binds a generic wrapper to an already deployed contract.
func bindAccessControl(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessControlABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControl *AccessControlRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControl.Contract.AccessControlCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControl *AccessControlRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControl.Contract.AccessControlTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControl *AccessControlRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControl.Contract.AccessControlTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessControl *AccessControlCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccessControl.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessControl *AccessControlTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControl.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessControl *AccessControlTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessControl.Contract.contract.Transact(opts, method, params...)
}

// Acl is a free data retrieval call binding the contract method 0x2010c034.
//
// Solidity: function acl( address) constant returns(key string, pending bool)
func (_AccessControl *AccessControlCaller) Acl(opts *bind.CallOpts, arg0 common.Address) (struct {
	Key     string
	Pending bool
}, error) {
	ret := new(struct {
		Key     string
		Pending bool
	})
	out := ret
	err := _AccessControl.contract.Call(opts, out, "acl", arg0)
	return *ret, err
}

// Acl is a free data retrieval call binding the contract method 0x2010c034.
//
// Solidity: function acl( address) constant returns(key string, pending bool)
func (_AccessControl *AccessControlSession) Acl(arg0 common.Address) (struct {
	Key     string
	Pending bool
}, error) {
	return _AccessControl.Contract.Acl(&_AccessControl.CallOpts, arg0)
}

// Acl is a free data retrieval call binding the contract method 0x2010c034.
//
// Solidity: function acl( address) constant returns(key string, pending bool)
func (_AccessControl *AccessControlCallerSession) Acl(arg0 common.Address) (struct {
	Key     string
	Pending bool
}, error) {
	return _AccessControl.Contract.Acl(&_AccessControl.CallOpts, arg0)
}

// Devices is a free data retrieval call binding the contract method 0xe7b4cac6.
//
// Solidity: function devices( address) constant returns(policy string)
func (_AccessControl *AccessControlCaller) Devices(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AccessControl.contract.Call(opts, out, "devices", arg0)
	return *ret0, err
}

// Devices is a free data retrieval call binding the contract method 0xe7b4cac6.
//
// Solidity: function devices( address) constant returns(policy string)
func (_AccessControl *AccessControlSession) Devices(arg0 common.Address) (string, error) {
	return _AccessControl.Contract.Devices(&_AccessControl.CallOpts, arg0)
}

// Devices is a free data retrieval call binding the contract method 0xe7b4cac6.
//
// Solidity: function devices( address) constant returns(policy string)
func (_AccessControl *AccessControlCallerSession) Devices(arg0 common.Address) (string, error) {
	return _AccessControl.Contract.Devices(&_AccessControl.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccessControl *AccessControlCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AccessControl.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccessControl *AccessControlSession) Owner() (common.Address, error) {
	return _AccessControl.Contract.Owner(&_AccessControl.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccessControl *AccessControlCallerSession) Owner() (common.Address, error) {
	return _AccessControl.Contract.Owner(&_AccessControl.CallOpts)
}

// PubKey is a free data retrieval call binding the contract method 0xac2a5dfd.
//
// Solidity: function pubKey() constant returns(string)
func (_AccessControl *AccessControlCaller) PubKey(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AccessControl.contract.Call(opts, out, "pubKey")
	return *ret0, err
}

// PubKey is a free data retrieval call binding the contract method 0xac2a5dfd.
//
// Solidity: function pubKey() constant returns(string)
func (_AccessControl *AccessControlSession) PubKey() (string, error) {
	return _AccessControl.Contract.PubKey(&_AccessControl.CallOpts)
}

// PubKey is a free data retrieval call binding the contract method 0xac2a5dfd.
//
// Solidity: function pubKey() constant returns(string)
func (_AccessControl *AccessControlCallerSession) PubKey() (string, error) {
	return _AccessControl.Contract.PubKey(&_AccessControl.CallOpts)
}

// GrantAccess is a paid mutator transaction binding the contract method 0x4ce65778.
//
// Solidity: function grantAccess(_user address, _key string, _attrs bytes32[]) returns()
func (_AccessControl *AccessControlTransactor) GrantAccess(opts *bind.TransactOpts, _user common.Address, _key string, _attrs [][32]byte) (*types.Transaction, error) {
	return _AccessControl.contract.Transact(opts, "grantAccess", _user, _key, _attrs)
}

// GrantAccess is a paid mutator transaction binding the contract method 0x4ce65778.
//
// Solidity: function grantAccess(_user address, _key string, _attrs bytes32[]) returns()
func (_AccessControl *AccessControlSession) GrantAccess(_user common.Address, _key string, _attrs [][32]byte) (*types.Transaction, error) {
	return _AccessControl.Contract.GrantAccess(&_AccessControl.TransactOpts, _user, _key, _attrs)
}

// GrantAccess is a paid mutator transaction binding the contract method 0x4ce65778.
//
// Solidity: function grantAccess(_user address, _key string, _attrs bytes32[]) returns()
func (_AccessControl *AccessControlTransactorSession) GrantAccess(_user common.Address, _key string, _attrs [][32]byte) (*types.Transaction, error) {
	return _AccessControl.Contract.GrantAccess(&_AccessControl.TransactOpts, _user, _key, _attrs)
}

// RemoveDevicePolicy is a paid mutator transaction binding the contract method 0x24781dae.
//
// Solidity: function removeDevicePolicy(_device address) returns()
func (_AccessControl *AccessControlTransactor) RemoveDevicePolicy(opts *bind.TransactOpts, _device common.Address) (*types.Transaction, error) {
	return _AccessControl.contract.Transact(opts, "removeDevicePolicy", _device)
}

// RemoveDevicePolicy is a paid mutator transaction binding the contract method 0x24781dae.
//
// Solidity: function removeDevicePolicy(_device address) returns()
func (_AccessControl *AccessControlSession) RemoveDevicePolicy(_device common.Address) (*types.Transaction, error) {
	return _AccessControl.Contract.RemoveDevicePolicy(&_AccessControl.TransactOpts, _device)
}

// RemoveDevicePolicy is a paid mutator transaction binding the contract method 0x24781dae.
//
// Solidity: function removeDevicePolicy(_device address) returns()
func (_AccessControl *AccessControlTransactorSession) RemoveDevicePolicy(_device common.Address) (*types.Transaction, error) {
	return _AccessControl.Contract.RemoveDevicePolicy(&_AccessControl.TransactOpts, _device)
}

// RequestAccess is a paid mutator transaction binding the contract method 0xeb2f4817.
//
// Solidity: function requestAccess() returns()
func (_AccessControl *AccessControlTransactor) RequestAccess(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessControl.contract.Transact(opts, "requestAccess")
}

// RequestAccess is a paid mutator transaction binding the contract method 0xeb2f4817.
//
// Solidity: function requestAccess() returns()
func (_AccessControl *AccessControlSession) RequestAccess() (*types.Transaction, error) {
	return _AccessControl.Contract.RequestAccess(&_AccessControl.TransactOpts)
}

// RequestAccess is a paid mutator transaction binding the contract method 0xeb2f4817.
//
// Solidity: function requestAccess() returns()
func (_AccessControl *AccessControlTransactorSession) RequestAccess() (*types.Transaction, error) {
	return _AccessControl.Contract.RequestAccess(&_AccessControl.TransactOpts)
}

// SetDevicePolicy is a paid mutator transaction binding the contract method 0x5901200f.
//
// Solidity: function setDevicePolicy(_device address, _policy string) returns()
func (_AccessControl *AccessControlTransactor) SetDevicePolicy(opts *bind.TransactOpts, _device common.Address, _policy string) (*types.Transaction, error) {
	return _AccessControl.contract.Transact(opts, "setDevicePolicy", _device, _policy)
}

// SetDevicePolicy is a paid mutator transaction binding the contract method 0x5901200f.
//
// Solidity: function setDevicePolicy(_device address, _policy string) returns()
func (_AccessControl *AccessControlSession) SetDevicePolicy(_device common.Address, _policy string) (*types.Transaction, error) {
	return _AccessControl.Contract.SetDevicePolicy(&_AccessControl.TransactOpts, _device, _policy)
}

// SetDevicePolicy is a paid mutator transaction binding the contract method 0x5901200f.
//
// Solidity: function setDevicePolicy(_device address, _policy string) returns()
func (_AccessControl *AccessControlTransactorSession) SetDevicePolicy(_device common.Address, _policy string) (*types.Transaction, error) {
	return _AccessControl.Contract.SetDevicePolicy(&_AccessControl.TransactOpts, _device, _policy)
}

// AccessControlAccessGrantedIterator is returned from FilterAccessGranted and is used to iterate over the raw logs and unpacked data for AccessGranted events raised by the AccessControl contract.
type AccessControlAccessGrantedIterator struct {
	Event *AccessControlAccessGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlAccessGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlAccessGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlAccessGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlAccessGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlAccessGranted represents a AccessGranted event raised by the AccessControl contract.
type AccessControlAccessGranted struct {
	Requester common.Address
	Key       string
	Attrs     [][32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAccessGranted is a free log retrieval operation binding the contract event 0xd682788c6c3dccd0663083b0dd0a64e2d58a5dad264fc64bbfccc2425dcbed53.
//
// Solidity: e AccessGranted(_requester indexed address, _key string, attrs bytes32[])
func (_AccessControl *AccessControlFilterer) FilterAccessGranted(opts *bind.FilterOpts, _requester []common.Address) (*AccessControlAccessGrantedIterator, error) {

	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}

	logs, sub, err := _AccessControl.contract.FilterLogs(opts, "AccessGranted", _requesterRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlAccessGrantedIterator{contract: _AccessControl.contract, event: "AccessGranted", logs: logs, sub: sub}, nil
}

// WatchAccessGranted is a free log subscription operation binding the contract event 0xd682788c6c3dccd0663083b0dd0a64e2d58a5dad264fc64bbfccc2425dcbed53.
//
// Solidity: e AccessGranted(_requester indexed address, _key string, attrs bytes32[])
func (_AccessControl *AccessControlFilterer) WatchAccessGranted(opts *bind.WatchOpts, sink chan<- *AccessControlAccessGranted, _requester []common.Address) (event.Subscription, error) {

	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}

	logs, sub, err := _AccessControl.contract.WatchLogs(opts, "AccessGranted", _requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlAccessGranted)
				if err := _AccessControl.contract.UnpackLog(event, "AccessGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AccessControlAccessRequestedIterator is returned from FilterAccessRequested and is used to iterate over the raw logs and unpacked data for AccessRequested events raised by the AccessControl contract.
type AccessControlAccessRequestedIterator struct {
	Event *AccessControlAccessRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlAccessRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlAccessRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlAccessRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlAccessRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlAccessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlAccessRequested represents a AccessRequested event raised by the AccessControl contract.
type AccessControlAccessRequested struct {
	From common.Address
	Time *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAccessRequested is a free log retrieval operation binding the contract event 0x383cc34b0f756994f1690edba5293974b3590b2f802f0eec957436ab88218a0f.
//
// Solidity: e AccessRequested(_from indexed address, _time uint256)
func (_AccessControl *AccessControlFilterer) FilterAccessRequested(opts *bind.FilterOpts, _from []common.Address) (*AccessControlAccessRequestedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _AccessControl.contract.FilterLogs(opts, "AccessRequested", _fromRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlAccessRequestedIterator{contract: _AccessControl.contract, event: "AccessRequested", logs: logs, sub: sub}, nil
}

// WatchAccessRequested is a free log subscription operation binding the contract event 0x383cc34b0f756994f1690edba5293974b3590b2f802f0eec957436ab88218a0f.
//
// Solidity: e AccessRequested(_from indexed address, _time uint256)
func (_AccessControl *AccessControlFilterer) WatchAccessRequested(opts *bind.WatchOpts, sink chan<- *AccessControlAccessRequested, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _AccessControl.contract.WatchLogs(opts, "AccessRequested", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlAccessRequested)
				if err := _AccessControl.contract.UnpackLog(event, "AccessRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AccessControlDevicePolicyDeletedIterator is returned from FilterDevicePolicyDeleted and is used to iterate over the raw logs and unpacked data for DevicePolicyDeleted events raised by the AccessControl contract.
type AccessControlDevicePolicyDeletedIterator struct {
	Event *AccessControlDevicePolicyDeleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlDevicePolicyDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlDevicePolicyDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlDevicePolicyDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlDevicePolicyDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlDevicePolicyDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlDevicePolicyDeleted represents a DevicePolicyDeleted event raised by the AccessControl contract.
type AccessControlDevicePolicyDeleted struct {
	Device common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDevicePolicyDeleted is a free log retrieval operation binding the contract event 0x3a48ef8a388b1e27a6e19b5a85532fb13d679d686486eecc9eed656a2222be7a.
//
// Solidity: e DevicePolicyDeleted(_device indexed address)
func (_AccessControl *AccessControlFilterer) FilterDevicePolicyDeleted(opts *bind.FilterOpts, _device []common.Address) (*AccessControlDevicePolicyDeletedIterator, error) {

	var _deviceRule []interface{}
	for _, _deviceItem := range _device {
		_deviceRule = append(_deviceRule, _deviceItem)
	}

	logs, sub, err := _AccessControl.contract.FilterLogs(opts, "DevicePolicyDeleted", _deviceRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlDevicePolicyDeletedIterator{contract: _AccessControl.contract, event: "DevicePolicyDeleted", logs: logs, sub: sub}, nil
}

// WatchDevicePolicyDeleted is a free log subscription operation binding the contract event 0x3a48ef8a388b1e27a6e19b5a85532fb13d679d686486eecc9eed656a2222be7a.
//
// Solidity: e DevicePolicyDeleted(_device indexed address)
func (_AccessControl *AccessControlFilterer) WatchDevicePolicyDeleted(opts *bind.WatchOpts, sink chan<- *AccessControlDevicePolicyDeleted, _device []common.Address) (event.Subscription, error) {

	var _deviceRule []interface{}
	for _, _deviceItem := range _device {
		_deviceRule = append(_deviceRule, _deviceItem)
	}

	logs, sub, err := _AccessControl.contract.WatchLogs(opts, "DevicePolicyDeleted", _deviceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlDevicePolicyDeleted)
				if err := _AccessControl.contract.UnpackLog(event, "DevicePolicyDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AccessControlDevicePolicyUpdatedIterator is returned from FilterDevicePolicyUpdated and is used to iterate over the raw logs and unpacked data for DevicePolicyUpdated events raised by the AccessControl contract.
type AccessControlDevicePolicyUpdatedIterator struct {
	Event *AccessControlDevicePolicyUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccessControlDevicePolicyUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccessControlDevicePolicyUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccessControlDevicePolicyUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccessControlDevicePolicyUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccessControlDevicePolicyUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccessControlDevicePolicyUpdated represents a DevicePolicyUpdated event raised by the AccessControl contract.
type AccessControlDevicePolicyUpdated struct {
	Device common.Address
	Policy string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDevicePolicyUpdated is a free log retrieval operation binding the contract event 0x81d8162d314cded5fb1930e8dd0a90f29252cdb994add44a133f7a64485d980f.
//
// Solidity: e DevicePolicyUpdated(_device indexed address, policy string)
func (_AccessControl *AccessControlFilterer) FilterDevicePolicyUpdated(opts *bind.FilterOpts, _device []common.Address) (*AccessControlDevicePolicyUpdatedIterator, error) {

	var _deviceRule []interface{}
	for _, _deviceItem := range _device {
		_deviceRule = append(_deviceRule, _deviceItem)
	}

	logs, sub, err := _AccessControl.contract.FilterLogs(opts, "DevicePolicyUpdated", _deviceRule)
	if err != nil {
		return nil, err
	}
	return &AccessControlDevicePolicyUpdatedIterator{contract: _AccessControl.contract, event: "DevicePolicyUpdated", logs: logs, sub: sub}, nil
}

// WatchDevicePolicyUpdated is a free log subscription operation binding the contract event 0x81d8162d314cded5fb1930e8dd0a90f29252cdb994add44a133f7a64485d980f.
//
// Solidity: e DevicePolicyUpdated(_device indexed address, policy string)
func (_AccessControl *AccessControlFilterer) WatchDevicePolicyUpdated(opts *bind.WatchOpts, sink chan<- *AccessControlDevicePolicyUpdated, _device []common.Address) (event.Subscription, error) {

	var _deviceRule []interface{}
	for _, _deviceItem := range _device {
		_deviceRule = append(_deviceRule, _deviceItem)
	}

	logs, sub, err := _AccessControl.contract.WatchLogs(opts, "DevicePolicyUpdated", _deviceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccessControlDevicePolicyUpdated)
				if err := _AccessControl.contract.UnpackLog(event, "DevicePolicyUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
