pragma solidity ^0.6.0;


contract ILockProxy {
    mapping(address => mapping(uint64 => bytes)) public assetHashMap;
    function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes memory toAssetHash) external returns (bool) {}
}