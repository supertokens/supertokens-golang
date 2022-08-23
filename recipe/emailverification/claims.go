package emailverification

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func initClaims() {
	claims.EmailVerificationClaim = claims.NewEmailVerificationClaim(
		"st-ev",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			instance, err := getRecipeInstanceOrThrowError()
			if err != nil {
				return nil, err
			}
			emailInfo, err := instance.GetEmailForUserID(userId, userContext)
			if err != nil {
				return false, err
			}
			if emailInfo.OK != nil {
				verified, err := (*instance.RecipeImpl.IsEmailVerified)(userId, emailInfo.OK.Email, userContext)
				if err != nil {
					return false, nil
				}
				return verified, nil
			} else if emailInfo.EmailDoesNotExistError != nil {
				// We consider people without email addresses as validated
				return true, nil
			} else if emailInfo.UnknownUserIDError != nil {
				return false, errors.New("should never come here: UnknownUserIdError or invalid result from getEmailForUserId")
			} else {
				return false, errors.New("should never come here")
			}
		},
	)
}
