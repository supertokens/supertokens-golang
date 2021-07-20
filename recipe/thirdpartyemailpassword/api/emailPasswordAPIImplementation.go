package api

import (
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func GetEmailPasswordIterfaceImpl(apiImplmentation models.APIImplementation) epm.APIImplementation {
	return epm.APIImplementation{
		EmailExistsGET:                 apiImplmentation.EmailExistsGET,
		GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
		PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
		SignInPOST:                     apiImplmentation.SignInPOST,
		SignUpPOST:                     apiImplmentation.SignUpPOST,
	}
}
