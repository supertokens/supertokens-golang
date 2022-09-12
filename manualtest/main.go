package main

import (
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:3567",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "Thirdparty Demo",
			WebsiteDomain: "localhost:3000",
			APIDomain:     "localhost:8000",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
			emailpassword.Init(nil),
			dashboard.Init(dashboardmodels.TypeInput{
				ApiKey: "abcd",
			}),
		},
	})
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(
		":8000",
		corsMiddleware(
			supertokens.Middleware(
				http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
					},
				),
			),
		),
	)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			// we add content-type + other headers used by SuperTokens
			response.Header().Set("Access-Control-Allow-Headers",
				strings.Join(append([]string{"Content-Type"},
					supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}
