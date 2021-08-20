package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/KSlashh/emergency-button/config"
	"github.com/KSlashh/emergency-button/log"
	"github.com/KSlashh/emergency-button/shutTools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ADDRESS_ZERO common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

var confFile string
var tokenFile string
var function string
var multiple float64
var force bool
var all bool
var poolId uint64
var chainId uint64

type Msg struct {
	ChainId uint64
	Err     error
}

type Token struct {
	ChainId uint64
	Address common.Address
	NetCfg  *config.Network
}

func init() {
	flag.Uint64Var(&poolId, "pool", 0, "pool id if needed")
	flag.Uint64Var(&chainId, "chain", 0, "chain id if single chainId needed")
	flag.StringVar(&tokenFile, "token", "./token.json", "token configuration file path")
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.Float64Var(&multiple, "mul", 1, "multiple of gasPrice, actual_gasPrice = suggested_gasPrice * mul ")
	flag.BoolVar(&force, "force", false, "need force send override bind or not")
	flag.BoolVar(&all, "all", false, "shut/restart all in config file")
	flag.StringVar(&function, "func", "", "choose function to run:\n"+
		"  -func shutCCM -mul {1} -conf {./config.json} [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  -func restartCCM -mul {1} -conf {./config.json} [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  -func shutToken -mul {1} -conf {./config.json} -token {./token.json} \n"+
		"  -func rebindToken -mul {1} -conf {./config.json} -token {./token.json} \n"+
		"  -func bindSingleToken -mul {1} -conf {./config.json} -token {./token.json} [fromChainId] [toChainId] \n"+
		"  -func shutSingleToken -mul {1} -conf {./config.json} -token {./token.json} [fromChainId] [toChainId] \n"+
		"  -func pauseSwapper \n"+
		"  -func unpauseSwapper \n"+
		"  -func unbindPool \n"+
		"  -func checkUnbind \n"+
		"  -func checkBind \n"+
		"  -func checkCCM \n"+
		"  -func checkSwapperPaused \n"+
		"  -func poolTokenMap \n"+
		"  {}contains default value")
	flag.Parse()
}

func main() {
	conf, err := config.LoadConfig(confFile)
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
			err = netCfg.PhraseCCMPrivateKey()
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
				err = shutTools.ShutCCM(multiple, client, netCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
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
			err = netCfg.PhraseCCMPrivateKey()
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
				err = shutTools.RestartCCM(multiple, client, netCfg)
				sig <- Msg{netCfg.PolyChainID, err}
			}()
			cnt += 1
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
			err = netCfg.PhraseLockProxyPrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			tokens = append(tokens, &Token{uint64(id), address, netCfg})
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
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].Address, tokens[j].ChainId)
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
	case "rebindToken":
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
			err = netCfg.PhraseLockProxyPrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			tokens = append(tokens, &Token{uint64(id), address, netCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Rebinding %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
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
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						err = fmt.Errorf(
							"fail to rebind %s from chain %d =>to=> chain %d , %s",
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
						tokens[i].Address,
						tokens[j].ChainId,
						tokens[j].Address.Bytes())
					if err != nil {
						err = fmt.Errorf(
							"fail to rebind %s from chain %d =>to=> chain %d , %s",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					log.Infof("%s: %d =>to=> %d pair has be rebind", tokenConfig.Name, tokens[i].ChainId, tokens[j].ChainId)
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
				log.Infof("%s at chain %d has been rebind.", tokenConfig.Name, msg.ChainId)
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
		err = netCfg.PhraseLockProxyPrivateKey()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fromAsset := &Token{uint64(id), address, netCfg}

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
		toAsset := &Token{uint64(id), address, netCfg}

		log.Infof("Binding %x and %x from %d to %d ...", fromAsset.Address, toAsset.Address, fromAsset.ChainId, toAsset.ChainId)
		client, err := ethclient.Dial(fromAsset.NetCfg.Provider)
		if err != nil {
			log.Fatal("fail to dial %s , %s", fromAsset.NetCfg.Provider, err)
		}
		mappedAsset, err := shutTools.TokenMap(client, fromAsset.NetCfg, fromAsset.Address, toAsset.ChainId)
		if err != nil {
			log.Fatalf(
				"fail to rebind %s from chain %d =>to=> chain %d , %s",
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
			fromAsset.Address,
			toAsset.ChainId,
			toAsset.Address.Bytes())
		if err != nil {
			log.Fatalf(
				"fail to rebind %s from chain %d =>to=> chain %d , %s",
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
		err = netCfg.PhraseLockProxyPrivateKey()
		if err != nil {
			log.Fatalf("%v", err)
		}
		fromAsset := &Token{uint64(id), address, netCfg}
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
		mappedAsset, err := shutTools.TokenMap(client, fromAsset.NetCfg, fromAsset.Address, toChainId)
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
	case "checkUnbind":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all || len(args) == 0 {
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
			tokens = append(tokens, &Token{uint64(id), address, netCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Checking %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
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
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						log.Fatalf(
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
							"token %s from chain %d =>to=> chain %d is still bind",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
					}
				}
				sig <- Msg{tokens[i].ChainId, nil}
			}(i)
		}
		cnt := len(tokens)
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
	case "checkBind":
		log.Info("Processing...")
		args := flag.Args()
		tokenConfig, err := config.LoadToken(tokenFile)
		if err != nil {
			log.Fatal("LoadToken fail", err)
		}
		if all || len(args) == 0 {
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
			tokens = append(tokens, &Token{uint64(id), address, netCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Checking %s at %s...", tokenConfig.Name, tokens[i].NetCfg.Name)
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
					toAsset, err := shutTools.TokenMap(client, tokens[i].NetCfg, tokens[i].Address, tokens[j].ChainId)
					if err != nil {
						log.Fatalf(
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
					} else {
						log.Infof(
							"token %s from chain %d =>to=> chain %d is binded",
							tokenConfig.Name,
							tokens[i].ChainId,
							tokens[j].ChainId)
					}
				}
				sig <- Msg{tokens[i].ChainId, nil}
			}(i)
		}
		cnt := len(tokens)
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
			err = netCfg.PhraseSwapperPrivateKey()
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
				paused, err := shutTools.SwapperPaused(client, netCfg)
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
				err = shutTools.PauseSwapper(multiple, client, netCfg)
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
			err = netCfg.PhraseSwapperPrivateKey()
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
				paused, err := shutTools.SwapperPaused(client, netCfg)
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
				err = shutTools.UnpauseSwapper(multiple, client, netCfg)
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
		err = netCfg.PhraseSwapperPrivateKey()
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
		err = shutTools.RegisterPool(multiple, client, netCfg, poolId, ADDRESS_ZERO)
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
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
				continue
			}
			go func() {
				log.Infof("Checking %s ...", netCfg.Name)
				time.Sleep(500 * time.Millisecond)
				paused, err := shutTools.SwapperPaused(client, netCfg)
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
	default:
		log.Fatal("unknown function", function)
	}

}
