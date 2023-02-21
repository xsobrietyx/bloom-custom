package bloom_custom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter_Generic(t *testing.T) {
	filter := &Filter{}
	filter = filter.New(1000000, 0.2).(*Filter)

	filter.Set([]byte("Hello"))
	filter.Set([]byte("Good"))
	filter.Set([]byte("Bad"))

	assert.Equal(t, true, filter.Verify([]byte("Hello")))
	assert.Equal(t, true, filter.Verify([]byte("Good")))
	assert.Equal(t, true, filter.Verify([]byte("Bad")))

	assert.Equal(t, false, filter.Verify([]byte("Nobody")))
	assert.Equal(t, false, filter.Verify([]byte("Nowhere")))
	assert.Equal(t, false, filter.Verify([]byte("Never")))
}
