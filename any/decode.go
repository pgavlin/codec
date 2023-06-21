package any

import (
	"reflect"

	"github.com/pgavlin/codec"
)

type Decoder struct {
	v any
}

func Decode[T any](v any) (t T, err error) {
	err = codec.GetDeserializer(&t).Deserialize(NewDecoder(v))
	return
}

func NewDecoder(v any) Decoder {
	return Decoder{v: v}
}

func (d Decoder) DecodeAny(v codec.Visitor) error {
	switch dv := d.v.(type) {
	case nil:
		return v.VisitNil()
	case bool:
		return v.VisitBool(dv)
	case int:
		return v.VisitInt(int(dv))
	case int8:
		return v.VisitInt8(int8(dv))
	case int16:
		return v.VisitInt16(int16(dv))
	case int32:
		return v.VisitInt32(int32(dv))
	case int64:
		return v.VisitInt64(dv)
	case uint:
		return v.VisitUint(uint(dv))
	case uint8:
		return v.VisitUint8(uint8(dv))
	case uint16:
		return v.VisitUint16(uint16(dv))
	case uint32:
		return v.VisitUint32(uint32(dv))
	case uint64:
		return v.VisitUint64(dv)
	case uintptr:
		return v.VisitUintptr(uintptr(dv))
	case float32:
		return v.VisitFloat32(float32(dv))
	case float64:
		return v.VisitFloat64(dv)
	case complex64:
		return v.VisitComplex64(complex64(dv))
	case complex128:
		return v.VisitComplex128(dv)
	case string:
		return v.VisitString(dv)
	case []any:
		return v.VisitSeq(&SeqDecoder{v: dv})
	case map[string]any:
		return v.VisitMap(&MapDecoder{size: len(dv), iter: reflect.ValueOf(dv).MapRange()})
	default:
		panic("unsupported")
	}
}

func (d Decoder) DecodeNil(v codec.Visitor) error                 { return d.DecodeAny(v) }
func (d Decoder) DecodeBool(v codec.Visitor) error                { return d.DecodeAny(v) }
func (d Decoder) DecodeInt(v codec.Visitor) error                 { return d.DecodeAny(v) }
func (d Decoder) DecodeInt8(v codec.Visitor) error                { return d.DecodeAny(v) }
func (d Decoder) DecodeInt16(v codec.Visitor) error               { return d.DecodeAny(v) }
func (d Decoder) DecodeInt32(v codec.Visitor) error               { return d.DecodeAny(v) }
func (d Decoder) DecodeInt64(v codec.Visitor) error               { return d.DecodeAny(v) }
func (d Decoder) DecodeUint(v codec.Visitor) error                { return d.DecodeAny(v) }
func (d Decoder) DecodeUint8(v codec.Visitor) error               { return d.DecodeAny(v) }
func (d Decoder) DecodeUint16(v codec.Visitor) error              { return d.DecodeAny(v) }
func (d Decoder) DecodeUint32(v codec.Visitor) error              { return d.DecodeAny(v) }
func (d Decoder) DecodeUint64(v codec.Visitor) error              { return d.DecodeAny(v) }
func (d Decoder) DecodeUintptr(v codec.Visitor) error             { return d.DecodeAny(v) }
func (d Decoder) DecodeFloat32(v codec.Visitor) error             { return d.DecodeAny(v) }
func (d Decoder) DecodeFloat64(v codec.Visitor) error             { return d.DecodeAny(v) }
func (d Decoder) DecodeComplex64(v codec.Visitor) error           { return d.DecodeAny(v) }
func (d Decoder) DecodeComplex128(v codec.Visitor) error          { return d.DecodeAny(v) }
func (d Decoder) DecodeString(v codec.Visitor) error              { return d.DecodeAny(v) }
func (d Decoder) DecodeBytes(v codec.Visitor) error               { return d.DecodeAny(v) }
func (d Decoder) DecodeOption(v codec.Visitor) error              { return d.DecodeAny(v) }
func (d Decoder) DecodeSeq(v codec.Visitor) error                 { return d.DecodeAny(v) }
func (d Decoder) DecodeMap(v codec.Visitor) error                 { return d.DecodeAny(v) }
func (d Decoder) DecodeStruct(name string, v codec.Visitor) error { return d.DecodeAny(v) }

type SeqDecoder struct {
	v []any
}

func (d *SeqDecoder) Size() (int, bool) {
	return len(d.v), true
}

func (d *SeqDecoder) NextElement(ds codec.Deserializer) (bool, error) {
	if len(d.v) == 0 {
		return false, nil
	}
	v := d.v[0]
	d.v = d.v[1:]
	return true, ds.Deserialize(Decoder{v: v})
}

type MapDecoder struct {
	size int
	iter *reflect.MapIter
}

func (d *MapDecoder) Size() (int, bool) {
	return d.size, true
}

func (d *MapDecoder) NextKey(ds codec.Deserializer) (bool, error) {
	if !d.iter.Next() {
		return false, nil
	}
	return true, ds.Deserialize(NewDecoder(d.iter.Key().String()))
}

func (d *MapDecoder) NextValue(ds codec.Deserializer) error {
	return ds.Deserialize(NewDecoder(d.iter.Value().Interface()))
}
