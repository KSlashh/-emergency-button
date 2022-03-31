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

// BusinessContractMetaData contains all meta data concerning the BusinessContract contract.
var BusinessContractMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"managerProxyContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ethCCMProxyAddr\",\"type\":\"address\"}],\"name\":\"setManagerProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"d798f881": "managerProxyContract()",
		"8da5cb5b": "owner()",
		"af9980f0": "setManagerProxy(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b5060e48061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80638da5cb5b146041578063af9980f0146063578063d798f881146088575b600080fd5b6047608e565b604080516001600160a01b039092168252519081900360200190f35b608660048036036020811015607757600080fd5b50356001600160a01b0316609d565b005b604760a0565b6000546001600160a01b031681565b50565b6001546001600160a01b03168156fea265627a7a7231582031f947c8729585d5ecef8b2637ece1d61b37b2981cb74ce3c9bcd8094482aadd64736f6c63430005110032",
}

// BusinessContractABI is the input ABI used to generate the binding from.
// Deprecated: Use BusinessContractMetaData.ABI instead.
var BusinessContractABI = BusinessContractMetaData.ABI

// Deprecated: Use BusinessContractMetaData.Sigs instead.
// BusinessContractFuncSigs maps the 4-byte function signature to its string representation.
var BusinessContractFuncSigs = BusinessContractMetaData.Sigs

// BusinessContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BusinessContractMetaData.Bin instead.
var BusinessContractBin = BusinessContractMetaData.Bin

// DeployBusinessContract deploys a new Ethereum contract, binding an instance of BusinessContract to it.
func DeployBusinessContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BusinessContract, error) {
	parsed, err := BusinessContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BusinessContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BusinessContract{BusinessContractCaller: BusinessContractCaller{contract: contract}, BusinessContractTransactor: BusinessContractTransactor{contract: contract}, BusinessContractFilterer: BusinessContractFilterer{contract: contract}}, nil
}

// BusinessContract is an auto generated Go binding around an Ethereum contract.
type BusinessContract struct {
	BusinessContractCaller     // Read-only binding to the contract
	BusinessContractTransactor // Write-only binding to the contract
	BusinessContractFilterer   // Log filterer for contract events
}

// BusinessContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type BusinessContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BusinessContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BusinessContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BusinessContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BusinessContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BusinessContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BusinessContractSession struct {
	Contract     *BusinessContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BusinessContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BusinessContractCallerSession struct {
	Contract *BusinessContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// BusinessContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BusinessContractTransactorSession struct {
	Contract     *BusinessContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// BusinessContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type BusinessContractRaw struct {
	Contract *BusinessContract // Generic contract binding to access the raw methods on
}

// BusinessContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BusinessContractCallerRaw struct {
	Contract *BusinessContractCaller // Generic read-only contract binding to access the raw methods on
}

// BusinessContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BusinessContractTransactorRaw struct {
	Contract *BusinessContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBusinessContract creates a new instance of BusinessContract, bound to a specific deployed contract.
func NewBusinessContract(address common.Address, backend bind.ContractBackend) (*BusinessContract, error) {
	contract, err := bindBusinessContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BusinessContract{BusinessContractCaller: BusinessContractCaller{contract: contract}, BusinessContractTransactor: BusinessContractTransactor{contract: contract}, BusinessContractFilterer: BusinessContractFilterer{contract: contract}}, nil
}

// NewBusinessContractCaller creates a new read-only instance of BusinessContract, bound to a specific deployed contract.
func NewBusinessContractCaller(address common.Address, caller bind.ContractCaller) (*BusinessContractCaller, error) {
	contract, err := bindBusinessContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BusinessContractCaller{contract: contract}, nil
}

// NewBusinessContractTransactor creates a new write-only instance of BusinessContract, bound to a specific deployed contract.
func NewBusinessContractTransactor(address common.Address, transactor bind.ContractTransactor) (*BusinessContractTransactor, error) {
	contract, err := bindBusinessContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BusinessContractTransactor{contract: contract}, nil
}

// NewBusinessContractFilterer creates a new log filterer instance of BusinessContract, bound to a specific deployed contract.
func NewBusinessContractFilterer(address common.Address, filterer bind.ContractFilterer) (*BusinessContractFilterer, error) {
	contract, err := bindBusinessContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BusinessContractFilterer{contract: contract}, nil
}

// bindBusinessContract binds a generic wrapper to an already deployed contract.
func bindBusinessContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BusinessContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BusinessContract *BusinessContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BusinessContract.Contract.BusinessContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BusinessContract *BusinessContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BusinessContract.Contract.BusinessContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BusinessContract *BusinessContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BusinessContract.Contract.BusinessContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BusinessContract *BusinessContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BusinessContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BusinessContract *BusinessContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BusinessContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BusinessContract *BusinessContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BusinessContract.Contract.contract.Transact(opts, method, params...)
}

// ManagerProxyContract is a free data retrieval call binding the contract method 0xd798f881.
//
// Solidity: function managerProxyContract() view returns(address)
func (_BusinessContract *BusinessContractCaller) ManagerProxyContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BusinessContract.contract.Call(opts, &out, "managerProxyContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ManagerProxyContract is a free data retrieval call binding the contract method 0xd798f881.
//
// Solidity: function managerProxyContract() view returns(address)
func (_BusinessContract *BusinessContractSession) ManagerProxyContract() (common.Address, error) {
	return _BusinessContract.Contract.ManagerProxyContract(&_BusinessContract.CallOpts)
}

// ManagerProxyContract is a free data retrieval call binding the contract method 0xd798f881.
//
// Solidity: function managerProxyContract() view returns(address)
func (_BusinessContract *BusinessContractCallerSession) ManagerProxyContract() (common.Address, error) {
	return _BusinessContract.Contract.ManagerProxyContract(&_BusinessContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BusinessContract *BusinessContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BusinessContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BusinessContract *BusinessContractSession) Owner() (common.Address, error) {
	return _BusinessContract.Contract.Owner(&_BusinessContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BusinessContract *BusinessContractCallerSession) Owner() (common.Address, error) {
	return _BusinessContract.Contract.Owner(&_BusinessContract.CallOpts)
}

// SetManagerProxy is a paid mutator transaction binding the contract method 0xaf9980f0.
//
// Solidity: function setManagerProxy(address ethCCMProxyAddr) returns()
func (_BusinessContract *BusinessContractTransactor) SetManagerProxy(opts *bind.TransactOpts, ethCCMProxyAddr common.Address) (*types.Transaction, error) {
	return _BusinessContract.contract.Transact(opts, "setManagerProxy", ethCCMProxyAddr)
}

// SetManagerProxy is a paid mutator transaction binding the contract method 0xaf9980f0.
//
// Solidity: function setManagerProxy(address ethCCMProxyAddr) returns()
func (_BusinessContract *BusinessContractSession) SetManagerProxy(ethCCMProxyAddr common.Address) (*types.Transaction, error) {
	return _BusinessContract.Contract.SetManagerProxy(&_BusinessContract.TransactOpts, ethCCMProxyAddr)
}

// SetManagerProxy is a paid mutator transaction binding the contract method 0xaf9980f0.
//
// Solidity: function setManagerProxy(address ethCCMProxyAddr) returns()
func (_BusinessContract *BusinessContractTransactorSession) SetManagerProxy(ethCCMProxyAddr common.Address) (*types.Transaction, error) {
	return _BusinessContract.Contract.SetManagerProxy(&_BusinessContract.TransactOpts, ethCCMProxyAddr)
}

