package claims

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestPrimitiveClaim(t *testing.T) {
	primClaim := PrimitiveClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		nil,
	)
	payload := map[string]interface{}{}
	payload = primClaim.AddToPayload_internal(payload, "world", nil)

	validators := primClaim.Validators
	assert.True(t, validators.HasValue("world", nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.HasValue("hello", nil, nil).Validate(payload, nil).IsValid)

	five := int64(5)
	one := int64(1)
	assert.True(t, validators.HasValue("world", &five, nil).Validate(payload, nil).IsValid)
	time.Sleep(2 * time.Second)
	assert.False(t, validators.HasValue("world", &one, nil).Validate(payload, nil).IsValid)
}
