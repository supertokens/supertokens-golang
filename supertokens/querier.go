package supertokens

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Querier struct {
	InitCalled     bool
	Hosts          []NormalisedURLDomain
	APIKey         string
	APIVersion     string
	LastTriedIndex int
	RIDToCore      string
}

func NewQuerier(hosts []NormalisedURLDomain, rIdToCore string) *Querier {
	return &Querier{
		Hosts:     hosts,
		RIDToCore: rIdToCore,
	}
}

func (q *Querier) getAPIVersion() (string, error) {
	if q.APIVersion != "" {
		return q.APIVersion, nil
	}
	response, err := q.sendRequestHelper(NormalisedURLPath{value: "/apiversion"}, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if q.APIKey != "" {
			req.Header.Set("api-key", q.APIKey)
		}
		client := &http.Client{}
		return client.Do(req)
	}, len(q.Hosts))

	if err != nil {
		return "", err
	}

	cdiSupportedByServer := response["versions"].([]string)
	supportedVersion := getLargestVersionFromIntersection(cdiSupportedByServer, cdiSupported)
	if supportedVersion == nil {
		return "", errors.New("The running SuperTokens core version is not compatible with this Golang SDK. Please visit https://supertokens.io/docs/community/compatibility to find the right version")
	}

	q.APIVersion = *supportedVersion

	return q.APIVersion, nil
}

func (q *Querier) GetNewInstanceOrThrowError(rIDToCore string) (*Querier, error) {
	if q.InitCalled == false {
		return nil, errors.New("Please call the supertokens.init function before using SuperTokens")
	}
	return &Querier{Hosts: q.Hosts, RIDToCore: rIDToCore}, nil
}

func (q *Querier) Init(hosts []NormalisedURLDomain, apiKey string) {
	if q.InitCalled == false {
		q.InitCalled = true
		q.Hosts = hosts
		q.APIKey = apiKey
		q.APIVersion = ""
		q.LastTriedIndex = 0
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

		apiVerion, apiVersionError := q.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cdi-version", apiVerion)
		if q.APIKey != "" {
			req.Header.Set("api-key", q.APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(q.Hosts))
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

		apiVerion, apiVersionError := q.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cdi-version", apiVerion)
		if q.APIKey != "" {
			req.Header.Set("api-key", q.APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(q.Hosts))
}

func (q *Querier) SendGetRequest(path NormalisedURLPath, params map[string]string) (map[string]interface{}, error) {
	return q.sendRequestHelper(path, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		query := req.URL.Query()

		for k, v := range params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()

		apiVerion, apiVersionError := q.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}
		req.Header.Set("cdi-version", apiVerion)
		if q.APIKey != "" {
			req.Header.Set("api-key", q.APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(q.Hosts))
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

		apiVerion, apiVersionError := q.getAPIVersion()
		if apiVersionError != nil {
			return nil, apiVersionError
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("cdi-version", apiVerion)
		if q.APIKey != "" {
			req.Header.Set("api-key", q.APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(q.Hosts))
}

type httpRequestFunction func(url string) (*http.Response, error)

func (q *Querier) sendRequestHelper(path NormalisedURLPath, httpRequest httpRequestFunction,
	numberOfTries int) (map[string]interface{}, error) {
	if numberOfTries == 0 {
		return nil, errors.New("No SuperTokens core available to query")
	}
	currentHost := q.Hosts[q.LastTriedIndex].GetAsStringDangerous()
	q.LastTriedIndex = (q.LastTriedIndex + 1) % len(q.Hosts)
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
		return nil, errors.New(strconv.Itoa(resp.StatusCode))
	}

	var body, readErr = ioutil.ReadAll(resp.Body)
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
