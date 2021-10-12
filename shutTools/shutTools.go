package shutTools

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/KSlashh/emergency-button/abi"
	"github.com/KSlashh/emergency-button/config"
	"github.com/KSlashh/emergency-button/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ok_id int64 = 66
var ok_test_id int64 = 65

var DefaultGasLimit uint64 = 0 // get gas limit via EstimateGas
var gasMultipleDecimal int64 = 8
var ADDRESS_ZERO common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

// CCM
func ShutCCM(gasMultiple float64, client *ethclient.Client, conf *config.Network) error {
	privateKey, err := crypto.HexToECDSA(conf.CCMPOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	CCMPContract, err := abi.NewICCMP(conf.CCMPAddress, client)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	tx, err := CCMPContract.PauseEthCrossChainManager(auth)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	return nil
}

func RestartCCM(gasMultiple float64, client *ethclient.Client, conf *config.Network) error {
	privateKey, err := crypto.HexToECDSA(conf.CCMPOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	CCMPContract, err := abi.NewICCMP(conf.CCMPAddress, client)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	tx, err := CCMPContract.UnpauseEthCrossChainManager(auth)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	return nil
}

func CCMPaused(client *ethclient.Client, conf *config.Network) (paused bool, err error) {
	CCMPContract, err := abi.NewICCMPCaller(conf.CCMPAddress, client)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("fail while request CCM of %s ,", conf.Name), err)
	}
	return CCMPContract.Paused(nil)
}

// LockProxy
func BindToken(gasMultiple float64, client *ethclient.Client, conf *config.Network, token common.Address, toChainId uint64, toAsset []byte) error {
	privateKey, err := crypto.HexToECDSA(conf.LockProxyOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token %s from chain %d =>to=> asset %x at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toAsset,
				toChainId),
			err)
	}
	LockProxyContract, err := abi.NewILockProxy(conf.LockProxyAddress, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token %s from chain %d =>to=> asset %x at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toAsset,
				toChainId),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token %s from chain %d =>to=> asset %x at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toAsset,
				toChainId),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token %s from chain %d =>to=> asset %x at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toAsset,
				toChainId),
			err)
	}
	tx, err := LockProxyContract.BindAssetHash(auth, token, toChainId, toAsset)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token %s from chain %d =>to=> asset %x at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toAsset,
				toChainId),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token %s from chain %d =>to=> asset %x at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toAsset,
				toChainId),
			err)
	}
	return nil
}

func BindProxyHash(gasMultiple float64, client *ethclient.Client, conf *config.Network, toChainId uint64, toProxy []byte) error {
	privateKey, err := crypto.HexToECDSA(conf.LockProxyOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	LockProxyContract, err := abi.NewILockProxy(conf.LockProxyAddress, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	tx, err := LockProxyContract.BindProxyHash(auth, toChainId, toProxy)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	return nil
}

func TokenMap(client *ethclient.Client, conf *config.Network, token common.Address, toChainId uint64) (targetToken []byte, err error) {
	LockProxyContract, err := abi.NewILockProxyCaller(conf.LockProxyAddress, client)
	if err != nil {
		return nil, fmt.Errorf(
			fmt.Sprintf(
				"fail while request Token %s from chain %d =>to=> asset at chain %d,",
				token.Hex(),
				conf.PolyChainID,
				toChainId),
			err)
	}
	return LockProxyContract.AssetHashMap(nil, token, toChainId)
}

func ProxyHashMap(client *ethclient.Client, conf *config.Network, toChainId uint64) (targetProxy []byte, err error) {
	LockProxyContract, err := abi.NewILockProxyCaller(conf.LockProxyAddress, client)
	if err != nil {
		return nil, fmt.Errorf(
			fmt.Sprintf(
				"fail while request Proxy from chain %d =>to=> asset at chain %d,",
				conf.PolyChainID,
				toChainId),
			err)
	}
	return LockProxyContract.ProxyHashMap(nil, toChainId)
}

// Swapper
func RegisterPool(gasMultiple float64, client *ethclient.Client, conf *config.Network, poolId uint64, poolTokenAddress common.Address) error {
	privateKey, err := crypto.HexToECDSA(conf.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.SwapperAddress, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	tx, err := SwapperContract.RegisterPool(auth, poolId, poolTokenAddress)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	return nil
}

func BindAsserAndPool(gasMultiple float64, client *ethclient.Client, conf *config.Network, fromAssetHash []byte, poolId uint64) error {
	privateKey, err := crypto.HexToECDSA(conf.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.SwapperAddress, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	tx, err := SwapperContract.BindAssetAndPool(auth, fromAssetHash, poolId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	return nil
}

func Bind3Asset(gasMultiple float64, client *ethclient.Client, conf *config.Network, asset1 []byte, asset2 []byte, asset3 []byte, poolId uint64) error {
	privateKey, err := crypto.HexToECDSA(conf.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.SwapperAddress, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	tx, err := SwapperContract.Bind3Asset(auth, asset1, asset2, asset3, poolId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	return nil
}

func RegisterPoolWith3Assets(gasMultiple float64, client *ethclient.Client, conf *config.Network, poolTokenAddress common.Address, asset1 []byte, asset2 []byte, asset3 []byte, poolId uint64) error {
	privateKey, err := crypto.HexToECDSA(conf.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.SwapperAddress, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	tx, err := SwapperContract.RegisterPoolWith3Assets(auth, poolId, poolTokenAddress, asset1, asset2, asset3)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	return nil
}

func PauseSwapper(gasMultiple float64, client *ethclient.Client, conf *config.Network) error {
	privateKey, err := crypto.HexToECDSA(conf.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	SwapperContract, err := abi.NewISwapper(conf.SwapperAddress, client)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	tx, err := SwapperContract.Pause(auth)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	return nil
}

func UnpauseSwapper(gasMultiple float64, client *ethclient.Client, conf *config.Network) error {
	privateKey, err := crypto.HexToECDSA(conf.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	SwapperContract, err := abi.NewISwapper(conf.SwapperAddress, client)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	tx, err := SwapperContract.Unpause(auth)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	return nil
}

func SwapperPaused(client *ethclient.Client, conf *config.Network) (paused bool, err error) {
	SwapperContract, err := abi.NewISwapperCaller(conf.SwapperAddress, client)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("fail while request Swapper of %s ,", conf.Name), err)
	}
	return SwapperContract.Paused(nil)
}

func PoolTokenMap(client *ethclient.Client, conf *config.Network, poolId uint64) (poolTokenAddress common.Address, err error) {
	SwapperContract, err := abi.NewISwapperCaller(conf.SwapperAddress, client)
	if err != nil {
		return ADDRESS_ZERO, fmt.Errorf(fmt.Sprintf("fail while request Swapper of %s ,", conf.Name), err)
	}
	return SwapperContract.PoolTokenMap(nil, poolId)
}

func AssetInPool(client *ethclient.Client, conf *config.Network, poolId uint64, assetHash []byte) (isIn bool, err error) {
	SwapperContract, err := abi.NewISwapperCaller(conf.SwapperAddress, client)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("fail while request Swapper of %s ,", conf.Name), err)
	}
	return SwapperContract.AssetInPool(nil, assetHash, poolId)
}

// Basic
func MakeAuth(client *ethclient.Client, key *ecdsa.PrivateKey, gasLimit uint64, gasMultiple float64, chainId *big.Int) (*bind.TransactOpts, error) {
	authAddress := crypto.PubkeyToAddress(*key.Public().(*ecdsa.PublicKey))
	nonce, err := client.PendingNonceAt(context.Background(), authAddress)
	if err != nil {
		return nil, fmt.Errorf("makeAuth, addr %s, err %v", authAddress.Hex(), err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("makeAuth, get suggest gas price err: %v", err)
	}
	res := gasPrice.Mul(gasPrice, big.NewInt(int64(gasMultiple*math.Pow(10, float64(gasMultipleDecimal)))))
	if res == nil {
		return nil, fmt.Errorf("calculate actual gas price error (at mul")
	}
	res = gasPrice.Div(gasPrice, big.NewInt(int64(math.Pow(10, float64(gasMultipleDecimal)))))
	if res == nil {
		return nil, fmt.Errorf("calculate actual gas price error (at div")
	}

	// auth := bind.NewKeyedTransactor(key)
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		return nil, fmt.Errorf("makeAuth, bind.NewKeyedTransactorWithChainID err: %v", err)
	}
	auth.From = authAddress
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(int64(0)) // in wei
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

func WaitTxConfirm(client *ethclient.Client, hash common.Hash) error {

	// cannot receive receipt at ok chain , so skip it
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("faild to get chainId %v", err)
	}
	if chainId.Int64() == ok_id || chainId.Int64() == ok_test_id {
		log.Info("Can not get receipt of txns at okex, check transaction %s on explorer yourself, make sure it's confirmed.", hash.Hex())
		return nil
	}

	ticker := time.NewTicker(time.Second * 1)
	end := time.Now().Add(60 * time.Second)
	for now := range ticker.C {
		_, pending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil {
			log.Debug("failed to call TransactionByHash: %v", err)
			continue
		}
		if !pending {
			break
		}
		if now.Before(end) {
			continue
		}
		log.Info("Transaction pending for more than 1 min, check transaction %s on explorer yourself, make sure it's confirmed.", hash.Hex())
		return nil
	}

	tx, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return fmt.Errorf("faild to get receipt %s", hash.Hex())
	}

	if tx.Status == 0 {
		return fmt.Errorf("receipt failed %s", hash.Hex())
	}

	return nil
}
