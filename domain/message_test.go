package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageStatusJSON(t *testing.T) {
	for _, s := range []MessageStatus{StatusSent, StatusDelivered, StatusRead} {
		b, err := json.Marshal(s)
		assert.NoError(t, err)
		var out MessageStatus
		assert.NoError(t, json.Unmarshal(b, &out))
		assert.Equal(t, s, out)
	}
	var bad MessageStatus
	assert.Error(t, json.Unmarshal([]byte(`"foo"`), &bad))
}
