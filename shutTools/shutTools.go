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

var DefaultGasLimit uint64 = 300000
var gasMultipleDecimal int64 = 8

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
