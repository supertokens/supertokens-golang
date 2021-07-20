package api

import (
	tpm "github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
)

func GetThirdPartyIterfaceImpl(apiImplmentation models.APIImplementation) tpm.APIImplementation {
	return tpm.APIImplementation{
		AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,
		SignInUpPOST:        apiImplmentation.SignInUpPOST,
	}
}
