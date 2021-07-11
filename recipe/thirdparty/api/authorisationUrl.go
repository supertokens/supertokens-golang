package api

import "github.com/supertokens/supertokens-golang/recipe/thirdpaty/models"

func AuthorisationUrlAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.AuthorisationUrlGET == nil {
		return nil
	}
	return nil
}
