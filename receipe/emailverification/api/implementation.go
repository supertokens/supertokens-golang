package api

import "github.com/supertokens/supertokens-golang/receipe/emailverification/schema"

type APIImplementation schema.APIImplementation

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
