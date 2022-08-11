package claims

type BooleanClaim struct {
	PrimitiveClaim
}

func (claim *BooleanClaim) GetValidators() BooleanClaimValidators {
	primitiveClaimValidators := claim.PrimitiveClaim.GetValidators()
	return BooleanClaimValidators{
		PrimitiveClaimValidators: primitiveClaimValidators,
		IsTrue: func(maxAgeInSeconds *int64) SessionClaimValidator {
			if maxAgeInSeconds != nil {
				return primitiveClaimValidators.HasFreshValue(true, *maxAgeInSeconds, nil)
			}
			return primitiveClaimValidators.HasValue(true, nil)
		},
		IsFalse: func(maxAgeInSeconds *int64) SessionClaimValidator {
			if maxAgeInSeconds != nil {
				return primitiveClaimValidators.HasFreshValue(false, *maxAgeInSeconds, nil)
			}
			return primitiveClaimValidators.HasValue(true, nil)
		},
	}
}

type BooleanClaimValidators struct {
	PrimitiveClaimValidators
	IsTrue  func(maxAgeInSeconds *int64) SessionClaimValidator
	IsFalse func(maxAgeInSeconds *int64) SessionClaimValidator
}
