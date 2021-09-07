package server

import (
	"log"

	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init() {
	config := config.GetConfig()

	// thirdpartyemailpasswordConfig := &models.TypeInput{
	// 	Providers: []tpm.TypeProvider{thirdparty.Github(providers.TypeThirdPartyProviderGithubConfig{
	// 		ClientID:     config.GetString("GITHUB_CLIENT_ID"),
	// 		ClientSecret: config.GetString("GITHUB_CLIENT_SECRET"),
	// 	}),
	// 	},
	// }

	// thirdpartyConfig := &tpm.TypeInput{
	// 	SignInAndUpFeature: tpm.TypeInputSignInAndUp{
	// 		Providers: []tpm.TypeProvider{thirdparty.Github(providers.TypeThirdPartyProviderGithubConfig{
	// 			ClientID:     config.GetString("GITHUB_CLIENT_ID"),
	// 			ClientSecret: config.GetString("GITHUB_CLIENT_SECRET"),
	// 		}),
	// 		},
	// 	},
	// }

	// TODO: general error handling?
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.SupertokenTypeInput{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     "http://localhost" + config.GetString("server.apiPort"),
			WebsiteDomain: "http://localhost" + config.GetString("server.websitePort"),
		},
		RecipeList: []supertokens.RecipeListFunction{
			emailpassword.EmailPasswordInit(nil),
			session.SessionInit(nil),
			// thirdparty.RecipeInit(thirdpartyConfig),
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
