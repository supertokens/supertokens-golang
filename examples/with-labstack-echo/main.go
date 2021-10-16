package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
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
			emailpassword.Init(nil),
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
		"sessionHandle": sessionContainer.GetHandle(),
		"userId":        sessionContainer.GetUserID(),
		"jwtPayload":    sessionContainer.GetJWTPayload(),
		"sessionData":   sessionData,
	})
	if err != nil {
		c.Response().WriteHeader(500)
		c.Response().Write([]byte("error in converting to json"))
	} else {
		c.Response().Write(bytes)
	}
	return nil
}
