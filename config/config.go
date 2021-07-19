package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	PolyChainId uint64
	Address     common.Address
}

type TokenConfig struct {
	Name   string
	Tokens []Token
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

func (n *Network) PhrasePrivateKey() (err error) {
	_, ok1 := crypto.HexToECDSA(n.CCMPOwnerPrivateKey)
	_, ok2 := crypto.HexToECDSA(n.LockProxyOwnerPrivateKey)
	reader := bufio.NewReader(os.Stdin)
	if ok1 == nil && ok2 == nil { // no need to do anything
	} else if ok1 == nil { // need to recover LockProxy owner privatekey
		fmt.Printf("Please type in password of %s: ", n.LockProxyOwnerKeyStore)
		password, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		password = strings.Replace(password, "\n", "", -1)
		err = n.LockProxyOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	} else if ok2 == nil { // need to recover CCMPowner privatekey
		fmt.Printf("Please type in password of %s: ", n.CCMPOwnerKeyStore)
		password, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		password = strings.Replace(password, "\n", "", -1)
		err = n.CCMPOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	} else { // both
		fmt.Printf("Please type in password of %s: ", n.CCMPOwnerKeyStore)
		password, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		password = strings.Replace(password, "\n", "", -1)
		password2 := password
		if n.LockProxyOwnerKeyStore != n.CCMPOwnerKeyStore {
			fmt.Printf("Please type in password of %s: ", n.LockProxyOwnerKeyStore)
			password2, err = reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("fail to phrase private key, %v", err)
			}
			password2 = strings.Replace(password2, "\n", "", -1)
		}
		err = n.CCMPOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		err = n.LockProxyOwnerFromKeyStore(password2)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	}
	return nil
}

func (n *Network) CCMPOwnerFromKeyStore(pswd string) (err error) {
	ks1, err := ioutil.ReadFile(n.CCMPOwnerKeyStore)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v", err)
	}
	key1, err := keystore.DecryptKey(ks1, pswd)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v", err)
	}
	n.CCMPOwnerPrivateKey = fmt.Sprintf("%x", crypto.FromECDSA(key1.PrivateKey))
	return nil
}

func (n *Network) LockProxyOwnerFromKeyStore(pswd string) (err error) {
	ks2, err := ioutil.ReadFile(n.LockProxyOwnerKeyStore)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v", err)
	}
	key2, err := keystore.DecryptKey(ks2, pswd)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v", err)
	}
	n.LockProxyOwnerPrivateKey = fmt.Sprintf("%x", crypto.FromECDSA(key2.PrivateKey))
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
