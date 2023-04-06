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
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ParsedJWTInfo struct {
	RawTokenString string
	RawPayload     string
	Header         string
	Payload        map[string]interface{}
	Signature      string
	Version        int
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

	// V1&V2 is functionally identical, plus all legacy tokens should be V2 now.
	version := 2
	// V2 or older tokens did not save the key id;
	err := checkHeader(splittedInput[0])
	if err != nil {
		parsedHeaderBytes, err := b64.StdEncoding.DecodeString(splittedInput[0])
		if err != nil {
			return ParsedJWTInfo{}, err
		}

		parsedHeader := map[string]interface{}{}
		err = json.Unmarshal(parsedHeaderBytes, &parsedHeader)
		if err != nil {
			return ParsedJWTInfo{}, err
		}

		versionInHeader := parsedHeader["version"]

		if reflect.TypeOf(version).Kind() != reflect.String {
			return ParsedJWTInfo{}, errors.New("JWT header mismatch")
		}

		versionNumber, parseError := strconv.Atoi(versionInHeader.(string))

		if parsedHeader["typ"].(string) != "JWT" || parseError != nil || versionNumber < 3 || parsedHeader["kid"] == nil {
			return ParsedJWTInfo{}, errors.New("JWT header mismatch")
		}

		version = versionNumber
	}

	payloadBytes, err := b64.StdEncoding.DecodeString(splittedInput[1])
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
		Version:        version,
	}, nil
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
