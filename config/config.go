package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ssh/terminal"
)

var emptyAddress common.Address = common.Address{}
var passwordCache string = ""

type Key struct {
	PublicKey  common.Address
	PrivateKey string
	KeyStore   string
}

type KeyConfig struct {
	Keys []Key
}

type Network struct {
	Name              string
	Provider          string
	NewCCMPAddress    common.Address
	BusinessContracts []common.Address
}

type Config struct {
	Networks []Network
}

// Standardized Config
type SNetwork struct {
	Name                 string
	Provider             string
	EthCrossChainManager common.Address
}

// Standardized Config
type SConfig struct {
	Networks []SNetwork
}

func MergeConfig(sconfigFile string, configFile string, force bool) (err error) {
	c, err := LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("Load config failed, error: %s", err.Error())
	}
	sc, err := LoadSConfig(sconfigFile)
	if err != nil {
		return fmt.Errorf("Load sconfig failed, error: %s", err.Error())
	}
	for i := 0; i < len(c.Networks); i++ {
		sn := sc.GetNetwork(c.Networks[i].Name)
		if sn == nil {
			continue
		}
		if force || c.Networks[i].Provider == "" {
			c.Networks[i].Provider = sn.Provider
		}
		if force || c.Networks[i].NewCCMPAddress == emptyAddress {
			c.Networks[i].NewCCMPAddress = sn.EthCrossChainManager
		}
	}
	return c.WriteConfig(configFile)
}

func LoadKeyConfig(file string) (kc *KeyConfig, err error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	kc = &KeyConfig{}
	err = json.Unmarshal(jsonBytes, kc)
	if err != nil {
		return nil, err
	}
	err = kc.PhraseKeys()
	if err != nil {
		return nil, err
	}
	return kc, nil
}

func LoadConfig(file string) (c *Config, err error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c = &Config{}
	err = json.Unmarshal(jsonBytes, c)
	return c, nil
}

func LoadSConfig(file string) (sc *SConfig, err error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	sc = &SConfig{}
	err = json.Unmarshal(jsonBytes, sc)
	return sc, nil
}

func (k *Key) PhraseKey() (err error) {
	_, hasPk := crypto.HexToECDSA(k.PrivateKey)
	hasCache := k.KeyFromKeyStore(passwordCache)
	ok := hasPk == nil || hasCache == nil
	if !ok { // need to recover WrapperFeeCollector privatekey
		fmt.Printf("Please type in password of %s: ", k.KeyStore)
		pass, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
		fmt.Println()
		password := string(pass)
		password = strings.Replace(password, "\n", "", -1)
		passwordCache = password
		err = k.KeyFromKeyStore(password)
		if err != nil {
			return fmt.Errorf("fail to phrase private key, %v", err)
		}
	}
	key, err := crypto.HexToECDSA(k.PrivateKey)
	if err != nil {
		return err
	}
	k.PublicKey = crypto.PubkeyToAddress(key.PublicKey)
	return nil
}

func (k *Key) KeyFromKeyStore(pswd string) (err error) {
	ks, err := ioutil.ReadFile(k.KeyStore)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v", err)
	}
	key1, err := keystore.DecryptKey(ks, pswd)
	if err != nil {
		return fmt.Errorf("fail to recover private key from keystore file, %v", err)
	}
	k.PrivateKey = fmt.Sprintf("%x", crypto.FromECDSA(key1.PrivateKey))
	return nil
}

func (kc *KeyConfig) PhraseKeys() (err error) {
	for i := 0; i < len(kc.Keys); i++ {
		err = kc.Keys[i].PhraseKey()
		if err != nil {
			return err
		}
	}
	return nil
}

func (kc *KeyConfig) GetKey(PublicKey common.Address) (key *Key) {
	for i := 0; i < len(kc.Keys); i++ {
		if kc.Keys[i].PublicKey == PublicKey {
			return &kc.Keys[i]
		}
	}
	return nil
}

func (sc *SConfig) GetNetwork(name string) (network *SNetwork) {
	for i := 0; i < len(sc.Networks); i++ {
		if sc.Networks[i].Name == name {
			return &sc.Networks[i]
		}
	}
	return nil
}

func (c *Config) GetNetwork(name string) (network *Network) {
	for i := 0; i < len(c.Networks); i++ {
		if c.Networks[i].Name == name {
			return &c.Networks[i]
		}
	}
	return nil
}

func (c *Config) WriteConfig(file string) (err error) {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0777)
}
