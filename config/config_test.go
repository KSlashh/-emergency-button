package config

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	Eth := Network{
		2,
		"eth",
		"https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		"{put CCM owner's sk here}",
		"{put LockProxy owner's sk here}",
		common.HexToAddress("0x5a51e2ebf8d136926b9ca7b59b60464e7c44d2eb"),
		common.HexToAddress("0x250e76987d838a75310c34bf422ea9f1AC4Cc906"),
	}
	Bsc := Network{
		6,
		"bsc",
		"https://bsc-dataseed.binance.org",
		"{put CCM owner's sk here}",
		"{put LockProxy owner's sk here}",
		common.HexToAddress("0xABD7f7B89c5fD5D0AEf06165f8173b1b83d7D5c9"),
		common.HexToAddress("0x2f7ac9436ba4B548f9582af91CA1Ef02cd2F1f03"),
	}
	Heco := Network{
		7,
		"heco",
		"https://http-mainnet-node.huobichain.com",
		"{put CCM owner's sk here}",
		"{put LockProxy owner's sk here}",
		common.HexToAddress("0xABD7f7B89c5fD5D0AEf06165f8173b1b83d7D5c9"),
		common.HexToAddress("0x020c15e7d08A8Ec7D35bCf3AC3CCbF0BBf2704e6"),
	}
	Ok := Network{
		12,
		"ok",
		"https://exchainrpc.okex.org/",
		"{put CCM owner's sk here}",
		"{put LockProxy owner's sk here}",
		common.HexToAddress("0x4739fe955BE4704BcB7d6a699823F5B29217Baf6"),
		common.HexToAddress("0x9a3658864Aa2Ccc63FA61eAAD5e4f65fA490cA7D"),
	}
	nets := []Network{Eth, Bsc, Heco, Ok}
	res, err := json.Marshal(&Config{nets})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", res)
	err = ioutil.WriteFile("./sampleConfig.json", res, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestReadConfig(t *testing.T) {
	conf, err := LoadConfig("sampleConfig.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	net1 := conf.GetNetwork(12)
	if net1 == nil {
		fmt.Println("1 nil")
		fmt.Println(conf.Networks[3])
		return
	}
	fmt.Printf("chain %d has name %s", net1.PolyChainID, net1.Name)
}

func TestGenerateTokenConfig(t *testing.T) {
	Eth := Token{
		2,
		common.HexToAddress("0x250e76987d838a75310c34bf422ea9f1AC4Cc906"),
	}
	Bsc := Token{
		6,
		common.HexToAddress("0x2f7ac9436ba4B548f9582af91CA1Ef02cd2F1f03"),
	}
	Heco := Token{
		7,
		common.HexToAddress("0x020c15e7d08A8Ec7D35bCf3AC3CCbF0BBf2704e6"),
	}
	Ok := Token{
		12,
		common.HexToAddress("0x9a3658864Aa2Ccc63FA61eAAD5e4f65fA490cA7D"),
	}
	tokens := []Token{Eth, Bsc, Heco, Ok}
	res, err := json.Marshal(&TokenConfig{"sampleToken",tokens})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", res)
	err = ioutil.WriteFile("./sampleToken.json", res, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestReadTokenConfig(t *testing.T) {
	conf, err := LoadToken("sampleToken.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("token %s at chain %d is %s", conf.Name, conf.Tokens[2].PolyChainId, conf.Tokens[2].Address.Hex())
}
