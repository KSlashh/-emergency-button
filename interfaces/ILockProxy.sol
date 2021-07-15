pragma solidity ^0.6.0;


contract ILockProxy {
    function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes memory toAssetHash) external returns (bool) {}
}