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

var querierLock sync.Mutex

type Querier struct {
	RIDToCore string
}

var (
	querierInitCalled     bool                  = false
	querierHosts          []NormalisedURLDomain = nil
	querierAPIKey         *string
	querierAPIVersion     string
	querierLastTriedIndex int
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
	// TODO: Why do we have locking here?
	querierLock.Lock()
	defer querierLock.Unlock()
	if !querierInitCalled {
		return nil, errors.New("please call the supertokens.init function before using SuperTokens")
	}
	return &Querier{RIDToCore: rIDToCore}, nil
}

func initQuerier(hosts []NormalisedURLDomain, APIKey *string) {
	// TODO: Why do we have locking here?
	querierLock.Lock()
	defer querierLock.Unlock()
	if !querierInitCalled {
		querierInitCalled = true
		querierHosts = hosts
		querierAPIKey = APIKey
		querierAPIVersion = ""
		querierLastTriedIndex = 0
	}
}

func (q *Querier) SendPostRequest(path NormalisedURLPath, data map[string]interface{}) (map[string]interface{}, error) {
	return q.sendRequestHelper(path, func(url string) (*http.Response, error) {
		if data == nil {
			data = map[string]interface{}{}
		}
		// TODO: what is the need to do this - since this is not being done in DELETE or any other place.
		for key, value := range data {
			switch value.(type) {
			case map[string]interface{}:
				if len(value.(map[string]interface{})) == 0 {
					data[key] = map[string]interface{}{}
				}
			}
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
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

func (q *Querier) SendDeleteRequest(path NormalisedURLPath, data map[string]interface{}) (map[string]interface{}, error) {
	return q.sendRequestHelper(path, func(url string) (*http.Response, error) {
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
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

func (q *Querier) SendGetRequest(path NormalisedURLPath, params map[string]interface{}) (map[string]interface{}, error) {
	return q.sendRequestHelper(path, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		query := req.URL.Query()

		for k, v := range params {
			query.Add(k, v.(string))
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
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

func (q *Querier) SendPutRequest(path NormalisedURLPath, data map[string]interface{}) (map[string]interface{}, error) {
	return q.sendRequestHelper(path, func(url string) (*http.Response, error) {
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
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(querierHosts))
}

type httpRequestFunction func(url string) (*http.Response, error)

// TODO: Add tests
func (q *Querier) sendRequestHelper(path NormalisedURLPath, httpRequest httpRequestFunction, numberOfTries int) (map[string]interface{}, error) {
	if numberOfTries == 0 {
		return nil, errors.New("no SuperTokens core available to query")
	}
	currentHost := querierHosts[querierLastTriedIndex].GetAsStringDangerous()
	// TODO: won't we need to apply some sort of locking here when updating querierLastTriedIndex?
	querierLastTriedIndex = (querierLastTriedIndex + 1) % len(querierHosts)
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
