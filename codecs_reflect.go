package codec

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/segmentio/asm/ascii"
	"github.com/segmentio/asm/keyset"
)

type codec interface {
	Visitor
	Deserializer
	Serializer

	new(v unsafe.Pointer) codec
}

type unsafeCodec struct {
	value unsafe.Pointer
}

func (unsafeCodec) VisitNil() error {
	return errors.New("unexpected nil")
}

func (unsafeCodec) VisitBool(v bool) error {
	return errors.New("unexpected bool")
}

func (unsafeCodec) VisitInt(v int) error {
	return errors.New("unexpected int")
}

func (unsafeCodec) VisitInt8(v int8) error {
	return errors.New("unexpected int8")
}

func (unsafeCodec) VisitInt16(v int16) error {
	return errors.New("unexpected int16")
}

func (unsafeCodec) VisitInt32(v int32) error {
	return errors.New("unexpected int32")
}

func (unsafeCodec) VisitInt64(v int64) error {
	return errors.New("unexpected int64")
}

func (unsafeCodec) VisitUint(v uint) error {
	return errors.New("unexpected uint")
}

func (unsafeCodec) VisitUint8(v uint8) error {
	return errors.New("unexpected uint8")
}

func (unsafeCodec) VisitUint16(v uint16) error {
	return errors.New("unexpected uint16")
}

func (unsafeCodec) VisitUint32(v uint32) error {
	return errors.New("unexpected uint32")
}

func (unsafeCodec) VisitUint64(v uint64) error {
	return errors.New("unexpected uint64")
}

func (unsafeCodec) VisitUintptr(v uintptr) error {
	return errors.New("unexpected uintptr")
}

func (unsafeCodec) VisitFloat32(v float32) error {
	return errors.New("unexpected float32")
}

func (unsafeCodec) VisitFloat64(v float64) error {
	return errors.New("unexpected float64")
}

func (unsafeCodec) VisitComplex64(v complex64) error {
	return errors.New("unexpected complex64")
}

func (unsafeCodec) VisitComplex128(v complex128) error {
	return errors.New("unexpected complex128")
}

func (unsafeCodec) VisitString(v string) error {
	return errors.New("unexpected string")
}

func (unsafeCodec) VisitBytes(v []byte) error {
	return errors.New("unexpected bytes")
}

func (unsafeCodec) VisitElem(d ElemDecoder) error {
	return errors.New("unexpected elem")
}

func (unsafeCodec) VisitSeq(d SeqDecoder) error {
	return errors.New("unexpected sequence")
}

func (unsafeCodec) VisitMap(d MapDecoder) error {
	return errors.New("unexpected map")
}

type nilCodec struct {
	unsafeCodec
}

func (c nilCodec) new(v unsafe.Pointer) codec {
	return nilCodec{}
}

func (c nilCodec) VisitNil() error {
	return nil
}

func (c nilCodec) Deserialize(d Decoder) error {
	return d.DecodeNil(c)
}

func (c nilCodec) Serialize(e Encoder) error {
	return e.EncodeNil()
}

type boolCodec struct{ unsafeCodec }

func (c boolCodec) new(v unsafe.Pointer) codec {
	return boolCodec{unsafeCodec: unsafeCodec{value: v}}
}

func (c boolCodec) VisitBool(v bool) error {
	*(*bool)(c.value) = v
	return nil
}

func (c boolCodec) Deserialize(d Decoder) error {
	return d.DecodeBool(c)
}

func (c boolCodec) Serialize(e Encoder) error {
	return e.EncodeBool(*(*bool)(c.value))
}

type intCodec struct{ unsafeCodec }

func (c intCodec) new(v unsafe.Pointer) codec {
	return intCodec{unsafeCodec: unsafeCodec{value: v}}
}

func (c intCodec) VisitInt(v int) error {
	*(*int)(c.value) = v
	return nil
}

func (c intCodec) Deserialize(d Decoder) error {
	return d.DecodeInt(c)
}

func (c intCodec) Serialize(e Encoder) error {
	return e.EncodeInt(*(*int)(c.value))
}

type int8Codec struct{ unsafeCodec }

func (c int8Codec) new(v unsafe.Pointer) codec {
	return int8Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c int8Codec) VisitInt8(v int8) error {
	*(*int8)(c.value) = v
	return nil
}

func (c int8Codec) Deserialize(d Decoder) error {
	return d.DecodeInt8(c)
}

func (c int8Codec) Serialize(e Encoder) error {
	return e.EncodeInt8(*(*int8)(c.value))
}

type int16Codec struct{ unsafeCodec }

func (c int16Codec) new(v unsafe.Pointer) codec {
	return int16Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c int16Codec) VisitInt16(v int16) error {
	*(*int16)(c.value) = v
	return nil
}

func (c int16Codec) Deserialize(d Decoder) error {
	return d.DecodeInt16(c)
}

func (c int16Codec) Serialize(e Encoder) error {
	return e.EncodeInt16(*(*int16)(c.value))
}

type int32Codec struct{ unsafeCodec }

func (c int32Codec) new(v unsafe.Pointer) codec {
	return int32Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c int32Codec) VisitInt32(v int32) error {
	*(*int32)(c.value) = v
	return nil
}

func (c int32Codec) Deserialize(d Decoder) error {
	return d.DecodeInt32(c)
}

func (c int32Codec) Serialize(e Encoder) error {
	return e.EncodeInt32(*(*int32)(c.value))
}

type int64Codec struct{ unsafeCodec }

func (c int64Codec) new(v unsafe.Pointer) codec {
	return int64Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c int64Codec) VisitInt64(v int64) error {
	*(*int64)(c.value) = v
	return nil
}

func (c int64Codec) Deserialize(d Decoder) error {
	return d.DecodeInt64(c)
}

func (c int64Codec) Serialize(e Encoder) error {
	return e.EncodeInt64(*(*int64)(c.value))
}

type uintCodec struct{ unsafeCodec }

func (c uintCodec) new(v unsafe.Pointer) codec {
	return uintCodec{unsafeCodec: unsafeCodec{value: v}}
}

func (c uintCodec) VisitUint(v uint) error {
	*(*uint)(c.value) = v
	return nil
}

func (c uintCodec) Deserialize(d Decoder) error {
	return d.DecodeUint(c)
}

func (c uintCodec) Serialize(e Encoder) error {
	return e.EncodeUint(*(*uint)(c.value))
}

type uint8Codec struct{ unsafeCodec }

func (c uint8Codec) new(v unsafe.Pointer) codec {
	return uint8Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c uint8Codec) VisitUint8(v uint8) error {
	*(*uint8)(c.value) = v
	return nil
}

func (c uint8Codec) Deserialize(d Decoder) error {
	return d.DecodeUint8(c)
}

func (c uint8Codec) Serialize(e Encoder) error {
	return e.EncodeUint8(*(*uint8)(c.value))
}

type uint16Codec struct{ unsafeCodec }

func (c uint16Codec) new(v unsafe.Pointer) codec {
	return uint16Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c uint16Codec) VisitUint16(v uint16) error {
	*(*uint16)(c.value) = v
	return nil
}

func (c uint16Codec) Deserialize(d Decoder) error {
	return d.DecodeUint16(c)
}

func (c uint16Codec) Serialize(e Encoder) error {
	return e.EncodeUint16(*(*uint16)(c.value))
}

type uint32Codec struct{ unsafeCodec }

func (c uint32Codec) new(v unsafe.Pointer) codec {
	return uint32Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c uint32Codec) VisitUint32(v uint32) error {
	*(*uint32)(c.value) = v
	return nil
}

func (c uint32Codec) Deserialize(d Decoder) error {
	return d.DecodeUint32(c)
}

func (c uint32Codec) Serialize(e Encoder) error {
	return e.EncodeUint32(*(*uint32)(c.value))
}

type uint64Codec struct{ unsafeCodec }

func (c uint64Codec) new(v unsafe.Pointer) codec {
	return uint64Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c uint64Codec) VisitUint64(v uint64) error {
	*(*uint64)(c.value) = v
	return nil
}

func (c uint64Codec) Deserialize(d Decoder) error {
	return d.DecodeUint64(c)
}

func (c uint64Codec) Serialize(e Encoder) error {
	return e.EncodeUint64(*(*uint64)(c.value))
}

type uintptrCodec struct{ unsafeCodec }

func (c uintptrCodec) new(v unsafe.Pointer) codec {
	return uintptrCodec{unsafeCodec: unsafeCodec{value: v}}
}

func (c uintptrCodec) VisitUintptr(v uintptr) error {
	*(*uintptr)(c.value) = v
	return nil
}

func (c uintptrCodec) Deserialize(d Decoder) error {
	return d.DecodeUintptr(c)
}

func (c uintptrCodec) Serialize(e Encoder) error {
	return e.EncodeUintptr(*(*uintptr)(c.value))
}

type float32Codec struct{ unsafeCodec }

func (c float32Codec) new(v unsafe.Pointer) codec {
	return float32Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c float32Codec) VisitFloat32(v float32) error {
	*(*float32)(c.value) = v
	return nil
}

func (c float32Codec) Deserialize(d Decoder) error {
	return d.DecodeFloat32(c)
}

func (c float32Codec) Serialize(e Encoder) error {
	return e.EncodeFloat32(*(*float32)(c.value))
}

type float64Codec struct{ unsafeCodec }

func (c float64Codec) new(v unsafe.Pointer) codec {
	return float64Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c float64Codec) VisitFloat64(v float64) error {
	*(*float64)(c.value) = v
	return nil
}

func (c float64Codec) Deserialize(d Decoder) error {
	return d.DecodeFloat64(c)
}

func (c float64Codec) Serialize(e Encoder) error {
	return e.EncodeFloat64(*(*float64)(c.value))
}

type complex64Codec struct{ unsafeCodec }

func (c complex64Codec) new(v unsafe.Pointer) codec {
	return complex64Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c complex64Codec) VisitComplex64(v complex64) error {
	*(*complex64)(c.value) = v
	return nil
}

func (c complex64Codec) Deserialize(d Decoder) error {
	return d.DecodeComplex64(c)
}

func (c complex64Codec) Serialize(e Encoder) error {
	return e.EncodeComplex64(*(*complex64)(c.value))
}

type complex128Codec struct{ unsafeCodec }

func (c complex128Codec) new(v unsafe.Pointer) codec {
	return complex128Codec{unsafeCodec: unsafeCodec{value: v}}
}

func (c complex128Codec) VisitComplex128(v complex128) error {
	*(*complex128)(c.value) = v
	return nil
}

func (c complex128Codec) Deserialize(d Decoder) error {
	return d.DecodeComplex128(c)
}

func (c complex128Codec) Serialize(e Encoder) error {
	return e.EncodeComplex128(*(*complex128)(c.value))
}

type stringCodec struct{ unsafeCodec }

func (c stringCodec) new(v unsafe.Pointer) codec {
	return stringCodec{unsafeCodec: unsafeCodec{value: v}}
}

func (c stringCodec) VisitString(v string) error {
	*(*string)(c.value) = v
	return nil
}

func (c stringCodec) Deserialize(d Decoder) error {
	return d.DecodeString(c)
}

func (c stringCodec) Serialize(e Encoder) error {
	return e.EncodeString(*(*string)(c.value))
}

type bytesCodec struct{ unsafeCodec }

func (c bytesCodec) new(v unsafe.Pointer) codec {
	return bytesCodec{unsafeCodec: unsafeCodec{value: v}}
}

func (c bytesCodec) VisitBytes(v []byte) error {
	*(*[]byte)(c.value) = v
	return nil
}

func (c bytesCodec) Deserialize(d Decoder) error {
	return d.DecodeBytes(c)
}

func (c bytesCodec) Serialize(e Encoder) error {
	return e.EncodeBytes(*(*[]byte)(c.value))
}

type ptrType struct {
	t    reflect.Type
	elem codec
}

type ptrCodec struct {
	unsafeCodec
	*ptrType
}

func (c ptrCodec) new(v unsafe.Pointer) codec {
	return ptrCodec{
		unsafeCodec: unsafeCodec{value: v},
		ptrType:     c.ptrType,
	}
}

func (c ptrCodec) VisitNil() error {
	*(*unsafe.Pointer)(c.value) = nil
	return nil
}

func (c ptrCodec) VisitElem(d ElemDecoder) error {
	v := reflect.New(c.t.Elem())
	p := v.UnsafePointer()
	if err := d.Element(v.Interface(), c.elem.new(p)); err != nil {
		return err
	}
	*(*unsafe.Pointer)(c.value) = p
	return nil
}

func (c ptrCodec) Deserialize(d Decoder) error {
	return d.DecodePtr(c)
}

func (c ptrCodec) Serialize(e Encoder) error {
	if (*unsafe.Pointer)(c.value) == nil {
		return e.EncodeNil()
	}
	return c.elem.new(c.value).Serialize(e)
}

type arrayType struct {
	n    int
	t    reflect.Type
	elem codec
}

type arrayCodec struct {
	unsafeCodec
	*arrayType
}

func (c arrayCodec) new(v unsafe.Pointer) codec {
	return arrayCodec{
		unsafeCodec: unsafeCodec{value: v},
		arrayType:   c.arrayType,
	}
}

func (c arrayCodec) VisitSeq(seq SeqDecoder) error {
	vals := reflect.NewAt(c.t, c.value)
	for i := 0; i < c.n; i++ {
		elem := vals.Index(i).Addr()
		ok, err := seq.NextElement(elem.Interface(), c.elem.new(elem.UnsafePointer()))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}

	discard := reflect.New(c.t.Elem())
	discardCodec := c.elem.new(discard.UnsafePointer())
	for {
		ok, err := seq.NextElement(discard.Interface(), discardCodec)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
}

func (c arrayCodec) Deserialize(d Decoder) error {
	return d.DecodeSeq(c)
}

func (c arrayCodec) Serialize(e Encoder) error {
	vals := reflect.NewAt(c.t, c.value)

	enc, err := e.EncodeSeq(vals.Len())
	if err != nil {
		return err
	}
	for i, n := 0, vals.Len(); i < n; i++ {
		elem := vals.Index(i)
		if err := enc.EncodeElement(elem, c.elem.new(elem.Addr().UnsafePointer())); err != nil {
			return err
		}
	}
	return enc.Close()
}

type sliceType struct {
	size uintptr
	t    reflect.Type
	elem codec
}

type sliceCodec struct {
	unsafeCodec
	*sliceType
}

func (c sliceCodec) new(v unsafe.Pointer) codec {
	return sliceCodec{
		unsafeCodec: unsafeCodec{value: v},
		sliceType:   c.sliceType,
	}
}

func (c sliceCodec) VisitSeq(seq SeqDecoder) error {
	s := (*slice)(c.value)
	for {
		if s.len == s.cap {
			cap := s.cap
			if cap == 0 {
				cap = 10
			} else {
				cap *= 2
			}
			*s = extendSlice(c.t, s, cap)
		}

		p := unsafe.Pointer(uintptr(s.data) + (uintptr(s.len) * c.size))
		elem := reflect.NewAt(c.t.Elem(), p)

		ok, err := seq.NextElement(elem.Interface(), c.elem.new(p))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		s.len++
	}
}

func (c sliceCodec) Deserialize(d Decoder) error {
	return d.DecodeSeq(c)
}

func (c sliceCodec) Serialize(e Encoder) error {
	s := (*slice)(c.value)

	enc, err := e.EncodeSeq(s.len)
	if err != nil {
		return err
	}
	for i, n, data := 0, s.len, s.data; i < n; i, data = i+1, unsafe.Pointer(uintptr(data)+c.size) {
		elem := reflect.NewAt(c.t.Elem(), data).Elem()
		if err := enc.EncodeElement(elem.Interface(), c.elem.new(unsafe.Pointer(data))); err != nil {
			return err
		}
	}
	return enc.Close()
}

type mapType struct {
	t        reflect.Type
	kt       reflect.Type
	vt       reflect.Type
	kz       reflect.Value
	vz       reflect.Value
	kc       codec
	vc       codec
	sortKeys sortFunc
}

type mapCodec struct {
	unsafeCodec
	*mapType
}

func (c mapCodec) new(v unsafe.Pointer) codec {
	return mapCodec{
		unsafeCodec: unsafeCodec{value: v},
		mapType:     c.mapType,
	}
}

func (c mapCodec) VisitMap(map_ MapDecoder) error {
	var m reflect.Value
	if len, ok := map_.Size(); ok {
		m = reflect.MakeMapWithSize(c.t, len)
	} else {
		m = reflect.MakeMap(c.t)
	}

	k := reflect.New(c.kt).Elem()
	v := reflect.New(c.vt).Elem()
	kptr := k.Addr().UnsafePointer()
	vptr := v.Addr().UnsafePointer()
	for {
		k.Set(c.kz)
		v.Set(c.vz)

		ok, err := map_.NextKey(k.Interface(), c.kc.new(kptr))
		if err != nil {
			return err
		}
		if !ok {
			reflect.NewAt(c.t, c.value).Elem().Set(m)
			return nil
		}

		if err = map_.NextValue(v.Interface(), c.vc.new(vptr)); err != nil {
			return err
		}
		m.SetMapIndex(k, v)
	}
}

func (c mapCodec) Deserialize(d Decoder) error {
	return d.DecodeMap(c)
}

func (c mapCodec) Serialize(e Encoder) error {
	p := c.value
	m := reflect.NewAt(c.t, unsafe.Pointer(&p)).Elem()

	keys := m.MapKeys()
	c.sortKeys(keys)

	enc, err := e.EncodeMap(m.Len())
	if err != nil {
		return err
	}
	for _, k := range keys {
		if err := enc.EncodeKey(k.Interface(), c.kc.new((*iface)(unsafe.Pointer(&k)).ptr)); err != nil {
			return err
		}
		v := m.MapIndex(k)
		if err := enc.EncodeValue(v.Interface(), c.vc.new((*iface)(unsafe.Pointer(&v)).ptr)); err != nil {
			return err
		}
	}
	return enc.Close()
}

type structType struct {
	name        string
	fields      []structField
	fieldsIndex map[string]*structField
	ficaseIndex map[string]*structField
	keyset      []byte
	typ         reflect.Type
	inlined     bool
}

type structField struct {
	codec     codec
	offset    uintptr
	empty     emptyFunc
	tag       bool
	omitempty bool
	embedded  *embeddedStructField
	name      string
	typ       reflect.Type
	zero      reflect.Value
	index     int
}

type embeddedStructField struct {
	unexported bool
	offset     uintptr
}

type structCodec struct {
	unsafeCodec
	*structType
}

func (c structCodec) new(v unsafe.Pointer) codec {
	return structCodec{
		unsafeCodec: unsafeCodec{value: v},
		structType:  c.structType,
	}
}

func (c structCodec) VisitMap(map_ MapDecoder) error {
	var keybuf []byte
	for {
		var k string
		ok, err := map_.NextKey(&k, NewString(&k))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		var f *structField
		if len(c.keyset) != 0 {
			if n := keyset.Lookup(c.keyset, stringBytes(k)); n < len(c.fields) {
				f = &c.fields[n]
			} else {
				f = c.fieldsIndex[k]
			}
		}
		if f == nil {
			// TODO: disallow case-insensitive match
			f = c.ficaseIndex[appendToLower(keybuf[:0], k)]
		}

		if f == nil {
			// TODO: disallow unknown fields
			if err := map_.NextValue(nil, SkipCodec{}); err != nil {
				return err
			}
			continue
		}

		v := unsafe.Pointer(uintptr(c.value) + f.offset)
		if f.embedded != nil {
			if p := (*unsafe.Pointer)(v); *p == nil {
				if f.embedded.unexported {
					return fmt.Errorf("codec: cannot set embedded pointer to unexported struct: %s", f.typ)
				}
				t := reflect.New(f.typ).UnsafePointer()
				*p = t
				v = unsafe.Pointer(uintptr(t) + f.embedded.offset)
			}
		}

		fv := reflect.NewAt(f.typ, v)
		if err := map_.NextValue(fv.Interface(), f.codec.new(v)); err != nil {
			return err
		}
	}
}

func (c structCodec) Deserialize(d Decoder) error {
	return d.DecodeStruct(c.name, c)
}

func (c structCodec) Serialize(e Encoder) error {
	enc, err := e.EncodeStruct(c.name)
	if err != nil {
		return err
	}

	for i := range c.fields {
		f := &c.fields[i]
		v := unsafe.Pointer(uintptr(c.value) + f.offset)

		if f.omitempty && f.empty(v) {
			continue
		}

		if f.embedded != nil {
			p := *(*unsafe.Pointer)(v)
			if p == nil {
				continue
			}
			v = unsafe.Pointer(uintptr(p) + f.embedded.offset)
		}

		fv := reflect.NewAt(f.typ, v).Elem()
		if err := enc.EncodeField(f.name, fv.Interface(), f.codec.new(v)); err != nil {
			return err
		}
	}

	return enc.Close()
}

type unsupportedTypeCodec struct {
	unsafeCodec
	t reflect.Type
}

func (c unsupportedTypeCodec) new(v unsafe.Pointer) codec {
	return c
}

func (c unsupportedTypeCodec) Deserialize(d Decoder) error {
	return &UnsupportedTypeError{Type: c.t}
}

func (c unsupportedTypeCodec) Serialize(e Encoder) error {
	return &UnsupportedTypeError{Type: c.t}
}

type serializerCodec struct {
	unsafeCodec
	t    reflect.Type
	next codec
}

func (c serializerCodec) new(v unsafe.Pointer) codec {
	return serializerCodec{
		unsafeCodec: unsafeCodec{value: v},
		t:           c.t,
		next:        c.next,
	}
}

func (c serializerCodec) Deserialize(d Decoder) error {
	return c.next.Deserialize(d)
}

func (c serializerCodec) Serialize(e Encoder) error {
	return reflect.NewAt(c.t, c.value).Interface().(Serializer).Serialize(e)
}

type deserializerCodec struct {
	unsafeCodec
	t    reflect.Type
	next codec
}

func (c deserializerCodec) new(v unsafe.Pointer) codec {
	return deserializerCodec{
		unsafeCodec: unsafeCodec{value: v},
		t:           c.t,
		next:        c.next,
	}
}

func (c deserializerCodec) Deserialize(d Decoder) error {
	return reflect.NewAt(c.t, c.value).Interface().(Deserializer).Deserialize(d)
}

func (c deserializerCodec) Serialize(e Encoder) error {
	return c.next.Serialize(e)
}

func stringBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func appendRune(b []byte, r rune) []byte {
	n := len(b)
	b = append(b, 0, 0, 0, 0)
	return b[:n+utf8.EncodeRune(b[n:], r)]
}

func appendToLower(b []byte, ss string) string {
	s := stringBytes(ss)

	if ascii.Valid(s) { // fast path for ascii strings
		i := 0

		for j := range s {
			c := s[j]

			if 'A' <= c && c <= 'Z' {
				b = append(b, s[i:j]...)
				b = append(b, c+('a'-'A'))
				i = j + 1
			}
		}

		b = append(b, s[i:]...)
		return unsafe.String(unsafe.SliceData(b), len(b))
	}

	for _, r := range string(s) {
		b = appendRune(b, foldRune(r))
	}

	return unsafe.String(unsafe.SliceData(b), len(b))
}

func foldRune(r rune) rune {
	if r = unicode.SimpleFold(r); 'A' <= r && r <= 'Z' {
		r = r + ('a' - 'A')
	}
	return r
}

func extendSlice(t reflect.Type, s *slice, n int) slice {
	arrayType := reflect.ArrayOf(n, t.Elem())
	arrayData := reflect.New(arrayType)
	reflect.Copy(arrayData.Elem(), reflect.NewAt(t, unsafe.Pointer(s)).Elem())
	return slice{
		data: unsafe.Pointer(arrayData.Pointer()),
		len:  s.len,
		cap:  n,
	}
}
