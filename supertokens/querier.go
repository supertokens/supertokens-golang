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
	RIDToCore string
}

// TODO: these are declared at global scope.. will they get reinit on every import?
var InitCalled bool = false
var Hosts []NormalisedURLDomain = nil
var APIKey *string
var APIVersion string
var LastTriedIndex int

func NewQuerier(rIdToCore string) Querier {
	return Querier{
		RIDToCore: rIdToCore,
	}
}

func (q *Querier) getAPIVersion() (string, error) {
	if APIVersion != "" {
		return APIVersion, nil
	}
	response, err := q.sendRequestHelper(NormalisedURLPath{value: "/apiversion"}, func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if APIKey != nil {
			req.Header.Set("api-key", *APIKey)
		}
		client := &http.Client{}
		return client.Do(req)
	}, len(Hosts))

	if err != nil {
		return "", err
	}

	cdiSupportedByServer := strings.Split(response["versions"], ",")
	supportedVersion := getLargestVersionFromIntersection(cdiSupportedByServer, cdiSupported)
	if supportedVersion == nil {
		return "", errors.New("The running SuperTokens core version is not compatible with this Golang SDK. Please visit https://supertokens.io/docs/community/compatibility to find the right version")
	}

	APIVersion = *supportedVersion

	return APIVersion, nil
}

func GetNewQuerierInstanceOrThrowError(rIDToCore string) (*Querier, error) {
	if InitCalled == false {
		return nil, errors.New("Please call the supertokens.init function before using SuperTokens")
	}
	return &Querier{RIDToCore: rIDToCore}, nil
}

func InitQuerier(hosts []NormalisedURLDomain, apiKey *string) {
	if InitCalled == false {
		InitCalled = true
		Hosts = hosts
		APIKey = apiKey
		APIVersion = ""
		LastTriedIndex = 0
	}
}

func (q *Querier) SendPostRequest(path NormalisedURLPath, data map[string]string) (map[string]string, error) {
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
		if APIKey != nil {
			req.Header.Set("api-key", *APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(Hosts))
}

func (q *Querier) SendDeleteRequest(path NormalisedURLPath, data map[string]string) (map[string]string, error) {
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
		if APIKey != nil {
			req.Header.Set("api-key", *APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(Hosts))
}

func (q *Querier) SendGetRequest(path NormalisedURLPath, params map[string]string) (map[string]string, error) {
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
		if APIKey != nil {
			req.Header.Set("api-key", *APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(Hosts))
}

func (q *Querier) SendPutRequest(path NormalisedURLPath, data map[string]string) (map[string]string, error) {
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
		if APIKey != nil {
			req.Header.Set("api-key", *APIKey)
		}
		if path.IsARecipePath() && q.RIDToCore != "" {
			req.Header.Set("rid", q.RIDToCore)
		}

		client := &http.Client{}
		return client.Do(req)
	}, len(Hosts))
}

type httpRequestFunction func(url string) (*http.Response, error)

func (q *Querier) sendRequestHelper(path NormalisedURLPath, httpRequest httpRequestFunction,
	numberOfTries int) (map[string]string, error) {
	if numberOfTries == 0 {
		return nil, errors.New("No SuperTokens core available to query")
	}
	currentHost := Hosts[LastTriedIndex].GetAsStringDangerous()
	LastTriedIndex = (LastTriedIndex + 1) % len(Hosts)
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

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	finalResult := make(map[string]string)
	jsonError := json.Unmarshal(body, &finalResult)
	if jsonError != nil {
		return map[string]string{
			"result": string(body),
		}, nil
	}
	return finalResult, nil
}
