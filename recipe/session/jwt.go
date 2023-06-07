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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"reflect"
	"strconv"
	"strings"
)

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

func ParseJWTWithoutSignatureVerification(token string) (sessmodels.ParsedJWTInfo, error) {
	splittedInput := strings.Split(token, ".")
	latestAccessTokenVersion := 3
	var kid *string
	if len(splittedInput) != 3 {
		errors.New("Invalid JWT")
	}

	// V1&V2 is functionally identical, plus all legacy tokens should be V2 now.
	version := 2
	// V2 or older tokens did not save the key id;
	err := checkHeader(splittedInput[0])
	payload := map[string]interface{}{}

	// If err != nil, it is a V3 token (or above)
	if err != nil {
		unverifiedToken, _, rawParseError := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if rawParseError != nil {
			return sessmodels.ParsedJWTInfo{}, rawParseError
		}

		parsedHeader := unverifiedToken.Header

		versionInHeader, ok := parsedHeader["version"]

		if !ok {
			versionInHeader = fmt.Sprint(latestAccessTokenVersion)
		}

		if reflect.TypeOf(versionInHeader).Kind() != reflect.String {
			return sessmodels.ParsedJWTInfo{}, errors.New("JWT header mismatch")
		}

		versionNumber, parseError := strconv.Atoi(versionInHeader.(string))

		kidInHeader, ok := parsedHeader["kid"]

		if !ok {
			return sessmodels.ParsedJWTInfo{}, errors.New("JWT header mismatch")
		}

		if reflect.TypeOf(kidInHeader).Kind() != reflect.String {
			return sessmodels.ParsedJWTInfo{}, errors.New("JWT header mismatch")
		}

		kidString := kidInHeader.(string)
		kid = &kidString

		if parsedHeader["typ"].(string) != "JWT" || parseError != nil || versionNumber < latestAccessTokenVersion || parsedHeader["kid"] == nil {
			return sessmodels.ParsedJWTInfo{}, errors.New("JWT header mismatch")
		}

		version = versionNumber

		claims, ok := unverifiedToken.Claims.(jwt.MapClaims)

		if ok {
			payload = claims
		} else {
			return sessmodels.ParsedJWTInfo{}, errors.New("Invalid JWT")
		}
	} else {
		bytes, err := base64.StdEncoding.DecodeString(splittedInput[1])

		if err != nil {
			return sessmodels.ParsedJWTInfo{}, err
		}

		decodedJson := map[string]interface{}{}
		err = json.Unmarshal(bytes, &decodedJson)

		if err != nil {
			return sessmodels.ParsedJWTInfo{}, err
		}

		payload = decodedJson

	}

	return sessmodels.ParsedJWTInfo{
		RawTokenString: token,
		RawPayload:     splittedInput[1],
		Header:         splittedInput[0],
		Payload:        payload,
		Signature:      splittedInput[2],
		Version:        version,
		KID:            kid,
	}, nil
}
