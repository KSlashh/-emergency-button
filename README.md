# Emergency Button
紧急情况下快速关停poly合约

## Compile
````
go build -o boom
````

## setup config file
config.json , 例如：
(注意，`PrivateKey` 和 `KeyStore` 至少提供一项)
````
{
  "Networks": [
    {
      "PolyChainID":2,
      "Name":"eth",
      "Provider":"https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "CCMPOwnerKeyStore":"{put CCM owner's sk keystore path here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "LockProxyOwnerKeyStore":"{put LockProxy owner's sk keystore path here}",
      "CCMPAddress":"0x5a51e2ebf8d136926b9ca7b59b60464e7c44d2eb",
      "LockProxyAddress":"0x250e76987d838a75310c34bf422ea9f1ac4cc906"
    },

    {
      "PolyChainID":6,
      "Name":"bsc",
      "Provider":"https://bsc-dataseed.binance.org",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "CCMPOwnerKeyStore":"{put CCM owner's sk keystore path here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "LockProxyOwnerKeyStore":"{put LockProxy owner's sk keystore path here}",
      "CCMPAddress":"0xabd7f7b89c5fd5d0aef06165f8173b1b83d7d5c9",
      "LockProxyAddress":"0x2f7ac9436ba4b548f9582af91ca1ef02cd2f1f03"
    },

    {
      "PolyChainID":7,
      "Name":"heco",
      "Provider":"https://http-mainnet-node.huobichain.com",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "CCMPOwnerKeyStore":"{put CCM owner's sk keystore path here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "LockProxyOwnerKeyStore":"{put LockProxy owner's sk keystore path here}",
      "CCMPAddress":"0xabd7f7b89c5fd5d0aef06165f8173b1b83d7d5c9",
      "LockProxyAddress":"0x020c15e7d08a8ec7d35bcf3ac3ccbf0bbf2704e6"
    },

    {
      "PolyChainID":12,
      "Name":"ok",
      "Provider":"https://exchainrpc.okex.org/",
      "CCMPOwnerPrivateKey":"{put CCM owner's sk here}",
      "CCMPOwnerKeyStore":"{put CCM owner's sk keystore path here}",
      "LockProxyOwnerPrivateKey":"{put LockProxy owner's sk here}",
      "LockProxyOwnerKeyStore":"{put LockProxy owner's sk keystore path here}",
      "CCMPAddress":"0x4739fe955be4704bcb7d6a699823f5b29217baf6",
      "LockProxyAddress":"0x9a3658864aa2ccc63fa61eaad5e4f65fa490ca7d"
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
      "PolyChainId":6,
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
例如 关掉 bsc(6) 和 ok(12) 的CCM：
````
./boom -func shutCCM -conf sampleTestnetConfig.json shutCCM 6 12
````

## restart CCM
````
./boom -func restartCCM -conf sampleTestnetConfig.json shutCCM 6 12
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
      "PolyChainId":6,
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
./boom -func shutToken -conf sampleTestnetConfig.json -token USDT.json 2 6 7
````

## rebind Token
````
./boom -func rebindToken -conf sampleTestnetConfig.json -token USDT.json 2 6 7
````

## multiple gas price
用`-mul`来提高gas price，例如 `-mul 6` 意思是按照推荐gas price的6倍发出交易
````
./boom -mul 6 -func shutCCM -conf sampleConfig.json shutCCM 6 12
````

## 使用keystore
上述例子为在config文件中直接配置了私钥的情况下运行。
如果需要使用keystore，在配置相应字段指定keystore文件的路径即可。
执行时会要求提供相应的password.