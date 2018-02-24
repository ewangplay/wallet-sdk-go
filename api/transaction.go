/*
Copyright ArxanFintech Technology Ltd. 2018 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/arxanchain/sdk-go-common/errors"
	"github.com/arxanchain/sdk-go-common/rest"
	restapi "github.com/arxanchain/sdk-go-common/rest/api"
	rtstructs "github.com/arxanchain/sdk-go-common/rest/structs"
	"github.com/arxanchain/sdk-go-common/structs"
)

// IssueCToken is used to issue colored token.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by caller using
// 'github.com/arxanchain/sdk-go-common/crypto/tools/sign-util' tool.
//
func (w *WalletClient) IssueCToken(header http.Header, body *structs.IssueBody, sign *structs.SignatureBody) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build http request
	r := w.c.NewRequest("POST", "/v1/transaction/tokens/issue")
	r.SetHeaders(header)

	// Build request payload
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}

	// build request body
	reqBody := &structs.WalletRequest{
		Payload:   string(reqPayload),
		Signature: sign,
	}
	r.SetBody(reqBody)

	// Do http request
	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse http response
	var respBody rtstructs.Response
	if err = restapi.DecodeBody(resp, &respBody); err != nil {
		return
	}

	if respBody.ErrCode != errors.SuccCode {
		err = rest.CodedError(respBody.ErrCode, respBody.ErrMessage)
		return
	}

	respPayload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(respPayload), &result)

	return
}

// IssueCTokenSign is used to issue colored token.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by SDK, need to pass the user private key to the SDK.
//
func (w *WalletClient) IssueCTokenSign(header http.Header, body *structs.IssueBody, signParams *structs.SignatureParam) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build request signature
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}
	sign, err := buildSignatureBody(signParams, reqPayload)
	if err != nil {
		return
	}

	return w.IssueCToken(header, body, sign)
}

// IssueAsset is used to issue digital asset.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by caller using
// 'github.com/arxanchain/sdk-go-common/crypto/tools/sign-util' tool.
//
func (w *WalletClient) IssueAsset(header http.Header, body *structs.IssueAssetBody, sign *structs.SignatureBody) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build http request
	r := w.c.NewRequest("POST", "/v1/transaction/assets/issue")
	r.SetHeaders(header)

	// Build request payload
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Build request body
	reqBody := &structs.WalletRequest{
		Payload:   string(reqPayload),
		Signature: sign,
	}
	r.SetBody(reqBody)

	// Do http request
	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse http response
	var respBody rtstructs.Response
	if err = restapi.DecodeBody(resp, &respBody); err != nil {
		return
	}

	if respBody.ErrCode != errors.SuccCode {
		err = rest.CodedError(respBody.ErrCode, respBody.ErrMessage)
		return
	}

	respPayload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(respPayload), &result)

	return
}

// IssueAssetSign is used to issue digital asset.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by SDK, need to pass the user private key to the SDK.
//
func (w *WalletClient) IssueAssetSign(header http.Header, body *structs.IssueAssetBody, signParams *structs.SignatureParam) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build request signature
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}
	sign, err := buildSignatureBody(signParams, reqPayload)
	if err != nil {
		return
	}

	return w.IssueAsset(header, body, sign)
}

// TransferCToken is used to transfer colored tokens from one user to another.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by caller using
// 'github.com/arxanchain/sdk-go-common/crypto/tools/sign-util' tool.
//
func (w *WalletClient) TransferCToken(header http.Header, body *structs.TransferBody, sign *structs.SignatureBody) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build http request
	r := w.c.NewRequest("POST", "/v1/transaction/tokens/transfer")
	r.SetHeaders(header)

	// Build request payload
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Build request body
	reqBody := &structs.WalletRequest{
		Payload:   string(reqPayload),
		Signature: sign,
	}
	r.SetBody(reqBody)

	// Do http request
	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse http response
	var respBody rtstructs.Response
	if err = restapi.DecodeBody(resp, &respBody); err != nil {
		return
	}

	if respBody.ErrCode != errors.SuccCode {
		err = rest.CodedError(respBody.ErrCode, respBody.ErrMessage)
		return
	}

	respPayload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(respPayload), &result)

	return
}

// TransferCTokenSign is used to transfer colored tokens from one user to another.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by SDK, need to pass the user private key to the SDK.
//
func (w *WalletClient) TransferCTokenSign(header http.Header, body *structs.TransferBody, signParams *structs.SignatureParam) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build request signature
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}
	sign, err := buildSignatureBody(signParams, reqPayload)
	if err != nil {
		return
	}

	return w.TransferCToken(header, body, sign)
}

// TransferAsset is used to transfer assets from one user to another.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by caller using
// 'github.com/arxanchain/sdk-go-common/crypto/tools/sign-util' tool.
//
func (w *WalletClient) TransferAsset(header http.Header, body *structs.TransferAssetBody, sign *structs.SignatureBody) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build http request
	r := w.c.NewRequest("POST", "/v1/transaction/assets/transfer")
	r.SetHeaders(header)

	// Build request payload
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}

	// Build request body
	reqBody := &structs.WalletRequest{
		Payload:   string(reqPayload),
		Signature: sign,
	}
	r.SetBody(reqBody)

	// Do http request
	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse http response
	var respBody rtstructs.Response
	if err = restapi.DecodeBody(resp, &respBody); err != nil {
		return
	}

	if respBody.ErrCode != errors.SuccCode {
		err = rest.CodedError(respBody.ErrCode, respBody.ErrMessage)
		return
	}

	respPayload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(respPayload), &result)

	return
}

// TransferAssetSign is used to transfer assets from one user to another.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// If you want to switch to synchronous invoking mode, set
// 'BC-Invoke-Mode' header to 'sync' value. In synchronous mode,
// it will not return until the blockchain transaction is confirmed.
//
// The signature is generated by SDK, need to pass the user private key to the SDK.
//
func (w *WalletClient) TransferAssetSign(header http.Header, body *structs.TransferAssetBody, signParams *structs.SignatureParam) (result *structs.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	// Build request signature
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}
	sign, err := buildSignatureBody(signParams, reqPayload)
	if err != nil {
		return
	}

	return w.TransferAsset(header, body, sign)
}

// QueryTransactionLogs is used to query transaction logs.
//
// txType:
// in: query income type transaction
// out: query spending type transaction
//
func (w *WalletClient) QueryTransactionLogs(header http.Header, id structs.Identifier, txType string) (result structs.TransactionLogs, err error) {
	if id == "" {
		err = fmt.Errorf("request id invalid")
		return
	}

	// Build http request
	r := w.c.NewRequest("GET", "/v1/transaction/logs")
	r.SetHeaders(header)
	r.SetParam("id", string(id))
	r.SetParam("type", txType)

	// Do http request
	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Parse http response
	var respBody rtstructs.Response
	if err = restapi.DecodeBody(resp, &respBody); err != nil {
		return
	}

	if respBody.ErrCode != errors.SuccCode {
		err = rest.CodedError(respBody.ErrCode, respBody.ErrMessage)
		return
	}

	respPayload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(respPayload), &result)

	return
}
