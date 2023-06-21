package codec

import "io"

type Encoder interface {
	EncodeNil() error
	EncodeBool(v bool) error
	EncodeInt(v int) error
	EncodeInt8(v int8) error
	EncodeInt16(v int16) error
	EncodeInt32(v int32) error
	EncodeInt64(v int64) error
	EncodeUint(v uint) error
	EncodeUint8(v uint8) error
	EncodeUint16(v uint16) error
	EncodeUint32(v uint32) error
	EncodeUint64(v uint64) error
	EncodeUintptr(v uintptr) error
	EncodeFloat32(v float32) error
	EncodeFloat64(v float64) error
	EncodeComplex64(v complex64) error
	EncodeComplex128(v complex128) error
	EncodeString(v string) error
	EncodeBytes(b []byte) error
	EncodeSeq(len int) (SeqEncoder, error)
	EncodeMap(len int) (MapEncoder, error)
	EncodeStruct(name string) (StructEncoder, error)
}

type SeqEncoder interface {
	io.Closer

	EncodeElement(v Serializer) error
}

type MapEncoder interface {
	io.Closer

	EncodeKey(v Serializer) error
	EncodeValue(v Serializer) error
}

type StructEncoder interface {
	io.Closer

	EncodeField(key string, v Serializer) error
}

type Serializer interface {
	Serialize(encoder Encoder) error
}
