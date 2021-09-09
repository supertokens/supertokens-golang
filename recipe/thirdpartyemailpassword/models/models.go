package models

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	evm "github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
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
	FormFields                 []epm.TypeFormField
	ThirdPartyAuthCodeResponse interface{}
}

type TypeInputSignUp struct {
	FormFields []epm.TypeInputFormField
}

type TypeNormalisedInputSignUp struct {
	FormFields []epm.NormalisedFormField
}

type TypeInputEmailVerificationFeature struct {
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string) error
}

type TypeInput struct {
	SignUpFeature                  *TypeInputSignUp
	Providers                      []tpm.TypeProvider
	ResetPasswordUsingTokenFeature *epm.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       *TypeInputEmailVerificationFeature
	Override                       *OverrideStruct
}

type TypeNormalisedInput struct {
	SignUpFeature                  TypeNormalisedInputSignUp
	Providers                      []tpm.TypeProvider
	ResetPasswordUsingTokenFeature *epm.TypeInputResetPasswordUsingTokenFeature
	EmailVerificationFeature       evm.TypeInput
	Override                       OverrideStruct
}

type OverrideStruct struct {
	Functions                func(originalImplementation RecipeInterface) RecipeInterface
	APIs                     func(originalImplementation APIInterface) APIInterface
	EmailVerificationFeature *evm.OverrideStruct
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}
