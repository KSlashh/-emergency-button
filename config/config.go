package config

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
)


type Network struct {
	PolyChainID              uint64
	Name                     string
	Provider                 string
	CCMPOwnerPrivateKey      string
	LockProxyOwnerPrivateKey string
	CCMPAddress              common.Address
	LockProxyAddress         common.Address
}

type Config struct {
	Networks []Network
}

// LoadConfig ...
func LoadConfig(confFile string) (config *Config, err error) {
	jsonBytes, err := ioutil.ReadFile(confFile)
	if err != nil {
		return
	}

	config = &Config{}
	err = json.Unmarshal(jsonBytes, config)
	return
}

// search by chainId or name
func (c *Config) GetNetwork(index uint64) (netConfig *Network) {
	for i := 0; i < len(c.Networks); i++ {
		if c.Networks[i].PolyChainID == index {
			return &c.Networks[i]
		}
	}
	return nil
}
