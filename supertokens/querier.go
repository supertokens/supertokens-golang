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

package supertokens

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Querier struct {
	RIDToCore string
}

var (
	querierInitCalled     bool                  = false
	querierHosts          []NormalisedURLDomain = nil
	querierAPIKey         *string
	querierAPIVersion     string
	querierLastTriedIndex int
	querierLock           sync.Mutex
)

func (q *Querier) getQuerierAPIVersion() (string, error) {
	querierLock.Lock()
	defer querierLock.Unlock()
	if querierAPIVersion != "" {
		return querierAPIVersion, nil
	}
	response, err := q.sendRequestHelper(NormalisedURLPath{value: "/apiversion"}, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if querierAPIKey != nil {
			req.Header.Set("api-key", *querierAPIKey)
		}
		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))

	if err != nil {
		return "", err
	}

	respJSON, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	var cdiSupportedByServer struct {
		Versions []string `json:"versions"`
	}
	err = json.Unmarshal(respJSON, &cdiSupportedByServer)
	if err != nil {
		return "", err
	}
	supportedVersion := getLargestVersionFromIntersection(cdiSupportedByServer.Versions, cdiSupported)
	if supportedVersion == nil {
		return "", errors.New("the running SuperTokens core version is not compatible with this Golang SDK. Please visit https://supertokens.io/docs/community/compatibility-table to find the right version")
	}

	querierAPIVersion = *supportedVersion

	return querierAPIVersion, nil
}

func GetNewQuerierInstanceOrThrowError(rIDToCore string) (*Querier, error) {
	if !querierInitCalled {
		return nil, errors.New("please call the supertokens.init function before using SuperTokens")
	}
	return &Querier{RIDToCore: rIDToCore}, nil
}

func initQuerier(hosts []NormalisedURLDomain, APIKey string) {
	if !querierInitCalled {
		querierInitCalled = true
		querierHosts = hosts
		if APIKey != "" {
			querierAPIKey = &APIKey
		}
		querierAPIVersion = ""
		querierLastTriedIndex = 0
	}
}

func (q *Querier) SendPostRequest(path string, data map[string]interface{}) (map[string]interface{}, error) {
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	return q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		if data == nil {
			data = map[string]interface{}{}
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, querierAPIVersionError := q.getQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("cdi-version", apiVerion)
		if querierAPIKey != nil {
			req.Header.Set("api-key", *querierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

func (q *Querier) SendDeleteRequest(path string, data map[string]interface{}) (map[string]interface{}, error) {
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	return q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, querierAPIVersionError := q.getQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("cdi-version", apiVerion)
		if querierAPIKey != nil {
			req.Header.Set("api-key", *querierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

func (q *Querier) SendGetRequest(path string, params map[string]string) (map[string]interface{}, error) {
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	return q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		query := req.URL.Query()

		for k, v := range params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()

		apiVerion, querierAPIVersionError := q.getQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}
		req.Header.Set("cdi-version", apiVerion)
		if querierAPIKey != nil {
			req.Header.Set("api-key", *querierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

func (q *Querier) SendPutRequest(path string, data map[string]interface{}) (map[string]interface{}, error) {
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	return q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, querierAPIVersionError := q.getQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("cdi-version", apiVerion)
		if querierAPIKey != nil {
			req.Header.Set("api-key", *querierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

type httpRequestFunction func(url string) (*http.Response, error)

func (q *Querier) sendRequestHelper(path NormalisedURLPath, httpRequest httpRequestFunction, numberOfTries int) (map[string]interface{}, error) {
	if numberOfTries == 0 {
		return nil, errors.New("no SuperTokens core available to query")
	}

	querierLock.Lock()
	currentHost := querierHosts[querierLastTriedIndex].GetAsStringDangerous()
	querierLastTriedIndex = (querierLastTriedIndex + 1) % len(querierHosts)
	querierLock.Unlock()

	resp, err := httpRequest(currentHost + path.GetAsStringDangerous())

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return q.sendRequestHelper(path, httpRequest, numberOfTries-1)
		}
		if resp != nil {
			resp.Body.Close()
		}
		return nil, err
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("SuperTokens core threw an error for a request to path: '%s' with status code: %v and message: %s", path.GetAsStringDangerous(), resp.StatusCode, body))
	}

	finalResult := make(map[string]interface{})
	jsonError := json.Unmarshal(body, &finalResult)
	if jsonError != nil {
		return map[string]interface{}{
			"result": string(body),
		}, nil
	}
	return finalResult, nil
}
