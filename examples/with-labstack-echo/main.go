package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
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

	e := echo.New()

	// CORS middleware
	e.Use(func(hf echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			if c.Request().Method == "OPTIONS" {
				c.Response().Header().Set("Access-Control-Allow-Headers", strings.Join(append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...), ","))
				c.Response().Header().Set("Access-Control-Allow-Methods", "*")
				c.Response().Write([]byte(""))
				return nil
			} else {
				return hf(c)
			}
		})
	})

	// SuperTokens Middleware
	e.Use(func(hf echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				hf(c)
			})).ServeHTTP(c.Response(), c.Request())
			return nil
		})
	})

	e.GET("/sessioninfo", sessioninfo, verifySession(nil))

	e.Start(":3001")
}

func verifySession(options *sessmodels.VerifySessionOptions) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
				c.Set("session", session.GetSessionFromRequestContext(r.Context()))
				hf(c)
			})(c.Response(), c.Request())
			return nil
		}
	}
}

func sessioninfo(c echo.Context) error {
	sessionContainer := c.Get("session").(*sessmodels.SessionContainer)

	if sessionContainer == nil {
		return errors.New("no session found")
	}
	sessionData, err := sessionContainer.GetSessionData()
	if err != nil {
		err = supertokens.ErrorHandler(err, c.Request(), c.Response())
		if err != nil {
			return err
		}
		return nil
	}
	c.Response().WriteHeader(200)
	c.Response().Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
	if err != nil {
		c.Response().WriteHeader(500)
		c.Response().Write([]byte("error in converting to json"))
	} else {
		c.Response().Write(bytes)
	}
	return nil
}
