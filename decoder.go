package codec

import (
	"errors"
)

type Visitor interface {
	VisitNil() error
	VisitBool(v bool) error
	VisitInt(v int) error
	VisitInt8(v int8) error
	VisitInt16(v int16) error
	VisitInt32(v int32) error
	VisitInt64(v int64) error
	VisitUint(v uint) error
	VisitUint8(v uint8) error
	VisitUint16(v uint16) error
	VisitUint32(v uint32) error
	VisitUint64(v uint64) error
	VisitUintptr(v uintptr) error
	VisitFloat32(v float32) error
	VisitFloat64(v float64) error
	VisitComplex64(v complex64) error
	VisitComplex128(v complex128) error
	VisitString(v string) error
	VisitBytes(v []byte) error
	VisitSeq(d SeqDecoder) error
	VisitMap(d MapDecoder) error
}

type Decoder interface {
	DecodeNil(v Visitor) error
	DecodeBool(v Visitor) error
	DecodeInt(v Visitor) error
	DecodeInt8(v Visitor) error
	DecodeInt16(v Visitor) error
	DecodeInt32(v Visitor) error
	DecodeInt64(v Visitor) error
	DecodeUint(v Visitor) error
	DecodeUint8(v Visitor) error
	DecodeUint16(v Visitor) error
	DecodeUint32(v Visitor) error
	DecodeUint64(v Visitor) error
	DecodeUintptr(v Visitor) error
	DecodeFloat32(v Visitor) error
	DecodeFloat64(v Visitor) error
	DecodeComplex64(v Visitor) error
	DecodeComplex128(v Visitor) error
	DecodeString(v Visitor) error
	DecodeBytes(v Visitor) error
	DecodeSeq(v Visitor) error
	DecodeMap(v Visitor) error
	DecodeStruct(name string, v Visitor) error
	DecodeAny(v Visitor) error
}

type SeqDecoder interface {
	Size() (int, bool)
	NextElement(de Deserializer) (bool, error)
}

type MapDecoder interface {
	Size() (int, bool)
	NextKey(de Deserializer) (bool, error)
	NextValue(de Deserializer) error
}

type Deserializer interface {
	Deserialize(decoder Decoder) error
}

type DefaultVisitor struct{}

func (DefaultVisitor) VisitNil() error {
	return errors.New("unexpected nil")
}

func (DefaultVisitor) VisitBool(v bool) error {
	return errors.New("unexpected bool")
}

func (DefaultVisitor) VisitInt(v int) error {
	return errors.New("unexpected int")
}

func (DefaultVisitor) VisitInt8(v int8) error {
	return errors.New("unexpected int8")
}

func (DefaultVisitor) VisitInt16(v int16) error {
	return errors.New("unexpected int16")
}

func (DefaultVisitor) VisitInt32(v int32) error {
	return errors.New("unexpected int32")
}

func (DefaultVisitor) VisitInt64(v int64) error {
	return errors.New("unexpected int64")
}

func (DefaultVisitor) VisitUint(v uint) error {
	return errors.New("unexpected uint")
}

func (DefaultVisitor) VisitUint8(v uint8) error {
	return errors.New("unexpected uint8")
}

func (DefaultVisitor) VisitUint16(v uint16) error {
	return errors.New("unexpected uint16")
}

func (DefaultVisitor) VisitUint32(v uint32) error {
	return errors.New("unexpected uint32")
}

func (DefaultVisitor) VisitUint64(v uint64) error {
	return errors.New("unexpected uint64")
}

func (DefaultVisitor) VisitUintptr(v uintptr) error {
	return errors.New("unexpected uintptr")
}

func (DefaultVisitor) VisitFloat32(v float32) error {
	return errors.New("unexpected float32")
}

func (DefaultVisitor) VisitFloat64(v float64) error {
	return errors.New("unexpected float64")
}

func (DefaultVisitor) VisitComplex64(v complex64) error {
	return errors.New("unexpected complex64")
}

func (DefaultVisitor) VisitComplex128(v complex128) error {
	return errors.New("unexpected complex128")
}

func (DefaultVisitor) VisitString(v string) error {
	return errors.New("unexpected string")
}

func (DefaultVisitor) VisitBytes(v []byte) error {
	return errors.New("unexpected bytes")
}

func (DefaultVisitor) VisitSeq(d SeqDecoder) error {
	return errors.New("unexpected sequence")
}

func (DefaultVisitor) VisitMap(d MapDecoder) error {
	return errors.New("unexpected map")
}
