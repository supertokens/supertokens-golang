/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package session

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"strings"
)

/*
{
	"alg":     "RS256",
	"typ":     "JWT",
	"version": "2",
}
*/
const header = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIyIn0="

func verifyJWTAndGetPayload(jwt string, jwtSigningPublicKey string) (map[string]interface{}, error) {
	var splitted = strings.Split(jwt, ".")
	if len(splitted) != 3 {
		return nil, errors.New("Invalid JWT")
	}
	if header != splitted[0] {
		return nil, errors.New("JWT header mismatch")
	}
	var payload = splitted[1]

	var publicKey, publicKeyError = getPublicKeyFromStr("-----BEGIN PUBLIC KEY-----\n" + jwtSigningPublicKey + "\n-----END PUBLIC KEY-----")
	if publicKeyError != nil {
		return nil, publicKeyError
	}

	h := sha256.New()
	h.Write([]byte(header + "." + payload))
	digest := h.Sum(nil)

	var decodedSignature, decodedSignatureError = b64.StdEncoding.DecodeString(splitted[2])
	if decodedSignatureError != nil {
		return nil, decodedSignatureError
	}

	verificationError := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, decodedSignature)
	if verificationError != nil {
		return nil, verificationError
	}

	var decodedPayload, base64Error = b64.StdEncoding.DecodeString(payload)
	if base64Error != nil {
		return nil, base64Error
	}

	var result map[string]interface{}
	jsonError := json.Unmarshal(decodedPayload, &result)
	if jsonError != nil {
		return nil, jsonError
	}
	return result, nil
}

func getPayloadWithoutVerifying(jwt string) (map[string]interface{}, error) {
	var splitted = strings.Split(jwt, ".")
	if len(splitted) != 3 {
		return nil, errors.New("Invalid JWT")
	}

	var payload = splitted[1]

	var decodedPayload, base64Error = b64.StdEncoding.DecodeString(payload)
	if base64Error != nil {
		return nil, base64Error
	}

	var result map[string]interface{}
	jsonError := json.Unmarshal(decodedPayload, &result)
	if jsonError != nil {
		return nil, jsonError
	}
	return result, nil
}

func getPublicKeyFromStr(str string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse DER encoded public key:" + err.Error())
	}

	return pub.(*rsa.PublicKey), nil
}
