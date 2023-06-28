package shutTools

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/KSlashh/emergency-button/abi"
	"github.com/KSlashh/emergency-button/config"
	"github.com/KSlashh/emergency-button/log"

	"github.com/ethereum/go-ethereum"
	eabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ok_id int64 = 66
var ok_test_id int64 = 65

var DefaultGasLimit uint64 = 0 // get gas limit via EstimateGas
var gasMultipleDecimal int64 = 8
var ADDRESS_ZERO common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

// CCM
func ShutCCM(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.CCMPOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while shut CCM of %s ,", conf.Name), err)
	}
	CCMPContract, err := abi.NewICCMP(conf.EthCrossChainManagerProxy, client)
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

func RestartCCM(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.CCMPOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while restart CCM of %s ,", conf.Name), err)
	}
	CCMPContract, err := abi.NewICCMP(conf.EthCrossChainManagerProxy, client)
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
	CCMPContract, err := abi.NewICCMPCaller(conf.EthCrossChainManagerProxy, client)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("fail while request CCM of %s ,", conf.Name), err)
	}
	return CCMPContract.Paused(nil)
}

///LockProxy
func BindToken(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, token common.Address, toChainId uint64, toAsset []byte) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyOwnerPrivateKey)
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
	LockProxyContract, err := abi.NewILockProxy(conf.LockProxy, client)
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

func BindProxyHash(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, toChainId uint64, toProxy []byte) error {

	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	lockProxyAbi, err := eabi.JSON(strings.NewReader(abi.ILockProxyABI))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	data, err := lockProxyAbi.Pack("bindProxyHash", toChainId, toProxy)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	gasLimit := auth.GasLimit
	if gasLimit == 0 {
		msg := ethereum.CallMsg{From: auth.From, To: &conf.LockProxy, GasPrice: auth.GasPrice, Value: auth.Value, Data: data}
		gasLimit, err = client.EstimateGas(context.Background(), msg)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
		}
	}

	tx := types.NewTransaction(auth.Nonce.Uint64(), conf.LockProxy, auth.Value, gasLimit, auth.GasPrice, data)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	signer := types.LatestSignerForChainID(chainId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	tx, err = types.SignTx(tx, signer, privateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while try to call bindProxyHash() from chain %d =>to=> asset %x at chain %d,", conf.PolyChainID, toProxy, toChainId), err)
	}
	return WaitTxConfirm(client, tx.Hash())
}

func BindProxyHashOld(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, toChainId uint64, toProxy []byte) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	LockProxyContract, err := abi.NewILockProxy(conf.LockProxy, client)
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

func TokenMap(client *ethclient.Client, conf *config.Network, lockproxy common.Address, token common.Address, toChainId uint64) (targetToken []byte, err error) {
	LockProxyContract, err := abi.NewILockProxyCaller(lockproxy, client)
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
	LockProxyContract, err := abi.NewILockProxyCaller(conf.LockProxy, client)
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
func ExtractFeeSwapper(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, tokenAddress common.Address) error {
	if conf.Swapper == ADDRESS_ZERO {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, swapper address in config is ZERO",
				conf.PolyChainID),
		)
	}
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperFeeCollectorPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewIWrapper(conf.Swapper, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	tx, err := SwapperContract.ExtractFee(auth, tokenAddress)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at Swapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	return nil
}

func RegisterPool(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, poolId uint64, poolTokenAddress common.Address) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.Swapper, client)
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

func BindAsserAndPool(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, fromAssetHash []byte, poolId uint64) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Asset %x at pool %d at chain %d, ",
				fromAssetHash,
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.Swapper, client)
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

func Bind3Asset(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, asset1 []byte, asset2 []byte, asset3 []byte, poolId uint64) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.Swapper, client)
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

func RegisterPoolWith3Assets(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, poolTokenAddress common.Address, asset1 []byte, asset2 []byte, asset3 []byte, poolId uint64) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while register and bind 3 Asset at pool %d at chain %d, ",
				poolId,
				conf.PolyChainID),
			err)
	}
	SwapperContract, err := abi.NewISwapper(conf.Swapper, client)
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

func PauseSwapper(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while pause swapper of %s ,", conf.Name), err)
	}
	SwapperContract, err := abi.NewISwapper(conf.Swapper, client)
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

func UnpauseSwapper(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.SwapperOwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("fail while unpause swapper of %s ,", conf.Name), err)
	}
	SwapperContract, err := abi.NewISwapper(conf.Swapper, client)
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

func SwapperPaused(client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey) (paused bool, err error) {
	SwapperContract, err := abi.NewISwapperCaller(conf.Swapper, client)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("fail while request Swapper of %s ,", conf.Name), err)
	}
	return SwapperContract.Paused(nil)
}

func PoolTokenMap(client *ethclient.Client, conf *config.Network, poolId uint64) (poolTokenAddress common.Address, err error) {
	SwapperContract, err := abi.NewISwapperCaller(conf.Swapper, client)
	if err != nil {
		return ADDRESS_ZERO, fmt.Errorf(fmt.Sprintf("fail while request Swapper of %s ,", conf.Name), err)
	}
	return SwapperContract.PoolTokenMap(nil, poolId)
}

func AssetInPool(client *ethclient.Client, conf *config.Network, poolId uint64, assetHash []byte) (isIn bool, err error) {
	SwapperContract, err := abi.NewISwapperCaller(conf.Swapper, client)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("fail while request Swapper of %s ,", conf.Name), err)
	}
	return SwapperContract.AssetInPool(nil, assetHash, poolId)
}

// Wrapper
func ExtractFeeWrapper(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, tokenAddress common.Address) error {
	if conf.Wrapper == ADDRESS_ZERO {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, wrapper address in config is ZERO",
				conf.PolyChainID),
		)
	}
	privateKey, err := crypto.HexToECDSA(pkCfg.WrapperFeeCollectorPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	WrapperContract, err := abi.NewIWrapper(conf.Wrapper, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	tx, err := WrapperContract.ExtractFee(auth, tokenAddress)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, ",
				conf.PolyChainID),
			err)
	}
	return nil
}

func ExtractFeeWrapperO3(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey) error {
	if conf.WrapperO3 == ADDRESS_ZERO {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at PolyWrapper at chain %d, wrapper address in config is ZERO",
				conf.PolyChainID),
		)
	}
	privateKey, err := crypto.HexToECDSA(pkCfg.WrapperO3FeeCollectorPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at WrapperO3 at chain %d, ",
				conf.PolyChainID),
			err)
	}
	WrapperContract, err := abi.NewWrapperO3(conf.WrapperO3, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at WrapperO3 at chain %d, ",
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at WrapperO3 at chain %d, ",
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at WrapperO3 at chain %d, ",
				conf.PolyChainID),
			err)
	}
	tx, err := WrapperContract.ExtractFee(auth)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at WrapperO3 at chain %d, ",
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while extract fee at WrapperO3 at chain %d, ",
				conf.PolyChainID),
			err)
	}
	return nil
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
	ticker := time.NewTicker(time.Second * 1)
	end := time.Now().Add(60 * time.Second)
	for now := range ticker.C {
		_, pending, err := client.TransactionByHash(context.Background(), hash)
		if err != nil {
			log.Info("failed to call TransactionByHash: %v", err)
			if now.After(end) {
				break
			}
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
