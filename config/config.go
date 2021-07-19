package config

import (
	"encoding/json"
	"fmt"
    "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
)


type Network struct {
	PolyChainID              uint64
	Name                     string
	Provider                 string
	CCMPOwnerPrivateKey      string
	CCMPOwnerKeyStore        string
	LockProxyOwnerPrivateKey string
	LockProxyOwnerKeyStore   string
	CCMPAddress              common.Address
	LockProxyAddress         common.Address
}

type Config struct {
	Networks []Network
}

type Token struct {
	PolyChainId   uint64
	Address       common.Address
}

type TokenConfig struct {
	Name    string
	Tokens  []Token
}

func LoadConfig(confFile string) (config *Config, err error) {
	jsonBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return
	}

	config = &Config{}
	err = json.Unmarshal(jsonBytes, config)
	return
}

func (c *Config) GetNetwork(index uint64) (netConfig *Network) {
	for i := 0; i < len(c.Networks); i++ {
		if c.Networks[i].PolyChainID == index {
			return &c.Networks[i]
		}
	}
	return nil
}

func (n *Network) CCMPOwnerFromKeyStore(pswd string) (err error) {
	ks1, err := ioutil.ReadFile(n.CCMPOwnerKeyStore)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v",err)
	}
	key1, err := keystore.DecryptKey(ks1, pswd)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v",err)
	}
	n.CCMPOwnerPrivateKey = fmt.Sprintf("%x",crypto.FromECDSA(key1.PrivateKey))
    return nil
}

func (n *Network) LockProxyFromKeyStore(pswd string) (err error) {
	ks2, err := ioutil.ReadFile(n.LockProxyOwnerKeyStore)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v",err)
	}
	key2, err := keystore.DecryptKey(ks2, pswd)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v",err)
	}
	n.CCMPOwnerPrivateKey = fmt.Sprintf("%x",crypto.FromECDSA(key2.PrivateKey))
	return nil
}

func LoadToken(tokenFile string) (tokens *TokenConfig, err error) {
	jsonBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return
	}

	tokens = &TokenConfig{}
	err = json.Unmarshal(jsonBytes, tokens)
	return
}

func (c *TokenConfig) GetToken(index uint64) (netConfig *Token) {
	for i := 0; i < len(c.Tokens); i++ {
		if c.Tokens[i].PolyChainId == index {
			return &c.Tokens[i]
		}
	}
	return nil
}
