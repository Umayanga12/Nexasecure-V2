// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/utils/Address.sol";

contract AuthToken is ERC721 {
    using Counters for Counters.Counter;
    using Address for address;

    Counters.Counter private _tokenIds;
    address public serverAccount;

    // Allowed values for user types and token types
    string[] private allowedUserTypes = ["admin", "user", "guest"];
    string[] private allowedTokenTypes = ["login", "session"];

    // Structure for token details
    struct TokenDetails {
        string username;
        string userType;  // e.g., "admin", "user"
        string tokenType; // e.g., "login", "session"
        bool active;      // To mark if the token is currently active
    }

    // Structure for registered user details.
    struct UserInfo {
        bool exists;
        string username;
        string userType;
    }

    // Mapping token ID to its details.
    mapping(uint256 => TokenDetails) public tokenMetadata;
    // Mapping account addresses to registered user info.
    mapping(address => UserInfo) public registeredUsers;

    // Events for logging actions.
    event UserRegistered(address indexed user, string username, string userType);
    event TokenMinted(uint256 indexed tokenId, address indexed recipient, string username, string userType, string tokenType);
    event TokenTransferred(uint256 indexed tokenId, address indexed from, address indexed to);
    event TokenBurned(uint256 indexed tokenId);

    constructor(address _serverAccount, string memory _name, string memory _symbol) ERC721(_name, _symbol) {
        require(_serverAccount != address(0), "Server account cannot be zero address");
        require(!_serverAccount.isContract(), "Server account cannot be a contract");
        serverAccount = _serverAccount;
    }

    // Internal helper to check if a string is in an allowed list.
    function _isAllowedValue(string memory value, string[] memory allowedValues) internal pure returns (bool) {
        bytes memory valueBytes = bytes(value);
        require(valueBytes.length > 0, "Value cannot be empty");
        for (uint256 i = 0; i < allowedValues.length; i++) {
            if (keccak256(bytes(allowedValues[i])) == keccak256(valueBytes)) {
                return true;
            }
        }
        return false;
    }

    // Register a user with their account address.
    // This simulates off‑chain identity verification by recording valid details on-chain.
    function registerUser(
        address user, 
        string memory username, 
        string memory userType
    ) public {
        require(msg.sender == serverAccount, "Only server can register users");
        require(user != address(0), "User cannot be zero address");
        require(!user.isContract(), "User address cannot be a contract");
        require(bytes(username).length > 0, "Username cannot be empty");
        require(_isAllowedValue(userType, allowedUserTypes), "Invalid user type");

        registeredUsers[user] = UserInfo(true, username, userType);
        emit UserRegistered(user, username, userType);
    }

    // Mint an authentication token with extended validations.
    // Validates that the recipient is a valid, registered account.
    function mintToken(
        address recipient, 
        string memory username, 
        string memory userType, 
        string memory tokenType
    ) public returns (uint256) {
        require(msg.sender == serverAccount, "Only server can mint tokens");
        require(recipient != address(0), "Recipient cannot be zero address");
        require(!recipient.isContract(), "Recipient address must be an externally owned account");
        require(bytes(username).length > 0, "Username cannot be empty");
        require(_isAllowedValue(userType, allowedUserTypes), "Invalid user type");
        require(_isAllowedValue(tokenType, allowedTokenTypes), "Invalid token type");

        // Validate that the recipient is registered and the provided details are correct.
        UserInfo memory user = registeredUsers[recipient];
        require(user.exists, "User is not registered");
        require(keccak256(bytes(username)) == keccak256(bytes(user.username)), "Provided username does not match registered username");
        require(keccak256(bytes(userType)) == keccak256(bytes(user.userType)), "Provided user type does not match registered user type");

        // Increment token counter and mint the token.
        _tokenIds.increment();
        uint256 newTokenId = _tokenIds.current();
        _mint(recipient, newTokenId);

        // Store token details and mark it as active.
        tokenMetadata[newTokenId] = TokenDetails(username, userType, tokenType, true);

        emit TokenMinted(newTokenId, recipient, username, userType, tokenType);
        return newTokenId;
    }

    // Retrieve token details ensuring the token exists.
    function getTokenDetails(uint256 tokenId) public view returns (string memory, string memory, string memory, bool) {
        require(_exists(tokenId), "Token does not exist");
        TokenDetails memory details = tokenMetadata[tokenId];
        return (details.username, details.userType, details.tokenType, details.active);
    }

    // Verify token ownership.
    function verifyTokenOwner(uint256 tokenId, address claimedOwner) public view returns (bool) {
        require(_exists(tokenId), "Token does not exist");
        return ownerOf(tokenId) == claimedOwner;
    }

    // Transfer the token back to the server. Only the token owner can transfer and the token must be active.
    function transferToServer(uint256 tokenId) public {
        require(_exists(tokenId), "Token does not exist");
        require(ownerOf(tokenId) == msg.sender, "Only the owner can transfer the token");
        require(tokenMetadata[tokenId].active, "Token is not active");

        _transfer(msg.sender, serverAccount, tokenId);
        emit TokenTransferred(tokenId, msg.sender, serverAccount);
    }

    // Burn the token after deactivating it. Only the server can burn tokens.
    function burnToken(uint256 tokenId) public {
        require(_exists(tokenId), "Token does not exist");
        require(ownerOf(tokenId) == serverAccount, "Only the server can burn the token");
        require(tokenMetadata[tokenId].active, "Token is already inactive");

        tokenMetadata[tokenId].active = false;
        _burn(tokenId);
        emit TokenBurned(tokenId);

        // Clean up stored metadata.
        delete tokenMetadata[tokenId];
    }
}

