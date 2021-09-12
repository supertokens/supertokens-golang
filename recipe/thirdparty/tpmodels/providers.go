package tpmodels

type GoogleConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
}

type GithubConfig struct {
	ClientID              string
	ClientSecret          string
	Scope                 []string
	AuthorisationRedirect *struct {
		Params map[string]interface{}
	}
}

type FacebookConfig struct {
	ClientID     string
	ClientSecret string
	Scope        []string
}

// type AppleConfig struct {
// 	ClientID              string
// 	ClientSecret          AppleClientSecret
// 	Scope                 []string
// 	AuthorisationRedirect *struct {
// 		Params map[string]interface{}
// 	}
// }

// type AppleClientSecret struct {
// 	KeyId      string
// 	PrivateKey string
// 	TeamId     string
// }
