package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KSlashh/emergency-button/config"
	"github.com/KSlashh/emergency-button/log"
	"github.com/KSlashh/emergency-button/shutTools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ADDRESS_ZERO common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

var confFile string
var pkconfFile string
var tokenFile string
var function string
var multiple float64
var force bool
var all bool
var poolId uint64
var chainId uint64
var pip4 bool

type Msg struct {
	ChainId uint64
	Err     error
}

type Token struct {
	ChainId uint64
	Address common.Address
	//LpAddress common.Address
	NetCfg *config.Network
	PkCfg  *config.PrivateKey
}

type LPToken struct {
	ChainId   uint64
	Address   common.Address
	LpAddress common.Address
	NetCfg    *config.Network
	PkCfg     *config.PrivateKey
}

type Tokenlist struct {
	ChainId   uint64
	Address   []common.Address
	LpAddress []common.Address
	toChainId []uint64
	toAsset   [][]byte
	toLPAsset [][]byte
	NetCfg    *config.Network
	PkCfg     *config.PrivateKey
}

func init() {
	flag.Uint64Var(&poolId, "pool", 0, "pool id if needed")
	flag.Uint64Var(&chainId, "chain", 0, "chain id if single chainId needed")
	flag.StringVar(&tokenFile, "token", "./configs/token.json", "token configuration file path")
	flag.StringVar(&confFile, "conf", "./configs/config.json", "configuration file path")
	flag.StringVar(&pkconfFile, "pkconf", "./configs/pk.json", "PrivateKey configuration file path")
	flag.Float64Var(&multiple, "mul", 1, "multiple of gasPrice, actual_gasPrice = suggested_gasPrice * mul ")
	flag.BoolVar(&force, "force", false, "need force send override bind or not")
	flag.BoolVar(&pip4, "pip4", false, "deploy lp or token")
	flag.BoolVar(&all, "all", false, "shut/restart all in config file")
	flag.StringVar(&function, "func", "", "choose function to run:\n"+
		"#### CCM \n"+
		"  -func shutCCM -mul {1} -conf {./config.json} -pkconf {./pkconfig.json} [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  -func restartCCM -mul {1} -conf {./config.json} -pkconf {./pkconfig.json} [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  -func checkCCM \n"+
		"#### LockProxy \n"+
		"  -func shutToken -mul {1} -conf {./config.json} -token {./token.json} -pkconf {./pkconfig.json}\n"+
		"  -func bindToken -mul {1} -conf {./config.json} -token {./token.json} -pkconf {./pkconfig.json}\n"+
		"  -func bindSingleToken -mul {1} -conf {./config.json} -token {./token.json} -pkconf {./pkconfig.json} [fromChainId] [toChainId] \n"+
		"  -func shutSingleToken -mul {1} -conf {./config.json} -token {./token.json} -pkconf {./pkconfig.json} [fromChainId] [toChainId] \n"+
		"  -func checkUnbindToken \n"+
		"  -func checkBindToken \n"+
		"  -func bindProxy -mul {1} -conf {./config.json} -pkconf {./pkconfig.json}\n"+
		"  -func unbindProxy -mul {1} -conf {./config.json} -pkconf {./pkconfig.json}\n"+
		"  -func transferOwner -mul {1} -conf {./config.json} -pkconf {./pkconfig.json} [ChainID-1] [ChainID-2] ... [ChainID-n]\n"+
		"#### LockProxyPip4  \n"+
		" -func deployToken -mul {1} -conf {./config.json} -token {./token.json} -pkconf {./pkconfig.json} -pip4 {false} [fromChainId]\n"+
		"#### Swapper \n"+
		"  -func pauseSwapper \n"+
		"  -func unpauseSwapper \n"+
		"  -func unbindPool -chain ** -pool **\n"+
		"  -func checkSwapperPaused \n"+
		"  -func poolTokenMap -chain ** -pool **\n"+
		"  -func checkFeeCollected\n"+
		"  -func extractFeeSwapper\n"+
		"#### Wrapper \n"+
		"  -func extractFeeWrapper\n"+
		"  -func extractFeeWrapperO3\n"+
		"  {}contains default value")
	flag.Parse()
}

func main() {
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		log.Fatal("LoadConfig fail", err)
	}

	PKconfig, err := config.LoadPrivateKeyConfig(pkconfFile)
	if err != nil {
		log.Fatal("LoadConfig fail", err)
	}

	switch function {
	case "shutCCM":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseCCMPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Shutting down %s ...", netCfg.Name)
				paused, err := shutTools.CCMPaused(client, netCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if paused && !force {
					log.Warnf("CCM at chain %d is already shut, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if paused && force {
					log.Warnf("CCM at chain %d is already shut, still force shut", netCfg.PolyChainID)
				}
				err = shutTools.ShutCCM(multiple, client, netCfg, pkCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		if cnt == 0 {
			log.Info("Done.")
			return
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("CCM at chain %d has been shut down.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "restartCCM":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseCCMPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Restarting %s ...", netCfg.Name)
				paused, err := shutTools.CCMPaused(client, netCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if !paused && !force {
					log.Warnf("CCM at chain %d is already running, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if !paused && force {
					log.Warnf("CCM at chain %d is already running, still force restart", netCfg.PolyChainID)
				}
				err = shutTools.RestartCCM(multiple, client, netCfg, pkCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		if cnt == 0 {
			log.Info("Done.")
			return
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("CCM at chain %d has been restarted.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "shutToken":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all {
			args = tokenConfig.GetTokenIds()
		}
		sig := make(chan Msg, 10)
		var tokens []*Token
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			token := tokenConfig.GetToken(uint64(id))
			if token == nil {
				log.Errorf("token with chainId %d not found in %s", id, tokenFile)
				continue
			}
			address := token.Address

			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseLockProxyPrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			tokens = append(tokens, &Token{uint64(id), address, netCfg, pkCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Shutting down %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
				client, err := ethclient.Dial(tokens[i].NetCfg.Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", tokens[i].NetCfg.Provider, err)
					sig <- Msg{tokens[i].ChainId, err}
					return
				}
				for j := 0; j < len(tokens); j++ {
					if i == j {
						continue
					}
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].NetCfg.LockProxy, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						err = fmt.Errorf(
							"fail to shut %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					if len(toAsset) == 0 && !force {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is already shut, ignored",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
						continue
					} else if len(toAsset) == 0 && force {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is already shut, still force shut",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
					}
					err = shutTools.BindToken(
						multiple,
						client,
						tokens[i].NetCfg,
						tokens[i].PkCfg,
						tokens[i].Address,
						tokens[j].ChainId,
						nil)
					if err != nil {
						err = fmt.Errorf(
							"fail to shut %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					log.Infof("%s: %d =>to=> %d pair has be unbind", tokenConfig.Name, tokens[i].ChainId, tokens[j].ChainId)
				}
				sig <- Msg{tokens[i].ChainId, nil}
			}(i)
		}
		cnt := len(tokens)
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("%s at chain %d has been shut down.", tokenConfig.Name, msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "bindToken":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all {
			args = tokenConfig.GetTokenIds()
		}
		sig := make(chan Msg, 10)
		var tokens []*Token
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			token := tokenConfig.GetToken(uint64(id))
			if token == nil {
				log.Errorf("token with chainId %d not found in %s", id, tokenFile)
				continue
			}
			address := token.Address

			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseLockProxyPrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			tokens = append(tokens, &Token{uint64(id), address, netCfg, pkCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Binding %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
				client, err := ethclient.Dial(tokens[i].NetCfg.Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", tokens[i].NetCfg.Provider, err)
					sig <- Msg{tokens[i].ChainId, err}
					return
				}
				for j := 0; j < len(tokens); j++ {
					if i == j {
						continue
					}
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].NetCfg.LockProxy, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						err = fmt.Errorf(
							"fail to bind %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					if len(toAsset) != 0 && !force {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , ignored",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							toAsset)
						continue
					} else if len(toAsset) != 0 && force {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , still force bind",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							toAsset)
					}
					err = shutTools.BindToken(
						multiple,
						client,
						tokens[i].NetCfg,
						tokens[i].PkCfg,
						tokens[i].Address,
						tokens[j].ChainId,
						tokens[j].Address.Bytes())
					if err != nil {
						err = fmt.Errorf(
							"fail to bind %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					log.Infof("%s: %d =>to=> %d pair has be bind", tokenConfig.Name, tokens[i].ChainId, tokens[j].ChainId)
				}
				sig <- Msg{tokens[i].ChainId, nil}
			}(i)
		}
		cnt := len(tokens)
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("%s at chain %d has been bind.", tokenConfig.Name, msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "bindSingleToken":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if len(args) != 2 {
			log.Fatal("Arg num not match")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
		}
		token := tokenConfig.GetToken(uint64(id))
		if token == nil {
			log.Errorf("token with chainId %d not found in %s", id, tokenFile)
		}
		address := token.Address
		netCfg := conf.GetNetwork(uint64(id))
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}

		pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
		if pkCfg == nil {
			log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
		}

		err = pkCfg.ParseLockProxyPrivateKey()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fromAsset := &Token{uint64(id), address, netCfg, pkCfg}

		id, err = strconv.Atoi(args[1])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 1, args[1], err)
		}
		token = tokenConfig.GetToken(uint64(id))
		if token == nil {
			log.Errorf("token with chainId %d not found in %s", id, tokenFile)
		}
		address = token.Address
		netCfg = conf.GetNetwork(uint64(id))
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}
		toAsset := &Token{uint64(id), address, netCfg, pkCfg}

		log.Infof("Binding %x and %x from %d to %d ...", fromAsset.Address, toAsset.Address, fromAsset.ChainId, toAsset.ChainId)
		client, err := ethclient.Dial(fromAsset.NetCfg.Provider)
		if err != nil {
			log.Fatal("fail to dial %s , %s", fromAsset.NetCfg.Provider, err)
		}
		mappedAsset, err := shutTools.TokenMap(client, fromAsset.NetCfg, fromAsset.NetCfg.LockProxy, fromAsset.Address, toAsset.ChainId)
		if err != nil {
			log.Fatalf(
				"fail to bind %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				err)
		}
		if len(mappedAsset) != 0 && !force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , ignored",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				mappedAsset)
			log.Info("Done.")
			return
		} else if len(mappedAsset) != 0 && force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , still force bind",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				mappedAsset)
		}
		err = shutTools.BindToken(
			multiple,
			client,
			fromAsset.NetCfg,
			fromAsset.PkCfg,
			fromAsset.Address,
			toAsset.ChainId,
			toAsset.Address.Bytes())
		if err != nil {
			log.Fatalf(
				"fail to bind %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				err)
		}
		log.Infof("%s: %d =>to=> %d pair has been bind", tokenConfig.Name, fromAsset.ChainId, toAsset.ChainId)
		log.Info("Done.")
	case "shutSingleToken":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if len(args) != 2 {
			log.Fatal("Arg num not match")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
		}
		token := tokenConfig.GetToken(uint64(id))
		if token == nil {
			log.Errorf("token with chainId %d not found in %s", id, tokenFile)
		}
		address := token.Address
		netCfg := conf.GetNetwork(uint64(id))
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}

		pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
		if pkCfg == nil {
			log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
		}

		err = pkCfg.ParseLockProxyPrivateKey()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fromAsset := &Token{uint64(id), address, netCfg, pkCfg}
		id, err = strconv.Atoi(args[1])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 1, args[1], err)
		}
		toChainId := uint64(id)

		log.Infof("Shuting %x from %d to %d ...", fromAsset.Address, fromAsset.ChainId, toChainId)
		client, err := ethclient.Dial(fromAsset.NetCfg.Provider)
		if err != nil {
			log.Fatal("fail to dial %s , %s", fromAsset.NetCfg.Provider, err)
		}
		mappedAsset, err := shutTools.TokenMap(client, fromAsset.NetCfg, fromAsset.NetCfg.LockProxy, fromAsset.Address, toChainId)
		if err != nil {
			log.Fatalf(
				"fail to shut %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toChainId,
				err)
		}
		if len(mappedAsset) == 0 && !force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already shut, ignored",
				tokenConfig.Name,
				fromAsset.ChainId,
				toChainId)
			log.Info("Done.")
			return
		} else if len(mappedAsset) == 0 && force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already shut, still force shut",
				tokenConfig.Name,
				fromAsset.ChainId,
				toChainId,
				mappedAsset)
		}
		err = shutTools.BindToken(
			multiple,
			client,
			fromAsset.NetCfg,
			fromAsset.PkCfg,
			fromAsset.Address,
			toChainId,
			nil)
		if err != nil {
			log.Fatalf(
				"fail to shut %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toChainId,
				err)
		}
		log.Infof("%s: %d =>to=> %d pair has been shut", tokenConfig.Name, fromAsset.ChainId, toChainId)
		log.Info("Done.")
	case "checkUnbindToken":
		log.Info("Processing...")
		args := flag.Args()
		flag := -1
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all || len(args) == 0 {
			args = tokenConfig.GetTokenIds()
		} else if len(args) == 1 {
			flag, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("can not parse arg %d : %s , %v", 0, args[0], err)
			}
			args = tokenConfig.GetTokenIds()
		}
		var tokens []*Token
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			token := tokenConfig.GetToken(uint64(id))
			if token == nil {
				log.Errorf("token with chainId %d not found in %s", id, tokenFile)
				continue
			}
			address := token.Address
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			tokens = append(tokens, &Token{uint64(id), address, netCfg, pkCfg})
		}
		for i := 0; i < len(tokens); i++ {
			if (flag != -1) && (int(tokens[i].ChainId) != flag) {
				continue
			}
			func(i int) {
				log.Infof("Checking %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
				client, err := ethclient.Dial(tokens[i].NetCfg.Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", tokens[i].NetCfg.Provider, err)
					log.Errorf(err.Error())
					return
				}
				for j := 0; j < len(tokens); j++ {
					if i == j {
						continue
					}
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].NetCfg.LockProxy, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						log.Errorf(
							"fail to check %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
					}
					if len(toAsset) == 0 {
						log.Infof(
							"token %s from chain %d =>to=> chain %d is unbind",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
						continue
					} else {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is still bind, bind at %x",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							toAsset)
					}
				}
				log.Infof("Check %s at %s done", tokenConfig.Name, tokens[i].NetCfg.Name)
				log.Info("-------------------------------------------------------------")
			}(i)
		}
		log.Info("All Done.")
	case "checkBindToken":
		log.Info("Processing...")
		args := flag.Args()
		flag := -1
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all || len(args) == 0 {
			args = tokenConfig.GetTokenIds()
		} else if len(args) == 1 {
			flag, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("can not parse arg %d : %s , %v", 0, args[0], err)
			}
			args = tokenConfig.GetTokenIds()
		}
		var tokens []*Token
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			token := tokenConfig.GetToken(uint64(id))
			if token == nil {
				log.Errorf("token with chainId %d not found in %s", id, tokenFile)
				continue
			}
			address := token.Address
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			tokens = append(tokens, &Token{uint64(id), address, netCfg, pkCfg})
		}
		for i := 0; i < len(tokens); i++ {
			if (flag != -1) && (int(tokens[i].ChainId) != flag) {
				continue
			}
			func(i int) {
				log.Infof("Checking %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
				client, err := ethclient.Dial(tokens[i].NetCfg.Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", tokens[i].NetCfg.Provider, err)
					log.Errorf(err.Error())
					return
				}
				for j := 0; j < len(tokens); j++ {
					if i == j {
						continue
					}
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].NetCfg.LockProxy, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						log.Errorf(
							"fail to check %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
					}
					if len(toAsset) == 0 {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d has not been bind",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
						continue
					} else if bytes.Equal(toAsset, tokens[j].Address[:]) {
						log.Infof(
							"token %s from chain %d =>to=> chain %d is binded",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
					} else {
						log.Infof(
							"token %s from chain %d =>to=> chain %d is binded unexpectedly at %x",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							toAsset)
					}
				}
				log.Infof("Check %s at %s done", tokenConfig.Name, tokens[i].NetCfg.Name)
				log.Info("-------------------------------------------------------------")
			}(i)
		}
		log.Info("All Done.")
	case "checkCCM":
		log.Info("Processing...")
		args := flag.Args()
		if all || len(args) == 0 {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Checking %s ...", netCfg.Name)
				time.Sleep(500 * time.Millisecond)
				paused, err := shutTools.CCMPaused(client, netCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if paused {
					log.Warnf("CCM at chain %d has been shut down", netCfg.PolyChainID)
				} else {
					log.Infof("CCM at chain %d is running", netCfg.PolyChainID)
				}
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "pauseSwapper":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseSwapperPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Pausing swapper at %s ...", netCfg.Name)
				paused, err := shutTools.SwapperPaused(client, netCfg, pkCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if paused && !force {
					log.Warnf("Swapper at chain %d is already paused, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if paused && force {
					log.Warnf("Swapper at chain %d is already paused, still force pause", netCfg.PolyChainID)
				}
				err = shutTools.PauseSwapper(multiple, client, netCfg, pkCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("Swapper at chain %d has been paused.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "unpauseSwapper":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseSwapperPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Unpausing swapper at %s ...", netCfg.Name)
				paused, err := shutTools.SwapperPaused(client, netCfg, pkCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if !paused && !force {
					log.Warnf("Swapper at chain %d is not paused, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if !paused && force {
					log.Warnf("Swapper at chain %d is not paused, still force unpause", netCfg.PolyChainID)
				}
				err = shutTools.UnpauseSwapper(multiple, client, netCfg, pkCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("Swapper at chain %d has been unpaused.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "unbindPool":
		log.Info("Processing...")
		netCfg := conf.GetNetwork(chainId)
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in config file", chainId)
		}

		pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
		if pkCfg == nil {
			log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
		}
		err = pkCfg.ParseSwapperPrivateKey()
		if err != nil {
			log.Fatalf(err.Error())
		}
		client, err := ethclient.Dial(netCfg.Provider)
		if err != nil {
			log.Fatalf("fail to dial client %s of network %d", netCfg.Provider, chainId)
		}
		log.Infof("Unbindind pool %d at %s ...", poolId, netCfg.Name)
		currentBind, err := shutTools.PoolTokenMap(client, netCfg, poolId)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if (currentBind == ADDRESS_ZERO) && !force {
			log.Warnf("Pool %d at chain %d is not registered, ignored", poolId, netCfg.PolyChainID)
			return
		} else if (currentBind == ADDRESS_ZERO) && force {
			log.Warnf("pool %d at chain %d is not registered, still force unbind", poolId, netCfg.PolyChainID)
		}
		err = shutTools.RegisterPool(multiple, client, netCfg, pkCfg, poolId, ADDRESS_ZERO)
		log.Info("Done")
	case "checkSwapperPaused":
		log.Info("Processing...")
		args := flag.Args()
		if all || len(args) == 0 {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Checking %s ...", netCfg.Name)
				time.Sleep(500 * time.Millisecond)
				paused, err := shutTools.SwapperPaused(client, netCfg, pkCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if paused {
					log.Warnf("Swapper at chain %d has been paused", netCfg.PolyChainID)
				} else {
					log.Infof("Swapper at chain %d is running", netCfg.PolyChainID)
				}
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "poolTokenMap":
		log.Info("Processing...")
		netCfg := conf.GetNetwork(chainId)
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in config file", chainId)
		}
		client, err := ethclient.Dial(netCfg.Provider)
		if err != nil {
			log.Fatalf("fail to dial client %s of network %d", netCfg.Provider, chainId)
		}
		log.Infof("Checking pool %d at %s ...", poolId, netCfg.Name)
		currentBind, err := shutTools.PoolTokenMap(client, netCfg, poolId)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if currentBind == ADDRESS_ZERO {
			log.Warnf("Pool %d at chain %d is not registered ", poolId, netCfg.PolyChainID)
			return
		} else {
			log.Infof("pool %d at chain %d is registered, currnet poolTokenAddress %x", poolId, netCfg.PolyChainID, currentBind)
		}
		log.Info("Done")
	case "bindProxy":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		var netCfgs []*config.Network
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}

			netCfgs = append(netCfgs, netCfg)
		}
		for i := 0; i < len(netCfgs); i++ {
			pkCfg := PKconfig.GetSenderPrivateKey(netCfgs[i].PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfgs[i].PrivateKeyNo)
			}

			err = pkCfg.ParseLockProxyPrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			go func(i int) {
				log.Infof("binding proxy at %s...", netCfgs[i].Name)
				client, err := ethclient.Dial(netCfgs[i].Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", netCfgs[i].Provider, err)
					sig <- Msg{netCfgs[i].PolyChainID, err}
					return
				}
				for j := 0; j < len(netCfgs); j++ {
					if i == j {
						continue
					}
					toProxy, err := shutTools.ProxyHashMap(client, netCfgs[i], netCfgs[j].PolyChainID)
					if err != nil {
						err = fmt.Errorf(
							"fail to bind proxy from chain %d =>to=> chain %d , %s",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							err)
						sig <- Msg{netCfgs[i].PolyChainID, err}
						return
					}
					if len(toProxy) != 0 && !force {
						log.Warnf(
							"proxy from chain %d =>to=> chain %d is already bind, current bind proxy: %x , ignored",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							toProxy)
						continue
					} else if len(toProxy) != 0 && force {
						log.Warnf(
							"proxy from chain %d =>to=> chain %d is already bind, current bind proxy: %x , still force bind",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							toProxy)
					}
					err = shutTools.BindProxyHash(
						multiple,
						client,
						netCfgs[i],
						pkCfg,
						netCfgs[j].PolyChainID,
						netCfgs[j].LockProxy.Bytes())
					if err != nil {
						err = fmt.Errorf(
							"fail to bind proxy from chain %d =>to=> chain %d , %s",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							err)
						sig <- Msg{netCfgs[i].PolyChainID, err}
						return
					}
					log.Infof("bindProxy : %d =>to=> %d proxy has be bind", netCfgs[i].PolyChainID, netCfgs[j].PolyChainID)
				}
				sig <- Msg{netCfgs[i].PolyChainID, nil}
			}(i)
		}
		cnt := len(netCfgs)
		if cnt == 0 {
			log.Info("Done.")
			return
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("proxy at chain %d has been bind.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "unbindProxy":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		var netCfgs []*config.Network
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}

			netCfgs = append(netCfgs, netCfg)
		}
		for i := 0; i < len(netCfgs); i++ {
			pkCfg := PKconfig.GetSenderPrivateKey(netCfgs[i].PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfgs[i].PrivateKeyNo)
			}

			err = pkCfg.ParseLockProxyPrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			go func(i int) {
				log.Infof("unbinding proxy at %s...", netCfgs[i].Name)
				client, err := ethclient.Dial(netCfgs[i].Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", netCfgs[i].Provider, err)
					sig <- Msg{netCfgs[i].PolyChainID, err}
					return
				}
				for j := 0; j < len(netCfgs); j++ {
					if i == j {
						continue
					}
					toProxy, err := shutTools.ProxyHashMap(client, netCfgs[i], netCfgs[j].PolyChainID)
					if err != nil {
						err = fmt.Errorf(
							"fail to unbind proxy from chain %d =>to=> chain %d , %s",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							err)
						sig <- Msg{netCfgs[i].PolyChainID, err}
						return
					}
					if len(toProxy) == 0 && !force {
						log.Warnf(
							"proxy from chain %d =>to=> chain %d is not bind , ignored",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID)
						continue
					} else if len(toProxy) == 0 && force {
						log.Warnf(
							"proxy from chain %d =>to=> chain %d is not bind, still force unbind",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID)
					}
					err = shutTools.BindProxyHash(
						multiple,
						client,
						netCfgs[i],
						pkCfg,
						netCfgs[j].PolyChainID,
						nil)
					if err != nil {
						err = fmt.Errorf(
							"fail to unbind proxy from chain %d =>to=> chain %d , %s",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							err)
						sig <- Msg{netCfgs[i].PolyChainID, err}
						return
					}
					log.Infof("unbindProxy : %d =>to=> %d proxy has be unbind", netCfgs[i].PolyChainID, netCfgs[j].PolyChainID)
				}
				sig <- Msg{netCfgs[i].PolyChainID, nil}
			}(i)
		}
		cnt := len(netCfgs)
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("proxy at chain %d has been unbind.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "checkBindProxy":
		log.Info("Processing...")
		args := flag.Args()
		flag := -1
		if all || len(args) == 0 {
			args = conf.GetNetworkIds()
		} else if len(args) == 1 {
			flag, err = strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("can not parse arg %d : %s , %v", 0, args[0], err)
			}
			args = conf.GetNetworkIds()
		}
		var netCfgs []*config.Network
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}
			netCfgs = append(netCfgs, netCfg)
		}
		for i := 0; i < len(netCfgs); i++ {
			if (flag != -1) && (int(netCfgs[i].PolyChainID) != flag) {
				continue
			}
			func(i int) {
				log.Infof("checking proxy at %s...", netCfgs[i].Name)
				client, err := ethclient.Dial(netCfgs[i].Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", netCfgs[i].Provider, err)
					log.Errorf(err.Error())
					return
				}
				for j := 0; j < len(netCfgs); j++ {
					if i == j {
						continue
					}
					toProxy, err := shutTools.ProxyHashMap(client, netCfgs[i], netCfgs[j].PolyChainID)
					if err != nil {
						err = fmt.Errorf(
							"fail to bind proxy from chain %d =>to=> chain %d , %s",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							err)
						log.Errorf(err.Error())
						return
					}
					if len(toProxy) == 0 {
						log.Warnf(
							"proxy from chain %d =>to=> chain %d has not been bind",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID)
						continue
					} else if bytes.Equal(toProxy, netCfgs[j].LockProxy.Bytes()) {
						log.Infof(
							"proxy from chain %d =>to=> chain %d is binded",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID)
					} else {
						log.Infof(
							"proxy from chain %d =>to=> chain %d is binded unexpectedly at %x",
							netCfgs[i].PolyChainID,
							netCfgs[j].PolyChainID,
							toProxy)
					}
				}
				log.Infof("check proxy at %s done", netCfgs[i].Name)
				log.Info("-------------------------------------------------------------")
			}(i)
		}
		log.Info("All Done.")
	case "bindSingleProxy":
		log.Info("Processing...")
		args := flag.Args()
		if len(args) != 2 {
			log.Fatal("Arg num not match")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
		}
		fromProxy := conf.GetNetwork(uint64(id))
		if fromProxy == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}

		pkCfg := PKconfig.GetSenderPrivateKey(fromProxy.PrivateKeyNo)
		if pkCfg == nil {
			log.Errorf("privatekey with chainId %d not found in PKconfig file", fromProxy.PrivateKeyNo)
		}

		err = pkCfg.ParseLockProxyPrivateKey()
		if err != nil {
			log.Fatalf("%v", err)
		}

		id, err = strconv.Atoi(args[1])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 1, args[1], err)
		}
		toProxy := conf.GetNetwork(uint64(id))
		if toProxy == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}

		log.Infof("Binding proxy from %d to %d ...", fromProxy.PolyChainID, toProxy.PolyChainID)
		client, err := ethclient.Dial(fromProxy.Provider)
		if err != nil {
			log.Fatal("fail to dial %s , %s", fromProxy.Provider, err)
		}
		mappedProxy, err := shutTools.ProxyHashMap(client, fromProxy, toProxy.PolyChainID)
		if err != nil {
			log.Fatalf(
				"fail to bind proxy from chain %d =>to=> chain %d , %s",
				fromProxy.PolyChainID,
				toProxy.PolyChainID,
				err)
		}
		if len(mappedProxy) != 0 && !force {
			log.Warnf(
				"proxy from chain %d =>to=> chain %d is already bind , ignored",
				fromProxy.PolyChainID,
				toProxy.PolyChainID)
			log.Info("Done.")
			return
		} else if len(mappedProxy) != 0 && force {
			log.Warnf(
				"proxy from chain %d =>to=> chain %d is already bind, current bind proxy: %x , still force bind",
				fromProxy.PolyChainID,
				toProxy.PolyChainID,
				mappedProxy)
		}
		err = shutTools.BindProxyHash(
			multiple,
			client,
			fromProxy,
			pkCfg,
			toProxy.PolyChainID,
			toProxy.LockProxy.Bytes())
		if err != nil {
			log.Fatalf(
				"fail to bind proxy from chain %d =>to=> chain %d , %s",
				fromProxy.PolyChainID,
				toProxy.PolyChainID,
				err)
		}
		log.Infof("bindProxy: %d =>to=> %d pair has been bind", fromProxy.PolyChainID, toProxy.PolyChainID)
		log.Info("Done.")
	case "checkFeeCollected":
		log.Info("Processing...")
		args := flag.Args()
		if all || len(args) == 0 {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 20)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			func() {
				log.Infof("Checking %s ...", netCfg.Name)
				time.Sleep(500 * time.Millisecond)
				bs, err := client.BalanceAt(context.Background(), netCfg.Swapper, nil)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				bw, err := client.BalanceAt(context.Background(), netCfg.Wrapper, nil)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				bo, err := client.BalanceAt(context.Background(), netCfg.WrapperO3, nil)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				if netCfg.Swapper == ADDRESS_ZERO {
					bs = big.NewInt(0)
				}
				if netCfg.Wrapper == ADDRESS_ZERO {
					bw = big.NewInt(0)
				}
				if netCfg.WrapperO3 == ADDRESS_ZERO {
					bo = big.NewInt(0)
				}
				balanceSwapper := big.NewFloat(0)
				balanceWrapper := big.NewFloat(0)
				balanceWrapperO3 := big.NewFloat(0)
				balanceSwapper.SetString(bs.String())
				balanceWrapper.SetString(bw.String())
				balanceWrapperO3.SetString(bo.String())
				balanceSwapper.Quo(balanceSwapper, big.NewFloat(math.Pow(10, 18)))
				balanceWrapper.Quo(balanceWrapper, big.NewFloat(math.Pow(10, 18)))
				balanceWrapperO3.Quo(balanceWrapperO3, big.NewFloat(math.Pow(10, 18)))
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				log.Infof("Balance of swapper (nativeToken) at %s is %f ", netCfg.Name, balanceSwapper)
				log.Infof("Balance of wrapper (nativeToken) at %s is %f ", netCfg.Name, balanceWrapper)
				log.Infof("Balance of wrapperO3 (nativeToken) at %s is %f ", netCfg.Name, balanceWrapperO3)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "extractFeeSwapper":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}
			err = pkCfg.ParseSwapperFeeCollectorPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Extract fee from swapper at %s ...", netCfg.Name)

				balance, err := client.BalanceAt(context.Background(), netCfg.Swapper, nil)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				zeroBalance := balance.Int64() == 0
				if zeroBalance && !force {
					log.Warnf("Swapper at chain %d do not have balance, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if zeroBalance && force {
					log.Warnf("Swapper at chain %d do not have balance, still force extractFee", netCfg.PolyChainID)
				}

				err = shutTools.ExtractFeeSwapper(multiple, client, netCfg, pkCfg, ADDRESS_ZERO)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("Fee has been taken from swapper at chain %d .", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "extractFeeWrapper":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}
			err = pkCfg.ParseWrapperFeeCollectorPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Extract fee from wrapper at %s ...", netCfg.Name)

				balance, err := client.BalanceAt(context.Background(), netCfg.Wrapper, nil)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				zeroBalance := balance.Int64() == 0
				if zeroBalance && !force {
					log.Warnf("Wrapper at chain %d do not have balance, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if zeroBalance && force {
					log.Warnf("Wrapper at chain %d do not have balance, still force extractFee", netCfg.PolyChainID)
				}

				err = shutTools.ExtractFeeWrapper(multiple, client, netCfg, pkCfg, ADDRESS_ZERO)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("Fee has been taken from wrapper at chain %d .", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "extractFeeWrapperO3":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseWrapperO3FeeCollectorPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Extract fee from wrapper at %s ...", netCfg.Name)

				balance, err := client.BalanceAt(context.Background(), netCfg.WrapperO3, nil)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				zeroBalance := balance.Int64() == 0
				if zeroBalance && !force {
					log.Warnf("Wrapper at chain %d do not have balance, ignored", netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if zeroBalance && force {
					log.Warnf("Wrapper at chain %d do not have balance, still force extractFee", netCfg.PolyChainID)
				}

				err = shutTools.ExtractFeeWrapperO3(multiple, client, netCfg, pkCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("Fee has been taken from wrapper at chain %d .", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "deployToken":
		log.Info("Token Depoly...")
		args := flag.Args()
		if len(args) != 1 {
			log.Fatal("Arg num not match")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
		}
		netCfg := conf.GetNetwork(uint64(id))
		if netCfg == nil {
			log.Errorf("network with chainId %d not found in config file", id)
		}

		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		client, err := ethclient.Dial(netCfg.Provider)
		if err != nil {
			log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
		}

		pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
		if pkCfg == nil {
			log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
		}
		fmt.Println("enter", netCfg.PrivateKeyNo)

		LptokenAddress, err := shutTools.DeployToken(multiple, client, netCfg, pkCfg, tokenConfig, pip4)
		if err != nil {
			log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
		}

		log.Info("token deploy success, hash is %s", LptokenAddress)

		token := tokenConfig.GetToken(uint64(id))
		if pip4 {
			token.LPAddress = LptokenAddress
		} else {
			token.Address = LptokenAddress
		}
		res2, err := json.MarshalIndent(tokenConfig, " ", "	")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile(tokenFile, res2, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Info("token write success")
	case "singleBindLPandAsset":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if len(args) != 2 {
			log.Fatal("Arg num not match")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
		}
		token := tokenConfig.GetToken(uint64(id))
		if token == nil {
			log.Errorf("token with chainId %d not found in %s", id, tokenFile)
		}
		address := token.Address
		lpaddress := token.LPAddress
		netCfg := conf.GetNetwork(uint64(id))
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}

		pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
		if pkCfg == nil {
			log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
		}

		err = pkCfg.ParseLockProxyPip4PrivateKey()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fromAsset := &Token{ChainId: uint64(id), Address: address, NetCfg: netCfg, PkCfg: pkCfg}

		id, err = strconv.Atoi(args[1])
		if err != nil {
			log.Errorf("can not parse arg %d : %s , %v", 1, args[1], err)
		}
		token = tokenConfig.GetToken(uint64(id))
		if token == nil {
			log.Errorf("token with chainId %d not found in %s", id, tokenFile)
		}
		address = token.Address
		tolpaddress := token.LPAddress
		netCfg = conf.GetNetwork(uint64(id))
		if netCfg == nil {
			log.Fatalf("network with chainId %d not found in %s", id, confFile)
		}
		toAsset := &Token{uint64(id), address, netCfg, pkCfg}

		log.Infof("Binding %x and %x from %d to %d ...", fromAsset.Address, toAsset.Address, fromAsset.ChainId, toAsset.ChainId)
		client, err := ethclient.Dial(fromAsset.NetCfg.Provider)
		if err != nil {
			log.Fatal("fail to dial %s , %s", fromAsset.NetCfg.Provider, err)
		}
		lpmappedAsset, err := shutTools.AssetLPMap(client, fromAsset.NetCfg, fromAsset.Address)
		if err != nil {
			log.Fatalf(
				"fail to bind %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				lpmappedAsset,
				err)
		}
		if len(lpmappedAsset) != 0 && !force && lpmappedAsset == lpaddress {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , ignored",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				lpmappedAsset)
			log.Info("Done.")
			return
		} else if len(lpmappedAsset) != 0 && force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , still force bind",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				lpmappedAsset)
		}
		mappedAsset, err := shutTools.TokenMap(client, fromAsset.NetCfg, fromAsset.NetCfg.LockProxyPip4, fromAsset.Address, toAsset.ChainId)
		if err != nil {
			log.Fatalf(
				"fail to bind %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				err)
		}
		if len(mappedAsset) != 0 && !force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , ignored",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				mappedAsset)
			log.Info("Done.")
			return
		} else if len(mappedAsset) != 0 && force {
			log.Warnf(
				"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , still force bind",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				mappedAsset)
		}
		err = shutTools.BindLPandAsset(
			multiple,
			client,
			fromAsset.NetCfg,
			fromAsset.PkCfg,
			fromAsset.Address,
			lpaddress,
			toAsset.ChainId,
			toAsset.Address.Bytes(),
			tolpaddress.Bytes())
		if err != nil {
			log.Fatalf(
				"fail to bind %s from chain %d =>to=> chain %d , %s",
				tokenConfig.Name,
				fromAsset.ChainId,
				toAsset.ChainId,
				err)
		}
		log.Infof("%s: %d =>to=> %d pair has been bind", tokenConfig.Name, fromAsset.ChainId, toAsset.ChainId)
		log.Info("Done.")
	case "bindLPandAsset":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all {
			args = tokenConfig.GetTokenIds()
		}
		sig := make(chan Msg, 10)
		var tokens []*LPToken
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
			}
			token := tokenConfig.GetToken(uint64(id))
			if token == nil {
				log.Errorf("token with chainId %d not found in %s", id, tokenFile)
			}
			address := token.Address
			lpaddress := token.LPAddress
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in %s", id, confFile)
			}
			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}
			err = pkCfg.ParseLockProxyPip4PrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			tokens = append(tokens, &LPToken{uint64(id), address, lpaddress, netCfg, pkCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Binding %x and %x at  %d ...", tokens[i].Address, tokens[i].LpAddress, tokens[i].ChainId)
				client, err := ethclient.Dial(tokens[i].NetCfg.Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", tokens[i].NetCfg.Provider, err)
					sig <- Msg{tokens[i].ChainId, err}
					return
				}
				lpmappedAsset, err := shutTools.AssetLPMap(client, tokens[i].NetCfg, tokens[i].Address)
				if err != nil {
					log.Fatalf(
						"fail to bind %s from chain %d, %s",
						tokenConfig.Name,
						tokens[i].ChainId,
						err)
				}
				if len(lpmappedAsset) != 0 && !force && lpmappedAsset == tokens[i].LpAddress {
					log.Warnf(
						"lptoken %s from chain %d is already bind, current bind token: %x , ignored",
						tokenConfig.Name,
						tokens[i].ChainId,
						lpmappedAsset)
					log.Info("Done.")
				} else if len(lpmappedAsset) != 0 && force {
					log.Warnf(
						"lptoken %s from chain %d is already bind, current bind token: %x , ignored",
						tokenConfig.Name,
						tokens[i].ChainId,
						lpmappedAsset)
				}
				for j := 0; j < len(tokens); j++ {
					if i == j {
						continue
					}
					mappedAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].NetCfg.LockProxyPip4, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						log.Fatalf(
							"fail to bind %s from chain %d =>to=> chain %d, %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
					}
					if len(mappedAsset) != 0 && !force {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , ignored",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							mappedAsset)
						log.Info("Done.")
						sig <- Msg{tokens[i].ChainId, nil}
						return
					} else if len(mappedAsset) != 0 && force {
						log.Warnf(
							"token %s from chain %d =>to=> chain %d is already bind, current bind token: %x , still force bind",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							mappedAsset)
					}
					err = shutTools.BindLPandAsset(
						multiple,
						client,
						tokens[i].NetCfg,
						tokens[i].PkCfg,
						tokens[i].Address,
						tokens[i].LpAddress,
						tokens[j].ChainId,
						tokens[j].Address.Bytes(),
						tokens[j].LpAddress.Bytes())
					if err != nil {
						log.Fatalf(
							"fail to bind %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
					}
					log.Infof("%s: %d =>to=> %d pair has been bind", tokenConfig.Name, tokens[i].ChainId, tokens[j].ChainId)
					log.Info("Done.")
				}
				sig <- Msg{tokens[i].ChainId, nil}
			}(i)
		}
		cnt := len(tokens)
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("%s at chain %d has been bind.", tokenConfig.Name, msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done")
				break
			}
		}
	case "bindLPandAssetBatch":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all {
			args = tokenConfig.GetTokenIds()
		}
		var tokenlist []*Tokenlist
		sig := make(chan Msg, 10)
		for i := 0; i < len(args); i++ {
			var tokens Tokenlist
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
			}
			token := tokenConfig.GetToken(uint64(id))
			if token == nil {
				log.Errorf("token with chainId %d not found in %s", id, tokenFile)
			}
			for j := 0; j < len(args); j++ {
				if i == j {
					continue
				}
				id_j, err := strconv.Atoi(args[j])
				if err != nil {
					log.Errorf("can not parse arg %d : %s , %v", 0, args[0], err)
				}
				tokenj := tokenConfig.GetToken(uint64(id_j))
				if tokenj == nil {
					log.Errorf("token with chainId %d not found in %s", id, tokenFile)
				}
				tokens.Address = append(tokens.Address, token.Address)
				tokens.LpAddress = append(tokens.LpAddress, token.LPAddress)
				tokens.toChainId = append(tokens.toChainId, tokenj.PolyChainId)
				tokens.toAsset = append(tokens.toAsset, tokenj.Address.Bytes())
				tokens.toLPAsset = append(tokens.toLPAsset, tokenj.LPAddress.Bytes())
			}
			tokenlist = append(tokenlist, &Tokenlist{tokens.ChainId, tokens.Address, tokens.LpAddress, tokens.toChainId, tokens.toAsset, tokens.toLPAsset, tokens.NetCfg, tokens.PkCfg})
		}
		for i := 0; i <= len(tokenlist); i++ {
			go func(i int) {
				log.Infof("Binding %x and %x at  %d ...", tokenlist[i].Address, tokenlist[i].LpAddress, tokenlist[i].ChainId)
				client, err := ethclient.Dial(tokenlist[i].NetCfg.Provider)
				if err != nil {
					err = fmt.Errorf("fail to dial %s , %s", tokenlist[i].NetCfg.Provider, err)
					sig <- Msg{tokenlist[i].ChainId, err}
					return
				}
				err = shutTools.BindLPandAeestBatch(
					multiple,
					client,
					tokenlist[i].NetCfg,
					tokenlist[i].PkCfg,
					tokenlist[i].Address,
					tokenlist[i].LpAddress,
					tokenlist[i].toChainId,
					tokenlist[i].toAsset,
					tokenlist[i].toLPAsset,
				)
				if err != nil {
					log.Fatalf(
						"fail to bind %s at chain %d",
						tokenConfig.Name,
						tokenlist[i].ChainId,
						err)
				}
				log.Infof("%s: %d =>to=> %d pair has been bind", tokenConfig.Name, tokenlist[i].ChainId, tokenlist[i].toChainId)
				log.Info("Done.")
				sig <- Msg{tokenlist[i].ChainId, nil}
			}(i)
		}
		cnt := len(tokenlist)
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("%s at chain %d has been bind.", tokenConfig.Name, msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done")
				break
			}
		}
	case "transferOwner":
		log.Info("Processing...")
		args := flag.Args()
		if all {
			args = conf.GetNetworkIds()
		}
		sig := make(chan Msg, 10)
		cnt := 0

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("\nPlease type in new owner for LockProxy:")
		newOwnerStr, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("can not read input, %v", err)
			return
		}
		newOwnerStr = strings.TrimSuffix(newOwnerStr, "\n")
		newOwner := common.FromHex(newOwnerStr)
		if !common.IsHexAddress(newOwnerStr) {
			log.Errorf("invalid input, not address")
			return
		}

		fmt.Println("\nPlease repeat new owner for LockProxy:")
		newOwnerStrRepeat, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("can not read input, %v", err)
			return
		}
		newOwnerStrRepeat = strings.TrimSuffix(newOwnerStrRepeat, "\n")
		newOwnerRepeat := common.FromHex(newOwnerStrRepeat)
		if !common.IsHexAddress(newOwnerStrRepeat) {
			log.Errorf("invalid input, not address")
			return
		}

		if bytes.Compare(newOwner, newOwnerRepeat) != 0 {
			log.Errorf("mismatched input!")
			return
		}

		fmt.Println(fmt.Sprintf("\nPlease make sure %x is the new owner of LockProxy, type in yes to continue:", newOwner))
		reply, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("can not read input, %v", err)
			return
		}
		reply = strings.TrimSuffix(reply, "\n")
		if strings.Compare(reply, "yes") != 0 {
			log.Warnf("user canceled Execution")
			return
		}

		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
				continue
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
				continue
			}

			pkCfg := PKconfig.GetSenderPrivateKey(netCfg.PrivateKeyNo)
			if pkCfg == nil {
				log.Errorf("privatekey with chainId %d not found in PKconfig file", netCfg.PrivateKeyNo)
			}

			err = pkCfg.ParseLockProxyPrivateKey()
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Transfer LockProxy Owner of chain %s ...", netCfg.Name)
				owner, err := shutTools.LockProxyOwner(client, netCfg)
				if err != nil {
					sig <- Msg{netCfg.PolyChainID, err}
					return
				}
				isOwner := bytes.Compare(owner.Bytes(), newOwner) == 0
				if isOwner && !force {
					log.Warnf("%x is already the LockProxyOwner in chain %d, ignored", newOwner, netCfg.PolyChainID)
					sig <- Msg{netCfg.PolyChainID, err}
					return
				} else if isOwner && force {
					log.Warnf("%x is already the LockProxyOwner in chain %d, still force shut", newOwner, netCfg.PolyChainID)
				}
				err = shutTools.TransferOwnership(multiple, client, netCfg, pkCfg, common.BytesToAddress(newOwner))
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
		}
		if cnt == 0 {
			log.Info("Done.")
			return
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("LockProxyOwnership at chain %d has been transfered.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "Approve":
		//todo
		log.Info("Approve Token Depoly...")

	default:
		log.Fatal("unknown function", function)
	}

}
