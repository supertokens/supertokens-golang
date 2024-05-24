package config

import (
	"log"

	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"

	"github.com/spf13/viper"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
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

	var providers []tpmodels.ProviderInput
	err = config.UnmarshalKey("Providers", &providers)
	if err != nil {
		log.Fatal("invalid 'Providers' config, ", err)
	}

	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: config.GetString("SuperTokens.ConnectionURI"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:       config.GetString("AppInfo.AppName"),
			APIDomain:     config.GetString("AppInfo.APIDomain"),
			WebsiteDomain: config.GetString("AppInfo.WebsiteDomain"),
		},
		RecipeList: []supertokens.Recipe{
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired,
			}),
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: providers,
				},
			}),
			emailpassword.Init(nil),
			session.Init(nil),
			dashboard.Init(nil),
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
