// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/supertokens/supertokens-golang/examples/with-twirp/haberdasher"
	"github.com/supertokens/supertokens-golang/examples/with-twirp/internal/haberdasherserver"
	"github.com/supertokens/supertokens-golang/examples/with-twirp/internal/hooks"
	"github.com/supertokens/supertokens-golang/examples/with-twirp/internal/interceptor"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/twitchtv/twirp"
)

func main() {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     "http://localhost:3001",
			WebsiteDomain: "http://localhost:3000",
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				Providers: []tpmodels.TypeProvider{
					// We have provided you with development keys which you can use for testsing.
					// IMPORTANT: Please replace them with your own OAuth keys for production use.
					thirdparty.Google(tpmodels.GoogleConfig{
						ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
						ClientSecret: "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",
					}),
					thirdparty.Github(tpmodels.GithubConfig{
						ClientID:     "467101b197249757c71f",
						ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
					}),
					thirdparty.Apple(tpmodels.AppleConfig{
						ClientID: "4398792-io.supertokens.example.service",
						ClientSecret: tpmodels.AppleClientSecret{
							KeyId:      "7M48Y4RYDL",
							PrivateKey: "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
							TeamId:     "YWQCXGJRJL",
						},
					}),
				},
			}),
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	hook := hooks.LoggingHooks(os.Stderr)
	service := haberdasherserver.New()
	server := haberdasher.NewHaberdasherServer(service,
		twirp.WithServerInterceptors(interceptor.NewSuperTokensErrorHandlerInterceptor()), hook)
	sessionRequired := false
	log.Fatal(http.ListenAndServe(":3001", handlers.CORS(
		handlers.AllowedHeaders(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...)),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowCredentials(),
	)(supertokens.Middleware(session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &sessionRequired,
	}, server.ServeHTTP)))))
}
