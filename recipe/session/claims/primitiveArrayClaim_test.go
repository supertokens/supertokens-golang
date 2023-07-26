package claims

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestPrimitiveArrayClaimValidators(t *testing.T) {
	primArrayClaim, validators := PrimitiveArrayClaim(
		"test",
		func(userId string, tenantId string, userContext supertokens.UserContext) (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		nil,
	)
	payload := map[string]interface{}{}
	payload = primArrayClaim.AddToPayload_internal(payload, []interface{}{
		100,
		"world",
		true,
	}, nil)

	assert.True(t, validators.Includes(100, nil, nil).Validate(payload, nil).IsValid)
	assert.True(t, validators.Includes("world", nil, nil).Validate(payload, nil).IsValid)
	assert.True(t, validators.Includes(true, nil, nil).Validate(payload, nil).IsValid)

	assert.False(t, validators.Includes(101, nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.Includes("hello", nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.Includes(false, nil, nil).Validate(payload, nil).IsValid)

	assert.True(t, validators.Excludes(101, nil, nil).Validate(payload, nil).IsValid)
	assert.True(t, validators.Excludes("hello", nil, nil).Validate(payload, nil).IsValid)
	assert.True(t, validators.Excludes(false, nil, nil).Validate(payload, nil).IsValid)

	assert.False(t, validators.Excludes(100, nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.Excludes("world", nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.Excludes(true, nil, nil).Validate(payload, nil).IsValid)

	assert.True(t, validators.IncludesAll([]interface{}{100, true}, nil, nil).Validate(payload, nil).IsValid)
	assert.True(t, validators.IncludesAll([]interface{}{true, "world"}, nil, nil).Validate(payload, nil).IsValid)

	assert.False(t, validators.IncludesAll([]interface{}{101, true}, nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.IncludesAll([]interface{}{false, "world"}, nil, nil).Validate(payload, nil).IsValid)

	assert.True(t, validators.ExcludesAll([]interface{}{101, false}, nil, nil).Validate(payload, nil).IsValid)
	assert.True(t, validators.ExcludesAll([]interface{}{false, "hello"}, nil, nil).Validate(payload, nil).IsValid)

	assert.False(t, validators.ExcludesAll([]interface{}{101, true}, nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.ExcludesAll([]interface{}{false, "world"}, nil, nil).Validate(payload, nil).IsValid)
}
