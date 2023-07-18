package pulumi

import (
	"testing"

	"github.com/pgavlin/codec"
	any_codec "github.com/pgavlin/codec/any"
	"github.com/pgavlin/codec/json"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type boolStruct struct {
	Field bool `codec:"field"`
}

type valueStruct struct {
	Bool   Value[bool]            `codec:bool`
	Array  Value[[]Value[string]] `codec:array`
	Struct Value[*valueStruct]    `codec:struct`
}

func TestDecode(t *testing.T) {
	b, err := Decode[bool](resource.NewBoolProperty(true))
	require.NoError(t, err)
	assert.Equal(t, true, b)

	bp, err := Decode[*bool](resource.NewBoolProperty(true))
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	bp, err = Decode[*bool](resource.NewNullProperty())
	require.NoError(t, err)
	require.Nil(t, bp)

	bools, err := Decode[[]bool](resource.NewArrayProperty([]resource.PropertyValue{resource.NewBoolProperty(true), resource.NewBoolProperty(false)}))
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, bools)

	boolMap, err := Decode[map[string]bool](resource.NewObjectProperty(resource.PropertyMap{"42": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, map[string]bool{"42": true}, boolMap)

	struct_, err := Decode[boolStruct](resource.NewObjectProperty(resource.PropertyMap{"field": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, boolStruct{Field: true}, struct_)

	any_, err := Decode[any](resource.NewObjectProperty(resource.PropertyMap{"field": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"field": true}, any_)

	vs, err := Decode[valueStruct](resource.NewObjectProperty(resource.PropertyMap{
		"bool":  resource.MakeComputed(resource.NewStringProperty("")),
		"array": resource.NewArrayProperty([]resource.PropertyValue{resource.NewStringProperty("hello"), resource.MakeComputed(resource.NewStringProperty(""))}),
		"struct": resource.NewObjectProperty(resource.PropertyMap{
			"bool":   resource.NewBoolProperty(true),
			"array":  resource.MakeSecret(resource.MakeComputed(resource.NewStringProperty(""))),
			"struct": resource.MakeComputed(resource.NewStringProperty("")),
		}),
	}))
	require.NoError(t, err)

	assert.Equal(t, NewUnknown[bool](), vs.Bool)
	assert.Equal(t, NewValue([]Value[string]{NewValue("hello"), NewUnknown[string]()}), vs.Array)
	assert.Equal(t, NewValue(&valueStruct{
		Bool:   NewValue(true),
		Array:  NewSecretUnknown[[]Value[string]](),
		Struct: NewUnknown[*valueStruct](),
	}), vs.Struct)
}

func ptr[T any](v T) *T {
	return &v
}

func TestEncode(t *testing.T) {
	v, err := Encode(true)
	require.NoError(t, err)
	assert.Equal(t, resource.NewBoolProperty(true), v)

	v, err = Encode(ptr(true))
	require.NoError(t, err)
	assert.Equal(t, resource.NewBoolProperty(true), v)

	v, err = Encode[*bool](nil)
	require.NoError(t, err)
	assert.Equal(t, resource.NewNullProperty(), v)

	v, err = Encode([]bool{true, false})
	require.NoError(t, err)
	assert.Equal(t, resource.NewArrayProperty([]resource.PropertyValue{resource.NewBoolProperty(true), resource.NewBoolProperty(false)}), v)

	v, err = Encode(map[string]bool{"42": true})
	require.NoError(t, err)
	assert.Equal(t, resource.NewObjectProperty(resource.PropertyMap{"42": resource.NewBoolProperty(true)}), v)

	v, err = Encode(boolStruct{Field: true})
	require.NoError(t, err)
	assert.Equal(t, resource.NewObjectProperty(resource.PropertyMap{"field": resource.NewBoolProperty(true)}), v)
}

func TestDeserialize(t *testing.T) {
	deserialize := func(v any) (resource.PropertyValue, error) {
		return Deserialize(any_codec.NewDecoder(v))
	}

	v, err := deserialize(true)
	require.NoError(t, err)
	assert.Equal(t, resource.NewBoolProperty(true), v)

	v, err = deserialize(nil)
	require.NoError(t, err)
	assert.Equal(t, resource.NewNullProperty(), v)

	v, err = deserialize([]any{true, false})
	require.NoError(t, err)
	assert.Equal(t, resource.NewArrayProperty([]resource.PropertyValue{resource.NewBoolProperty(true), resource.NewBoolProperty(false)}), v)

	v, err = deserialize(map[string]any{"42": true})
	require.NoError(t, err)
	assert.Equal(t, resource.NewObjectProperty(resource.PropertyMap{"42": resource.NewBoolProperty(true)}), v)

	asset, err := resource.NewTextAsset("hello")
	require.NoError(t, err)
	asset.Sig = ""

	v, err = deserialize(asset.Serialize())
	require.NoError(t, err)
	assert.Equal(t, resource.NewAssetProperty(asset), v)
}

func TestDeserializeJSON(t *testing.T) {
	deserialize := func(v any) (pv resource.PropertyValue, err error) {
		bytes, err := json.Append(nil, v, codec.GetSerializer(v), 0)
		require.NoError(t, err)
		_, err = json.Parse(bytes, nil, NewDeserializer(&pv), 0)
		return
	}

	v, err := deserialize(true)
	require.NoError(t, err)
	assert.Equal(t, resource.NewBoolProperty(true), v)

	v, err = deserialize(nil)
	require.NoError(t, err)
	assert.Equal(t, resource.NewNullProperty(), v)

	v, err = deserialize([]any{true, false})
	require.NoError(t, err)
	assert.Equal(t, resource.NewArrayProperty([]resource.PropertyValue{resource.NewBoolProperty(true), resource.NewBoolProperty(false)}), v)

	v, err = deserialize(map[string]any{"42": true})
	require.NoError(t, err)
	assert.Equal(t, resource.NewObjectProperty(resource.PropertyMap{"42": resource.NewBoolProperty(true)}), v)

	asset, err := resource.NewTextAsset("hello")
	require.NoError(t, err)
	asset.Sig = ""

	v, err = deserialize(asset.Serialize())
	require.NoError(t, err)
	assert.Equal(t, resource.NewAssetProperty(asset), v)

}

func TestSerialize(t *testing.T) {
	serialize := func(v resource.PropertyValue) (res any, err error) {
		err = Serialize(v, any_codec.NewEncoder(&res))
		return
	}

	v, err := serialize(resource.NewBoolProperty(true))
	require.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = serialize(resource.NewNullProperty())
	require.NoError(t, err)
	require.Nil(t, v)

	v, err = serialize(resource.NewArrayProperty([]resource.PropertyValue{resource.NewBoolProperty(true), resource.NewBoolProperty(false)}))
	require.NoError(t, err)
	require.Equal(t, []any{true, false}, v)

	v, err = serialize(resource.NewObjectProperty(resource.PropertyMap{"42": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"42": true}, v)

	v, err = serialize(resource.NewObjectProperty(resource.PropertyMap{"field": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"field": true}, v)

	asset, err := resource.NewTextAsset("hello")
	require.NoError(t, err)

	v, err = serialize(resource.NewAssetProperty(asset))
	require.NoError(t, err)
	require.Equal(t, asset.Serialize(), v)
}

func TestSerializeJSON(t *testing.T) {
	serialize := func(v resource.PropertyValue) (res any, err error) {
		bytes, err := json.Append(nil, nil, NewSerializer(v), 0)
		require.NoError(t, err)
		_, err = json.Parse(bytes, nil, codec.GetDeserializer(&res), 0)
		return
	}

	v, err := serialize(resource.NewBoolProperty(true))
	require.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = serialize(resource.NewNullProperty())
	require.NoError(t, err)
	require.Nil(t, v)

	v, err = serialize(resource.NewArrayProperty([]resource.PropertyValue{resource.NewBoolProperty(true), resource.NewBoolProperty(false)}))
	require.NoError(t, err)
	require.Equal(t, []any{true, false}, v)

	v, err = serialize(resource.NewObjectProperty(resource.PropertyMap{"42": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"42": true}, v)

	v, err = serialize(resource.NewObjectProperty(resource.PropertyMap{"field": resource.NewBoolProperty(true)}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"field": true}, v)

	asset, err := resource.NewTextAsset("hello")
	require.NoError(t, err)

	v, err = serialize(resource.NewAssetProperty(asset))
	require.NoError(t, err)
	require.Equal(t, asset.Serialize(), v)
}
