package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/KSlashh/emergency-button/config"
	"github.com/KSlashh/emergency-button/log"
	"github.com/KSlashh/emergency-button/shutTools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var confFile string
var tokenFile string
var function string
var multiple float64
var force bool

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
	flag.StringVar(&tokenFile, "token", "./token.json", "token configuration file path")
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.Float64Var(&multiple, "mul", 1, "multiple of gasPrice, actual_gasPrice = suggested_gasPrice * mul ")
	flag.BoolVar(&force, "force", false, "need force send override bind or not")
	flag.StringVar(&function, "func", "", "choose function to run:\n"+
		"  -func shutCCM -mul {1} -conf {./config.json} [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  -func restartCCM -mul {1} -conf {./config.json} [ChainID-1] [ChainID-2] ... [ChainID-n] \n"+
		"  -func shutToken -mul {1} -conf {./config.json} -token {./token.json} \n"+
		"  -func rebindToken -mul {1} -conf {./config.json} -token {./token.json} \n"+
<<<<<<< Updated upstream
		"  ## {} contains default value")
=======
		"  {}contains default value")
>>>>>>> Stashed changes
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
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
			}
			log.Infof("Shutting down %s ...", netCfg.Name)
			err = netCfg.PhrasePrivateKey()
			if err != nil {
				log.Fatalf("%v", err)
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
			}
			go func() {
<<<<<<< Updated upstream
=======
				log.Infof("Shutting down %s ...", netCfg.Name)
>>>>>>> Stashed changes
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
		sig := make(chan Msg, 10)
		cnt := 0
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
			}
			netCfg := conf.GetNetwork(uint64(id))
			if netCfg == nil {
				log.Errorf("network with chainId %d not found in config file", id)
			}
			log.Infof("Restarting %s ...", netCfg.Name)
			err = netCfg.PhrasePrivateKey()
			if err != nil {
				log.Errorf("%v", err)
			}
			client, err := ethclient.Dial(netCfg.Provider)
			if err != nil {
				log.Errorf("fail to dial client %s of network %d", netCfg.Provider, id)
			}
			go func() {
<<<<<<< Updated upstream
=======
				log.Infof("Restarting %s ...", netCfg.Name)
>>>>>>> Stashed changes
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
		sig := make(chan Msg, 10)
		var tokens []*Token
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
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
			err = netCfg.PhrasePrivateKey()
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
		sig := make(chan Msg, 10)
		var tokens []*Token
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			if err != nil {
				log.Errorf("can not parse arg %d : %s , %v", i, args[i], err)
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
			err = netCfg.PhrasePrivateKey()
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
<<<<<<< Updated upstream
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
=======
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
		/*
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
						log.Fatalf("can not parse arg %d : %s , %v", i, args[i], err)
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
								multiple,
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
						log.Fatalf("can not parse arg %d : %s , %v", i, args[i], err)
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
								multiple,
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
		*/
>>>>>>> Stashed changes
	default:
		log.Fatal("unknown function", function)
	}
}
