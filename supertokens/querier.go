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
	"sort"
	"strings"
	"sync"
	"time"
)

type Querier struct {
	RIDToCore string
}

type QuerierHost struct {
	Domain   NormalisedURLDomain
	BasePath NormalisedURLPath
}

var (
	querierInitCalled     bool          = false
	QuerierHosts          []QuerierHost = nil
	QuerierAPIKey         *string
	querierAPIVersion     string
	querierLastTriedIndex int
	querierLock           sync.Mutex
	querierHostLock       sync.Mutex
	querierInterceptor    func(*http.Request, UserContext) *http.Request
	querierGlobalCacheTag uint64
	querierDisableCache   bool
)

func SetQuerierApiVersionForTests(version string) {
	querierAPIVersion = version
}

func (q *Querier) GetQuerierAPIVersion() (string, error) {
	querierLock.Lock()
	defer querierLock.Unlock()
	if querierAPIVersion != "" {
		return querierAPIVersion, nil
	}
	response, _, err := q.sendRequestHelper(NormalisedURLPath{value: "/apiversion"}, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if QuerierAPIKey != nil {
			req.Header.Set("api-key", *QuerierAPIKey)
		}
		client := &http.Client{}
		return client.Do(req)
	}, len(QuerierHosts), nil)

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

func initQuerier(hosts []QuerierHost, APIKey string, interceptor func(*http.Request, UserContext) *http.Request, disableCache bool) {
	if !querierInitCalled {
		querierInitCalled = true
		QuerierHosts = hosts
		if APIKey != "" {
			QuerierAPIKey = &APIKey
		}
		querierAPIVersion = ""
		querierLastTriedIndex = 0
		querierInterceptor = interceptor
		querierGlobalCacheTag = GetCurrTimeInMS()
		querierDisableCache = disableCache
	}
}

func (q *Querier) SendPostRequest(path string, data map[string]interface{}, userContext UserContext) (map[string]interface{}, error) {
	q.InvalidateCoreCallCache(userContext, true)
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	resp, _, err := q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
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

		apiVersion, querierAPIVersionError := q.GetQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("content-type", "application/json; charset=utf-8")
		req.Header.Set("cdi-version", apiVersion)
		if QuerierAPIKey != nil {
			req.Header.Set("api-key", *QuerierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		if querierInterceptor != nil {
			req = querierInterceptor(req, userContext)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(QuerierHosts), nil)
	return resp, err
}

func (q *Querier) SendDeleteRequest(path string, data map[string]interface{}, params map[string]string, userContext UserContext) (map[string]interface{}, error) {
	q.InvalidateCoreCallCache(userContext, true)
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	resp, _, err := q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		query := req.URL.Query()

		for k, v := range params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()

		apiVersion, querierAPIVersionError := q.GetQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("content-type", "application/json; charset=utf-8")
		req.Header.Set("cdi-version", apiVersion)
		if QuerierAPIKey != nil {
			req.Header.Set("api-key", *QuerierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		if querierInterceptor != nil {
			req = querierInterceptor(req, userContext)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(QuerierHosts), nil)
	return resp, err
}

func (q *Querier) SendGetRequest(path string, params map[string]string, userContext UserContext) (map[string]interface{}, error) {
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	resp, _, err := q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		query := req.URL.Query()

		// Sort the keys for deterministic order
		sortedKeys := make([]string, 0, len(params))
		for k := range params {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		// Start with the path as the unique key
		uniqueKey := nP.GetAsStringDangerous()

		// Append sorted params to the unique key
		for _, key := range sortedKeys {
			value := params[key]
			uniqueKey += ";" + key + "=" + value
		}

		// Append a separator for headers
		uniqueKey += ";hdrs"

		// Append sorted headers to the unique key
		headers := make(map[string]string)

		apiVersion, querierAPIVersionError := q.GetQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}
		headers["cdi-version"] = apiVersion

		if QuerierAPIKey != nil {
			headers["api-key"] = *QuerierAPIKey
		}

		if nP.IsARecipePath() && q.RIDToCore != "" {
			headers["rid"] = q.RIDToCore
		}

		sortedHeaderKeys := make([]string, 0, len(headers))
		for k := range headers {
			sortedHeaderKeys = append(sortedHeaderKeys, k)
		}
		sort.Strings(sortedHeaderKeys)

		for _, key := range sortedHeaderKeys {
			value := headers[key]
			uniqueKey += ";" + key + "=" + value
		}

		for k, v := range params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		if userContext != nil {
			defaultContext, ok := (*userContext)["_default"].(map[string]interface{})
			if !ok {
				defaultContext = make(map[string]interface{})
			}

			globalCacheTag, ok := defaultContext["globalCacheTag"].(uint64)
			if !ok || globalCacheTag != querierGlobalCacheTag {
				q.InvalidateCoreCallCache(userContext, false)
			}

			coreCallCache, ok := defaultContext["coreCallCache"].(map[string]interface{})
			if !ok {
				coreCallCache = make(map[string]interface{})
			}

			if !querierDisableCache && coreCallCache[uniqueKey] != nil {
				return coreCallCache[uniqueKey].(*http.Response), nil
			}
		}

		if querierInterceptor != nil {
			req = querierInterceptor(req, userContext)
		}

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if response.StatusCode == 200 && !querierDisableCache && userContext != nil {
			defaultContext, ok := (*userContext)["_default"].(map[string]interface{})
			if !ok {
				defaultContext = make(map[string]interface{})
			}

			coreCallCache, ok := defaultContext["coreCallCache"].(map[string]interface{})
			if !ok {
				coreCallCache = make(map[string]interface{})
			}

			coreCallCache[uniqueKey] = response
			defaultContext["coreCallCache"] = coreCallCache
			defaultContext["globalCacheTag"] = querierGlobalCacheTag

			(*userContext)["_default"] = defaultContext
		}

		return response, nil
	}, len(QuerierHosts), nil)
	return resp, err
}

func (q *Querier) SendGetRequestWithResponseHeaders(path string, params map[string]string, userContext UserContext) (map[string]interface{}, http.Header, error) {
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, nil, err
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

		apiVersion, querierAPIVersionError := q.GetQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}
		req.Header.Set("cdi-version", apiVersion)
		if QuerierAPIKey != nil {
			req.Header.Set("api-key", *QuerierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		if querierInterceptor != nil {
			req = querierInterceptor(req, userContext)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(QuerierHosts), nil)
}

func (q *Querier) SendPutRequest(path string, data map[string]interface{}, userContext UserContext) (map[string]interface{}, error) {
	q.InvalidateCoreCallCache(userContext, true)
	nP, err := NewNormalisedURLPath(path)
	if err != nil {
		return nil, err
	}
	resp, _, err := q.sendRequestHelper(nP, func(url string) (*http.Response, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVersion, querierAPIVersionError := q.GetQuerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("content-type", "application/json; charset=utf-8")
		req.Header.Set("cdi-version", apiVersion)
		if QuerierAPIKey != nil {
			req.Header.Set("api-key", *QuerierAPIKey)
		}
		if nP.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		if querierInterceptor != nil {
			req = querierInterceptor(req, userContext)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(QuerierHosts), nil)
	return resp, err
}

func (q *Querier) InvalidateCoreCallCache(userContext UserContext, updGlobalCacheTagIfNecessary bool) {
	if userContext == nil {
		// Create an empty map to avoid nil pointer dereference
		emptyMap := make(map[string]interface{})
		userContext = &emptyMap
	}

	if updGlobalCacheTagIfNecessary {
		defaultContext, ok := (*userContext)["_default"].(map[string]interface{})
		if !ok {
			defaultContext = make(map[string]interface{})
		}

		keepCacheAlive, ok := defaultContext["keepCacheAlive"].(bool)
		if !ok || !keepCacheAlive {
			// Update the global cache tag to invalidate the cache
			querierGlobalCacheTag = GetCurrTimeInMS()
		}
	}

	defaultContext, ok := (*userContext)["_default"].(map[string]interface{})
	if !ok {
		defaultContext = make(map[string]interface{})
	}

	// Clear the core call cache
	defaultContext["coreCallCache"] = make(map[string]interface{})

	(*userContext)["_default"] = defaultContext
}

type httpRequestFunction func(url string) (*http.Response, error)

func GetAllCoreUrlsForPath(path string) []string {
	if QuerierHosts == nil {
		return []string{}
	}

	normalisedPath := NormalisedURLPath{value: path}
	result := []string{}

	for _, host := range QuerierHosts {
		currentDomain := host.Domain.GetAsStringDangerous()
		currentBasePath := host.BasePath.GetAsStringDangerous()

		result = append(result, currentDomain+currentBasePath+normalisedPath.GetAsStringDangerous())
	}

	return result
}

func (q *Querier) sendRequestHelper(path NormalisedURLPath, httpRequest httpRequestFunction, numberOfTries int, retryInfoMap *map[string]int) (map[string]interface{}, http.Header, error) {
	if numberOfTries == 0 {
		return nil, nil, errors.New("no SuperTokens core available to query")
	}

	querierHostLock.Lock()
	currentDomain := QuerierHosts[querierLastTriedIndex].Domain.GetAsStringDangerous()
	currentBasePath := QuerierHosts[querierLastTriedIndex].BasePath.GetAsStringDangerous()
	url := currentDomain + currentBasePath + path.GetAsStringDangerous()

	maxRetries := 5
	var _retryInfoMap map[string]int

	if retryInfoMap != nil {
		_retryInfoMap = *retryInfoMap
	} else {
		_retryInfoMap = map[string]int{}
	}

	_, ok := _retryInfoMap[url]

	if !ok {
		_retryInfoMap[url] = maxRetries
	}

	querierLastTriedIndex = (querierLastTriedIndex + 1) % len(QuerierHosts)
	querierHostLock.Unlock()

	resp, err := httpRequest(url)

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return q.sendRequestHelper(path, httpRequest, numberOfTries-1, &_retryInfoMap)
		}
		if resp != nil {
			resp.Body.Close()
		}
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, nil, readErr
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == RateLimitStatusCode {
			retriesLeft := _retryInfoMap[url]

			if retriesLeft > 0 {
				_retryInfoMap[url] = retriesLeft - 1

				attemptsMade := maxRetries - retriesLeft
				delay := 10 + (250 * attemptsMade)

				time.Sleep(time.Millisecond * time.Duration(delay))

				return q.sendRequestHelper(path, httpRequest, numberOfTries, &_retryInfoMap)
			}
		}

		return nil, nil, fmt.Errorf("SuperTokens core threw an error for a request to path: '%s' with status code: %v and message: %s", path.GetAsStringDangerous(), resp.StatusCode, body)
	}

	headers := resp.Header.Clone()
	finalResult := make(map[string]interface{})
	jsonError := json.Unmarshal(body, &finalResult)
	if jsonError != nil {
		return map[string]interface{}{
			"result": string(body),
		}, headers, nil
	}
	return finalResult, headers, nil
}

func ResetQuerierForTest() {
	querierInitCalled = false
}

func (q *Querier) SetApiVersionForTests(apiVersion string) {
	querierAPIVersion = apiVersion
}
