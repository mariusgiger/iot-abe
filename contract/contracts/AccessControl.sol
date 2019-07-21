pragma solidity ^0.5.0;

contract AccessControl {
    address public owner;
    string public pubKey;

    event AccessRequested(
        address indexed _from,
        uint _time
    );

    event AccessGranted(
        address indexed _requester,
        string _key,
        bytes32[] attrs
    );

    event DevicePolicyUpdated(
        address indexed _device,
        string policy
    );

    event DevicePolicyDeleted(
        address indexed _device
    );

    struct AccessPolicy {
        string key;
        bool pending;
        bytes32[] attrs; //one attribute must be smaller than 32 bytes
    }

    struct DevicePolicy {
        string policy;
    }

    mapping(address => AccessPolicy) public acl;
    mapping(address => DevicePolicy) public devices;

    constructor(string memory _pubKey) public {
        owner = msg.sender;
        pubKey = _pubKey;
    }

    // isOwner checks if caller is owner of the contract
    modifier isOwner() {
        require(owner == msg.sender, "Caller is not owner");
        _;
    }

    function requestAccess() public {
        //TODO check if access has been granted already
        acl[msg.sender].pending = true;
        emit AccessRequested(msg.sender, block.timestamp);
    }

    function grantAccess(address _user, string memory _key, bytes32[] memory _attrs) public isOwner() {
        acl[_user].key = _key;
        acl[_user].attrs = _attrs;
        acl[_user].pending = false;

        emit AccessGranted(_user, _key, _attrs);
    }

    function setDevicePolicy(address _device, string memory _policy) public isOwner() {
        devices[_device].policy = _policy; 

        emit DevicePolicyUpdated(_device, _policy);
    }

    function removeDevicePolicy(address _device) public isOwner() {
        delete devices[_device];
        
        emit DevicePolicyDeleted(_device);
    }
}