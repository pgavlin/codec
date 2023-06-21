package codec

import (
	"unsafe"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type Codec[T any] interface {
	Visitor
	Deserializer
	Serializer

	New(v *T) Codec[T]
}

type SkipCodec struct{}

func (SkipCodec) VisitNil() error                    { return nil }
func (SkipCodec) VisitBool(v bool) error             { return nil }
func (SkipCodec) VisitInt(v int) error               { return nil }
func (SkipCodec) VisitInt8(v int8) error             { return nil }
func (SkipCodec) VisitInt16(v int16) error           { return nil }
func (SkipCodec) VisitInt32(v int32) error           { return nil }
func (SkipCodec) VisitInt64(v int64) error           { return nil }
func (SkipCodec) VisitUint(v uint) error             { return nil }
func (SkipCodec) VisitUint8(v uint8) error           { return nil }
func (SkipCodec) VisitUint16(v uint16) error         { return nil }
func (SkipCodec) VisitUint32(v uint32) error         { return nil }
func (SkipCodec) VisitUint64(v uint64) error         { return nil }
func (SkipCodec) VisitUintptr(v uintptr) error       { return nil }
func (SkipCodec) VisitFloat32(v float32) error       { return nil }
func (SkipCodec) VisitFloat64(v float64) error       { return nil }
func (SkipCodec) VisitComplex64(v complex64) error   { return nil }
func (SkipCodec) VisitComplex128(v complex128) error { return nil }
func (SkipCodec) VisitString(v string) error         { return nil }

func (SkipCodec) VisitSeq(d SeqDecoder) error {
	for {
		ok, err := d.NextElement(SkipCodec{})
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
}

func (SkipCodec) VisitMap(d MapDecoder) error {
	for {
		ok, err := d.NextKey(SkipCodec{})
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := d.NextValue(SkipCodec{}); err != nil {
			return err
		}
	}
}

func (SkipCodec) Deserialize(d Decoder) error {
	return d.DecodeAny(SkipCodec{})
}

type AnyCodec struct {
	value *any
}

func NewAny(v *any) AnyCodec {
	return AnyCodec{value: v}
}

func (c AnyCodec) new(v unsafe.Pointer) codec {
	return AnyCodec{value: (*any)(v)}
}

func (c AnyCodec) New(v *any) Codec[any] {
	return AnyCodec{value: v}
}

func (c AnyCodec) VisitNil() error {
	*c.value = nil
	return nil
}

func (c AnyCodec) VisitBool(v bool) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitInt(v int) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitInt8(v int8) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitInt16(v int16) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitInt32(v int32) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitInt64(v int64) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitUint(v uint) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitUint8(v uint8) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitUint16(v uint16) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitUint32(v uint32) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitUint64(v uint64) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitUintptr(v uintptr) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitFloat32(v float32) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitFloat64(v float64) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitComplex64(v complex64) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitComplex128(v complex128) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitString(v string) error {
	*c.value = v
	return nil
}

func (c AnyCodec) VisitSeq(seq SeqDecoder) error {
	var vals []any
	if len, ok := seq.Size(); ok {
		vals = make([]any, 0, len)
	}
	for {
		var v any
		ok, err := seq.NextElement(AnyCodec{value: &v})
		if err != nil {
			return err
		}
		if !ok {
			*c.value = vals
			return nil
		}
		vals = append(vals, v)
	}
}

func (c AnyCodec) VisitMap(map_ MapDecoder) error {
	var m map[string]any
	if len, ok := map_.Size(); ok {
		m = make(map[string]any, len)
	} else {
		m = make(map[string]any)
	}
	for {
		var k string
		ok, err := map_.NextKey(NewString(&k))
		if err != nil {
			return err
		}
		if !ok {
			*c.value = m
			return nil
		}

		var v any
		if err = map_.NextValue(AnyCodec{value: &v}); err != nil {
			return err
		}
		m[k] = v
	}
}

func (c AnyCodec) Deserialize(d Decoder) error {
	return d.DecodeAny(c)
}

func (c AnyCodec) Serialize(e Encoder) error {
	if *c.value == nil {
		return e.EncodeNil()
	}
	return GetSerializer(*c.value).Serialize(e)
}

type NilCodec struct {
	DefaultVisitor
}

func (NilCodec) VisitNil() error {
	return nil
}

func (NilCodec) Deserialize(d Decoder) error {
	return d.DecodeNil(NilCodec{})
}

func (NilCodec) Serialize(e Encoder) error {
	return e.EncodeNil()
}

type BoolCodec[T ~bool] struct {
	DefaultVisitor
	value *T
}

func NewBool[T ~bool](v *T) BoolCodec[T] {
	return BoolCodec[T]{value: v}
}

func (BoolCodec[T]) New(v *T) Codec[T] {
	return BoolCodec[T]{value: v}
}

func (c BoolCodec[T]) VisitBool(v bool) error {
	*c.value = T(v)
	return nil
}

func (c BoolCodec[T]) Deserialize(d Decoder) error {
	return d.DecodeBool(c)
}

func (c BoolCodec[T]) Serialize(e Encoder) error {
	return e.EncodeBool(bool(*c.value))
}

type IntCodec[T ~int] struct {
	DefaultVisitor
	value *T
}

func NewInt[T ~int](v *T) IntCodec[T] {
	return IntCodec[T]{value: v}
}

func (IntCodec[T]) New(v *T) Codec[T] {
	return IntCodec[T]{value: v}
}

func (c IntCodec[T]) VisitInt(v int) error {
	*c.value = T(v)
	return nil
}

func (c IntCodec[T]) Deserialize(d Decoder) error {
	return d.DecodeInt(c)
}

func (c IntCodec[T]) Serialize(e Encoder) error {
	return e.EncodeInt(int(*c.value))
}

type Int8Codec[T ~int8] struct {
	DefaultVisitor
	value *T
}

func NewInt8[T ~int8](v *T) Int8Codec[T] {
	return Int8Codec[T]{value: v}
}

func (Int8Codec[T]) New(v *T) Codec[T] {
	return Int8Codec[T]{value: v}
}

func (c Int8Codec[T]) VisitInt8(v int8) error {
	*c.value = T(v)
	return nil
}

func (c Int8Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeInt8(c)
}

func (c Int8Codec[T]) Serialize(e Encoder) error {
	return e.EncodeInt8(int8(*c.value))
}

type Int16Codec[T ~int16] struct {
	DefaultVisitor
	value *T
}

func NewInt16[T ~int16](v *T) Int16Codec[T] {
	return Int16Codec[T]{value: v}
}

func (Int16Codec[T]) New(v *T) Codec[T] {
	return Int16Codec[T]{value: v}
}

func (c Int16Codec[T]) VisitInt16(v int16) error {
	*c.value = T(v)
	return nil
}

func (c Int16Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeInt16(c)
}

func (c Int16Codec[T]) Serialize(e Encoder) error {
	return e.EncodeInt16(int16(*c.value))
}

type Int32Codec[T ~int32] struct {
	DefaultVisitor
	value *T
}

func NewInt32[T ~int32](v *T) Int32Codec[T] {
	return Int32Codec[T]{value: v}
}

func (Int32Codec[T]) New(v *T) Codec[T] {
	return Int32Codec[T]{value: v}
}

func (c Int32Codec[T]) VisitInt32(v int32) error {
	*c.value = T(v)
	return nil
}

func (c Int32Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeInt32(c)
}

func (c Int32Codec[T]) Serialize(e Encoder) error {
	return e.EncodeInt32(int32(*c.value))
}

type Int64Codec[T ~int64] struct {
	DefaultVisitor
	value *T
}

func NewInt64[T ~int64](v *T) Int64Codec[T] {
	return Int64Codec[T]{value: v}
}

func (Int64Codec[T]) New(v *T) Codec[T] {
	return Int64Codec[T]{value: v}
}

func (c Int64Codec[T]) VisitInt64(v int64) error {
	*c.value = T(v)
	return nil
}

func (c Int64Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeInt64(c)
}

func (c Int64Codec[T]) Serialize(e Encoder) error {
	return e.EncodeInt64(int64(*c.value))
}

type UintCodec[T ~uint] struct {
	DefaultVisitor
	value *T
}

func NewUint[T ~uint](v *T) UintCodec[T] {
	return UintCodec[T]{value: v}
}

func (UintCodec[T]) New(v *T) Codec[T] {
	return UintCodec[T]{value: v}
}

func (c UintCodec[T]) VisitUint(v uint) error {
	*c.value = T(v)
	return nil
}

func (c UintCodec[T]) Deserialize(d Decoder) error {
	return d.DecodeUint(c)
}

func (c UintCodec[T]) Serialize(e Encoder) error {
	return e.EncodeUint(uint(*c.value))
}

type Uint8Codec[T ~uint8] struct {
	DefaultVisitor
	value *T
}

func NewUint8[T ~uint8](v *T) Uint8Codec[T] {
	return Uint8Codec[T]{value: v}
}

func (Uint8Codec[T]) New(v *T) Codec[T] {
	return Uint8Codec[T]{value: v}
}

func (c Uint8Codec[T]) VisitUint8(v uint8) error {
	*c.value = T(v)
	return nil
}

func (c Uint8Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeUint8(c)
}

func (c Uint8Codec[T]) Serialize(e Encoder) error {
	return e.EncodeUint8(uint8(*c.value))
}

type Uint16Codec[T ~uint16] struct {
	DefaultVisitor
	value *T
}

func NewUint16[T ~uint16](v *T) Uint16Codec[T] {
	return Uint16Codec[T]{value: v}
}

func (Uint16Codec[T]) New(v *T) Codec[T] {
	return Uint16Codec[T]{value: v}
}

func (c Uint16Codec[T]) VisitUint16(v uint16) error {
	*c.value = T(v)
	return nil
}

func (c Uint16Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeUint16(c)
}

func (c Uint16Codec[T]) Serialize(e Encoder) error {
	return e.EncodeUint16(uint16(*c.value))
}

type Uint32Codec[T ~uint32] struct {
	DefaultVisitor
	value *T
}

func NewUint32[T ~uint32](v *T) Uint32Codec[T] {
	return Uint32Codec[T]{value: v}
}

func (Uint32Codec[T]) New(v *T) Codec[T] {
	return Uint32Codec[T]{value: v}
}

func (c Uint32Codec[T]) VisitUint32(v uint32) error {
	*c.value = T(v)
	return nil
}

func (c Uint32Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeUint32(c)
}

func (c Uint32Codec[T]) Serialize(e Encoder) error {
	return e.EncodeUint32(uint32(*c.value))
}

type Uint64Codec[T ~uint64] struct {
	DefaultVisitor
	value *T
}

func NewUint64[T ~uint64](v *T) Uint64Codec[T] {
	return Uint64Codec[T]{value: v}
}

func (Uint64Codec[T]) New(v *T) Codec[T] {
	return Uint64Codec[T]{value: v}
}

func (c Uint64Codec[T]) VisitUint64(v uint64) error {
	*c.value = T(v)
	return nil
}

func (c Uint64Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeUint64(c)
}

func (c Uint64Codec[T]) Serialize(e Encoder) error {
	return e.EncodeUint64(uint64(*c.value))
}

type UintptrCodec[T ~uintptr] struct {
	DefaultVisitor
	value *T
}

func NewUintptr[T ~uintptr](v *T) UintptrCodec[T] {
	return UintptrCodec[T]{value: v}
}

func (UintptrCodec[T]) New(v *T) Codec[T] {
	return UintptrCodec[T]{value: v}
}

func (c UintptrCodec[T]) VisitUintptr(v uintptr) error {
	*c.value = T(v)
	return nil
}

func (c UintptrCodec[T]) Deserialize(d Decoder) error {
	return d.DecodeUintptr(c)
}

func (c UintptrCodec[T]) Serialize(e Encoder) error {
	return e.EncodeUintptr(uintptr(*c.value))
}

type Float32Codec[T ~float32] struct {
	DefaultVisitor
	value *T
}

func NewFloat32[T ~float32](v *T) Float32Codec[T] {
	return Float32Codec[T]{value: v}
}

func (Float32Codec[T]) New(v *T) Codec[T] {
	return Float32Codec[T]{value: v}
}

func (c Float32Codec[T]) VisitFloat32(v float32) error {
	*c.value = T(v)
	return nil
}

func (c Float32Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeFloat32(c)
}

func (c Float32Codec[T]) Serialize(e Encoder) error {
	return e.EncodeFloat32(float32(*c.value))
}

type Float64Codec[T ~float64] struct {
	DefaultVisitor
	value *T
}

func NewFloat64[T ~float64](v *T) Float64Codec[T] {
	return Float64Codec[T]{value: v}
}

func (Float64Codec[T]) New(v *T) Codec[T] {
	return Float64Codec[T]{value: v}
}

func (c Float64Codec[T]) VisitFloat64(v float64) error {
	*c.value = T(v)
	return nil
}

func (c Float64Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeFloat64(c)
}

func (c Float64Codec[T]) Serialize(e Encoder) error {
	return e.EncodeFloat64(float64(*c.value))
}

type Complex64Codec[T ~complex64] struct {
	DefaultVisitor
	value *T
}

func NewComplex64[T ~complex64](v *T) Complex64Codec[T] {
	return Complex64Codec[T]{value: v}
}

func (Complex64Codec[T]) New(v *T) Codec[T] {
	return Complex64Codec[T]{value: v}
}

func (c Complex64Codec[T]) VisitComplex64(v complex64) error {
	*c.value = T(v)
	return nil
}

func (c Complex64Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeComplex64(c)
}

func (c Complex64Codec[T]) Serialize(e Encoder) error {
	return e.EncodeComplex64(complex64(*c.value))
}

type Complex128Codec[T ~complex128] struct {
	DefaultVisitor
	value *T
}

func NewComplex128[T ~complex128](v *T) Complex128Codec[T] {
	return Complex128Codec[T]{value: v}
}

func (Complex128Codec[T]) New(v *T) Codec[T] {
	return Complex128Codec[T]{value: v}
}

func (c Complex128Codec[T]) VisitComplex128(v complex128) error {
	*c.value = T(v)
	return nil
}

func (c Complex128Codec[T]) Deserialize(d Decoder) error {
	return d.DecodeComplex128(c)
}

func (c Complex128Codec[T]) Serialize(e Encoder) error {
	return e.EncodeComplex128(complex128(*c.value))
}

type StringCodec[T ~string] struct {
	DefaultVisitor
	value *T
}

func NewString[T ~string](v *T) StringCodec[T] {
	return StringCodec[T]{value: v}
}

func (StringCodec[T]) New(v *T) Codec[T] {
	return StringCodec[T]{value: v}
}

func (c StringCodec[T]) VisitString(v string) error {
	*c.value = T(v)
	return nil
}

func (c StringCodec[T]) Deserialize(d Decoder) error {
	return d.DecodeString(c)
}

func (c StringCodec[T]) Serialize(e Encoder) error {
	return e.EncodeString(string(*c.value))
}

type PtrCodec[P ~*T, T any, C Codec[T]] struct {
	value *P
}

func NewPtr[C Codec[T], P ~*T, T any](v *P) PtrCodec[P, T, C] {
	return PtrCodec[P, T, C]{value: v}
}

func (c PtrCodec[P, T, C]) codec(v *T) Codec[T] {
	var codec C
	return codec.New(v)
}

func (c PtrCodec[P, T, C]) New(v *P) PtrCodec[P, T, C] {
	return PtrCodec[P, T, C]{value: v}
}

func (c PtrCodec[P, T, C]) VisitNil() error {
	*c.value = nil
	return nil
}

func (c PtrCodec[P, T, C]) VisitBool(dv bool) error {
	var v T
	if err := c.codec(&v).VisitBool(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitInt(dv int) error {
	var v T
	if err := c.codec(&v).VisitInt(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitInt8(dv int8) error {
	var v T
	if err := c.codec(&v).VisitInt8(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitInt16(dv int16) error {
	var v T
	if err := c.codec(&v).VisitInt16(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitInt32(dv int32) error {
	var v T
	if err := c.codec(&v).VisitInt32(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitInt64(dv int64) error {
	var v T
	if err := c.codec(&v).VisitInt64(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitUint(dv uint) error {
	var v T
	if err := c.codec(&v).VisitUint(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitUint8(dv uint8) error {
	var v T
	if err := c.codec(&v).VisitUint8(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitUint16(dv uint16) error {
	var v T
	if err := c.codec(&v).VisitUint16(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitUint32(dv uint32) error {
	var v T
	if err := c.codec(&v).VisitUint32(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitUint64(dv uint64) error {
	var v T
	if err := c.codec(&v).VisitUint64(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitUintptr(dv uintptr) error {
	var v T
	if err := c.codec(&v).VisitUintptr(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitFloat32(dv float32) error {
	var v T
	if err := c.codec(&v).VisitFloat32(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitFloat64(dv float64) error {
	var v T
	if err := c.codec(&v).VisitFloat64(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitComplex64(dv complex64) error {
	var v T
	if err := c.codec(&v).VisitComplex64(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitComplex128(dv complex128) error {
	var v T
	if err := c.codec(&v).VisitComplex128(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitString(dv string) error {
	var v T
	if err := c.codec(&v).VisitString(dv); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitSeq(d SeqDecoder) error {
	var v T
	if err := c.codec(&v).VisitSeq(d); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) VisitMap(d MapDecoder) error {
	var v T
	if err := c.codec(&v).VisitMap(d); err != nil {
		return err
	}
	*c.value = &v
	return nil
}

func (c PtrCodec[P, T, C]) Deserialize(d Decoder) error {
	return d.DecodeOption(c)
}

func (c PtrCodec[P, T, C]) Serialize(e Encoder) error {
	if *c.value == nil {
		return e.EncodeNil()
	}
	return c.codec((*T)(*c.value)).Serialize(e)
}

type SeqCodec[Q ~[]T, T any, C Codec[T]] struct {
	DefaultVisitor
	value *Q
}

func NewSeq[C Codec[T], Q ~[]T, T any](v *Q) SeqCodec[Q, T, C] {
	return SeqCodec[Q, T, C]{value: v}
}

func (c SeqCodec[Q, T, C]) InitCodec(v *Q) {
	c.value = v
}

func (c SeqCodec[Q, T, C]) VisitNil() error {
	*c.value = nil
	return nil
}

func (c SeqCodec[Q, T, C]) VisitSeq(seq SeqDecoder) error {
	var vals Q
	if len, ok := seq.Size(); ok {
		vals = make(Q, 0, len)
	}
	var codec C
	for {
		var v T
		ok, err := seq.NextElement(codec.New(&v))
		if err != nil {
			return err
		}
		if !ok {
			*c.value = vals
			return nil
		}
		vals = append(vals, v)
	}
}

func (c SeqCodec[Q, T, C]) Deserialize(d Decoder) error {
	return d.DecodeSeq(c)
}

func (c SeqCodec[Q, T, C]) Serialize(e Encoder) error {
	vals := *c.value

	enc, err := e.EncodeSeq(len(vals))
	if err != nil {
		return err
	}
	var codec C
	for _, v := range vals {
		if err := enc.EncodeElement(codec.New(&v)); err != nil {
			return err
		}
	}
	return enc.Close()
}

type MapCodec[M ~map[K]V, K constraints.Ordered, V any, CK Codec[K], CV Codec[V]] struct {
	DefaultVisitor
	value *M
}

func NewMap[CK Codec[K], CV Codec[V], M ~map[K]V, K constraints.Ordered, V any](v *M) MapCodec[M, K, V, CK, CV] {
	return MapCodec[M, K, V, CK, CV]{value: v}
}

func (c MapCodec[M, K, V, CK, CV]) VisitNil() error {
	*c.value = nil
	return nil
}

func (c MapCodec[M, K, V, CK, CV]) VisitMap(map_ MapDecoder) error {
	var m M
	if len, ok := map_.Size(); ok {
		m = make(M, len)
	} else {
		m = make(M)
	}
	var keyCodec CK
	var valueCodec CV
	for {
		var k K
		ok, err := map_.NextKey(keyCodec.New(&k))
		if err != nil {
			return err
		}
		if !ok {
			*c.value = m
			return nil
		}

		var v V
		if err = map_.NextValue(valueCodec.New(&v)); err != nil {
			return err
		}
		m[k] = v
	}
}

func (c MapCodec[M, K, V, CK, CV]) Deserialize(d Decoder) error {
	return d.DecodeMap(c)
}

func (c MapCodec[M, K, V, CK, CV]) Serialize(e Encoder) error {
	m := *c.value

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	enc, err := e.EncodeMap(len(m))
	if err != nil {
		return err
	}
	var keyCodec CK
	var valueCodec CV
	for _, k := range keys {
		if err := enc.EncodeKey(keyCodec.New(&k)); err != nil {
			return err
		}
		v := m[k]
		if err := enc.EncodeValue(valueCodec.New(&v)); err != nil {
			return err
		}
	}
	return enc.Close()
}
