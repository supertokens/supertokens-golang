package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	emailPasswordBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	emailVerificationBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl tpepmodels.RecipeInterface, emailPasswordRecipeInterfaceImpl epmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext), sendResetPasswordEmail func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext)) emaildelivery.EmailDeliveryInterface {
	if sendEmailVerificationEmail != nil {
		inputSendEmailVerificationEmail := sendEmailVerificationEmail
		sendEmailVerificationEmail = func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
			userInfo, err := (*recipeInterfaceImpl.GetUserByID)(user.ID, userContext)
			if err != nil {
				return // FIXME: No error handling here
			}
			inputSendEmailVerificationEmail(evmodels.User{ID: userInfo.ID, Email: userInfo.Email}, emailVerificationURLWithToken, userContext)
		}
	}
	emailVerificationService := emailVerificationBackwardsCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendEmailVerificationEmail)
	emailPasswordService := emailPasswordBackwardsCompatibilityService.MakeBackwardCompatibilityService(emailPasswordRecipeInterfaceImpl, appInfo, sendResetPasswordEmail, sendEmailVerificationEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationService.SendEmail)(input, userContext)

		} else if input.PasswordReset != nil {
			return (*emailPasswordService.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}