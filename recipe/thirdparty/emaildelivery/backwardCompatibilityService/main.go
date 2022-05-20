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
