&#x20;

# Nexasecure V2

> Enterprise-grade NFT-based authentication framework for microservices

**Description:** This system combines traditional authentication mechanisms with Web3 technologies to create a secure and modern authentication framework. It utilizes a dual-blockchain architecture to establish a physical and logical barrier between the user and the core system, enhancing both security and decentralization.

---

## Table of Contents

* [🚀 Features](#-features)
* [🏛 Architecture](#-architecture)
* [🔑 Key Components](#-key-components)
* [⚙️ How It Works](#️-how-it-works)
* [🔧 Prerequisites](#-prerequisites)
* [⚙️ Installation & Setup](#️-installation--setup)

  * [1. Clone Repository](#1-clone-repository)
  * [2. Configure Environment](#2-configure-environment)
  * [3. Start Services](#3-start-services)
  * [4. Deploy Smart Contracts](#4-deploy-smart-contracts)
  * [5. Flash Hardware Wallet](#5-flash-hardware-wallet)
  * [6. Run Client Driver](#6-run-client-driver)
* [💡 Usage](#-usage)
* [🛠 Production Considerations](#-production-considerations)
* [🔗 Integration](#-integration)
* [🔒 Security Best Practices](#-security-best-practices)
* [📄 License](#-license)

---

## 🚀 Features

* **Physical Barrier:** Requires an ESP32-S3 hardware wallet for login, preventing software-only breaches.
* **Ephemeral Tokens:** Uses NFTs instead of JWTs; tokens self-destruct after each session.
* **Dual Blockchains:** Isolated ledgers for authentication and request token lifecycles.
* **Seamless Integration:** Exposes RESTful APIs and socket endpoints for easy incorporation.
* **Enterprise-Ready:** Designed for scalability, high security, and compliance.

## 🏛 Architecture

![image](https://github.com/user-attachments/assets/c9c03e1a-3c56-456f-9376-a3da457fbaaf)


## 🔑 Key Components

| Component           | Description                                                    |
| ------------------- | -------------------------------------------------------------- |
| **Client Driver**   | Python-based desktop app for wallet connectivity & requests.   |
| **ESP32-S3 Wallet** | Secure device for NFT storage & signing with biometric option. |
| **Auth Server**     | Validates credentials + wallet, mints/burns AuthNFTs.          |
| **Socket Server**   | Encrypted socket channel for real-time wallet interactions.    |
| **Backend API**     | Middleware validates ReqNFT for each API request.              |
| **Private Chains**  | Two Hyperledger ledgers for AuthNFT & ReqNFT lifecycles.       |
| **Database**        | PostgreSQL (dev) / AWS QLDB (prod) for immutable logs.         |

## ⚙️ How It Works

1. **Login**: Connect ESP32-S3 & authenticate credentials.
2. **Wallet Check**: Server ensures hardware wallet is present.
3. **Credential Validation**: Verifies user identity.
4. **AuthNFT Transfer**: Ownership shifts to company on Auth Blockchain.
5. **Burn AuthNFT**: Prevents token reuse by burning.
6. **ReqNFT Issue**: New session-specific NFT minted on Req Blockchain.
7. **Session Requests**: Client auto-attaches ReqNFT; API middleware verifies.
8. **Logout**: New AuthNFT minted, transferred, & burned; logs stored.

## 🔧 Prerequisites

* **Hardware**: ESP32-S3 module
* **Software**:

  * Docker & Docker Compose
  * Go ≥ v1.20
  * Python ≥ 3.9
  * PostgreSQL / AWS QLDB
  * AWS CLI with IAM permissions

## ⚙️ Installation & Setup

### 1. Clone Repository

```bash
git clone https://github.com/Umayanga12/Nexasecure-V2.git
cd nexasecure-v2
```

### 2. Configure Environment

1. Copy `.env.example` to `.env`.
2. Set blockchain endpoints, DB credentials, OTP service keys.

### 3. Start Services

```bash
docker-compose up -d
```

### 5. Flash Hardware Wallet

```bash
arduino-cli compile --fqbn esp32:esp32:esp32s3 wallet
arduino-cli upload --fqbn esp32:esp32:esp32s3 wallet
```
After Uploading the code into ESP32 It need to setup(Need to set password and AuthNFT) using serial monitor.

### 6. Run Client Driver

```bash
cd client-driver
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py
```

## 💡 Usage

* **Login:** `POST /login` with wallet connected.
* **Request:** based on the connection address secket server can rectrive the ReqNFT and Validate it.
* **Logout:** `POST /logout` triggers NFT teardown.

## 🛠 Production Considerations

* Swap PostgreSQL for **AWS QLDB**
* Use **AWS Managed Blockchain** with Hyperledger Fabric
* Deploy biometric hardware wallets in secure enclaves
* Enforce cloud IAM, VPN, and proxy configurations

## 🔗 Integration

Use behind API Gateway (e.g., AWS API Gateway) as a standalone auth microservice. Leverage health-check endpoints and rate-limiting.

## 🔒 Security Best Practices

* Audit contracts (MythX, Slither)
* TLS for all communications
* Enforce MFA (OTP)
* Rotate IAM roles & credentials regularly
* Implement monitoring & alerting (Prometheus + Grafana)

## 📄 License

Licensed under [MIT](LICENSE). For details, see the [LICENSE](LICENSE) file.
