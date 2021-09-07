package server

import (
	"log"
	"net/http"

	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init() {
	config := config.GetConfig()

	thirdpartyemailpasswordConfig := &models.TypeInput{
		Providers: []tpm.TypeProvider{thirdparty.Github(providers.TypeThirdPartyProviderGithubConfig{
			ClientID:     config.GetString("GITHUB_CLIENT_ID"),
			ClientSecret: config.GetString("GITHUB_CLIENT_SECRET"),
		}),
		},
	}

	// thirdpartyConfig := &tpm.TypeInput{
	// 	SignInAndUpFeature: tpm.TypeInputSignInAndUp{
	// 		Providers: []tpm.TypeProvider{thirdparty.Github(providers.TypeThirdPartyProviderGithubConfig{
	// 			ClientID:     config.GetString("GITHUB_CLIENT_ID"),
	// 			ClientSecret: config.GetString("GITHUB_CLIENT_SECRET"),
	// 		}),
	// 		},
	// 	},
	// }

	err := supertokens.SupertokensInit(supertokens.TypeInput{
		Supertokens: &supertokens.SupertokenTypeInput{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     "http://localhost" + config.GetString("server.apiPort"),
			WebsiteDomain: "http://localhost" + config.GetString("server.websitePort"),
		},
		RecipeList: []supertokens.RecipeListFunction{
			// emailpassword.RecipeInit(nil),
			session.RecipeInit(nil),
			thirdpartyemailpassword.RecipeInit(thirdpartyemailpasswordConfig),
			// thirdparty.RecipeInit(thirdpartyConfig),
		},
		OnGeneralError: func(err error, req *http.Request, res http.ResponseWriter) {
			res.WriteHeader(500)
			res.Write([]byte("Internal error: " + err.Error()))
		},
	})
	if err != nil {
		panic(err.Error())
	}

	r := newRouter()
	err = r.Run(config.GetString("server.apiPort"))
	if err != nil {
		log.Println("error running server => ", err)
	}
}

// http://localhost:3000/auth/callback/github
