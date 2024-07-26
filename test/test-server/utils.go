package main

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

var testClaimSetups map[string]*claims.TypeSessionClaim
var testClaimValidators map[string]map[string]func(args ...any) claims.SessionClaimValidator

func init() {
	testClaimSetups = make(map[string]*claims.TypeSessionClaim)
	testClaimValidators = make(map[string]map[string]func(args ...any) claims.SessionClaimValidator)

	addBuiltinClaimValidators()
}

func addBuiltinClaimValidators() {
	addEmailVerificationClaimAndValidators()
}

func addEmailVerificationClaimAndValidators() {
	evClaim := evclaims.EmailVerificationClaim
	testClaimSetups[evClaim.Key] = evClaim
	testClaimValidators[evClaim.Key] = map[string]func(args ...any) claims.SessionClaimValidator{}

	testClaimValidators[evClaim.Key]["isVerified"] = func(args ...any) claims.SessionClaimValidator {
		if len(args) == 0 {
			return evclaims.EmailVerificationClaimValidators.IsVerified(nil, nil)
		} else if len(args) == 1 {
			arg1 := int64(args[0].(float64))
			return evclaims.EmailVerificationClaimValidators.IsVerified(&arg1, nil)
		} else {
			arg1 := int64(args[0].(float64))
			arg2 := int64(args[1].(float64))
			return evclaims.EmailVerificationClaimValidators.IsVerified(&arg1, &arg2)
		}
	}
}

func deserializeValidator(serializedValidator map[string]interface{}) func(args ...any) claims.SessionClaimValidator {
	key := serializedValidator["key"].(string)
	if validators, ok := testClaimValidators[key]; ok {
		validatorName := serializedValidator["validatorName"].(string)
		if validator, ok := validators[validatorName]; ok {
			return validator
		}
	}
	return nil
}
