package any

import (
	"errors"

	"github.com/pgavlin/codec"
)

type Encoder struct {
	v *any
}

func Encode[T any](v T) (res any, err error) {
	err = codec.GetSerializer(v).Serialize(NewEncoder(&res))
	return
}

func NewEncoder(v *any) *Encoder {
	return &Encoder{v: v}
}

func (e *Encoder) EncodeNil() error {
	*e.v = nil
	return nil
}

func (e *Encoder) EncodeBool(v bool) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeInt(v int) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeInt8(v int8) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeInt16(v int16) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeInt32(v int32) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeInt64(v int64) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeUint(v uint) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeUint8(v uint8) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeUint16(v uint16) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeUint32(v uint32) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeUint64(v uint64) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeUintptr(v uintptr) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeFloat32(v float32) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeFloat64(v float64) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeComplex64(v complex64) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeComplex128(v complex128) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeString(v string) error {
	*e.v = v
	return nil
}

func (e *Encoder) EncodeBytes(b []byte) error {
	*e.v = b
	return nil
}

func (e *Encoder) EncodeSeq(len int) (codec.SeqEncoder, error) {
	var vs []any
	if len != 0 {
		vs = make([]any, 0, len)
	}
	return &SeqEncoder{v: e.v, vs: vs}, nil
}

func (e *Encoder) EncodeMap(len int) (codec.MapEncoder, error) {
	var m map[string]any
	if len != 0 {
		m = make(map[string]any, len)
	} else {
		m = make(map[string]any)
	}
	return &MapEncoder{v: e.v, m: m}, nil
}

func (e *Encoder) EncodeStruct(name string) (codec.StructEncoder, error) {
	return &StructEncoder{v: e.v, m: make(map[string]any)}, nil
}

type SeqEncoder struct {
	v  *any
	vs []any
}

func (e *SeqEncoder) Close() error {
	if len(e.vs) == 0 {
		*e.v = nil
	} else {
		*e.v = e.vs
	}
	return nil
}

func (e *SeqEncoder) EncodeElement(ser codec.Serializer) error {
	var v any
	if err := ser.Serialize(NewEncoder(&v)); err != nil {
		return err
	}
	e.vs = append(e.vs, v)
	return nil
}

type MapEncoder struct {
	v   *any
	m   map[string]any
	key string
}

func (e *MapEncoder) Close() error {
	if len(e.m) == 0 {
		*e.v = nil
	} else {
		*e.v = e.m
	}
	return nil
}

func (e *MapEncoder) EncodeKey(ser codec.Serializer) error {
	return ser.Serialize(mapKeyEncoder{key: &e.key})
}

func (e *MapEncoder) EncodeValue(ser codec.Serializer) error {
	var v any
	if err := ser.Serialize(NewEncoder(&v)); err != nil {
		return err
	}
	e.m[e.key] = v
	return nil
}

type StructEncoder struct {
	v *any
	m map[string]any
}

func (e *StructEncoder) Close() error {
	*e.v = e.m
	return nil
}

func (e *StructEncoder) EncodeField(key string, ser codec.Serializer) error {
	var v any
	if err := ser.Serialize(NewEncoder(&v)); err != nil {
		return err
	}
	e.m[key] = v
	return nil
}

type mapKeyEncoder struct {
	key *string
}

func (e mapKeyEncoder) EncodeString(v string) error {
	*e.key = v
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
func (e mapKeyEncoder) EncodeSeq(len int) (codec.SeqEncoder, error) {
	return nil, errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeMap(len int) (codec.MapEncoder, error) {
	return nil, errors.New("map key must be a string")
}
func (e mapKeyEncoder) EncodeStruct(name string) (codec.StructEncoder, error) {
	return nil, errors.New("map key must be a string")
}
