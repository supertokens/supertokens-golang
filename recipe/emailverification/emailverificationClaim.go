package emailverification

import (
	"errors"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// key string, fetchValue claims.FetchValueFunc
func NewEmailVerificationClaim() (claims.TypeSessionClaim, evclaims.TypeEmailVerificationClaimValidators) {
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

	evClaim, booleanClaimValidators := claims.BooleanClaim("st-ev", fetchValue, nil)

	getValueFromPayload := func(payload map[string]interface{}, userContext supertokens.UserContext) interface{} {
		if value, ok := evClaim.GetValueFromPayload(payload, userContext).(map[string]interface{}); ok {
			return value["v"]
		}
		return nil
	}

	getLastRefetchTime := func(payload map[string]interface{}, userContext supertokens.UserContext) *int64 {
		if value, ok := evClaim.GetValueFromPayload(payload, userContext).(map[string]interface{}); ok {
			val := value["t"].(int64)
			return &val
		}
		return nil
	}

	validators := evclaims.TypeEmailVerificationClaimValidators{
		BooleanClaimValidators: booleanClaimValidators,
		IsVerified: func(refetchTimeOnFalseInSeconds *int64) claims.SessionClaimValidator {
			if refetchTimeOnFalseInSeconds == nil {
				var defaultTimeout int64 = 10
				refetchTimeOnFalseInSeconds = &defaultTimeout
			}

			claimValidator := booleanClaimValidators.HasValue(true, nil, nil)
			claimValidator.ShouldRefetch = func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
				value := getValueFromPayload(payload, userContext)
				return value == nil || (value == false && *getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*refetchTimeOnFalseInSeconds*1000)
			}
			return claimValidator
		},
	}
	return evClaim, validators
}

func init() {
	// this function is called automatically when the package is imported
	evclaims.EmailVerificationClaim, evclaims.EmailVerificationClaimValidators = NewEmailVerificationClaim()
}
