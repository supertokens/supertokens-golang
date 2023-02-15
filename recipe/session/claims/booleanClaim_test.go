package claims

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestBooleanClaim(t *testing.T) {
	boolClaim, validators := BooleanClaim(
		"test",
		func(userId string, tenantId *string, userContext supertokens.UserContext) (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		nil,
	)

	payload := map[string]interface{}{}
	payload = boolClaim.AddToPayload_internal(payload, true, nil)

	assert.True(t, validators.IsTrue(nil, nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.IsFalse(nil, nil).Validate(payload, nil).IsValid)

	maxAge := int64(1)
	assert.True(t, validators.IsTrue(&maxAge, nil).Validate(payload, nil).IsValid)
	time.Sleep(2 * time.Second)
	assert.False(t, validators.IsTrue(&maxAge, nil).Validate(payload, nil).IsValid)
}
