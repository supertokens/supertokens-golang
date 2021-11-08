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
				/*
				   We use different credentials for different platforms when required. For example the redirect URI for Github
				   is different for Web and mobile. In such a case we can provide multiple providers with different client Ids.

				   When the frontend makes a request and wants to use a specific clientId, it needs to send the clientId to use in the
				   request. In the absence of a clientId in the request the SDK uses the default provider, indicated by `isDefault: true`.
				   When adding multiple providers for the same type (Google, Github etc), make sure to set `isDefault: true`.
				*/
				Providers: []tpmodels.TypeProvider{
					// We have provided you with development keys which you can use for testsing.
					// IMPORTANT: Please replace them with your own OAuth keys for production use.
					thirdparty.Google(tpmodels.GoogleConfig{
						// We use this for websites
						IsDefault:    true,
						ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
						ClientSecret: "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",
					}),
					thirdparty.Google(tpmodels.GoogleConfig{
						// we use this for mobile apps
						ClientID:     "1060725074195-c7mgk8p0h27c4428prfuo3lg7ould5o7.apps.googleusercontent.com",
						ClientSecret: "", // this is empty because we follow Authorization code grant flow via PKCE for mobile apps (Google doesn't issue a client secret for mobile apps).
					}),
					thirdparty.Github(tpmodels.GithubConfig{
						// We use this for websites
						IsDefault:    true,
						ClientID:     "467101b197249757c71f",
						ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
					}),
					thirdparty.Github(tpmodels.GithubConfig{
						// We use this for mobile apps
						ClientID:     "8a9152860ce869b64c44",
						ClientSecret: "00e841f10f288363cd3786b1b1f538f05cfdbda2",
					}),
					/*
					   For Apple signin, iOS apps always use the bundle identifier as the client ID when communicating with Apple. Android, Web and other platforms
					   need to configure a Service ID on the Apple developer dashboard and use that as client ID.
					   In the example below 4398792-io.supertokens.example.service is the client ID for Web. Android etc and thus we mark it as default. For iOS
					   the frontend for the demo app sends the clientId in the request which is then used by the SDK.
					*/
					thirdparty.Apple(tpmodels.AppleConfig{
						// For Android and website apps
						IsDefault: true,
						ClientID:  "4398792-io.supertokens.example.service",
						ClientSecret: tpmodels.AppleClientSecret{
							KeyId:      "7M48Y4RYDL",
							PrivateKey: "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
							TeamId:     "YWQCXGJRJL",
						},
					}),
					thirdparty.Apple(tpmodels.AppleConfig{
						// For iOS Apps
						ClientID: "4398792-io.supertokens.example",
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
