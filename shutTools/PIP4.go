package shutTools

import (
	"context"
	"fmt"
	"github.com/KSlashh/emergency-button/abi"
	"github.com/KSlashh/emergency-button/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func AssetLPMap(client *ethclient.Client, conf *config.Network, token common.Address) (targetToken common.Address, err error) {
	LockProxyContract, err := abi.NewILockProxyWithLPCaller(conf.LockProxyPip4, client)
	if err != nil {
		return ADDRESS_ZERO, fmt.Errorf(
			fmt.Sprintf(
				"fail while request Asset %s =>to=> LPtoken at chain %d,",
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	return LockProxyContract.AssetLPMap(nil, token)
}

func DeployToken(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, token *config.TokenConfig, f bool) (common.Address, error) {
	fmt.Println("suuuuu")
	privateKey, err := crypto.HexToECDSA(pkCfg.SenderPrivateKey)
	if err != nil {
		return ADDRESS_ZERO, fmt.Errorf(
			fmt.Sprintf("fail while deploy Token %s from chain %d,", token.Name, conf.PolyChainID), err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return ADDRESS_ZERO, fmt.Errorf(
			fmt.Sprintf("fail while deploy Token %s from chain %d,", token.Name, conf.PolyChainID), err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return ADDRESS_ZERO, fmt.Errorf(
			fmt.Sprintf("fail while deploy Token %s from chain %d,", token.Name, conf.PolyChainID), err)
	}
	fmt.Println("suuuuu")
	var tx common.Address
	var txhash *types.Transaction
	if f {
		fmt.Println("pip", token.LPName, token.LPSymbol, token.Decimal)
		tx, txhash, _, err = abi.DeployERC20PreMint(auth, client, token.LPName, token.LPSymbol, token.Decimal, conf.LockProxyPip4, token.InitSupply)
		if err != nil {
			return ADDRESS_ZERO, fmt.Errorf(
				fmt.Sprintf("fail while deploy Token %s from chain %d,", token.Name, conf.PolyChainID), err)
		}
	} else {
		fmt.Println("lll", token.Name, token.Symbol, token.Decimal)

		tx, txhash, _, err = abi.DeployERC20PreMint(auth, client, token.Name, token.Symbol, token.Decimal, pkCfg.SenderPublicKey, token.InitSupply)
		if err != nil {
			return ADDRESS_ZERO, fmt.Errorf(
				fmt.Sprintf("fail while deploy Token %s from chain %d,", token.Name, conf.PolyChainID), err)
		}
	}
	fmt.Println("wait")

	err = WaitTxConfirm(client, txhash.Hash())
	if err != nil {

		return ADDRESS_ZERO, fmt.Errorf(
			fmt.Sprintf("fail while deploy Token %s from chain %d,", token.Name, conf.PolyChainID), err)
	}
	fmt.Println("hhhhhh", token.Name, token.Symbol, token.Decimal)
	return tx, err
}

func BindLPandAsset(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, token common.Address, lptoken common.Address, toChainId uint64, toAsset []byte, toLPAsset []byte) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyPip4OwnerPrivateKey)
	fmt.Println("pk")
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s and from chain %d =>to=> asset %x at chain %d",
				lptoken.Hex(),
				conf.PolyChainID,
				toLPAsset,
				toChainId),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s and from chain %d =>to=> asset %x at chain %d",
				lptoken.Hex(),
				conf.PolyChainID,
				toLPAsset,
				toChainId),
			err)
	}
	fmt.Println("45454545", chainId, conf.LockProxyPip4)
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s and from chain %d =>to=> asset %x at chain %d",
				lptoken.Hex(),
				conf.PolyChainID,
				toLPAsset,
				toChainId),
			err)
	}
	fmt.Println("222222")
	LockProxyContract, err := abi.NewILockProxyWithLP(conf.LockProxyPip4, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s and from chain %d =>to=> asset %x at chain %d,",
				lptoken.Hex(),
				conf.PolyChainID,
				toLPAsset,
				toChainId),
			err)
	}
	fmt.Println("3333", token, lptoken, toChainId, toAsset, toLPAsset)
	tx, err := LockProxyContract.BindLPAndAsset(auth, token, lptoken, toChainId, toAsset, toLPAsset)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s and from chain %d =>to=> asset %x at chain %d",
				lptoken.Hex(),
				conf.PolyChainID,
				toLPAsset,
				toChainId),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s and from chain %d =>to=> asset %x at chain %d,",
				lptoken.Hex(),
				conf.PolyChainID,
				toLPAsset,
				toChainId),
			err)
	}
	return nil
}

func BindLPToAsset(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, token common.Address, lptoken common.Address) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyPip4OwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s =>to=> asset %s at chain %d,",
				lptoken.Hex(),
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s =>to=> asset %s at chain %d,",
				lptoken.Hex(),
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s =>to=> asset %s at chain %d,",
				lptoken.Hex(),
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	LockProxyContract, err := abi.NewILockProxyWithLP(conf.LockProxyPip4, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s =>to=> asset %s at chain %d,",
				lptoken.Hex(),
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	tx, err := LockProxyContract.BindLPToAsset(auth, token, lptoken)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s =>to=> asset %s at chain %d,",
				lptoken.Hex(),
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("fail while bindLPToken %s =>to=> asset %s at chain %d,",
				lptoken.Hex(),
				token.Hex(),
				conf.PolyChainID),
			err)
	}
	return nil
}
func BindTokenPip4(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, token common.Address, toChainId uint64, toAsset []byte) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyPip4OwnerPrivateKey)
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
	LockProxyContract, err := abi.NewILockProxyWithLP(conf.LockProxyPip4, client)
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

func BindProxyHashPip4(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey, toChainId uint64, toProxy []byte) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyPip4OwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Proxy from chain %d =>to=> asset %x at chain %d,",
				conf.PolyChainID,
				toProxy,
				toChainId),
			err)
	}
	LockProxyContract, err := abi.NewILockProxyWithLP(conf.LockProxyPip4, client)
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

func BindLPandAeestBatch(gasMultiple float64, client *ethclient.Client, conf *config.Network, pkCfg *config.PrivateKey,
	fromAddress []common.Address, fromLpAddress []common.Address, toChainId []uint64, toAsset [][]byte, toLPAddress [][]byte) error {
	privateKey, err := crypto.HexToECDSA(pkCfg.LockProxyPip4OwnerPrivateKey)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token from chain %d",
				conf.PolyChainID),
			err)
	}
	LockProxyContract, err := abi.NewILockProxyWithLP(conf.LockProxyPip4, client)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token from chain %d",
				conf.PolyChainID),
			err)
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token from chain %d",
				conf.PolyChainID),
			err)
	}
	auth, err := MakeAuth(client, privateKey, DefaultGasLimit, gasMultiple, chainId)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token from chain %d",
				conf.PolyChainID),
			err)
	}
	tx, err := LockProxyContract.BindLPAndAssetBatch(auth, fromAddress, fromAddress, toChainId, toAsset, toLPAddress)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token from chain %d",
				conf.PolyChainID),
			err)
	}
	err = WaitTxConfirm(client, tx.Hash())
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf(
				"fail while bind Token from chain %d",
				conf.PolyChainID),
			err)
	}
	return nil
}
