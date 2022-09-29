package emailverification

import (
	"errors"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// key string, fetchValue claims.FetchValueFunc
func NewEmailVerificationClaim() (*claims.TypeSessionClaim, evclaims.TypeEmailVerificationClaimValidators) {
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
			return false, errors.New("UNKNOWN_USER_ID")
		}
	}

	evClaim, booleanClaimValidators := claims.BooleanClaim("st-ev", fetchValue, nil)

	getLastRefetchTime := func(payload map[string]interface{}, userContext supertokens.UserContext) *int64 {
		if value, ok := payload[evClaim.Key].(map[string]interface{}); ok {
			switch t := value["t"].(type) {
			case int64:
				return &t
			case float64:
				it := int64(t)
				return &it
			}
		}
		return nil
	}

	validators := evclaims.TypeEmailVerificationClaimValidators{
		BooleanClaimValidators: booleanClaimValidators,
		IsVerified: func(refetchTimeOnFalseInSeconds *int64, maxAgeInSeconds *int64) claims.SessionClaimValidator {
			if refetchTimeOnFalseInSeconds == nil {
				var defaultTimeout int64 = 10
				refetchTimeOnFalseInSeconds = &defaultTimeout
			}
			if maxAgeInSeconds == nil {
				var defaultTimeout int64 = 300
				maxAgeInSeconds = &defaultTimeout
			}

			claimValidator := booleanClaimValidators.HasValue(true, nil, nil)
			claimValidator.ShouldRefetch = func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
				value := evClaim.GetValueFromPayload(payload, userContext)
				return value == nil || (*getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*maxAgeInSeconds*1000) || (value == false && *getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*refetchTimeOnFalseInSeconds*1000)
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
