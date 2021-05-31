package supertokens

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	if appInfo.apiDomain == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your apiDomain inside the appInfo object when calling supertokens.init")
	}
	if appInfo.appName == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your appName inside the appInfo object when calling supertokens.init")
	}
	if appInfo.websiteDomain == "" {
		return NormalisedAppinfo{}, errors.New("Please provide your websiteDomain inside the appInfo object when calling supertokens.init")
	}
	return NormalisedAppinfo{}, nil
}

func getDataFromFileForServerlessCache(filePath string) string {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(byteValue)
}

func containsHost(hostsAlive []string, host string) bool {
	if len(hostsAlive) == 0 {
		return false
	}
	for _, value := range hostsAlive {
		if value == host {
			return true
		}
	}
	return false
}

func getLargestVersionFromIntersection(v1 []string, v2 []string) *string {
	var intersection = []string{}
	for i := 0; i < len(v1); i++ {
		for y := 0; y < len(v2); y++ {
			if v1[i] == v2[y] {
				intersection = append(intersection, v1[i])
			}
		}
	}
	if len(intersection) == 0 {
		return nil
	}
	maxVersionSoFar := intersection[0]
	for i := 1; i < len(intersection); i++ {
		maxVersionSoFar = MaxVersion(intersection[i], maxVersionSoFar)
	}
	return &maxVersionSoFar
}

// MaxVersion returns max of v1 and v2
func MaxVersion(version1 string, version2 string) string {
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

func storeIntoTempFolderForServerlessCache(filePath, data string) {
	// todo
}

// func normaliseHttpMethod(method string) string {
// 	return strings.ToLower(method)
// }

func getRIDFromRequest(r *http.Request) string {
	return r.Header.Get("HEADER_RID")
}
