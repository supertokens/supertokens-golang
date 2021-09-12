package tpepmodels

import "github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"

type RecipeInterface struct {
	GetUserByID              func(userID string) (*User, error)
	GetUsersByEmail          func(email string) ([]User, error)
	GetUserByThirdPartyInfo  func(thirdPartyID string, thirdPartyUserID string) (*User, error)
	SignInUp                 func(thirdPartyID string, thirdPartyUserID string, email EmailStruct) (SignInUpResponse, error)
	SignUp                   func(email string, password string) (SignUpResponse, error)
	SignIn                   func(email string, password string) (SignInResponse, error)
	CreateResetPasswordToken func(userID string) (epmodels.CreateResetPasswordTokenResponse, error)
	ResetPasswordUsingToken  func(token string, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error)
	UpdateEmailOrPassword    func(userId string, email *string, password *string) (epmodels.UpdateEmailOrPasswordResponse, error)
}

type SignInUpResponse struct {
	OK *struct {
		CreatedNewUser bool
		User           User
	}
	FieldError *struct{ Error string }
}

type SignUpResponse struct {
	OK *struct {
		User User
	}
	EmailAlreadyExistsError *struct{}
}

type SignInResponse struct {
	OK *struct {
		User User
	}
	WrongCredentialsError *struct{}
}
