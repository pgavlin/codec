package any

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type boolStruct struct {
	Field bool `codec:"field"`
}

func TestDecode(t *testing.T) {
	b, err := Decode[bool](true)
	require.NoError(t, err)
	assert.Equal(t, true, b)

	bp, err := Decode[*bool](true)
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	bp, err = Decode[*bool](nil)
	require.NoError(t, err)
	require.Nil(t, bp)

	bools, err := Decode[[]bool]([]any{true, false})
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, bools)

	boolMap, err := Decode[map[string]bool](map[string]any{"42": true})
	require.NoError(t, err)
	require.Equal(t, map[string]bool{"42": true}, boolMap)

	struct_, err := Decode[boolStruct](map[string]any{"field": true})
	require.NoError(t, err)
	require.Equal(t, boolStruct{Field: true}, struct_)

	any_, err := Decode[any](map[string]any{"field": true})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"field": true}, any_)
}

func ptr[T any](v T) *T {
	return &v
}

func TestEncode(t *testing.T) {
	v, err := Encode(true)
	require.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = Encode(ptr(true))
	require.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = Encode[*bool](nil)
	require.NoError(t, err)
	assert.Nil(t, v)

	v, err = Encode([]bool{true, false})
	require.NoError(t, err)
	assert.Equal(t, []any{true, false}, v)

	v, err = Encode(map[string]bool{"42": true})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"42": true}, v)

	v, err = Encode(boolStruct{Field: true})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"field": true}, v)
}
