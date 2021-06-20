package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
)

type UserInfo struct {
	ID    string
	Email *struct {
		ID         string
		IsVerified bool
	}
}

type TypeProviderGetResponse struct {
	AccessTokenAPI struct {
		URL    string
		Params map[string]string
	}
	AuthorisationRedirect struct {
		URL    string
		Params map[string]string
	}
	GetProfileInfo func(authCodeResponse interface{}) UserInfo
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

type TypeInputSetJwtPayloadForSession func(User User, thirdPartyAuthCodeResponse interface{}, action string) map[string]string

type TypeInputSetSessionDataForSession func(User User, thirdPartyAuthCodeResponse interface{}, action string) map[string]string

type TypeInputSessionFeature struct {
	SetJwtPayload  TypeInputSetJwtPayloadForSession
	SetSessionData TypeInputSetJwtPayloadForSession
}

type TypeNormalisedInputSessionFeature struct {
	setJwtPayload  TypeInputSetJwtPayloadForSession
	setSessionData TypeInputSetSessionDataForSession
}

type TypeInputEmailVerificationFeature struct {
	getEmailVerificationURL  func(user User) string
	createAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
}

type TypeInputSignInAndUp struct {
	providers []TypeProvider
}

type TypeNormalisedInputSignInAndUp struct {
	providers []TypeProvider
}

type TypeInput struct {
	SessionFeature           *TypeInputSessionFeature
	SignInAndUpFeature       TypeInputSignInAndUp
	EmailVerificationFeature *TypeInputEmailVerificationFeature
	Override                 *struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		Apis                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature *struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			aApis     func(originalImplementation models.APIImplementation) models.APIImplementation
		}
	}
}

type TypeNormalisedInput struct {
	SessionFeature           TypeInputSessionFeature
	SignInAndUpFeature       TypeInputSignInAndUp
	EmailVerificationFeature TypeInputEmailVerificationFeature
	Override                 struct {
		Functions                func(originalImplementation RecipeImplementation) RecipeImplementation
		Apis                     func(originalImplementation APIImplementation) APIImplementation
		EmailVerificationFeature struct {
			Functions func(originalImplementation models.RecipeImplementation) models.RecipeImplementation
			aApis     func(originalImplementation models.APIImplementation) models.APIImplementation
		}
	}
}

type RecipeImplementation struct {
	GetUserById             func(userId string) *User
	GetUserByThirdPartyInfo func(thirdPartyId string, thirdPartyUserId string) *User
	GetUsersOldestFirst     func(limit *int, nextPaginationToken *string) struct {
		users               []User
		nextPaginationToken *string
	}
	GetUsersNewestFirst func(limit *int, nextPaginationToken *string) struct {
		users               []User
		nextPaginationToken *string
	}
	GetUserCount func() int
	SignInUp     func(
		thirdPartyId string,
		thirdPartyUserId string,
		email struct {
			ID         string
			IsVerified bool
		},
	) struct {
		CreatedNewUser bool
		User           User
	}
}

type APIOptions struct {
	RecipeImplementation RecipeImplementation
	Config               TypeNormalisedInput
	RecipeId             string
	Providers            []TypeProvider
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type APIImplementation struct {
	AuthorisationUrlGET func(provider TypeProvider, options APIOptions) error
	SignInUpPOST        func(provider TypeProvider, code string, redirectURI string, options APIOptions) error
}
