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
	"io"
	"io/ioutil"
	"math"
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
	websiteDomain, err := NewNormalisedURLDomain(appInfo.WebsiteDomain)
	if err != nil {
		return NormalisedAppinfo{}, err
	}
	apiDomain, err := NewNormalisedURLDomain(appInfo.APIDomain)
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
	apiBasePath := apiGatewayPath.AppendPath(APIBasePathURL)

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
		APIGatewayPath:  apiGatewayPath,
		WebsiteDomain:   websiteDomain,
		APIDomain:       apiDomain,
		APIBasePath:     apiBasePath,
		WebsiteBasePath: websiteBasePath,
	}, nil
}

// TODO: Add tests
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

func Send200Response(res http.ResponseWriter, responseJson interface{}) error {
	LogDebugMessage("Sending response to client with status code: 200")
	dw := MakeDoneWriter(res)
	if !dw.IsDone() {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(200)
		bytes, err := json.Marshal(responseJson)
		if err != nil {
			return err
		} else {
			res.Write(bytes)
		}
	}
	return nil
}

func SendHTMLResponse(res http.ResponseWriter, statusCode int, htmlString string) error {
	LogDebugMessage(fmt.Sprintf("Sending HTML response to client with status code: %d", statusCode))
	dw := MakeDoneWriter(res)
	if !dw.IsDone() {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		res.WriteHeader(200)
		_, err := fmt.Fprint(res, htmlString)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendNon200ResponseWithMessage(res http.ResponseWriter, message string, statusCode int) error {
	return SendNon200Response(res, statusCode, map[string]interface{}{"message": message})
}

func SendNon200Response(res http.ResponseWriter, statusCode int, body map[string]interface{}) error {
	dw := MakeDoneWriter(res)
	if !dw.IsDone() {
		if statusCode < 300 {
			return errors.New("calling SendNon200Response with status code < 300")
		}

		LogDebugMessage("Sending response to client with status code: " + strconv.Itoa(statusCode))

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(statusCode)

		bytes, err := json.Marshal(body)
		if err != nil {
			return err
		} else {
			res.Write(bytes)
		}
	}
	return nil
}

func SendUnauthorisedAccess(res http.ResponseWriter) error {
	return SendNon200ResponseWithMessage(res, "unauthorised access", 401)
}

func ReadFromRequest(r *http.Request) ([]byte, error) {
	f := r.Body
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return buf, err
	}

	r.Body = io.NopCloser(bytes.NewReader(buf))

	return buf, nil
}

func formatOneDecimalFloat(n float64) (string, string) {
	n = math.Floor(n*10) / 10
	if float64(int(n)) == n {
		if n == 1.0 {
			return fmt.Sprintf("%d", int(n)), ""
		} else {
			return fmt.Sprintf("%d", int(n)), "s"
		}
	}
	return fmt.Sprintf("%.1f", n), "s"
}

func HumaniseMilliseconds(m uint64) string {
	t := m / 1000
	var suffix string = ""
	if t < 60 {
		if t > 1 {
			suffix = "s"
		}
		return fmt.Sprintf("%d second%s", t, suffix)
	} else if t < 3600 {
		if t/60 > 1 {
			suffix = "s"
		}
		return fmt.Sprintf("%d minute%s", t/60, suffix)
	}
	if t/3600 > 1 {
		suffix = "s"
	}
	h := float64(t) / 3600
	hStr, suffix := formatOneDecimalFloat(h)
	return fmt.Sprintf("%s hour%s", hStr, suffix)
}

func ConvertGeneralErrorToJsonResponse(resp GeneralErrorResponse) map[string]interface{} {
	return map[string]interface{}{
		"status":  "GENERAL_ERROR",
		"message": resp.Message,
	}
}

func ErrorIfNoResponse(res http.ResponseWriter) error {
	dw := MakeDoneWriter(res)
	if !dw.IsDone() {
		return errors.New("invalid return from API interface function")
	}
	return nil
}

func MakeDefaultUserContextFromAPI(r *http.Request) UserContext {
	return &map[string]interface{}{
		"_default": map[string]interface{}{
			"request": r,
		},
	}
}
