package emailverification

type APIImplementation struct {
	verifyEmailPOST              func(token string, options APIOptions) map[string]interface{}
	isEmailVerifiedGET           func(options APIOptions) map[string]interface{}
	generateEmailVerifyTokenPOST func(options APIOptions) map[string]interface{}
}
