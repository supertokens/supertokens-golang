package api

import "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"

// TODO:
func AuthorisationUrlAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.AuthorisationUrlGET == nil {
		return nil
	}
	return nil
}
