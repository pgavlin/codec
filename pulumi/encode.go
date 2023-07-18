package pulumi

import (
	"errors"

	"github.com/pgavlin/codec"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Encoder struct {
	v *resource.PropertyValue
}

func Encode[T any](v T) (res resource.PropertyValue, err error) {
	err = codec.GetSerializer(v).Serialize(NewEncoder(&res))
	return
}

func NewEncoder(v *resource.PropertyValue) Encoder {
	return Encoder{v: v}
}

func (e Encoder) encode(v any, s codec.Serializer) error {
	switch v := v.(type) {
	case resource.PropertyValue:
		*e.v = v
		return nil
	case *resource.Asset:
		*e.v = resource.NewAssetProperty(v)
		return nil
	case *resource.Archive:
		*e.v = resource.NewArchiveProperty(v)
		return nil
	case resource.ResourceReference:
		*e.v = resource.NewResourceReferenceProperty(v)
		return nil
	case Marshaler:
		p, err := v.MarshalPropertyValue()
		*e.v = p
		return err
	default:
		return s.Serialize(e)
	}
}

func (e Encoder) EncodeNil() error {
	*e.v = resource.NewNullProperty()
	return nil
}

func (e Encoder) EncodeBool(v bool) error {
	*e.v = resource.NewBoolProperty(v)
	return nil
}

func (e Encoder) EncodeInt(v int) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeInt8(v int8) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeInt16(v int16) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeInt32(v int32) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeInt64(v int64) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeUint(v uint) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeUint8(v uint8) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeUint16(v uint16) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeUint32(v uint32) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeUint64(v uint64) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeUintptr(v uintptr) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeFloat32(v float32) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeFloat64(v float64) error {
	*e.v = resource.NewNumberProperty(float64(v))
	return nil
}

func (e Encoder) EncodeComplex64(v complex64) error {
	return errors.New("unsupported")
}

func (e Encoder) EncodeComplex128(v complex128) error {
	return errors.New("unsupported")
}

func (e Encoder) EncodeString(v string) error {
	*e.v = resource.NewStringProperty(v)
	return nil
}

func (e Encoder) EncodeBytes(b []byte) error {
	return errors.New("unsupported")
}

func (e Encoder) EncodeElem(v any, s codec.Serializer) error {
	return e.encode(v, s)
}

func (e Encoder) EncodeSeq(len int) (codec.SeqEncoder, error) {
	var vs []resource.PropertyValue
	if len != 0 {
		vs = make([]resource.PropertyValue, 0, len)
	}
	return &SeqEncoder{v: e.v, vs: vs}, nil
}

func (e Encoder) EncodeMap(len int) (codec.MapEncoder, error) {
	var m resource.PropertyMap
	if len != 0 {
		m = make(resource.PropertyMap, len)
	} else {
		m = make(resource.PropertyMap)
	}
	return &MapEncoder{v: e.v, m: m}, nil
}

func (e Encoder) EncodeStruct(name string) (codec.StructEncoder, error) {
	return &StructEncoder{v: e.v, m: make(resource.PropertyMap)}, nil
}

type SeqEncoder struct {
	v  *resource.PropertyValue
	vs []resource.PropertyValue
}

func (e *SeqEncoder) Close() error {
	*e.v = resource.NewArrayProperty(e.vs)
	return nil
}

func (e *SeqEncoder) EncodeElement(x any, ser codec.Serializer) error {
	var v resource.PropertyValue
	if err := NewEncoder(&v).encode(x, ser); err != nil {
		return err
	}
	e.vs = append(e.vs, v)
	return nil
}

type MapEncoder struct {
	v   *resource.PropertyValue
	m   resource.PropertyMap
	key resource.PropertyKey
}

func (e *MapEncoder) Close() error {
	*e.v = resource.NewObjectProperty(e.m)
	return nil
}

func (e *MapEncoder) EncodeKey(_ any, ser codec.Serializer) error {
	return ser.Serialize(mapKeyEncoder{key: &e.key})
}

func (e *MapEncoder) EncodeValue(x any, ser codec.Serializer) error {
	var v resource.PropertyValue
	if err := NewEncoder(&v).encode(x, ser); err != nil {
		return err
	}
	e.m[e.key] = v
	return nil
}

type StructEncoder struct {
	v *resource.PropertyValue
	m resource.PropertyMap
}

func (e *StructEncoder) Close() error {
	*e.v = resource.NewObjectProperty(e.m)
	return nil
}

func (e *StructEncoder) EncodeField(key string, x any, ser codec.Serializer) error {
	var v resource.PropertyValue
	if err := NewEncoder(&v).encode(x, ser); err != nil {
		return err
	}
	e.m[resource.PropertyKey(key)] = v
	return nil
}

type mapKeyEncoder struct {
	key *resource.PropertyKey
}

func (e mapKeyEncoder) EncodeString(v string) error {
	*e.key = resource.PropertyKey(v)
	return nil
}

func (e mapKeyEncoder) EncodeNil() error              { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeBool(v bool) error       { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeInt(v int) error         { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeInt8(v int8) error       { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeInt16(v int16) error     { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeInt32(v int32) error     { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeInt64(v int64) error     { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeUint(v uint) error       { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeUint8(v uint8) error     { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeUint16(v uint16) error   { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeUint32(v uint32) error   { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeUint64(v uint64) error   { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeUintptr(v uintptr) error { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeFloat32(v float32) error { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeFloat64(v float64) error { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeComplex64(v complex64) error {
	return errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeComplex128(v complex128) error {
	return errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeBytes(b []byte) error { return errors.New("map key must be a string") }
func (e mapKeyEncoder) EncodeElem(v any, s codec.Serializer) error {
	return errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeSeq(len int) (codec.SeqEncoder, error) {
	return nil, errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeMap(len int) (codec.MapEncoder, error) {
	return nil, errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeStruct(name string) (codec.StructEncoder, error) {
	return nil, errors.New("map key must be a string")
}
