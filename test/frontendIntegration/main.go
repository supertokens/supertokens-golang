/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var noOfTimesGetSessionCalledDuringTest int = 0
var noOfTimesRefreshCalledDuringTest int = 0
var noOfTimesRefreshAttemptedDuringTest int = 0
var lastAntiCsrfSetting bool
var lastEnableJWTSetting bool

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

var routes *http.Handler

func callSTInit(enableAntiCsrf bool, enableJWT bool, jwtPropertyName string) {
	if maxVersion(supertokens.VERSION, "0.3.1") == supertokens.VERSION && enableJWT {
		port := "8080"
		if len(os.Args) == 2 {
			port = os.Args[1]
		}
		antiCsrf := "NONE"
		if enableAntiCsrf {
			antiCsrf = "VIA_TOKEN"
		}
		err := supertokens.Init(supertokens.TypeInput{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:9000",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				APIDomain:     "0.0.0.0:" + port,
				WebsiteDomain: "http://localhost.org:8080",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					Jwt: &sessmodels.JWTInputConfig{
						Enable:                           true,
						PropertyNameInAccessTokenPayload: &jwtPropertyName,
					},
					ErrorHandlers: &sessmodels.ErrorHandlers{
						OnUnauthorised: func(message string, req *http.Request, res http.ResponseWriter) error {
							res.Header().Set("Content-Type", "text/html; charset=utf-8")
							res.WriteHeader(401)
							res.Write([]byte(""))
							return nil
						},
					},
					AntiCsrf: &antiCsrf,
					Override: &sessmodels.OverrideStruct{
						Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
							ogCNS := *originalImplementation.CreateNewSession
							(*originalImplementation.CreateNewSession) = func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
								if accessTokenPayload == nil {
									accessTokenPayload = map[string]interface{}{}
								}
								accessTokenPayload["customClaim"] = "customValue"

								return ogCNS(req, res, userID, accessTokenPayload, sessionData, userContext)
							}
							return originalImplementation
						},
						APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
							originalImplementation.RefreshPOST = nil
							return originalImplementation
						},
					},
				}),
			},
		})

		if err != nil {
			panic(err.Error())
		}
	} else {
		port := "8080"
		if len(os.Args) == 2 {
			port = os.Args[1]
		}
		antiCsrf := "NONE"
		if enableAntiCsrf {
			antiCsrf = "VIA_TOKEN"
		}
		err := supertokens.Init(supertokens.TypeInput{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:9000",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				APIDomain:     "0.0.0.0:" + port,
				WebsiteDomain: "http://localhost.org:8080",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					ErrorHandlers: &sessmodels.ErrorHandlers{
						OnUnauthorised: func(message string, req *http.Request, res http.ResponseWriter) error {
							res.Header().Set("Content-Type", "text/html; charset=utf-8")
							res.WriteHeader(401)
							res.Write([]byte(""))
							return nil
						},
					},
					AntiCsrf: &antiCsrf,
					Override: &sessmodels.OverrideStruct{
						APIs: func(originalImplementation sessmodels.APIInterface) sessmodels.APIInterface {
							originalImplementation.RefreshPOST = nil
							return originalImplementation
						},
					},
				}),
			},
		})

		if err != nil {
			panic(err.Error())
		}
	}

	middleware := supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/setAntiCsrf" && r.Method == "POST" {
			setAntiCsrf(rw, r)
		} else if r.URL.Path == "/setEnableJWT" && r.Method == "POST" {
			setEnableJWT(rw, r)
		} else if r.URL.Path == "/login" && r.Method == "POST" {
			login(rw, r)
		} else if r.URL.Path == "/beforeeach" && r.Method == "POST" {
			beforeeach(rw, r)
		} else if r.URL.Path == "/testUserConfig" && r.Method == "POST" {
			testUserConfig(rw, r)
		} else if r.URL.Path == "/multipleInterceptors" && r.Method == "POST" {
			multipleInterceptors(rw, r)
		} else if r.URL.Path == "/" && r.Method == "GET" {
			session.VerifySession(nil, simpleGet).ServeHTTP(rw, r)
		} else if r.URL.Path == "/check-rid" && r.Method == "GET" {
			session.VerifySession(nil, checkRID).ServeHTTP(rw, r)
		} else if r.URL.Path == "/update-jwt" && r.Method == "GET" {
			session.VerifySession(nil, getJWT).ServeHTTP(rw, r)
		} else if r.URL.Path == "/update-jwt" && r.Method == "POST" {
			session.VerifySession(nil, updateJwt).ServeHTTP(rw, r)
		} else if r.URL.Path == "/update-jwt-with-handle" && r.Method == "POST" {
			session.VerifySession(nil, updateJwtWithHandle).ServeHTTP(rw, r)
		} else if r.URL.Path == "/testing" {
			testing(rw, r)
		} else if r.URL.Path == "/logout" && r.Method == "POST" {
			session.VerifySession(nil, logout).ServeHTTP(rw, r)
		} else if r.URL.Path == "/revokeAll" && r.Method == "POST" {
			session.VerifySession(nil, revokeAll).ServeHTTP(rw, r)
		} else if r.URL.Path == "/auth/session/refresh" && r.Method == "POST" {
			refresh(rw, r)
		} else if r.URL.Path == "/refreshCalledTime" && r.Method == "GET" {
			rw.Write([]byte(strconv.Itoa(noOfTimesRefreshCalledDuringTest)))
		} else if r.URL.Path == "/refreshAttemptedTime" && r.Method == "GET" {
			rw.Write([]byte(strconv.Itoa(noOfTimesRefreshAttemptedDuringTest)))
		} else if r.URL.Path == "/getSessionCalledTime" && r.Method == "GET" {
			rw.Write([]byte(strconv.Itoa(noOfTimesGetSessionCalledDuringTest)))
		} else if r.URL.Path == "/ping" && r.Method == "GET" {
			rw.Write([]byte(""))
		} else if r.URL.Path == "/testHeader" && r.Method == "GET" {
			testHeader(rw, r)
		} else if r.URL.Path == "/checkAllowCredentials" && r.Method == "POST" {
			rw.Write([]byte(strconv.FormatBool(r.Header.Get("allow-credentials") != "")))
		} else if r.URL.Path == "/index.html" && r.Method == "GET" {
			index(rw, r)
		} else if r.URL.Path == "/featureFlags" && r.Method == "GET" {
			featureFlag(rw, r)
		} else if r.URL.Path == "/testError" && r.Method == "GET" {
			code := 500
			codeStr := r.URL.Query().Get("code")
			if codeStr != "" {
				code, _ = strconv.Atoi(codeStr)
			}
			rw.WriteHeader(code)
			rw.Write([]byte("test error message"))
		} else if r.URL.Path == "/reinitialiseBackendConfig" && r.Method == "POST" {
			reinitialiseBackendConfig(rw, r)
		} else if strings.HasPrefix(r.URL.Path, "/angular") {
			angular(rw, r)
		} else {
			fail(rw, r)
		}
	}))

	routes = &middleware
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost.org:8080")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.WriteHeader(204)
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

func main() {
	supertokens.IsTestFlag = true
	port := "8080"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	callSTInit(true, false, "jwt")
	err := http.ListenAndServe(":"+port, corsMiddleware(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			(*routes).ServeHTTP(rw, r)
		})))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func reinitialiseBackendConfig(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)

	var jwtPropertyName string
	if val, ok := body["jwtPropertyName"]; ok {
		jwtPropertyName = val.(string)
	}
	supertokens.ResetForTest()
	session.ResetForTest()
	callSTInit(lastAntiCsrfSetting, lastEnableJWTSetting, jwtPropertyName)
	w.Write([]byte(""))
}

func featureFlag(response http.ResponseWriter, request *http.Request) {
	json.NewEncoder(response).Encode(map[string]interface{}{
		"sessionJwt": maxVersion(supertokens.VERSION, "0.3.1") == supertokens.VERSION && lastEnableJWTSetting,
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	dat, _ := ioutil.ReadFile("./static/index.html")
	w.Header().Set("Content-Type", "text/html")
	w.Write(dat)
}

func angular(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./angular"))
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/angular")
	fs.ServeHTTP(w, r)
}

func testHeader(response http.ResponseWriter, request *http.Request) {
	testheader := request.Header.Get("st-custom-header")
	success := testheader != ""
	json.NewEncoder(response).Encode(map[string]interface{}{
		"success": success,
	})
}

func refresh(response http.ResponseWriter, request *http.Request) {
	noOfTimesRefreshAttemptedDuringTest++
	session.VerifySession(nil, func(rw http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("rid")
		if rid == "" {
			response.Write([]byte("refresh failed"))
		} else {
			noOfTimesRefreshCalledDuringTest++
			response.Write([]byte("refresh success"))
		}
	}).ServeHTTP(response, request)
}

func revokeAll(response http.ResponseWriter, request *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(request.Context())
	userID := sessionContainer.GetUserID()
	session.RevokeAllSessionsForUser(userID)
	response.Write([]byte("success"))
}

func logout(response http.ResponseWriter, request *http.Request) {
	session := session.GetSessionFromRequestContext(request.Context())
	err := session.RevokeSession()
	if err != nil {
		err = supertokens.ErrorHandler(err, request, response)
		if err != nil {
			response.WriteHeader(500)
			response.Write([]byte(""))
		}
		return
	}
	response.Write([]byte("success"))
}

func testing(response http.ResponseWriter, request *http.Request) {
	value := request.Header.Get("testing")
	if value != "" {
		response.Header().Set("testing", value)
	}
	response.Write([]byte("success"))
}

func getJWT(response http.ResponseWriter, request *http.Request) {
	session := session.GetSessionFromRequestContext(request.Context())
	json.NewEncoder(response).Encode(session.GetAccessTokenPayload())
}

func updateJwt(response http.ResponseWriter, request *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(request.Body).Decode(&body)
	userSession := session.GetSessionFromRequestContext(request.Context())
	userSession.UpdateAccessTokenPayload(body)
	json.NewEncoder(response).Encode(userSession.GetAccessTokenPayload())
}

func updateJwtWithHandle(response http.ResponseWriter, request *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(request.Body).Decode(&body)
	userSession := session.GetSessionFromRequestContext(request.Context())
	session.UpdateAccessTokenPayload(userSession.GetHandle(), body)
	json.NewEncoder(response).Encode(userSession.GetAccessTokenPayload())
}

func checkRID(response http.ResponseWriter, request *http.Request) {
	rid := request.Header.Get("rid")
	if rid == "" {
		response.Write([]byte("fail"))
	} else {
		response.Write([]byte("success"))
	}
}

func setAntiCsrf(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)

	enableAntiCsrf := true
	if val, ok := body["enableAntiCsrf"]; ok {
		enableAntiCsrf = val.(bool)
	}
	lastAntiCsrfSetting = enableAntiCsrf
	supertokens.ResetForTest()
	session.ResetForTest()
	callSTInit(enableAntiCsrf, lastEnableJWTSetting, "jwt")
	w.Write([]byte("success"))
}

func setEnableJWT(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)

	enableJWT := false
	if val, ok := body["enableJWT"]; ok {
		enableJWT = val.(bool)
	}
	lastEnableJWTSetting = enableJWT
	supertokens.ResetForTest()
	session.ResetForTest()
	callSTInit(lastAntiCsrfSetting, enableJWT, "jwt")
	w.Write([]byte("success"))
}

func login(response http.ResponseWriter, request *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(request.Body).Decode(&body)
	userID := body["userId"].(string)
	sess, _ := session.CreateNewSession(request, response, userID, nil, nil)
	response.Write([]byte(sess.GetUserID()))
}

func fail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(""))
}

func beforeeach(response http.ResponseWriter, request *http.Request) {
	noOfTimesRefreshCalledDuringTest = 0
	noOfTimesGetSessionCalledDuringTest = 0
	noOfTimesRefreshAttemptedDuringTest = 0
	response.Write([]byte(""))
}

func testUserConfig(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(""))
}

func multipleInterceptors(response http.ResponseWriter, request *http.Request) {
	interceptorheader2 := request.Header.Get("interceptorheader2")
	interceptorheader1 := request.Header.Get("interceptorheader1")

	var resp string
	if interceptorheader2 != "" && interceptorheader1 != "" {
		resp = "success"
	} else {
		resp = "failure"
	}
	response.Write([]byte(resp))
}

func simpleGet(response http.ResponseWriter, request *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(request.Context())
	noOfTimesGetSessionCalledDuringTest += 1
	response.Write([]byte(sessionContainer.GetUserID()))
}
