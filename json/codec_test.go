package json

import (
	"encoding/json"
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
		ok, err := map_.NextKey(&k, codec.NewString(&k))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		var f any
		var d codec.Deserializer
		switch k {
		case "field":
			f, d = &v.struct_.Field, codec.NewBool(&v.struct_.Field)
		default:
			f, d = nil, codec.SkipCodec{}
		}

		if err := map_.NextValue(f, d); err != nil {
			return err
		}
	}
}

func (v boolStructVisitor) Deserialize(d codec.Decoder) error {
	return d.DecodeStruct("boolStruct", v)
}

type SecretValue struct {
	Value      string // plaintext
	Ciphertext []byte // ciphertext
	Secret     bool
}

type secretWorkflowValue struct {
	Secret     string `json:"secret,omitempty"`
	Ciphertext []byte `json:"ciphertext,omitempty"`
}

func (v SecretValue) MarshalJSON() ([]byte, error) {
	switch {
	case len(v.Ciphertext) != 0:
		return json.Marshal(secretWorkflowValue{Ciphertext: v.Ciphertext})
	case v.Secret:
		return json.Marshal(secretWorkflowValue{Secret: v.Value})
	default:
		return json.Marshal(v.Value)
	}
}

func (v *SecretValue) UnmarshalJSON(bytes []byte) error {
	var secret secretWorkflowValue
	if err := json.Unmarshal(bytes, &secret); err == nil {
		v.Value, v.Ciphertext, v.Secret = secret.Secret, secret.Ciphertext, true
		return nil
	}

	var plaintext string
	if err := json.Unmarshal(bytes, &plaintext); err != nil {
		return err
	}
	v.Value, v.Secret = plaintext, false
	return nil
}

func TestCodec(t *testing.T) {
	var b bool
	_, err := Parse([]byte("true"), &b, codec.NewBool(&b), 0)
	require.NoError(t, err)
	assert.Equal(t, true, b)

	var bp *bool
	_, err = Parse([]byte("true"), &bp, codec.NewPtr[codec.BoolCodec[bool]](&bp), 0)
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	_, err = Parse([]byte("null"), &bp, codec.NewPtr[codec.BoolCodec[bool]](&bp), 0)
	require.NoError(t, err)
	require.Nil(t, bp)

	var bools []bool
	_, err = Parse([]byte("[true, false]"), &bools, codec.NewSeq[codec.BoolCodec[bool]](&bools), 0)
	require.NoError(t, err)
	assert.Equal(t, []bool{true, false}, bools)

	var boolMap map[string]bool
	_, err = Parse([]byte(`{"42": true}`), &boolMap, codec.NewMap[codec.StringCodec[string], codec.BoolCodec[bool]](&boolMap), 0)
	require.NoError(t, err)
	assert.Equal(t, map[string]bool{"42": true}, boolMap)

	var struct_ boolStruct
	_, err = Parse([]byte(`{"field": true}`), &struct_, boolStructVisitor{struct_: &struct_}, 0)
	require.NoError(t, err)
	assert.Equal(t, boolStruct{Field: true}, struct_)

	var any_ any
	_, err = Parse([]byte(`{"field": true}`), &any_, codec.NewAny(&any_), 0)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"field": true}, any_)

	var secret SecretValue
	_, err = Parse([]byte(`"plaintext"`), &secret, nil, 0)
	require.NoError(t, err)
	assert.Equal(t, SecretValue{Value: "plaintext"}, secret)
}

func TestCodecReflect(t *testing.T) {
	var b bool
	_, err := Parse([]byte("true"), &b, codec.GetDeserializer(&b), 0)
	require.NoError(t, err)
	assert.Equal(t, true, b)

	var bp *bool
	_, err = Parse([]byte("true"), &bp, codec.GetDeserializer(&bp), 0)
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	_, err = Parse([]byte("null"), &bp, codec.GetDeserializer(&bp), 0)
	require.NoError(t, err)
	require.Nil(t, bp)

	var bools []bool
	_, err = Parse([]byte("[true, false]"), &bools, codec.GetDeserializer(&bools), 0)
	require.NoError(t, err)
	assert.Equal(t, []bool{true, false}, bools)

	var boolMap map[string]bool
	_, err = Parse([]byte(`{"42": true}`), &boolMap, codec.GetDeserializer(&boolMap), 0)
	require.NoError(t, err)
	assert.Equal(t, map[string]bool{"42": true}, boolMap)

	var struct_ boolStruct
	_, err = Parse([]byte(`{"field": true}`), &struct_, codec.GetDeserializer(&struct_), 0)
	require.NoError(t, err)
	assert.Equal(t, boolStruct{Field: true}, struct_)

	var any_ any
	_, err = Parse([]byte(`{"field": true}`), &any_, codec.GetDeserializer(&any_), 0)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"field": true}, any_)

	var secret SecretValue
	_, err = Parse([]byte(`"plaintext"`), &secret, codec.GetDeserializer(&secret), 0)
	require.NoError(t, err)
	assert.Equal(t, SecretValue{Value: "plaintext"}, secret)
}
