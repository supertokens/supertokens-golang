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

type ParsedJWTInfo struct {
	RawTokenString string
	RawPayload     string
	Header         string
	Payload        map[string]interface{}
	Signature      string
}

var HEADERS = []string{
	"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIxIn0=", // {"alg":"RS256","typ":"JWT","version":"1"}
	"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIyIn0=", // {"alg":"RS256","typ":"JWT","version":"2"}
}

func checkHeader(header string) error {
	for _, h := range HEADERS {
		if h == header {
			return nil
		}
	}
	return errors.New("Invalid JWT header")
}

func parseJWTWithoutSignatureVerification(jwt string) (ParsedJWTInfo, error) {
	splittedInput := strings.Split(jwt, ".")
	if len(splittedInput) != 3 {
		errors.New("Invalid JWT")
	}

	err := checkHeader(splittedInput[0])
	if err != nil {
		return ParsedJWTInfo{}, err
	}

	payloadBytes, err := b64.RawStdEncoding.DecodeString(splittedInput[1])
	if err != nil {
		return ParsedJWTInfo{}, err
	}
	payload := map[string]interface{}{}
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return ParsedJWTInfo{}, err
	}

	return ParsedJWTInfo{
		RawTokenString: jwt,
		RawPayload:     splittedInput[1],
		Header:         splittedInput[0],
		Payload:        payload,
		Signature:      splittedInput[2],
	}, nil
}

func verifyJWT(jwtInfo ParsedJWTInfo, jwtSigningPublicKey string) error {
	var publicKey, publicKeyError = getPublicKeyFromStr("-----BEGIN PUBLIC KEY-----\n" + jwtSigningPublicKey + "\n-----END PUBLIC KEY-----")
	if publicKeyError != nil {
		return publicKeyError
	}

	h := sha256.New()
	h.Write([]byte(jwtInfo.Header + "." + jwtInfo.RawPayload))
	digest := h.Sum(nil)

	var decodedSignature, decodedSignatureError = b64.StdEncoding.DecodeString(jwtInfo.Signature)
	if decodedSignatureError != nil {
		return decodedSignatureError
	}

	verificationError := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, decodedSignature)
	if verificationError != nil {
		return verificationError
	}

	return nil
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
