package models

type RecipeInterface struct {
	CreateEmailVerificationToken  func(userID, email string) (CreateEmailVerificationTokenResponse, error)
	VerifyEmailUsingToken         func(token string) (VerifyEmailUsingTokenResponse, error)
	IsEmailVerified               func(userID, email string) (bool, error)
	RevokeEmailVerificationTokens func(userId, email string) (RevokeEmailVerificationTokensResponse, error)
	UnverifyEmail                 func(userId, email string) (UnverifyEmailResponse, error)
}

type CreateEmailVerificationTokenResponse struct {
	OK *struct {
		Token string
	}
	EmailAlreadyVerifiedError *struct{}
}

type VerifyEmailUsingTokenResponse struct {
	OK *struct {
		User User
	}
	EmailVerificationInvalidTokenError *struct{}
}

type RevokeEmailVerificationTokensResponse struct {
	OK *struct{}
}

type UnverifyEmailResponse struct {
	OK *struct{}
}
