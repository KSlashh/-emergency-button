// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ICCMPABI is the input ABI used to generate the binding from.
const ICCMPABI = "[{\"inputs\":[],\"name\":\"pauseEthCrossChainManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpauseEthCrossChainManager\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ICCMPFuncSigs maps the 4-byte function signature to its string representation.
var ICCMPFuncSigs = map[string]string{
	"3b9a80b8": "pauseEthCrossChainManager()",
	"4390c707": "unpauseEthCrossChainManager()",
}

// ICCMPBin is the compiled bytecode used for deploying new contracts.
var ICCMPBin = "0x6080604052348015600f57600080fd5b50608c8061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80633b9a80b81460375780634390c707146037575b600080fd5b603d6051565b604080519115158252519081900360200190f35b60009056fea2646970667358221220a15b3a487ae2fcdb24c05e227993944b8ed5b5bf447af38a6a289fceb9f3298064736f6c634300060c0033"

// DeployICCMP deploys a new Ethereum contract, binding an instance of ICCMP to it.
func DeployICCMP(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ICCMP, error) {
	parsed, err := abi.JSON(strings.NewReader(ICCMPABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ICCMPBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ICCMP{ICCMPCaller: ICCMPCaller{contract: contract}, ICCMPTransactor: ICCMPTransactor{contract: contract}, ICCMPFilterer: ICCMPFilterer{contract: contract}}, nil
}

// ICCMP is an auto generated Go binding around an Ethereum contract.
type ICCMP struct {
	ICCMPCaller     // Read-only binding to the contract
	ICCMPTransactor // Write-only binding to the contract
	ICCMPFilterer   // Log filterer for contract events
}

// ICCMPCaller is an auto generated read-only Go binding around an Ethereum contract.
type ICCMPCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICCMPTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ICCMPTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICCMPFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ICCMPFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICCMPSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ICCMPSession struct {
	Contract     *ICCMP            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ICCMPCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ICCMPCallerSession struct {
	Contract *ICCMPCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ICCMPTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ICCMPTransactorSession struct {
	Contract     *ICCMPTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ICCMPRaw is an auto generated low-level Go binding around an Ethereum contract.
type ICCMPRaw struct {
	Contract *ICCMP // Generic contract binding to access the raw methods on
}

// ICCMPCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ICCMPCallerRaw struct {
	Contract *ICCMPCaller // Generic read-only contract binding to access the raw methods on
}

// ICCMPTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ICCMPTransactorRaw struct {
	Contract *ICCMPTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICCMP creates a new instance of ICCMP, bound to a specific deployed contract.
func NewICCMP(address common.Address, backend bind.ContractBackend) (*ICCMP, error) {
	contract, err := bindICCMP(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICCMP{ICCMPCaller: ICCMPCaller{contract: contract}, ICCMPTransactor: ICCMPTransactor{contract: contract}, ICCMPFilterer: ICCMPFilterer{contract: contract}}, nil
}

// NewICCMPCaller creates a new read-only instance of ICCMP, bound to a specific deployed contract.
func NewICCMPCaller(address common.Address, caller bind.ContractCaller) (*ICCMPCaller, error) {
	contract, err := bindICCMP(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICCMPCaller{contract: contract}, nil
}

// NewICCMPTransactor creates a new write-only instance of ICCMP, bound to a specific deployed contract.
func NewICCMPTransactor(address common.Address, transactor bind.ContractTransactor) (*ICCMPTransactor, error) {
	contract, err := bindICCMP(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICCMPTransactor{contract: contract}, nil
}

// NewICCMPFilterer creates a new log filterer instance of ICCMP, bound to a specific deployed contract.
func NewICCMPFilterer(address common.Address, filterer bind.ContractFilterer) (*ICCMPFilterer, error) {
	contract, err := bindICCMP(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICCMPFilterer{contract: contract}, nil
}

// bindICCMP binds a generic wrapper to an already deployed contract.
func bindICCMP(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ICCMPABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICCMP *ICCMPRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICCMP.Contract.ICCMPCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICCMP *ICCMPRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICCMP.Contract.ICCMPTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICCMP *ICCMPRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICCMP.Contract.ICCMPTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICCMP *ICCMPCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICCMP.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICCMP *ICCMPTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICCMP.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICCMP *ICCMPTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICCMP.Contract.contract.Transact(opts, method, params...)
}

// PauseEthCrossChainManager is a paid mutator transaction binding the contract method 0x3b9a80b8.
//
// Solidity: function pauseEthCrossChainManager() returns(bool)
func (_ICCMP *ICCMPTransactor) PauseEthCrossChainManager(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICCMP.contract.Transact(opts, "pauseEthCrossChainManager")
}

// PauseEthCrossChainManager is a paid mutator transaction binding the contract method 0x3b9a80b8.
//
// Solidity: function pauseEthCrossChainManager() returns(bool)
func (_ICCMP *ICCMPSession) PauseEthCrossChainManager() (*types.Transaction, error) {
	return _ICCMP.Contract.PauseEthCrossChainManager(&_ICCMP.TransactOpts)
}

// PauseEthCrossChainManager is a paid mutator transaction binding the contract method 0x3b9a80b8.
//
// Solidity: function pauseEthCrossChainManager() returns(bool)
func (_ICCMP *ICCMPTransactorSession) PauseEthCrossChainManager() (*types.Transaction, error) {
	return _ICCMP.Contract.PauseEthCrossChainManager(&_ICCMP.TransactOpts)
}

// UnpauseEthCrossChainManager is a paid mutator transaction binding the contract method 0x4390c707.
//
// Solidity: function unpauseEthCrossChainManager() returns(bool)
func (_ICCMP *ICCMPTransactor) UnpauseEthCrossChainManager(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICCMP.contract.Transact(opts, "unpauseEthCrossChainManager")
}

// UnpauseEthCrossChainManager is a paid mutator transaction binding the contract method 0x4390c707.
//
// Solidity: function unpauseEthCrossChainManager() returns(bool)
func (_ICCMP *ICCMPSession) UnpauseEthCrossChainManager() (*types.Transaction, error) {
	return _ICCMP.Contract.UnpauseEthCrossChainManager(&_ICCMP.TransactOpts)
}

// UnpauseEthCrossChainManager is a paid mutator transaction binding the contract method 0x4390c707.
//
// Solidity: function unpauseEthCrossChainManager() returns(bool)
func (_ICCMP *ICCMPTransactorSession) UnpauseEthCrossChainManager() (*types.Transaction, error) {
	return _ICCMP.Contract.UnpauseEthCrossChainManager(&_ICCMP.TransactOpts)
}

