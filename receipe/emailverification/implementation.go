package emailverification

type APIImplementation struct{}

func (a *APIImplementation) verifyEmailPOST(token string, options APIOptions) map[string]interface{} {
	return options.recipeImplementation.verifyEmailUsingToken(token)
}

func (a *APIImplementation) isEmailVerifiedGET(options APIOptions) map[string]interface{} {
	return nil
}

func (a *APIImplementation) generateEmailVerifyTokenPOST(options APIOptions) map[string]interface{} {
	return nil
}
