package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ssh/terminal"
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

var passwordCache string = ""

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

func (c *Config) GetNetworkIds() []string {
	var res []string
	for i := 0; i < len(c.Networks); i++ {
		res = append(res, strconv.Itoa(int(c.Networks[i].PolyChainID)))
	}
	return res
}

func (n *Network) PhraseCCMPrivateKey() (err error) {
	_, hasPk1 := crypto.HexToECDSA(n.CCMPOwnerPrivateKey)
	hasCache1 := n.CCMPOwnerFromKeyStore(passwordCache)
	ok1 := hasPk1 == nil || hasCache1 == nil
	if !ok1 {
		fmt.Printf("Please type in password of %s: ", n.CCMPOwnerKeyStore)
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		fmt.Println()
		password := string(pass)
		password = strings.Replace(password, "\n", "", -1)
		passwordCache = password
		err = n.CCMPOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	}
	return nil
}

func (n *Network) PhraseLockProxyPrivateKey() (err error) {
	_, hasPk2 := crypto.HexToECDSA(n.LockProxyOwnerPrivateKey)
	hasCache2 := n.LockProxyOwnerFromKeyStore(passwordCache)
	ok2 := hasPk2 == nil || hasCache2 == nil
	if !ok2 { // need to recover LockProxy owner privatekey
		fmt.Printf("Please type in password of %s: ", n.LockProxyOwnerKeyStore)
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		fmt.Println()
		password := string(pass)
		password = strings.Replace(password, "\n", "", -1)
		passwordCache = password
		err = n.LockProxyOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	}
	return nil
}

func (n *Network) PhrasePrivateKey() (err error) {
	_, hasPk1 := crypto.HexToECDSA(n.CCMPOwnerPrivateKey)
	_, hasPk2 := crypto.HexToECDSA(n.LockProxyOwnerPrivateKey)
	hasCache1 := n.CCMPOwnerFromKeyStore(passwordCache)
	hasCache2 := n.LockProxyOwnerFromKeyStore(passwordCache)
	ok1 := hasPk1 == nil || hasCache1 == nil
	ok2 := hasPk2 == nil || hasCache2 == nil
	if ok1 && ok2 { // no need to do anything
	} else if ok1 { // need to recover LockProxy owner privatekey
		fmt.Printf("Please type in password of %s: ", n.LockProxyOwnerKeyStore)
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		fmt.Println()
		password := string(pass)
		password = strings.Replace(password, "\n", "", -1)
		passwordCache = password
		err = n.LockProxyOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	} else if ok2 { // need to recover CCMPowner privatekey
		fmt.Printf("Please type in password of %s: ", n.CCMPOwnerKeyStore)
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		fmt.Println()
		password := string(pass)
		password = strings.Replace(password, "\n", "", -1)
		passwordCache = password
		err = n.CCMPOwnerFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	} else { // both
		fmt.Printf("Please type in password of %s: ", n.CCMPOwnerKeyStore)
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		fmt.Println()
		password := string(pass)
		password2 := password
		passwordCache = password
		if n.LockProxyOwnerKeyStore != n.CCMPOwnerKeyStore {
			fmt.Printf("Please type in password of %s: ", n.LockProxyOwnerKeyStore)
			pass, err = terminal.ReadPassword(0)
			if err != nil {
				return fmt.Errorf("fail to phrase private key, %v", err)
			}
			fmt.Println()
			password2 = string(pass)
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

func (c *TokenConfig) GetTokenIds() []string {
	var res []string
	for i := 0; i < len(c.Tokens); i++ {
		res = append(res, strconv.Itoa(int(c.Tokens[i].PolyChainId)))
	}
	return res
}
