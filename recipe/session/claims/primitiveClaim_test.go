package claims

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestPrimitiveClaimFetchAndSetClaim(t *testing.T) {
	val := map[string]interface{}{
		"a": 1,
	}
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return val, nil
		},
		nil,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)

	assert.Equal(t, val, payload["test"].(map[string]interface{})["v"])
}

func TestPrimitiveClaimAddToPayloadInternal(t *testing.T) {
	val := map[string]interface{}{
		"a": 1,
	}
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		nil,
	)

	payload := primClaim.AddToPayload_internal(map[string]interface{}{}, val, nil)
	assert.Equal(t, val, payload["test"].(map[string]interface{})["v"])
}

func TestPrimitiveClaimFetchValue(t *testing.T) {
	val := map[string]interface{}{
		"a": 1,
	}
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return val, nil
		},
		nil,
	)
	fval, err := primClaim.FetchValue("userId", nil)
	assert.NoError(t, err)
	assert.Equal(t, val, fval)
}

func TestPrimitiveClaimFetchValueReturningEmpty(t *testing.T) {
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return nil, nil
		},
		nil,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{}, payload)
}

func TestPrimitiveClaimGetValueFromPayloadEmptyPayload(t *testing.T) {
	val := map[string]interface{}{
		"a": 1,
	}
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return val, nil
		},
		nil,
	)

	assert.Equal(t, nil, primClaim.GetValueFromPayload(map[string]interface{}{}, nil))
}

func TestPrimitiveClaimGetValueFromPayloadUsingBuild(t *testing.T) {
	val := map[string]interface{}{
		"a": 1,
	}
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return val, nil
		},
		nil,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, val, primClaim.GetValueFromPayload(payload, nil))
}

func TestPrimitiveClaimGetValueFromPayloadUsingAddToPayloadInternal(t *testing.T) {
	val := map[string]interface{}{
		"a": 1,
	}
	primClaim, _ := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		nil,
	)

	payload := primClaim.AddToPayload_internal(map[string]interface{}{}, val, nil)
	assert.Equal(t, val, primClaim.GetValueFromPayload(payload, nil))
}

func TestPrimitiveClaimValidateWithEmptyPayload(t *testing.T) {
	_, validator := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		nil,
	)

	validationResult := validator.HasValue("hello", nil, nil).Validate(map[string]interface{}{}, nil)
	assert.Equal(t, false, validationResult.IsValid)
	assert.Equal(t, map[string]interface{}{
		"actualValue":   nil,
		"expectedValue": "hello",
		"message":       "value does not exist",
	}, validationResult.Reason)
}

func TestPrimitiveClaimValidateWithMismatchingPayload(t *testing.T) {
	primClaim, validator := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		nil,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)

	validationResult := validator.HasValue("world", nil, nil).Validate(payload, nil)
	assert.Equal(t, false, validationResult.IsValid)
	assert.Equal(t, map[string]interface{}{
		"actualValue":   "hello",
		"expectedValue": "world",
		"message":       "wrong value",
	}, validationResult.Reason)
}

func TestPrimitiveClaimValidateWithMatchingPayload(t *testing.T) {
	primClaim, validator := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		nil,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)

	validationResult := validator.HasValue("hello", nil, nil).Validate(payload, nil)
	assert.Equal(t, true, validationResult.IsValid)
	assert.Equal(t, nil, validationResult.Reason)
}

func TestPrimitiveClaimValidateExpiry(t *testing.T) {
	primClaim, validator := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		nil,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)

	maxAgeInSec := int64(1)

	time.Sleep(2 * time.Second)

	validationResult := validator.HasValue("hello", &maxAgeInSec, nil).Validate(payload, nil)
	assert.Equal(t, false, validationResult.IsValid)
	assert.Equal(t, map[string]interface{}{
		"ageInSeconds":    int64(2),
		"maxAgeInSeconds": int64(1),
		"message":         "expired",
	}, validationResult.Reason)
}

func TestPrimitiveClaimValidateDefaultAgeExpiry(t *testing.T) {
	maxAgeInSec := int64(1)

	primClaim, validator := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		&maxAgeInSec,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	validationResult := validator.HasValue("hello", nil, nil).Validate(payload, nil)
	assert.Equal(t, false, validationResult.IsValid)
	assert.Equal(t, map[string]interface{}{
		"ageInSeconds":    int64(2),
		"maxAgeInSeconds": int64(1),
		"message":         "expired",
	}, validationResult.Reason)
}

func TestPrimitiveClaimValidateMaxAgeOverride(t *testing.T) {
	maxAgeInSec := int64(1)

	primClaim, validator := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return "hello", nil
		},
		&maxAgeInSec,
	)

	payload, err := primClaim.Build("userId", nil, nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	{
		validationResult := validator.HasValue("hello", nil, nil).Validate(payload, nil)
		assert.Equal(t, false, validationResult.IsValid)
		assert.Equal(t, map[string]interface{}{
			"ageInSeconds":    int64(2),
			"maxAgeInSeconds": int64(1),
			"message":         "expired",
		}, validationResult.Reason)
	}

	{
		maxAgeInSec = 10
		validationResult := validator.HasValue("hello", &maxAgeInSec, nil).Validate(payload, nil)
		assert.Equal(t, true, validationResult.IsValid)
		assert.Equal(t, nil, validationResult.Reason)
	}
}
