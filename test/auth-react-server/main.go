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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var latestURLWithToken string = ""
var apiPort string = "8080"
var webPort string = "3031"

func callSTInit() {
	countryOptional := true
	formFields := []epmodels.TypeInputFormField{
		{
			ID: "name",
		},
		{
			ID: "age",
			Validate: func(value interface{}) *string {
				age, _ := strconv.Atoi(value.(string))
				if age >= 18 {
					// return nil to indicate success
					return nil
				}
				err := "You must be over 18 to register"
				return &err
			},
		},
		{
			ID:       "country",
			Optional: &countryOptional,
		},
	}
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:9000",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "localhost:" + apiPort,
			WebsiteDomain: "http://localhost:" + webPort,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: formFields,
				},
				ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string) {
						latestURLWithToken = passwordResetURLWithToken
					},
				},
				EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
					CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string) {
						latestURLWithToken = emailVerificationURLWithToken
					},
				},
			}),
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.TypeProvider{
						thirdparty.Google(tpmodels.GoogleConfig{
							ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
							ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
						}),
						thirdparty.Github(tpmodels.GithubConfig{
							ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
							ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
						}),
						thirdparty.Facebook(tpmodels.FacebookConfig{
							ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
							ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
						}),
					},
				},
			}),
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: formFields,
				},
				Providers: []tpmodels.TypeProvider{
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
						ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
					}),
					thirdparty.Github(tpmodels.GithubConfig{
						ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
						ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
					}),
					thirdparty.Facebook(tpmodels.FacebookConfig{
						ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
						ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
					}),
				},
			}),
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+webPort)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

func main() {
	godotenv.Load()
	if len(os.Args) >= 2 {
		apiPort = os.Args[1]
	}
	if len(os.Args) >= 3 {
		webPort = os.Args[2]
	}
	supertokens.IsTestFlag = true
	port := "8080"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	callSTInit()

	http.ListenAndServe(":"+port, corsMiddleware(
		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/sessioninfo" && r.Method == "GET" {
				session.VerifySession(nil, sessioninfo).ServeHTTP(rw, r)
			} else if r.URL.Path == "/token" && r.Method == "GET" {
				rw.WriteHeader(200)
				rw.Header().Add("content-type", "application/json")
				bytes, _ := json.Marshal(map[string]interface{}{
					"latestURLWithToken": latestURLWithToken,
				})
				rw.Write(bytes)
			}
		}))))
}

func sessioninfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	sessionData, err := sessionContainer.GetSessionData()
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	w.WriteHeader(200)
	w.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(map[string]interface{}{
		"sessionHandle": sessionContainer.GetHandle(),
		"userId":        sessionContainer.GetUserID(),
		"jwtPayload":    sessionContainer.GetJWTPayload(),
		"sessionData":   sessionData,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error in converting to json"))
	} else {
		w.Write(bytes)
	}
}
