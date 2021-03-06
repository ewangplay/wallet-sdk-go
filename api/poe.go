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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"

	"strconv"

	"github.com/arxanchain/sdk-go-common/errors"
	"github.com/arxanchain/sdk-go-common/rest"
	restapi "github.com/arxanchain/sdk-go-common/rest/api"
	rtstructs "github.com/arxanchain/sdk-go-common/rest/structs"
	"github.com/arxanchain/sdk-go-common/structs/did"
	"github.com/arxanchain/sdk-go-common/structs/pki"
	"github.com/arxanchain/sdk-go-common/structs/wallet"
)

// CreatePOE is used to create POE digital asset.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// The default key pair trust mode does not trust, it will required key pair.
// If you had trust the key pair, it will required security code.
//
func (w *WalletClient) CreatePOE(header http.Header, body *wallet.POEBody, signParams *pki.SignatureParam) (result *wallet.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	if w.s != nil {
		signParams, err = w.queryPrivateKey(header, signParams)
		if err != nil {
			return
		}
	}

	// Build request signature
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}
	sign, err := buildSignatureBody(signParams, reqPayload)
	if err != nil {
		return nil, err
	}

	// Build http request
	r := w.c.NewRequest("POST", "/v1/poe/create")
	r.SetHeaders(header)

	// Build request body
	reqBody := &wallet.WalletRequest{
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

	payload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(payload), &result)

	return
}

// UpdatePOE is used to update POE digital asset.
//
// The default invoking mode is asynchronous, it will return
// without waiting for blockchain transaction confirmation.
//
// The default key pair trust mode does not trust, it will required key pair.
// If you had trust the key pair, it will required security code.
//
func (w *WalletClient) UpdatePOE(header http.Header, body *wallet.POEBody, signParams *pki.SignatureParam) (result *wallet.WalletResponse, err error) {
	if body == nil {
		err = fmt.Errorf("request payload invalid")
		return
	}

	if w.s != nil {
		signParams, err = w.queryPrivateKey(header, signParams)
		if err != nil {
			return
		}
	}

	// Build request signature
	reqPayload, err := json.Marshal(body)
	if err != nil {
		return
	}
	sign, err := buildSignatureBody(signParams, reqPayload)
	if err != nil {
		return nil, err
	}

	// Build http request
	r := w.c.NewRequest("PUT", "/v1/poe/update")
	r.SetHeaders(header)

	// Build request body
	reqBody := &wallet.WalletRequest{
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

	payload, ok := respBody.Payload.(string)
	if !ok {
		err = fmt.Errorf("response payload type invalid: %v", reflect.TypeOf(respBody.Payload))
		return
	}

	err = json.Unmarshal([]byte(payload), &result)

	return
}

// QueryPOE is used to query POE digital asset.
//
func (w *WalletClient) QueryPOE(header http.Header, id did.Identifier) (result *wallet.POEPayload, err error) {
	r := w.c.NewRequest("GET", "/v1/poe")
	r.SetHeaders(header)
	r.SetParam("id", string(id))

	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// parse http response
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

// UploadPOEFile is used to upload file for specified POE digital asset
//
// poeID parameter is the POE digital asset ID pre-created using CreatePOE API.
//
// poeFile parameter is the path to file to be uploaded.
//
func (w *WalletClient) UploadPOEFile(header http.Header, poeID string, poeFile string, readOnly bool) (result *wallet.UploadResponse, err error) {
	log.Println("Call UploadPOEFile...")

	if poeID == "" {
		err = fmt.Errorf("poe id must be set when uploading poe file")
		return
	}
	if poeFile == "" {
		err = fmt.Errorf("poe file must be set when uploading poe file")
		return
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// Create poeID form field
	err = writer.WriteField(wallet.OffchainPOEID, poeID)
	if err != nil {
		log.Printf("Write %s field to form fail: %v", wallet.OffchainPOEID, err)
		return
	}

	log.Printf("Write %s field to form succ", wallet.OffchainPOEID)

	// Create readOnly form field
	err = writer.WriteField(wallet.OffchainReadOnly, strconv.FormatBool(readOnly))
	if err != nil {
		log.Printf("Write %s field to form fail: %v", wallet.OffchainReadOnly, err)
		return
	}

	log.Printf("Write %s field to form succ", wallet.OffchainReadOnly)

	// Create poeFile form field
	formFile, err := writer.CreateFormFile(wallet.OffchainPOEFile, poeFile)
	if err != nil {
		log.Printf("Create form file handler for %s fail: %v", poeFile, err)
		return
	}

	log.Printf("Create form file handler for %s succ", poeFile)

	// Read data from file and Write to form
	srcFile, err := os.Open(poeFile)
	if err != nil {
		log.Printf("Open %s file fail: %v", poeFile, err)
		return
	}
	defer srcFile.Close()

	log.Printf("Open %s file succ", poeFile)

	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Printf("Write file contents to form fail: %v", err)
		return
	}

	log.Printf("Write file contents to form succ")

	// Send form
	contentType := writer.FormDataContentType()
	log.Printf("Content-Type: %s", contentType)
	// Must call Close() before http post to write EOF flag.
	writer.Close()

	// New request
	r := w.c.NewRequest("POST", "/v1/poe/upload")
	r.SetHeaders(header)
	r.SetHeader("Content-Type", contentType)
	r.SetBody(buf.Bytes())

	// Do upload
	_, resp, err := restapi.RequireOK(w.c.DoRequest(r))
	if err != nil {
		log.Printf("Request to upload file fail: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Request to upload file succ")

	// Parse http response
	var respBody rtstructs.Response
	if err = restapi.DecodeBody(resp, &respBody); err != nil {
		log.Printf("Parse the http response fail: %v", err)
		return
	}

	log.Printf("Parse the http response succ")

	if respBody.ErrCode != errors.SuccCode {
		err = rest.CodedError(respBody.ErrCode, respBody.ErrMessage)
		log.Printf("Upload file(%s) fail: %v", poeFile, err)
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
