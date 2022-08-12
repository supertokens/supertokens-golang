package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init() {

}

type TypeEmailVerificationClaim struct {
	*claims.TypeBooleanClaim

	Validators *EmailVerificationClaimValidators
}

func NewEmailVerificationClaim(key string, fetchValue claims.FetchValueFunc) *TypeEmailVerificationClaim {
	booleanClaim := claims.BooleanClaim(key, fetchValue)

	emailVerificationClaim := &TypeEmailVerificationClaim{
		TypeBooleanClaim: booleanClaim,
	}
	emailVerificationClaim.Validators = &EmailVerificationClaimValidators{
		BooleanClaimValidators: booleanClaim.Validators,
		IsVerified: func(refetchTimeOnFalseInSeconds int64) *claims.SessionClaimValidator {
			if refetchTimeOnFalseInSeconds == 0 {
				refetchTimeOnFalseInSeconds = 10
			}

			id := "st-ev-isVerified"
			claimValidator := booleanClaim.Validators.HasValue(true, &id)
			claimValidator.ShouldRefetch = func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
				value := emailVerificationClaim.GetValueFromPayload(payload, userContext)
				return value == nil || value == false || emailVerificationClaim.GetLastRefetchTime(payload, userContext) < time.Now().Unix()-int64(refetchTimeOnFalseInSeconds)
			}
			return claimValidator
		},
	}
	return emailVerificationClaim
}

var EmailVerificationClaim *TypeEmailVerificationClaim

type EmailVerificationClaimValidators struct {
	*claims.BooleanClaimValidators
	IsVerified func(refetchTimeOnFalseInSeconds int64) *claims.SessionClaimValidator
}
