package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	emailVerificationBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl tpmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext)) emaildelivery.EmailDeliveryInterface {
	// We are using evmodels.User as opposed to tpmodels.User because TypeInput of thirdparty accepts evmodels.TypeInput for EmailVerificationFeature
	emailVerificationService := emailVerificationBackwardsCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendEmailVerificationEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationService.SendEmail)(input, userContext)
		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
