# Emergency Button
紧急情况下快速关停poly合约

## Compile
````
go build -o boom
````

## setup config file
config.json , 例如：
````
{
  "Networks": [
    {
      "PolyChainID":2,
      "Name":"eth-test",
      "Provider":"https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "CCMPAddress":"0xb600c8a2e8852832B75DB9Da1A3A1c173eAb28d8",
      "LockProxyAddress":"0xD8aE73e06552E270340b63A8bcAbf9277a1aac99"
    },

    {
      "PolyChainID":79,
      "Name":"bsc-test",
      "Provider":"https://data-seed-prebsc-1-s1.binance.org:8545/",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "CCMPAddress":"0x441C035446c947a97bD36b425B67907244576990",
      "LockProxyAddress":"0x097Ae585BfEf78DDC8E266ABCb840dAF7265130c"
    },

    {
      "PolyChainID":7,
      "Name":"heco-test",
      "Provider":"https://http-testnet.hecochain.com",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "CCMPAddress":"0xc5757b5d22984E534004cC7Fb1D59eD14EC510a5",
      "LockProxyAddress":"0x4a76E52600C6285029c8f7c52183cf86282cA5b8"
    },

    {
      "PolyChainID":200,
      "Name":"ok-test",
      "Provider":"https://exchaintestrpc.okex.org/",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "CCMPAddress":"0x38917884b397447227fb45cbA0342F1bFf7A3470",
      "LockProxyAddress":"0x74cE7D56cd1b5AEe9A3345A490b5Ed768134C7D4"
    }
  ]
}
````
token.json , 例如：
````
{
  "Name":"sampleToken",
  "Tokens":[
    {
      "PolyChainId":2,
      "Address":"0xad3f96ae966ad60347f31845b7e4b333104c52fb"
    },
    {
      "PolyChainId":79,
      "Address":"0xB60e03E6973B1d0b90a763f5B64C48ca7cB8c2d1"
    },
    {
      "PolyChainId":7,
      "Address":"0xc4e419CC0945dC9860A73B3A2cAcAA12aD7CF3B8"
    },
    {
      "PolyChainId":12,
      "Address":"0x9a3658864aa2ccc63fa61eaad5e4f65fa490ca7d"
    }
  ]
}
````
## shut down CCM
例如 关掉 bsc(79) 和 ok(200) 的CCM：
````
./boom -func shutCCM -conf sampleTestnetConfig.json shutCCM 79 200
````

## restart CCM
````
./boom -func restartCCM -conf sampleTestnetConfig.json shutCCM 79 200
````

## shut down Token
例如 停掉eth，bsc及heco之间的USDT
配置如下 USDT.json
````
{
  "Name":"USDT",
  "Tokens":[
    {
      "PolyChainId":2,
      "Address":"0xad3f96ae966ad60347f31845b7e4b333104c52fb"
    },
    {
      "PolyChainId":79,
      "Address":"0xB60e03E6973B1d0b90a763f5B64C48ca7cB8c2d1"
    },
    {
      "PolyChainId":7,
      "Address":"0xc4e419CC0945dC9860A73B3A2cAcAA12aD7CF3B8"
    }
  ]
}
````
执行：
````
./boom -func shutToken -conf sampleTestnetConfig.json -token USDT.json 
````

## rebind Token
````
./boom -func rebindToken -conf sampleTestnetConfig.json -token USDT.json
````

## multiple gas price
用`-mul`来提高gas price，例如 `-mul 6` 意思是按照推荐gas price的6倍发出交易
````
./boom -mul 6 -func shutCCM -conf sampleConfig.json shutCCM 79 200
````