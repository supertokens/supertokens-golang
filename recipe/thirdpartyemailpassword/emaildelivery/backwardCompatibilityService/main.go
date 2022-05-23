package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	emailPasswordBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl tpepmodels.RecipeInterface, emailPasswordRecipeInterfaceImpl epmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext), sendResetPasswordEmail func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext)) emaildelivery.EmailDeliveryInterface {
	// We are using evmodels.User as opposed to tpmodels.User because TypeInput of thirdparty accepts evmodels.TypeInput for EmailVerificationFeature
	// Similarly with epmodels.User as well
	emailPasswordService := emailPasswordBackwardsCompatibilityService.MakeBackwardCompatibilityService(emailPasswordRecipeInterfaceImpl, appInfo, sendResetPasswordEmail, sendEmailVerificationEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil || input.PasswordReset != nil {
			return (*emailPasswordService.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
