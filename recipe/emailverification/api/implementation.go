package api

import "github.com/supertokens/supertokens-golang/recipe/emailverification/schema"

type APIImplementation struct{}

func (a *APIImplementation) VerifyEmailPOST(token string, options schema.APIOptions) map[string]interface{} {
	return options.RecipeImplementation.VerifyEmailUsingToken(token)
}

func (a *APIImplementation) IsEmailVerifiedGET(options schema.APIOptions) map[string]interface{} {
	// todo
	return nil
}

func (a *APIImplementation) GenerateEmailVerifyTokenPOST(options schema.APIOptions) map[string]interface{} {
	// todo
	return nil
}
