package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type OverrideLog struct {
	T    int64       `json:"t"`
	Name string      `json:"name"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	overrideLogs []OverrideLog
	logMutex     sync.Mutex
)

func resetOverrideLogs() {
	logMutex.Lock()
	defer logMutex.Unlock()
	overrideLogs = []OverrideLog{}
}

func getOverrideLogs() []OverrideLog {
	logMutex.Lock()
	defer logMutex.Unlock()
	return overrideLogs
}

func logOverrideEvent(name string, logType string, data interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	overrideLogs = append(overrideLogs, OverrideLog{
		T:    time.Now().UnixNano() / int64(time.Millisecond),
		Name: name,
		Type: logType,
		Data: transformLoggedData(data),
	})
}

func transformLoggedData(data interface{}) interface{} {
	visited := make(map[interface{}]bool)
	return transformLoggedDataRecursive(data, visited)
}

func transformLoggedDataRecursive(data interface{}, visited map[interface{}]bool) interface{} {
	if data == nil {
		return nil
	}

	transformMap := func(data map[string][]string) map[string]interface{} {
		res := make(map[string]interface{})
		for k, v := range data {
			if len(v) == 1 {
				res[k] = v[0]
			} else {
				res[k] = v
			}
		}
		return res
	}

	if req, ok := data.(*http.Request); ok {
		return map[string]interface{}{
			"method":  req.Method,
			"url":     req.URL.Scheme + "://" + req.URL.Host + req.URL.Path,
			"headers": transformMap(req.Header),
			"params":  transformMap(req.URL.Query()),
		}
	}

	if _, ok := data.([]interface{}); ok {
		res := make([]interface{}, len(data.([]interface{})))
		for i, item := range data.([]interface{}) {
			res[i] = transformLoggedDataRecursive(item, visited)
		}
		return res
	}

	if v, ok := data.(supertokens.UserContext); ok {
		if v == nil {
			return nil
		}
		return transformLoggedDataRecursive(*v, visited)
	}

	if v, ok := data.(map[string]interface{}); ok {
		res := make(map[string]interface{})
		for k, item := range v {
			if k == "_default" {
				continue
			}
			res[k] = transformLoggedDataRecursive(item, visited)
		}
		return res
	}

	// if visited[data] {
	// 	return "VISITED"
	// }
	// visited[data] = true

	// switch v := data.(type) {
	// case []interface{}:
	// 	result := make([]interface{}, len(v))
	// 	for i, item := range v {
	// 		result[i] = transformLoggedDataRecursive(item, visited)
	// 	}
	// 	return result
	// case map[string]interface{}:
	// 	result := make(map[string]interface{})
	// 	for key, value := range v {
	// 		result[key] = transformLoggedDataRecursive(value, visited)
	// 	}
	// 	return result
	// case json.Marshaler:
	// 	jsonData, _ := v.MarshalJSON()
	// 	var unmarshaled interface{}
	// 	json.Unmarshal(jsonData, &unmarshaled)
	// 	return transformLoggedDataRecursive(unmarshaled, visited)
	// }

	return data
}
