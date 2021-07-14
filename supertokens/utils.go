package supertokens

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func IsAnIPAddress(ipaddress string) (bool, error) {
	return regexp.MatchString(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, ipaddress)
}

func NormaliseInputAppInfoOrThrowError(appInfo AppInfo) (NormalisedAppinfo, error) {
	if reflect.DeepEqual(appInfo, AppInfo{}) {
		return NormalisedAppinfo{}, errors.New("Please provide the appInfo object when calling supertokens.init")
	}
	if appInfo.APIDomain == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your apiDomain inside the appInfo object when calling supertokens.init")
	}
	if appInfo.AppName == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your appName inside the appInfo object when calling supertokens.init")
	}
	if appInfo.WebsiteDomain == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your websiteDomain inside the appInfo object when calling supertokens.init")
	}
	apiGatewayPath, err := NewNormalisedURLPath("")
	if err != nil {
		return NormalisedAppinfo{}, err
	}
	if appInfo.APIGatewayPath != nil {
		apiGatewayPath, err = NewNormalisedURLPath(*appInfo.APIGatewayPath)
		if err != nil {
			return NormalisedAppinfo{}, err
		}
	}
	websiteDomain, err := NewNormalisedURLDomain(appInfo.WebsiteDomain, false)
	if err != nil {
		return NormalisedAppinfo{}, err
	}
	apiDomain, err := NewNormalisedURLDomain(appInfo.WebsiteDomain, false)
	if err != nil {
		return NormalisedAppinfo{}, err
	}

	APIBasePathStr := "/auth"
	if appInfo.APIBasePath != nil {
		APIBasePathStr = *appInfo.APIBasePath
	}
	APIBasePathURL, err := NewNormalisedURLPath(APIBasePathStr)
	if err != nil {
		return NormalisedAppinfo{}, err
	}
	apiBasePath := apiGatewayPath.AppendPath(*APIBasePathURL)

	websiteBasePathStr := "/auth"
	if appInfo.WebsiteBasePath != nil {
		websiteBasePathStr = *appInfo.WebsiteBasePath
	}
	websiteBasePath, err := NewNormalisedURLPath(websiteBasePathStr)
	if err != nil {
		return NormalisedAppinfo{}, err
	}
	return NormalisedAppinfo{
		AppName:         appInfo.AppName,
		APIGatewayPath:  *apiGatewayPath,
		WebsiteDomain:   *websiteDomain,
		APIDomain:       *apiDomain,
		APIBasePath:     apiBasePath,
		WebsiteBasePath: *websiteBasePath,
	}, nil
}

func getLargestVersionFromIntersection(v1 []string, v2 []string) *string {
	var intersection = []string{}
	for _, i := range v1 {
		for _, j := range v2 {
			if i == j {
				intersection = append(intersection, i)
			}
		}
	}
	if len(intersection) == 0 {
		return nil
	}
	maxVersionSoFar := intersection[0]
	for i := 1; i < len(intersection); i++ {
		maxVersionSoFar = maxVersion(intersection[i], maxVersionSoFar)
	}
	return &maxVersionSoFar
}

// MaxVersion returns max of v1 and v2
func maxVersion(version1 string, version2 string) string {
	var splittedv1 = strings.Split(version1, ".")
	var splittedv2 = strings.Split(version2, ".")
	var minLength = len(splittedv1)
	if minLength > len(splittedv2) {
		minLength = len(splittedv2)
	}
	for i := 0; i < minLength; i++ {
		var v1, _ = strconv.Atoi(splittedv1[i])
		var v2, _ = strconv.Atoi(splittedv2[i])
		if v1 > v2 {
			return version1
		} else if v2 > v1 {
			return version2
		}
	}
	if len(splittedv1) >= len(splittedv2) {
		return version1
	}
	return version2
}

func getRIDFromRequest(r *http.Request) string {
	return r.Header.Get(HeaderRID)
}

func Send200Response(res http.ResponseWriter, responseJson interface{}) {
	res.WriteHeader(200)
	res.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(responseJson)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("error in json marshalling"))
	}
	res.Write(bytes)
}

func SendNon200Response(res http.ResponseWriter, message string, statusCode int) error {
	if statusCode < 300 {
		return errors.New("Calling sendNon200Response with status code < 300")
	}
	res.WriteHeader(statusCode)
	res.Header().Add("content-type", "application/json")
	response := map[string]interface{}{
		"message": message,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("error in json marshalling"))
	}
	res.Write(bytes)
	return nil
}
