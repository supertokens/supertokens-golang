package supertokens

// This is there so that the supertokens package can use the email verification
// recipe since account linking requires it. We cannot directly use the ev recipe from here
// cause of cyclic dependencies.

var InternalUseEmailVerificationRecipeProxyInstance *InternalUseEmailVerificationRecipeProxy = nil

type InternalUseEmailVerificationRecipeProxy struct {
	CreateEmailVerificationToken  func(recipeUserID RecipeUserID, email string, tenantId string, userContext UserContext) (InternalUseCreateEmailVerificationTokenResponse, error)
	VerifyEmailUsingToken         func(token string, tenantId string, attemptAccountLinking bool, userContext UserContext) (InternalUseVerifyEmailUsingTokenResponse, error)
	IsEmailVerified               func(userID, email string, userContext UserContext) (bool, error)
	RevokeEmailVerificationTokens func(userId, email string, tenantId string, userContext UserContext) (InternalUseRevokeEmailVerificationTokensResponse, error)
	UnverifyEmail                 func(userId, email string, userContext UserContext) (InternalUseUnverifyEmailResponse, error)
}

type InternalUseCreateEmailVerificationTokenResponse struct {
	OK *struct {
		Token string
	}
	EmailAlreadyVerifiedError *struct{}
}

type InternalUseVerifyEmailUsingTokenResponse struct {
	OK *struct {
		User InternalUseEmailVerificationUser
	}
	EmailVerificationInvalidTokenError *struct{}
}

type InternalUseRevokeEmailVerificationTokensResponse struct {
	OK *struct{}
}

type InternalUseUnverifyEmailResponse struct {
	OK *struct{}
}

type InternalUseEmailVerificationUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
