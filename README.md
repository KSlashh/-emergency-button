# Emergency Button
紧急情况下快速关停poly合约

## Compile
````
go build -o boom
````

## setup config file
例如：
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

## shut down CCM
例如 关掉 bsc(79) 和 ok(200) 的CCM：
````
./boom -func shutCCM -conf sampleConfig.json shutCCM 79 200
````

## restart CCM
````
./boom -func restartCCM -conf sampleConfig.json shutCCM 79 200
````

## shut down Token
例如 停掉eth，bsc及heco之间的USDT
````
eth USDT : 0xad3f96ae966ad60347f31845b7e4b333104c52fb
bsc USDT : 0x23F5075740c2C99C569FfD0768c383A92d1a4aD7
heco USDT : 0x7698Da475B3132F37E40DE8C222d7D74d3A4172d	
````
执行：
````
./boom -func shutToken -conf sampleConfig.json 2 0xad3f96ae966ad60347f31845b7e4b333104c52fb 79 0x23F5075740c2C99C569FfD0768c383A92d1a4aD7 6 0x7698Da475B3132F37E40DE8C222d7D74d3A4172d 
````

## rebind Token
````
./boom -func rebindToken -conf sampleConfig.json 2 0xad3f96ae966ad60347f31845b7e4b333104c52fb 79 0x23F5075740c2C99C569FfD0768c383A92d1a4aD7 6 0x7698Da475B3132F37E40DE8C222d7D74d3A4172d 
````

## multiple gas price
用`-mul`来提高gas price，例如 `-mul 6` 意思是按照推荐gas price的十倍发出交易
````
./boom  -mul 6 -func shutCCM -conf sampleConfig.json shutCCM 79 200
````