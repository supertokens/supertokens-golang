package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	emailVerificationBackwardsCompatibilityService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(recipeInterfaceImpl epmodels.RecipeInterface, appInfo supertokens.NormalisedAppinfo, sendResetPasswordEmail func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext), sendEmailVerificationEmail func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext)) emaildelivery.EmailDeliveryInterface {

	emailVerificationService := emailVerificationBackwardsCompatibilityService.MakeBackwardCompatibilityService(appInfo, sendEmailVerificationEmail)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationService.SendEmail)(input, userContext)
		} else if input.PasswordReset != nil {
			user, err := (*recipeInterfaceImpl.GetUserByID)(input.PasswordReset.User.ID, userContext)
			if err != nil {
				return err
			}
			if user == nil {
				return errors.New("should never come here")
			}
			sendResetPasswordEmail(*user, input.PasswordReset.PasswordResetLink, userContext)
		} else {
			return errors.New("should never come here")
		}
		return nil
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
