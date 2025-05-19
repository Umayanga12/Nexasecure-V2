#include <Preferences.h>
#include <AES.h>
#include <AESLib.h>
#include <SHA256.h>
#include <uECC.h>
const char* AUTH_WALLET_PRIV = "307702010104209c748bbfb4892adc891cc0fe4af81ea7dd0806beefa2e8b23ec16a37d81c86aea00a06082a8648ce3d030107a14403420004f85cc5e5b94e5f223150210d9c7a7be7be6f371bcc8ebb07dc4e8290af7419eecba844e2d60c7ec6644f1a186bef89aea10ddb175cc5f301c1351583b9c759bd";
const char* AUTH_WALLET_PUB = "3059301306072a8648ce3d020106082a8648ce3d03010703420004f85cc5e5b94e5f223150210d9c7a7be7be6f371bcc8ebb07dc4e8290af7419eecba844e2d60c7ec6644f1a186bef89aea10ddb175cc5f301c1351583b9c759bd";
const char* REQ_WALLET_PRIV = "3077020101042004fdefddd650d8ebef42d7744025178adc3ec0bdaebb2a1656e8478650253719a00a06082a8648ce3d030107a144034200042a7ed6e5de2228b01add17dc157aaac53bad659bb4d7816a41cfa40e12b1d5b339cbd69b1fd68d50de93a5f93a94decd83c5063b477daa544330c086fd66f590";
const char* REQ_WALLET_PUB = "3059301306072a8648ce3d020106082a8648ce3d030107034200042a7ed6e5de2228b01add17dc157aaac53bad659bb4d7816a41cfa40e12b1d5b339cbd69b1fd68d50de93a5f93a94decd83c5063b477daa544330c086fd66f590";

Preferences preferences;
AES aes;
AESLib aesLib;

byte aesKey[16];
bool isAuthenticated = false;
bool isFirstTimeSetup = false;

void hexToBin(const char* hex, uint8_t* bin, size_t binLen) {
  size_t len = strlen(hex);
  for (size_t i = 0; i < len / 2 && i < binLen; i++) {
    char byteStr[3] = {hex[2 * i], hex[2 * i + 1], '\0'};
    bin[i] = (uint8_t)strtol(byteStr, NULL, 16);
  }
}

void deriveKey(const char* password, uint8_t* key, size_t keyLen) {
  uint8_t salt[16];
  preferences.getBytes("salt", salt, 16);
  memset(key, 0, keyLen);
  size_t passLen = strlen(password);
  for (size_t i = 0; i < keyLen && i < passLen; i++) {
    key[i] = (uint8_t)password[i] ^ salt[i % 16];
  }
}

bool verifyPassword(const char* password) {
  uint8_t storedKey[32];
  uint8_t inputKey[32];
  preferences.getBytes("key_hash", storedKey, 32);
  deriveKey(password, inputKey, 32);
  return memcmp(storedKey, inputKey, 32) == 0;
}

void encryptPrivateKey(const char* privHex, const uint8_t* key, uint8_t* iv, uint8_t* enc, size_t* encLen) {
  size_t binLen = strlen(privHex) / 2;
  uint8_t* bin = (uint8_t*)malloc(binLen);
  hexToBin(privHex, bin, binLen);

  size_t paddedLen = ((binLen + 15) / 16) * 16;
  uint8_t* padded = (uint8_t*)malloc(paddedLen);
  memcpy(padded, bin, binLen);
  uint8_t paddingValue = paddedLen - binLen;
  for (size_t i = binLen; i < paddedLen; i++) {
    padded[i] = paddingValue;
  }

  aes.set_key(key, 256);
  aes.cbc_encrypt(padded, enc, paddedLen / 16, iv);
  *encLen = paddedLen;

  free(bin);
  free(padded);
}

void StoreWalletDetails(uint8_t* key) {
  preferences.putBytes("key_hash", key, 32);

  if (!preferences.isKey("auth_priv_enc")) {
    uint8_t authIV[16];
    for (int i = 0; i < 16; i++) authIV[i] = random(256);
    size_t authEncLen;
    uint8_t authEnc[512];
    encryptPrivateKey(AUTH_WALLET_PRIV, key, authIV, authEnc, &authEncLen);
    preferences.putBytes("auth_priv_iv", authIV, 16);
    preferences.putBytes("auth_priv_enc", authEnc, authEncLen);
    preferences.putString("auth_pub", AUTH_WALLET_PUB);
  }

  if (!preferences.isKey("req_priv_enc")) {
    uint8_t reqIV[16];
    for (int i = 0; i < 16; i++) reqIV[i] = random(256);
    size_t reqEncLen;
    uint8_t reqEnc[512];
    encryptPrivateKey(REQ_WALLET_PRIV, key, reqIV, reqEnc, &reqEncLen);
    preferences.putBytes("req_priv_iv", reqIV, 16);
    preferences.putBytes("req_priv_enc", reqEnc, reqEncLen);
    preferences.putString("req_pub", REQ_WALLET_PUB);
  }

  Serial.println("Wallet setup complete. Authenticate with 'AUTH <password>'.");
}

// Store NFT data (optimized: no base64)
// Store NFT data (plaintext, no encryption)
void storeNFT(const String& wallet, const String& nftData) {
  if (nftData.length() > 256) {
    Serial.println("NFT data too large");
    return;
  }
  String keyName = "nft-" + wallet;
  preferences.putString(keyName.c_str(), nftData);
 // Serial.println("NFT stored for " + wallet);
}
// Get NFT data
// Get NFT data (plaintext)
String getNFT(const String& wallet) {
  String keyName = "nft-" + wallet;
  String nftData = preferences.getString(keyName.c_str());
  if (nftData == "") {
    return "No NFT";
  }
  return nftData;
}
// Remove NFT data
// Remove NFT data
void removeNFT(const String& wallet) {
  String keyName = "nft-" + wallet;
  preferences.remove(keyName.c_str());
}


String signMessage(const String& wallet, const String& message) {
  String keyName = wallet == "AUTH" ? "auth_priv_enc" : "req_priv_enc";;
  size_t storedLen = preferences.getBytesLength(keyName.c_str());
  if (storedLen == 0) return "Wallet not set";
  uint8_t* storedData = (uint8_t*)malloc(storedLen);
  preferences.getBytes(keyName.c_str(), storedData, storedLen);
  if (storedLen < 16) {
    free(storedData);
    return "Invalid data";
  }
  uint8_t iv[16];
  memcpy(iv, storedData, 16);
  size_t ciphertextLen = storedLen - 16;
  uint8_t* ciphertext = storedData + 16;
  uint8_t* decrypted = (uint8_t*)malloc(ciphertextLen);
  aes.set_key(aesKey, 256); // if using 256-bit keys
aes.cbc_decrypt(ciphertext, decrypted, ciphertextLen / 16, iv);
  uint8_t padValue = decrypted[ciphertextLen - 1];
//  if (padValue > 16 || padValue == 0) {
//    free(storedData);
//    free(decrypted);
//    return "Invalid padding";
//  }
  size_t privKeyLen = ciphertextLen - padValue;
//  if (privKeyLen != 32) {
//    free(storedData);
//    free(decrypted);
//    return "Invalid private key length";
//  }
  uint8_t privKey[32];
  memcpy(privKey, decrypted, 32);
  free(storedData);
  free(decrypted);

  SHA256 sha256;
  sha256.update((const uint8_t*)message.c_str(), message.length());
  uint8_t hash[32];
  sha256.finalize(hash, 32);

  uint8_t signature[64];
//  if (!uECC_sign(privKey, hash, 32, signature, uECC_secp256k1())) {
//    memset(privKey, 0, 32);
//    return "Sign failed";
//  }

  char sigStr[129];
  for (int i = 0; i < 64; i++) {
    snprintf(sigStr + i * 2, sizeof(sigStr) - i * 2, "%02x", signature[i]);
  }

  memset(privKey, 0, 32);
  return String(sigStr);
}


void setup() {
  Serial.begin(115200);
  delay(1000);
  Serial.println("Secure Wallet Starting...");
  preferences.begin("wallet", false);

  if (preferences.getBytesLength("salt") == 0) {
    isFirstTimeSetup = true;
    Serial.println("First-time setup. Set password with 'SET_PASS <password>'.");
  } else {
    Serial.println("Wallet ready. Authenticate with 'PASS <password>'.");
  }
}

void loop() {
  if (Serial.available()) {
    String command = Serial.readStringUntil('\n');
    command.trim();

    if (isFirstTimeSetup && command.startsWith("SET_PASS ")) {
      String password = command.substring(9);
      uint8_t salt[16];
      for (int i = 0; i < 16; i++) {
        salt[i] = random(256);
      }
      preferences.putBytes("salt", salt, 16);

      uint8_t key[32];
      deriveKey(password.c_str(), key, 32);
      StoreWalletDetails(key);
      isFirstTimeSetup = false;
      Serial.println("Password set successfully. Authenticate with 'PASS <password>'.");
    } else if (isFirstTimeSetup) {
      Serial.println("First-time setup required. Use 'SET_PASS <password>'.");
    } else if (command.startsWith("PASS ")) {
      String password = command.substring(5);
      if (verifyPassword(password.c_str())) {
        isAuthenticated = true;
        Serial.println("PASSWORD_OK");
      } else {
        Serial.println("FAIL");
      }
    } else if (!isAuthenticated) {
      Serial.println("Please authenticate first with 'PASS <password>'.");
    } else if (command == "GET_ADDR_AUTH") {
      Serial.println(preferences.getString("auth_pub"));
    } else if (command == "GET_ADDR_REQ") {
      Serial.println(preferences.getString("req_pub"));
    } else if (command == "LOGOUT" && isAuthenticated) {
        isAuthenticated = false;
        memset(aesKey, 0, 16);
        Serial.println("Logged out");
    } else if (command.startsWith("SIGN_MSG_AUTH ")) {
      String msg = command.substring(14);
      Serial.println("SIG_AUTH " + signMessage("AUTH", msg));
    }
    else if (command.startsWith("SIGN_MSG_REQ ")) {
      String msg = command.substring(13);
      Serial.println("SIG_REQ " + signMessage("REQ", msg));
    }
    else if (command.startsWith("SET_NFT_AUTH ")) {
        String nft = command.substring(13);
        storeNFT("AUTH", nft);
        Serial.println("NFT_AUTH stored");
      }
      else if (command.startsWith("SET_NFT_REQ ")) {
        String nft = command.substring(12);
        storeNFT("REQ", nft);
        Serial.println("NFT_REQ stored");
      }
      else if (command == "GET_NFT_AUTH") {
        Serial.println(getNFT("AUTH"));
      }
      else if (command == "GET_NFT_REQ") {
        Serial.println(getNFT("REQ"));
      }
      else if (command == "REMOVE_NFT_AUTH") {
        removeNFT("AUTH");
        Serial.println("NFT_REMOVED");
      }
      else if (command == "REMOVE_NFT_REQ") {
        removeNFT("REQ");
        Serial.println("NFT_REMOVED");
      }else {
        Serial.println("INVALID");
      }
  }
}
