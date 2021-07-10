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

func NewQuerier(rIdToCore string) Querier {
	return Querier{
		RIDToCore: rIdToCore,
	}
}

func (q *Querier) getquerierAPIVersion() (string, error) {
	querierLock.Lock()
	defer querierLock.Unlock()
	if querierAPIVersion != "" {
		return querierAPIVersion, nil
	}
	response, err := q.sendRequestHelper(NormalisedURLPath{value: "/querierAPIVersion"}, func(url string) (*http.Response, error) {
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

	cdiSupportedByServer := strings.Split(response["versions"].(string), ",")
	supportedVersion := getLargestVersionFromIntersection(cdiSupportedByServer, cdiSupported)
	if supportedVersion == nil {
		return "", errors.New("The running SuperTokens core version is not compatible with this Golang SDK. Please visit https://supertokens.io/docs/community/compatibility to find the right version")
	}

	querierAPIVersion = *supportedVersion

	return querierAPIVersion, nil
}

func GetNewQuerierInstanceOrThrowError(rIDToCore string) (*Querier, error) {
	querierLock.Lock()
	defer querierLock.Unlock()
	if querierInitCalled == false {
		return nil, errors.New("Please call the supertokens.init function before using SuperTokens")
	}
	return &Querier{RIDToCore: rIDToCore}, nil
}

func initQuerier(hosts []NormalisedURLDomain, APIKey *string) {
	querierLock.Lock()
	defer querierLock.Unlock()
	if querierInitCalled == false {
		querierInitCalled = true
		querierHosts = hosts
		querierAPIKey = APIKey
		querierAPIVersion = ""
		querierLastTriedIndex = 0
	}
}

func (q *Querier) SendPostRequest(path NormalisedURLPath, data map[string]interface{}) (map[string]interface{}, error) {
	return q.sendRequestHelper(path, func(url string) (*http.Response, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		apiVerion, querierAPIVersionError := q.getquerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("Content-Type", "application/json")
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

		apiVerion, querierAPIVersionError := q.getquerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("Content-Type", "application/json")
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

		apiVerion, querierAPIVersionError := q.getquerierAPIVersion()
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

		apiVerion, querierAPIVersionError := q.getquerierAPIVersion()
		if querierAPIVersionError != nil {
			return nil, querierAPIVersionError
		}

		req.Header.Set("Content-Type", "application/json")
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

func (q *Querier) sendRequestHelper(path NormalisedURLPath, httpRequest httpRequestFunction, numberOfTries int) (map[string]interface{}, error) {
	if numberOfTries == 0 {
		return nil, errors.New("No SuperTokens core available to query")
	}
	currentHost := querierHosts[querierLastTriedIndex].GetAsStringDangerous()
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

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("%v", resp.StatusCode))
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
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
