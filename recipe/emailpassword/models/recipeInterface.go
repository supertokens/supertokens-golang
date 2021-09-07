package models

type RecipeInterface struct {
	SignUp                   func(email string, password string) (SignUpResponse, error)
	SignIn                   func(email string, password string) (SignInResponse, error)
	GetUserByID              func(userID string) (*User, error)
	GetUserByEmail           func(email string) (*User, error)
	CreateResetPasswordToken func(userID string) (CreateResetPasswordTokenResponse, error)
	ResetPasswordUsingToken  func(token string, newPassword string) (ResetPasswordUsingTokenResponse, error)
	UpdateEmailOrPassword    func(userId string, email *string, password *string) (UpdateEmailOrPasswordResponse, error)
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

type CreateResetPasswordTokenResponse struct {
	OK *struct {
		Token string
	}
	UnknownUserIdError *struct{}
}

type ResetPasswordUsingTokenResponse struct {
	OK                             *struct{}
	ResetPasswordInvalidTokenError *struct{}
}

type UpdateEmailOrPasswordResponse struct {
	OK                      *struct{}
	UnknownUserIdError      *struct{}
	EmailAlreadyExistsError *struct{}
}
