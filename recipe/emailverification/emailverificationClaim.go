package emailverification

import (
	"errors"
	"time"

	evclaims "github.com/supertokens/supertokens-golang/recipe/emailverification/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// key string, fetchValue claims.FetchValueFunc
func NewEmailVerificationClaim() *evclaims.TypeEmailVerificationClaim {
	fetchValue := func(userId string, userContext supertokens.UserContext) (interface{}, error) {
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
		} else {
			return false, errors.New("should never come here: UnknownUserIdError or invalid result from getEmailForUserId")
		}
	}

	booleanClaim := claims.BooleanClaim("st-ev", fetchValue, nil)

	emailVerificationClaim := &evclaims.TypeEmailVerificationClaim{
		TypeBooleanClaim: booleanClaim,
	}
	emailVerificationClaim.Validators = &evclaims.EmailVerificationClaimValidators{
		BooleanClaimValidators: booleanClaim.Validators,
		IsVerified: func(refetchTimeOnFalseInSeconds *int64) *claims.SessionClaimValidator {
			if refetchTimeOnFalseInSeconds == nil {
				var defaultTimeout int64 = 10
				refetchTimeOnFalseInSeconds = &defaultTimeout
			}

			claimValidator := booleanClaim.Validators.HasValue(true, refetchTimeOnFalseInSeconds, nil)
			claimValidator.ShouldRefetch = func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
				value := emailVerificationClaim.GetValueFromPayload(payload, userContext)
				return value == nil || (value == false && *emailVerificationClaim.GetLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*refetchTimeOnFalseInSeconds*1000)
			}
			return claimValidator
		},
	}
	return emailVerificationClaim
}

func init() {
	evclaims.EmailVerificationClaim = NewEmailVerificationClaim()
}
