package api

import "github.com/supertokens/supertokens-golang/recipe/session/models"

func SignOutAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
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
