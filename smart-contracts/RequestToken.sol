// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/utils/Address.sol";

contract RequestToken is ERC721 {
    using Counters for Counters.Counter;
    using Address for address;

    Counters.Counter private _tokenIds;
    address public serverAccount;

    // Mapping to track whether a token is active.
    mapping(uint256 => bool) public tokenActive;

    // Events for logging actions.
    event TokenMinted(uint256 indexed tokenId, address indexed recipient);
    event TokenTransferred(uint256 indexed tokenId, address indexed from, address indexed to);
    event TokenBurned(uint256 indexed tokenId);

    constructor(address _serverAccount) ERC721("RequestToken", "RTKN") {
        require(_serverAccount != address(0), "Server account cannot be zero address");
        require(!_serverAccount.isContract(), "Server account must be an externally owned account");
        serverAccount = _serverAccount;
    }

    // Mint a token with additional validations.
    function mintToken(address recipient) public returns (uint256) {
        require(msg.sender == serverAccount, "Only server can mint tokens");
        require(recipient != address(0), "Recipient cannot be zero address");
        require(!recipient.isContract(), "Recipient must be an externally owned account");

        _tokenIds.increment();
        uint256 newTokenId = _tokenIds.current();

        _mint(recipient, newTokenId);
        tokenActive[newTokenId] = true; // Mark the token as active

        emit TokenMinted(newTokenId, recipient);
        return newTokenId;
    }

    // Verify token ownership. Requires that the token exists.
    function verifyTokenOwner(uint256 tokenId, address claimedOwner) public view returns (bool) {
        require(_exists(tokenId), "Token does not exist");
        return ownerOf(tokenId) == claimedOwner;
    }

    // Transfer the token back to the server.
    // Ensures token exists, is active, and that the caller is the token owner.
    function transferToServer(uint256 tokenId) public {
        require(_exists(tokenId), "Token does not exist");
        require(tokenActive[tokenId], "Token is not active");
        require(ownerOf(tokenId) == msg.sender, "Only the owner can transfer the token");

        _transfer(msg.sender, serverAccount, tokenId);
        emit TokenTransferred(tokenId, msg.sender, serverAccount);
    }

    // Burn the token.
    // Only the server can burn the token, and it must be active.
    function burnToken(uint256 tokenId) public {
        require(_exists(tokenId), "Token does not exist");
        require(tokenActive[tokenId], "Token is already inactive");
        require(ownerOf(tokenId) == serverAccount, "Only the server can burn the token");

        tokenActive[tokenId] = false; // Mark token as inactive
        _burn(tokenId);
        emit TokenBurned(tokenId);
    }
}
