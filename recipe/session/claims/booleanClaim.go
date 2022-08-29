package claims

func BooleanClaim(key string, fetchValue FetchValueFunc, defaultMaxAgeInSeconds *int64) (*TypeSessionClaim, *BooleanClaimValidators) {
	claim, primitiveClaimValidators := PrimitiveClaim(key, fetchValue, defaultMaxAgeInSeconds)

	validators := &BooleanClaimValidators{
		PrimitiveClaimValidators: primitiveClaimValidators,

		IsTrue: func(maxAgeInSeconds *int64, id *string) *SessionClaimValidator {
			return primitiveClaimValidators.HasValue(true, maxAgeInSeconds, id)
		},

		IsFalse: func(maxAgeInSeconds *int64, id *string) *SessionClaimValidator {
			return primitiveClaimValidators.HasValue(false, maxAgeInSeconds, id)
		},
	}

	return claim, validators
}

type BooleanClaimValidators struct {
	*PrimitiveClaimValidators
	IsTrue  func(maxAgeInSeconds *int64, id *string) *SessionClaimValidator
	IsFalse func(maxAgeInSeconds *int64, id *string) *SessionClaimValidator
}
