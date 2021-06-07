package api

import "github.com/supertokens/supertokens-golang/recipe/emailverification/schema"

func GenerateEmailVerifyToken(apiImplementation schema.APIInterface, options schema.APIOptions) {
	var result map[string]string

	response := apiImplementation.GenerateEmailVerifyTokenPOST(options)
	for k, v := range response {
		result[k] = v.(string)
	}
	// todo: send200Response(options.res, result);
}
