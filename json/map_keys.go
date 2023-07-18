package json

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/pgavlin/codec"
	"github.com/pgavlin/codec/typecache"
)

var mapKeyCodecs typecache.Cache[jsonCodec]

func getMapKeyCodec(v any) jsonCodec {
	return mapKeyCodecs.GetOrCreate(reflect.TypeOf(v), func(t reflect.Type) jsonCodec {
		var c jsonCodec
		if t.Implements(textMarshalerType) {
			c.encode = encoder.encodeTextMarshaler
		}

		p := reflect.PtrTo(t)
		if p.Implements(textUnmarshalerType) {
			c.decode = decoder.decodeTextUnmarshaler
		}

		return c
	})
}

type mapKeyEncoder struct {
	enc *Encoder
}

func (e mapKeyEncoder) encode(v any, s codec.Serializer) (err error) {
	codec := getMapKeyCodec(v)
	if codec.encode != nil {
		e.enc.out, err = codec.encode(e.enc.enc, e.enc.out, v)
		return
	}
	return s.Serialize(e)
}

func (e mapKeyEncoder) EncodeNil() error {
	return &UnsupportedTypeError{Type: nilType}
}

func (e mapKeyEncoder) EncodeBool(v bool) error {
	return &UnsupportedTypeError{Type: boolType}
}

func (e mapKeyEncoder) EncodeInt(v int) error {
	return e.enc.EncodeString(strconv.FormatInt(int64(v), 10))
}

func (e mapKeyEncoder) EncodeInt8(v int8) error {
	return e.enc.EncodeString(strconv.FormatInt(int64(v), 10))
}

func (e mapKeyEncoder) EncodeInt16(v int16) error {
	return e.enc.EncodeString(strconv.FormatInt(int64(v), 10))
}

func (e mapKeyEncoder) EncodeInt32(v int32) error {
	return e.enc.EncodeString(strconv.FormatInt(int64(v), 10))
}

func (e mapKeyEncoder) EncodeInt64(v int64) error {
	return e.enc.EncodeString(strconv.FormatInt(v, 10))
}

func (e mapKeyEncoder) EncodeUint(v uint) error {
	return e.enc.EncodeString(strconv.FormatUint(uint64(v), 10))
}

func (e mapKeyEncoder) EncodeUintptr(v uintptr) error {
	return e.enc.EncodeString(strconv.FormatUint(uint64(v), 10))
}

func (e mapKeyEncoder) EncodeUint8(v uint8) error {
	return e.enc.EncodeString(strconv.FormatUint(uint64(v), 10))
}

func (e mapKeyEncoder) EncodeUint16(v uint16) error {
	return e.enc.EncodeString(strconv.FormatUint(uint64(v), 10))
}

func (e mapKeyEncoder) EncodeUint32(v uint32) error {
	return e.enc.EncodeString(strconv.FormatUint(uint64(v), 10))
}

func (e mapKeyEncoder) EncodeUint64(v uint64) error {
	return e.enc.EncodeString(strconv.FormatUint(v, 10))
}

func (e mapKeyEncoder) EncodeFloat32(v float32) error {
	return &UnsupportedTypeError{Type: float32Type}
}

func (e mapKeyEncoder) EncodeFloat64(v float64) error {
	return &UnsupportedTypeError{Type: float64Type}
}

func (e mapKeyEncoder) EncodeComplex64(v complex64) error {
	return &UnsupportedTypeError{Type: complex64Type}
}

func (e mapKeyEncoder) EncodeComplex128(v complex128) error {
	return &UnsupportedTypeError{Type: complex128Type}
}

func (e mapKeyEncoder) EncodeString(v string) error {
	return e.enc.EncodeString(v)
}

func (e mapKeyEncoder) EncodeBytes(v []byte) error {
	return &UnsupportedTypeError{Type: bytesType}
}

func (e mapKeyEncoder) EncodeElem(v any, s codec.Serializer) error {
	return &UnsupportedTypeError{Type: reflect.TypeOf(v)}
}

func (e mapKeyEncoder) EncodeSeq(count int) (codec.SeqEncoder, error) {
	return nil, &UnsupportedTypeError{Type: sliceType}
}

func (e mapKeyEncoder) EncodeMap(len int) (codec.MapEncoder, error) {
	return nil, &UnsupportedTypeError{Type: mapType}
}

func (e mapKeyEncoder) EncodeStruct(name string) (codec.StructEncoder, error) {
	return nil, &UnsupportedTypeError{Type: structType}
}

type mapKeyDecoder struct {
	dec *Decoder
}

func (d mapKeyDecoder) decode(v any, ds codec.Deserializer) (err error) {
	codec := getMapKeyCodec(v)
	if codec.decode != nil {
		d.dec.rest, err = codec.decode(decoder{flags: d.dec.flags}, d.dec.rest, v)
		return
	}
	return ds.Deserialize(d)
}

func (d mapKeyDecoder) decodeString() (string, error) {
	dec := decoder{flags: d.dec.flags}

	s, r, _, err := dec.parseStringUnquote(d.dec.rest, nil)
	if err != nil {
		d.dec.rest = r
		return "", err
	}
	str := *(*string)(unsafe.Pointer(&s))
	d.dec.rest = r
	return str, nil
}

func (d mapKeyDecoder) decodeInt(v codec.Visitor, t reflect.Type) (n int64, err error) {
	s, err := d.decodeString()
	if err != nil {
		return 0, err
	}
	n, err = strconv.ParseInt(s, 10, 64)
	if err != nil || reflect.Zero(t).OverflowInt(n) {
		return 0, &UnmarshalTypeError{Value: "number " + s, Type: t}
	}
	return n, nil
}

func (d mapKeyDecoder) decodeUint(v codec.Visitor, t reflect.Type) (n uint64, err error) {
	s, err := d.decodeString()
	if err != nil {
		return 0, err
	}
	n, err = strconv.ParseUint(s, 10, 64)
	if err != nil || reflect.Zero(t).OverflowUint(n) {
		return 0, &UnmarshalTypeError{Value: "number " + s, Type: t}
	}
	return n, nil
}

func (d mapKeyDecoder) DecodeAny(v codec.Visitor) (err error) {
	return &UnsupportedTypeError{Type: nilType}
}

func (d mapKeyDecoder) DecodeNil(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: nilType}
}

func (d mapKeyDecoder) DecodeBool(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: boolType}
}

func (d mapKeyDecoder) DecodeInt(v codec.Visitor) error {
	n, err := d.decodeInt(v, intType)
	if err != nil {
		return err
	}
	return v.VisitInt(int(n))
}

func (d mapKeyDecoder) DecodeInt8(v codec.Visitor) error {
	n, err := d.decodeInt(v, int8Type)
	if err != nil {
		return err
	}
	return v.VisitInt8(int8(n))
}

func (d mapKeyDecoder) DecodeInt16(v codec.Visitor) error {
	n, err := d.decodeInt(v, int16Type)
	if err != nil {
		return err
	}
	return v.VisitInt16(int16(n))
}

func (d mapKeyDecoder) DecodeInt32(v codec.Visitor) error {
	n, err := d.decodeInt(v, int32Type)
	if err != nil {
		return err
	}
	return v.VisitInt32(int32(n))
}

func (d mapKeyDecoder) DecodeInt64(v codec.Visitor) error {
	n, err := d.decodeInt(v, int64Type)
	if err != nil {
		return err
	}
	return v.VisitInt64(n)
}

func (d mapKeyDecoder) DecodeUint(v codec.Visitor) error {
	n, err := d.decodeUint(v, uintType)
	if err != nil {
		return err
	}
	return v.VisitUint(uint(n))
}

func (d mapKeyDecoder) DecodeUint8(v codec.Visitor) error {
	n, err := d.decodeUint(v, uint8Type)
	if err != nil {
		return err
	}
	return v.VisitUint8(uint8(n))
}

func (d mapKeyDecoder) DecodeUint16(v codec.Visitor) error {
	n, err := d.decodeUint(v, uint16Type)
	if err != nil {
		return err
	}
	return v.VisitUint16(uint16(n))
}

func (d mapKeyDecoder) DecodeUint32(v codec.Visitor) error {
	n, err := d.decodeUint(v, uint32Type)
	if err != nil {
		return err
	}
	return v.VisitUint32(uint32(n))
}

func (d mapKeyDecoder) DecodeUint64(v codec.Visitor) error {
	n, err := d.decodeUint(v, uint64Type)
	if err != nil {
		return err
	}
	return v.VisitUint64(n)
}

func (d mapKeyDecoder) DecodeUintptr(v codec.Visitor) error {
	n, err := d.decodeUint(v, uintptrType)
	if err != nil {
		return err
	}
	return v.VisitUintptr(uintptr(n))
}

func (d mapKeyDecoder) DecodeFloat32(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: float32Type}
}

func (d mapKeyDecoder) DecodeFloat64(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: float64Type}
}

func (d mapKeyDecoder) DecodeComplex64(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: complex64Type}
}

func (d mapKeyDecoder) DecodeComplex128(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: complex128Type}
}

func (d mapKeyDecoder) DecodeString(v codec.Visitor) error {
	return d.dec.DecodeString(v)
}

func (d mapKeyDecoder) DecodeBytes(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: bytesType}
}

func (d mapKeyDecoder) DecodePtr(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: ptrType}
}

func (d mapKeyDecoder) DecodeSeq(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: sliceType}
}

func (d mapKeyDecoder) DecodeMap(v codec.Visitor) error {
	return &UnsupportedTypeError{Type: mapType}
}

func (d mapKeyDecoder) DecodeStruct(name string, v codec.Visitor) error {
	return &UnsupportedTypeError{Type: structType}
}
