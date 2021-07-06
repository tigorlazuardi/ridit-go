package pkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_jsonize(t *testing.T) {
	const errorText = "testing"
	err := errors.New(errorText)
	val := jsonize(err)
	if v, ok := val.(string); !ok {
		require.FailNow(t, "expected string error")
	} else {
		require.NotEqual(t, "{}", v)
		assert.Equal(t, errorText, v)
	}

	ss := []string{"a", "b", "c"}
	val = jsonize(ss)
	assert.Equal(t, `["a","b","c"]`, val)
}
