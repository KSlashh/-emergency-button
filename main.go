package main

import (
	"flag"
	"fmt"
	"github.com/KSlashh/emergency-button/config"
	"github.com/KSlashh/emergency-button/log"
	"github.com/KSlashh/emergency-button/shutTools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"strconv"
)

var confFile string
var function string

type Msg struct {
	ChainId uint64
	Err     error
}

type token struct {
	ChainId uint64
	Address string
	NetCfg  *config.Network
}

func init() {
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.StringVar(&function, "func", "", "choose function to run:\n"+
		"  shutCCM [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  restartCCM [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  shutToken [chainID-1] [tokenAddress-1] ... [chainId-n] [tokenAddress-n] \n"+
		"  rebindToken [chainID-1] [tokenAddress-1] ... [chainId-n] [tokenAddress-n]")
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
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %s", i, args[i], err)
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
			}
			go func() {
				log.Infof("Shutting down %s ...",netCfg.Name)
				err = shutTools.ShutCCM(client, netCfg)
				sig<-Msg{netCfg.PolyChainID,err}
			}()
			cnt += 1
		}
		for msg := range sig {
			cnt -= 1
			if msg.Err != nil {
				log.Error(msg.Err)
			} else {
				log.Infof("CCM at chain %d has been shut down.", msg.ChainId)
				if cnt == 0 {
					log.Info("Done.")
					break
				}
			}
		}
	case "restartCCM":
		log.Info("Processing...")
		args := flag.Args()
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %s", i, args[i], err)
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
			}
			go func() {
				log.Infof("Restarting %s ...",netCfg.Name)
				err = shutTools.RestartCCM(client, netCfg)
				sig<-Msg{netCfg.PolyChainID,err}
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
		sig := make(chan Msg, 10)
		if len(args)%2 != 0 {
			log.Fatalf("invalid arg amount %d ,must be even", len(args))
		}
		var tokens []*token
		for i := 0; i < len(args); i += 2 {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Fatalf("can not parse arg %d : %s , %s", i, args[i], err)
			}
			address := args[i+1]
			if !common.IsHexAddress(address) {
				log.Fatalf("%s is not a valid address", address)
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in config file", id)
			}
			tokens = append(tokens, &token{uint64(id), address, netCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Shutting down token for %s...",tokens[i].NetCfg.Name)
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
					err = shutTools.BindToken(
						client,
						tokens[i].NetCfg,
						common.HexToAddress(tokens[i].Address),
						tokens[j].ChainId,
						nil)
					if err != nil {
						err = fmt.Errorf(
							"fail to shut bind from chain %d =>to=> chain %d , %s",
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					log.Infof("%d =>to=> %d pair has be unbind",tokens[i].ChainId,tokens[j].ChainId)
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
				log.Infof("Token at chain %d has been shut down.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	case "rebindToken":
		log.Info("Processing...")
		args := flag.Args()
		sig := make(chan Msg, 10)
		if len(args)%2 != 0 {
			log.Fatalf("invalid arg amount %d ,must be even", len(args))
		}
		var tokens []*token
		for i := 0; i < len(args); i += 2 {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Fatalf("can not parse arg %d : %s , %s", i, args[i], err)
			}
			address := args[i+1]
			if !common.IsHexAddress(address) {
				log.Fatalf("%s is not a valid address", address)
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Fatalf("network with chainId %d not found in config file", id)
			}
			tokens = append(tokens, &token{uint64(id), address, netCfg})
		}
		for i := 0; i < len(tokens); i++ {
			go func(i int) {
				log.Infof("Rebinding for %s...",tokens[i].NetCfg.Name)
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
					err = shutTools.BindToken(
						client,
						tokens[i].NetCfg,
						common.HexToAddress(tokens[i].Address),
						tokens[j].ChainId,
						common.FromHex(tokens[j].Address))
					if err != nil {
						err = fmt.Errorf(
							"fail to shut bind from chain %d =>to=> chain %d , %s",
							tokens[i].ChainId,
							tokens[j].ChainId,
							err)
						sig <- Msg{tokens[i].ChainId, err}
						return
					}
					log.Infof("%d =>to=> %d pair has be rebind",tokens[i].ChainId,tokens[j].ChainId)
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
				log.Infof("Token at chain %d has been rebind.", msg.ChainId)
			}
			if cnt == 0 {
				log.Info("Done.")
				break
			}
		}
	default:
		log.Fatal("unknown function", function)
	}
}
