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

package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
)

func doGetRequest(req *http.Request) (interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Provider API returned response with status `%s` and body `%s`", resp.Status, string(body)))
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var jwksKeys = map[string]*keyfunc.JWKS{}
var jwksKeysLock = sync.Mutex{}

func getJWKSFromURL(url string) (*keyfunc.JWKS, error) {
	if jwks, ok := jwksKeys[url]; ok {
		return jwks, nil
	}

	jwksKeysLock.Lock()
	defer jwksKeysLock.Unlock()

	// Check again to see if it was added while we were waiting for the lock
	if jwks, ok := jwksKeys[url]; ok {
		return jwks, nil
	}

	options := keyfunc.Options{
		RefreshInterval: time.Hour,
	}
	jwks, err := keyfunc.Get(url, options)
	if err != nil {
		return nil, err
	}
	jwksKeys[url] = jwks
	return jwks, nil
}
