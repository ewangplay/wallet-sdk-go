# Status
[![Build Status](https://travis-ci.org/arxanchain/wallet-sdk-go.svg?branch=master)](https://travis-ci.org/arxanchain/wallet-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/arxanchain/wallet-sdk-go)](https://goreportcard.com/report/github.com/arxanchain/wallet-sdk-go)
[![GoDoc](https://godoc.org/github.com/arxanchain/wallet-sdk-go?status.svg)](https://godoc.org/github.com/arxanchain/wallet-sdk-go)

# wallet-go-sdk

Blockchain Wallet SDK includes APIs for managing wallet accounts (DID),
digital assets (POE), colored tokens etc.

You need not care about how the backend blockchain runs or the unintelligible
techniques, such as consensus, endorsement and decentralization. Simply use
the SDK we provide to implement your business logics, we will handle caching,
tagging, compressing, encrypting and high availability.

Please refer to the API document: [Blockchain wallet platform](http://www.arxanfintech.com/infocenter/html/development/wallet.html)

# Usage

## Install

Run the following command to download the Go SDK:

```code
go get github.com/arxanchain/wallet-sdk-go/api
```

## New wallet client

To invoke the SDK API, you first need to create a wallet client as follows:

```code
import (
	pw "github.com/arxanchain/sdk-go-common/protos/wallet"
	restapi "github.com/arxanchain/sdk-go-common/rest/api"
	"github.com/arxanchain/sdk-go-common/structs/wallet"
	"github.com/arxanchain/sdk-go-common/structs/pki"
	walletapi "github.com/arxanchain/wallet-sdk-go/api"
)

// Create wallet client
config := &restapi.Config{
	Address:    "https://API-Proxy-Gateway:Port",
	ApiKey:     "Your-API-Access-Key",
	TLSConfig: &restapi.TLSConfig:{
		CAFile: "path/to/tls/ca/cert",
		CertFile: "path/to/tls/user/cert",
		KeyFile: "path/to/tls/user/key",
	},
	CallbackUrl: "http://callback-url",
	EnterpriseSignParam: &restapi.EnterpriseSignParam{
		Creator: "did:axn:09e2fc68-f51e-4aff-b6e4-427cce3ed1af",
		Nonce: "nonce",
		PrivateKey: "RiQ+oEuaelf2aecUZvG7xrWr+p43ZfjGZYfDCXfQD+ku0xY5BXP8kIKhiqzKRvfyKBKM3y7V9O1bF7X3M9mxkQ==",
	},
}
walletClient, err := walletapi.NewWalletClient(config)
if err != nil {
	fmt.Printf("New wallet client fail: %v\n", err)
	return
}
fmt.Printf("New wallet client succ\n")
```

* When building the client configuration, the **Address**, **ApiKey** and **TLSConfig** fields must
be set. The **Address** is set to the address of BaaS API proxy gateway, the **ApiKey** is set to 
the API access key obtained on `ChainConsole` management page, and the **TLSConfig** is set to the 
real TLS config.

* `Callback-Url` is optional. You only need to set it if you need to receive blockchain transaction events.

* `Enterprisesignparam`: Enterprise signature parameter, used to sign UTXO records for AXT fee.
	- Creator: Enterprise wallet did
	- Nonce: Signature random nonce string
	- PrivateKey: The ed25519 private key of enterprise wallet

About how to apply API-Key, please refer to [Apikey Application](http://www.arxanfintech.com/infocenter/html/baas/enterprise/v1.2/api-access.html#api-access-ref)

## Register wallet account

After creating wallet client, you can use this client to register wallet account
as follows:

```code
// Build request header
header := http.Header{}
// If you use synchronous invoking mode, set following header
header.Set("Bc-Invoke-Mode", "sync")

// Register wallet account
registerBody := &wallet.RegisterWalletBody{
	Type:   pw.DidType_ORGANIZATION,
	Access: "alice0001",
	Secret: "Alice#123456",
}
resp, err = walletClient.Register(header, registerBody)
if err != nil {
	fmt.Printf("Register wallet fail: %v\n", err)
	return
}
walletID := resp.Id
keyPair := resp.KeyPair
fmt.Printf("Register wallet succ.\nwallet id: %v\nED25519 public key: %v\nED25519 private key: %v", walletID, keyPair.PublicKey, keyPair.PrivateKey)
```

## Create POE digital asset and upload file

After creating the wallet account, you can create POE assets for this account as follows:

```code
// Create poe asset
poeBody := &wallet.POEBody{
	Name:     "TestPOE",
	Owner:    walletID,
	Metadata: []byte("poe metadata"),
}
signParam := &pki.SignatureParam{
	Creator:    walletID,
	Nonce:      "nonce",
	PrivateKey: keyPair.PrivateKey,
}
resp, err = walletClient.CreatePOE(header, poeBody, signParam)
if err != nil {
	fmt.Printf("CreatePOE fail: %v\n", err)
	return
}
fmt.Printf("Create POE succ. Response: %+v\n", resp)

// Upload poe file
poeID := string(resp.Id)
poeFile := "./test-upload-file"
resp, err = walletClient.UploadPOEFile(header, poeID, poeFile, false)
if err != nil {
	fmt.Printf("UploadPOEFail fail: %v\n", err)
	return
}
fmt.Printf("Upload POE file succ. Response: %+v\n", resp)
```

* When creating POE assets, the **Name** and **Owner** fields must be set, and the
**Owner** field must be set to the wallet account ID.

* When building the signature parameter, use the ed25519 private key returned
when registering wallet to do ed25519 signing.

* `UploadPOEFile` API uploads the file to **Offchain** storage, generates SHA256
hash value for this file, and saves this hash value into blockchain.

## Issue colored token using digital asset

Once you have possessed assets, you can use a specific asset to issue colored
token as follows:

```code
// Issue colored token
issueBody = &wallet.IssueBody{
	Issuer:  string(issuerID),
	Owner:   string(walletID),
	AssetId: string(poeID),
	Amount:  1000,
}
signParam = &pki.SignatureParam{
	Creator:    issuerID,
	Nonce:      "nonce",
	PrivateKey: issuerKeyPair.PrivateKey,
}
resp, err = walletClient.IssueCToken(header, issueBody, signParam)
if err != nil {
	log.Fatalf("Issue colored token fail: %v\n", err)
	return
}
log.Printf("Issue colored token succ. Response: %+v", resp)
```

* When issuing colored token, you need to specify an issuer (one wallet account ID),
an asset to issue token, and the asset owner (another wallet account ID).

## Transfer colored token

After issuing colored token, the asset owner's wallet account will own these
colored tokens, and can transfer some of them to other wallet accounts.

```code
// Transfer colored token
transferBody = &wallet.TransferCTokenBody{
	From: string(walletID),
	To:   string(toID),
	Tokens: []*wallet.TokenAmount{
		&wallet.TokenAmount{
			TokenId: tokenId,
			Amount:  100,
		},
	},
}
signParam = &pki.SignatureParam{
	Creator:    walletID,
	Nonce:      "nonce",
	PrivateKey: keyPair.PrivateKey,
}
resp, err = walletClient.TransferCToken(header, transferBody, signParam)
if err != nil {
	log.Fatalf("Transfer colored token fail: %v\n", err)
	return
}
log.Printf("Transfer colored token succ.\nResponse: %+v", resp)
```

## Query colored token balance

You can use the `GetWalletBalance` API to get the balance of the specified wallet
account as follows:

```code
// Query wallet balance
balance, err = walletClient.GetWalletBalance(header, walletID)
if err != nil {
	fmt.Printf("Get wallet(%s) balance fail: %v\n", walletID, err)
	return
}
if balance.ColoredTokens != nil {
	fmt.Printf("Get wallet(%s) colored tokens succ\n", walletID)
	for ctokenId, ctoken := range balance.ColoredTokens {
		fmt.Printf("===> CTokenID: %v, Amount: %v\n", ctokenId, ctoken.Amount)
	}
}
if balance.DigitalAssets != nil {
	fmt.Printf("Get wallet(%s) digital assets succ\n", walletID)
	for assetId, asset := range balance.DigitalAssets {
		fmt.Printf("===> AssetID: %v, Amount: %v\n", assetId, asset.Amount)
	}
}
```

## Query transaction logs
You can use the `QueryTransactionLogs` API to get the transaction logs of the
specified wallet account as follows:

```
// Query wallet tx logs
txType := "in" // tx type: transfer in/out, other: all
var num int32 = 1
var page int32 = 1
logs, err = walletClient.QueryTransactionLogs(header, walletID, txType, num, page)
if err != nil {
	fmt.Printf("Get wallet(%s) tx logs fail: %v\n", walletID, err)
	return
}
if logs != nil {
	fmt.Printf("Get wallet(%s) tx logs succ: %+v\n", walletID, logs)
	}
}
```

## Query transaction UTXO logs
You can use the `QueryTransactionUTXO` API to get the transaction UTXOs of the
specified wallet account as follows:

```
// Query wallet UTXO
var num int32 = 1
var page int32 = 1
logs, err = walletClient.QueryTransactionUTXO(header, walletID, num, page)
if err != nil {
	fmt.Printf("Get wallet(%s) UTXOs fail: %v\n", walletID, err)
	return
}
if logs != nil {
	fmt.Printf("Get wallet(%s) UTXOs succ: %+v\n", walletID, logs)
	}
}
```

## Query transaction STXO logs
You can use the `QueryTransactionSTXO` API to get the transaction STXOs of the
specified wallet account as follows:

```
// Query wallet STXO
var num int32 = 1
var page int32 = 1
logs, err = walletClient.QueryTransactionSTXO(header, walletID, num, page)
if err != nil {
	fmt.Printf("Get wallet(%s) STXOs fail: %v\n", walletID, err)
	return
}
if logs != nil {
	fmt.Printf("Get wallet(%s) STXOs succ: %+v\n", walletID, logs)
	}
}
```

## Using callback URL to receive blockchain transaction events

Each of the APIs for invoking blockchain has two invoking modes, one is `sync`
mode, the other is `async` mode.

The default invoking mode is asynchronous, it will return without waiting for
blockchain transaction confirmation. In asynchronous mode, you should set
`Callback-Url` in the http header to receive blockchain transaction events.

The blockchain transaction event structure is defined as follows:

```code
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp

// Blockchain transaction event payload
type BcTxEventPayload struct {
	BlockNumber   uint64                     `json:"block_number"`   // Block number
	BlockHash     []byte                     `json:"block_hash"`     // Block hash
	ChannelId     string                     `json:"channel_id"`     // Channel ID
	ChaincodeId   string                     `json:"chaincode_id"`   // Chaincode ID
	TransactionId string                     `json:"transaction_id"` // Transaction ID
	Timestamp     *google_protobuf.Timestamp `json:"timestamp"`      // Transaction timestamp
	IsInvalid     bool                       `json:"is_invalid"`     // Is transaction invalid
	Payload       interface{}                `json:"payload"`        // Transaction Payload
}
```

One blockchain transaction event sample as follows:

```code
{
	"block_number":63,
	"block_hash":"vTRmfHZ3aaecbbw2A5zPcuzekUC42Lid3w+i6dOU5C0=",
	"channel_id":"pubchain",
	"chaincode_id":"pubchain-c4:",
	"transaction_id":"243eaa6e695cc4ce736e765395a64b8b917ff13a6c6500a11558b5e94e02556a",
	"timestamp":{
		"seconds":1521189855,
		"nanos":192203115
	},
	"is_invalid":false,
	"payload":{
		"id":"4debe20b-ca00-49b0-9130-026a1aefcf2d",
		"metadata":{
			"member_id_value":"3714811988020512",
			"member_mobile":"6666",
			"member_name":"8777896121269017",
			"member_truename":"Tony"
		}
	}
}
```

**NOTE** Please make sure that you response with http status code 200 when you have received one
blockchain transaction event. If you don't respond, You might get the same event multiple times, 
because the sender cannot confirm that you have received the event, so it will resend.

If you don't care the blockchain transaction event, you can switch to synchronous invoking mode, 
set `Bc-Invoke-Mode` header to `sync` value. In synchronous mode, it will not return until the blockchain
transaction is confirmed.
