package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type UserInfo struct {
	ID    string
	Email *EmailStruct
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

type TypeProviderGetResponse struct {
	AccessTokenAPI        AccessTokenAPI
	AuthorisationRedirect AuthorisationRedirect
	GetProfileInfo        func(authCodeResponse interface{}) (UserInfo, error)
}

type AccessTokenAPI struct {
	URL    string
	Params map[string]string
}

type AuthorisationRedirect struct {
	URL    string
	Params map[string]interface{}
}

type TypeProvider struct {
	ID  string
	Get func(redirectURI *string, authCodeFromRequest *string) TypeProviderGetResponse
}

type User struct {
	ID         string
	TimeJoined uint64
	Email      string
	ThirdParty struct {
		ID     string
		UserID string
	}
}

type TypeInputSetJwtPayloadForSession func(User User, thirdPartyAuthCodeResponse interface{}, action string) map[string]interface{}

type TypeInputSetSessionDataForSession func(User User, thirdPartyAuthCodeResponse interface{}, action string) map[string]interface{}

type TypeNormalisedInputSessionFeature struct {
	SetJwtPayload  TypeInputSetJwtPayloadForSession
	SetSessionData TypeInputSetSessionDataForSession
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
}

type TypeInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeNormalisedInputSignInAndUp struct {
	Providers []TypeProvider
}

type TypeInput struct {
	SessionFeature           *TypeNormalisedInputSessionFeature
	SignInAndUpFeature       TypeInputSignInAndUp
	EmailVerificationFeature *TypeInputEmailVerificationFeature
	Override                 *struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
		}
	}
}

type TypeNormalisedInput struct {
	SessionFeature           TypeNormalisedInputSessionFeature
	SignInAndUpFeature       TypeNormalisedInputSignInAndUp
	EmailVerificationFeature models.TypeInput
	Override                 struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			APIs      func(originalImplementation models.APIImplementation) models.APIImplementation
		}
	}
}

type RecipeImplementation struct {
	GetUserByID             func(userID string) *User
	GetUserByThirdPartyInfo func(thirdPartyID string, thirdPartyUserID string) *User
	SignInUp                func(thirdPartyID string, thirdPartyUserID string, email EmailStruct) SignInUpResponse
}

type SignInUpResponse struct {
	Status         string `json:"status"`
	CreatedNewUser bool   `json:"createdNewUser"`
	User           User   `json:"user"`
	Error          error  `json:"error"`
}

type APIOptions struct {
	RecipeImplementation RecipeImplementation
	Config               TypeNormalisedInput
	RecipeID             string
	Providers            []TypeProvider
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type APIImplementation struct {
	AuthorisationUrlGET func(provider TypeProvider, options APIOptions) AuthorisationUrlGETResponse
	SignInUpPOST        func(provider TypeProvider, code string, redirectURI string, options APIOptions) SignInUpPOSTResponse
}

type AuthorisationUrlGETResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

type SignInUpPOSTResponse struct {
	Status           string      `json:"status"`
	CreatedNewUser   bool        `json:"createdNewUser"`
	User             User        `json:"user"`
	AuthCodeResponse interface{} `json:"authCodeResponse"`
	Error            error       `json:"error"`
}
