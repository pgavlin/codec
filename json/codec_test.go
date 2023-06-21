package json

import (
	"testing"

	"github.com/pgavlin/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type boolStruct struct {
	Field bool `codec:"field"`
}

type boolStructVisitor struct {
	codec.DefaultVisitor
	struct_ *boolStruct
}

func (v boolStructVisitor) VisitMap(map_ codec.MapDecoder) error {
	for {
		var k string
		ok, err := map_.NextKey(codec.NewString(&k))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		var d codec.Deserializer
		switch k {
		case "field":
			d = codec.NewBool(&v.struct_.Field)
		default:
			d = codec.SkipCodec{}
		}

		if err := map_.NextValue(d); err != nil {
			return err
		}
	}
}

func (v boolStructVisitor) Deserialize(d codec.Decoder) error {
	return d.DecodeStruct("boolStruct", v)
}

func TestCodec(t *testing.T) {
	var b bool
	_, err := Parse([]byte("true"), codec.NewBool(&b), 0)
	require.NoError(t, err)
	assert.Equal(t, true, b)

	var bp *bool
	_, err = Parse([]byte("true"), codec.NewPtr[codec.BoolCodec[bool]](&bp), 0)
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	_, err = Parse([]byte("null"), codec.NewPtr[codec.BoolCodec[bool]](&bp), 0)
	require.NoError(t, err)
	require.Nil(t, bp)

	var bools []bool
	_, err = Parse([]byte("[true, false]"), codec.NewSeq[codec.BoolCodec[bool]](&bools), 0)
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, bools)

	var boolMap map[string]bool
	_, err = Parse([]byte(`{"42": true}`), codec.NewMap[codec.StringCodec[string], codec.BoolCodec[bool]](&boolMap), 0)
	require.NoError(t, err)
	require.Equal(t, map[string]bool{"42": true}, boolMap)

	var struct_ boolStruct
	_, err = Parse([]byte(`{"field": true}`), boolStructVisitor{struct_: &struct_}, 0)
	require.NoError(t, err)
	require.Equal(t, boolStruct{Field: true}, struct_)

	var any_ any
	_, err = Parse([]byte(`{"field": true}`), codec.NewAny(&any_), 0)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"field": true}, any_)
}

func TestCodecReflect(t *testing.T) {
	var b bool
	_, err := Parse([]byte("true"), codec.GetDeserializer(&b), 0)
	require.NoError(t, err)
	assert.Equal(t, true, b)

	var bp *bool
	_, err = Parse([]byte("true"), codec.GetDeserializer(&bp), 0)
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	_, err = Parse([]byte("null"), codec.GetDeserializer(&bp), 0)
	require.NoError(t, err)
	require.Nil(t, bp)

	var bools []bool
	_, err = Parse([]byte("[true, false]"), codec.GetDeserializer(&bools), 0)
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, bools)

	var boolMap map[string]bool
	_, err = Parse([]byte(`{"42": true}`), codec.GetDeserializer(&boolMap), 0)
	require.NoError(t, err)
	require.Equal(t, map[string]bool{"42": true}, boolMap)

	var struct_ boolStruct
	_, err = Parse([]byte(`{"field": true}`), codec.GetDeserializer(&struct_), 0)
	require.NoError(t, err)
	require.Equal(t, boolStruct{Field: true}, struct_)

	var any_ any
	_, err = Parse([]byte(`{"field": true}`), codec.GetDeserializer(&any_), 0)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"field": true}, any_)
}
