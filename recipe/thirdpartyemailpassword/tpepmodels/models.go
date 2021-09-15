package tpepmodels

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

type User struct {
	ID         string
	TimeJoined uint64
	Email      string
	ThirdParty *struct {
		ID     string
		UserID string
	}
}

type TypeContext struct {
	FormFields                 []epmodels.TypeFormField
	ThirdPartyAuthCodeResponse interface{}
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
}

type TypeInput struct {
	SignUpFeature                  *epmodels.TypeInputSignUp
	Providers                      []tpmodels.TypeProvider
	ResetPasswordUsingTokenFeature *epmodels.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *OverrideStruct
}

type TypeNormalisedInput struct {
	SignUpFeature                  *epmodels.TypeInputSignUp
	Providers                      []tpmodels.TypeProvider
	ResetPasswordUsingTokenFeature *epmodels.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       evmodels.TypeInput
	Override                       OverrideStruct
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evmodels.OverrideStruct
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}
