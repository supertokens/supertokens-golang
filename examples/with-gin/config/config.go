package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init() {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("dev")
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}

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

	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.SupertokenTypeInput{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     "http://localhost" + config.GetString("server.apiPort"),
			WebsiteDomain: "http://localhost" + config.GetString("server.websitePort"),
		},
		RecipeList: []supertokens.RecipeListFunction{
			emailpassword.Init(&models.TypeInput{
				Override: &models.OverrideStruct{
					APIs: func(originalImplementation models.APIInterface) models.APIInterface {
						return models.APIInterface{
							EmailExistsGET:                 nil,
							GeneratePasswordResetTokenPOST: originalImplementation.GeneratePasswordResetTokenPOST,
							PasswordResetPOST:              originalImplementation.PasswordResetPOST,
							SignInPOST:                     originalImplementation.SignInPOST,
							SignUpPOST:                     originalImplementation.SignUpPOST,
						}
					},
				},
			}),
			session.Init(nil),
			// thirdparty.Init(thirdpartyConfig),
		},
	})
	if err != nil {
		panic(err.Error())
	}
}

func GetConfig() *viper.Viper {
	return config
}
