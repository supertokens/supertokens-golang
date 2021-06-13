package api

import "github.com/supertokens/supertokens-golang/recipe/session/schema"

func SignOutAPI(apiImplementation schema.APIImplementation, options schema.APIOptions) error {
	if apiImplementation.SignOutPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}
	_, err := apiImplementation.SignOutPOST(options)
	if err != nil {
		return err
	}
	// TODO: response 200
	return nil
}
