// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ILockProxyMetaData contains all meta data concerning the ILockProxy contract.
var ILockProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"assetHashMap\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"fromAssetHash\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"toAssetHash\",\"type\":\"bytes\"}],\"name\":\"bindAssetHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"toChainId\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"targetProxyHash\",\"type\":\"bytes\"}],\"name\":\"bindProxyHash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"proxyHashMap\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4f7d9808": "assetHashMap(address,uint64)",
		"3348f63b": "bindAssetHash(address,uint64,bytes)",
		"379b98f6": "bindProxyHash(uint64,bytes)",
		"8f32d59b": "isOwner()",
		"8da5cb5b": "owner()",
		"9e5767aa": "proxyHashMap(uint64)",
		"f2fde38b": "transferOwnership(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b5061048f806100206000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80638da5cb5b1161005b5780638da5cb5b146102bd5780638f32d59b146102e15780639e5767aa146102e9578063f2fde38b146103105761007d565b80633348f63b14610082578063379b98f61461015b5780634f7d980814610212575b600080fd5b6101476004803603606081101561009857600080fd5b6001600160a01b038235169167ffffffffffffffff602082013516918101906060810160408201356401000000008111156100d257600080fd5b8201836020820111156100e457600080fd5b8035906020019184600183028401116401000000008311171561010657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610338945050505050565b604080519115158252519081900360200190f35b6101476004803603604081101561017157600080fd5b67ffffffffffffffff823516919081019060408101602082013564010000000081111561019d57600080fd5b8201836020820111156101af57600080fd5b803590602001918460018302840111640100000000831117156101d157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610341945050505050565b6102486004803603604081101561022857600080fd5b5080356001600160a01b0316906020013567ffffffffffffffff16610349565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561028257818101518382015260200161026a565b50505050905090810190601f1680156102af5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102c56103ea565b604080516001600160a01b039092168252519081900360200190f35b6101476103ea565b610248600480360360208110156102ff57600080fd5b503567ffffffffffffffff166103ef565b6103366004803603602081101561032657600080fd5b50356001600160a01b0316610456565b005b60009392505050565b600092915050565b60006020818152928152604080822084529181528190208054825160026001831615610100026000190190921691909104601f8101859004850282018501909352828152929091908301828280156103e25780601f106103b7576101008083540402835291602001916103e2565b820191906000526020600020905b8154815290600101906020018083116103c557829003601f168201915b505050505081565b600090565b60016020818152600092835260409283902080548451600294821615610100026000190190911693909304601f81018390048302840183019094528383529192908301828280156103e25780601f106103b7576101008083540402835291602001916103e2565b5056fea2646970667358221220916eddfa49322b179324d1f923a5303983ea6c56bc2f25de0a3597ebb6b9eb1664736f6c634300060c0033",
}

// ILockProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ILockProxyMetaData.ABI instead.
var ILockProxyABI = ILockProxyMetaData.ABI

// Deprecated: Use ILockProxyMetaData.Sigs instead.
// ILockProxyFuncSigs maps the 4-byte function signature to its string representation.
var ILockProxyFuncSigs = ILockProxyMetaData.Sigs

// ILockProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ILockProxyMetaData.Bin instead.
var ILockProxyBin = ILockProxyMetaData.Bin

// DeployILockProxy deploys a new Ethereum contract, binding an instance of ILockProxy to it.
func DeployILockProxy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ILockProxy, error) {
	parsed, err := ILockProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ILockProxyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ILockProxy{ILockProxyCaller: ILockProxyCaller{contract: contract}, ILockProxyTransactor: ILockProxyTransactor{contract: contract}, ILockProxyFilterer: ILockProxyFilterer{contract: contract}}, nil
}

// ILockProxy is an auto generated Go binding around an Ethereum contract.
type ILockProxy struct {
	ILockProxyCaller     // Read-only binding to the contract
	ILockProxyTransactor // Write-only binding to the contract
	ILockProxyFilterer   // Log filterer for contract events
}

// ILockProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ILockProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILockProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ILockProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILockProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ILockProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILockProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ILockProxySession struct {
	Contract     *ILockProxy       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ILockProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ILockProxyCallerSession struct {
	Contract *ILockProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ILockProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ILockProxyTransactorSession struct {
	Contract     *ILockProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ILockProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ILockProxyRaw struct {
	Contract *ILockProxy // Generic contract binding to access the raw methods on
}

// ILockProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ILockProxyCallerRaw struct {
	Contract *ILockProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ILockProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ILockProxyTransactorRaw struct {
	Contract *ILockProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewILockProxy creates a new instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxy(address common.Address, backend bind.ContractBackend) (*ILockProxy, error) {
	contract, err := bindILockProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ILockProxy{ILockProxyCaller: ILockProxyCaller{contract: contract}, ILockProxyTransactor: ILockProxyTransactor{contract: contract}, ILockProxyFilterer: ILockProxyFilterer{contract: contract}}, nil
}

// NewILockProxyCaller creates a new read-only instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxyCaller(address common.Address, caller bind.ContractCaller) (*ILockProxyCaller, error) {
	contract, err := bindILockProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ILockProxyCaller{contract: contract}, nil
}

// NewILockProxyTransactor creates a new write-only instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ILockProxyTransactor, error) {
	contract, err := bindILockProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ILockProxyTransactor{contract: contract}, nil
}

// NewILockProxyFilterer creates a new log filterer instance of ILockProxy, bound to a specific deployed contract.
func NewILockProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ILockProxyFilterer, error) {
	contract, err := bindILockProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ILockProxyFilterer{contract: contract}, nil
}

// bindILockProxy binds a generic wrapper to an already deployed contract.
func bindILockProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ILockProxyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILockProxy *ILockProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILockProxy.Contract.ILockProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILockProxy *ILockProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILockProxy.Contract.ILockProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILockProxy *ILockProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILockProxy.Contract.ILockProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILockProxy *ILockProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILockProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILockProxy *ILockProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILockProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILockProxy *ILockProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILockProxy.Contract.contract.Transact(opts, method, params...)
}

// AssetHashMap is a free data retrieval call binding the contract method 0x4f7d9808.
//
// Solidity: function assetHashMap(address , uint64 ) view returns(bytes)
func (_ILockProxy *ILockProxyCaller) AssetHashMap(opts *bind.CallOpts, arg0 common.Address, arg1 uint64) ([]byte, error) {
	var out []interface{}
	err := _ILockProxy.contract.Call(opts, &out, "assetHashMap", arg0, arg1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// AssetHashMap is a free data retrieval call binding the contract method 0x4f7d9808.
//
// Solidity: function assetHashMap(address , uint64 ) view returns(bytes)
func (_ILockProxy *ILockProxySession) AssetHashMap(arg0 common.Address, arg1 uint64) ([]byte, error) {
	return _ILockProxy.Contract.AssetHashMap(&_ILockProxy.CallOpts, arg0, arg1)
}

// AssetHashMap is a free data retrieval call binding the contract method 0x4f7d9808.
//
// Solidity: function assetHashMap(address , uint64 ) view returns(bytes)
func (_ILockProxy *ILockProxyCallerSession) AssetHashMap(arg0 common.Address, arg1 uint64) ([]byte, error) {
	return _ILockProxy.Contract.AssetHashMap(&_ILockProxy.CallOpts, arg0, arg1)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ILockProxy *ILockProxyCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ILockProxy.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ILockProxy *ILockProxySession) IsOwner() (bool, error) {
	return _ILockProxy.Contract.IsOwner(&_ILockProxy.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_ILockProxy *ILockProxyCallerSession) IsOwner() (bool, error) {
	return _ILockProxy.Contract.IsOwner(&_ILockProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ILockProxy *ILockProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ILockProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ILockProxy *ILockProxySession) Owner() (common.Address, error) {
	return _ILockProxy.Contract.Owner(&_ILockProxy.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ILockProxy *ILockProxyCallerSession) Owner() (common.Address, error) {
	return _ILockProxy.Contract.Owner(&_ILockProxy.CallOpts)
}

// ProxyHashMap is a free data retrieval call binding the contract method 0x9e5767aa.
//
// Solidity: function proxyHashMap(uint64 ) view returns(bytes)
func (_ILockProxy *ILockProxyCaller) ProxyHashMap(opts *bind.CallOpts, arg0 uint64) ([]byte, error) {
	var out []interface{}
	err := _ILockProxy.contract.Call(opts, &out, "proxyHashMap", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ProxyHashMap is a free data retrieval call binding the contract method 0x9e5767aa.
//
// Solidity: function proxyHashMap(uint64 ) view returns(bytes)
func (_ILockProxy *ILockProxySession) ProxyHashMap(arg0 uint64) ([]byte, error) {
	return _ILockProxy.Contract.ProxyHashMap(&_ILockProxy.CallOpts, arg0)
}

// ProxyHashMap is a free data retrieval call binding the contract method 0x9e5767aa.
//
// Solidity: function proxyHashMap(uint64 ) view returns(bytes)
func (_ILockProxy *ILockProxyCallerSession) ProxyHashMap(arg0 uint64) ([]byte, error) {
	return _ILockProxy.Contract.ProxyHashMap(&_ILockProxy.CallOpts, arg0)
}

// BindAssetHash is a paid mutator transaction binding the contract method 0x3348f63b.
//
// Solidity: function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes toAssetHash) returns(bool)
func (_ILockProxy *ILockProxyTransactor) BindAssetHash(opts *bind.TransactOpts, fromAssetHash common.Address, toChainId uint64, toAssetHash []byte) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "bindAssetHash", fromAssetHash, toChainId, toAssetHash)
}

// BindAssetHash is a paid mutator transaction binding the contract method 0x3348f63b.
//
// Solidity: function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes toAssetHash) returns(bool)
func (_ILockProxy *ILockProxySession) BindAssetHash(fromAssetHash common.Address, toChainId uint64, toAssetHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindAssetHash(&_ILockProxy.TransactOpts, fromAssetHash, toChainId, toAssetHash)
}

// BindAssetHash is a paid mutator transaction binding the contract method 0x3348f63b.
//
// Solidity: function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes toAssetHash) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) BindAssetHash(fromAssetHash common.Address, toChainId uint64, toAssetHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindAssetHash(&_ILockProxy.TransactOpts, fromAssetHash, toChainId, toAssetHash)
}

// BindProxyHash is a paid mutator transaction binding the contract method 0x379b98f6.
//
// Solidity: function bindProxyHash(uint64 toChainId, bytes targetProxyHash) returns(bool)
func (_ILockProxy *ILockProxyTransactor) BindProxyHash(opts *bind.TransactOpts, toChainId uint64, targetProxyHash []byte) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "bindProxyHash", toChainId, targetProxyHash)
}

// BindProxyHash is a paid mutator transaction binding the contract method 0x379b98f6.
//
// Solidity: function bindProxyHash(uint64 toChainId, bytes targetProxyHash) returns(bool)
func (_ILockProxy *ILockProxySession) BindProxyHash(toChainId uint64, targetProxyHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindProxyHash(&_ILockProxy.TransactOpts, toChainId, targetProxyHash)
}

// BindProxyHash is a paid mutator transaction binding the contract method 0x379b98f6.
//
// Solidity: function bindProxyHash(uint64 toChainId, bytes targetProxyHash) returns(bool)
func (_ILockProxy *ILockProxyTransactorSession) BindProxyHash(toChainId uint64, targetProxyHash []byte) (*types.Transaction, error) {
	return _ILockProxy.Contract.BindProxyHash(&_ILockProxy.TransactOpts, toChainId, targetProxyHash)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ILockProxy *ILockProxyTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ILockProxy.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ILockProxy *ILockProxySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ILockProxy.Contract.TransferOwnership(&_ILockProxy.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ILockProxy *ILockProxyTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ILockProxy.Contract.TransferOwnership(&_ILockProxy.TransactOpts, newOwner)
}
