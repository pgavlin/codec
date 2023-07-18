package pulumi

import (
	"fmt"
	"reflect"

	"github.com/pgavlin/codec"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Decoder struct {
	v resource.PropertyValue
}

func Decode[T any](v resource.PropertyValue) (t T, err error) {
	err = codec.GetDeserializer(&t).Deserialize(NewDecoder(v))
	return
}

func NewDecoder(v resource.PropertyValue) Decoder {
	return Decoder{v: v}
}

func (d Decoder) DecodeAny(v codec.Visitor) error {
	switch {
	case d.v.IsNull():
		return v.VisitNil()
	case d.v.IsBool():
		return v.VisitBool(d.v.BoolValue())
	case d.v.IsNumber():
		return v.VisitFloat64(d.v.NumberValue())
	case d.v.IsString():
		return v.VisitString(d.v.StringValue())
	case d.v.IsArray():
		return v.VisitSeq(&SeqDecoder{v: d.v.ArrayValue()})
	case d.v.IsObject():
		dv := d.v.ObjectValue()
		return v.VisitMap(&MapDecoder{size: len(dv), iter: reflect.ValueOf(dv).MapRange()})
	default:
		return fmt.Errorf("cannot decode %v", d.v.TypeString())
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
func (d Decoder) DecodeSeq(v codec.Visitor) error                 { return d.DecodeAny(v) }
func (d Decoder) DecodeMap(v codec.Visitor) error                 { return d.DecodeAny(v) }
func (d Decoder) DecodeStruct(name string, v codec.Visitor) error { return d.DecodeAny(v) }

func (d Decoder) DecodePtr(v codec.Visitor) error {
	if d.v.IsNull() {
		return v.VisitNil()
	}
	return v.VisitElem(ElemDecoder{d})
}

func (d Decoder) decode(v any, ds codec.Deserializer) error {
	switch v := v.(type) {
	case *resource.PropertyValue:
		*v = d.v
		return nil
	case **resource.Asset:
		if !d.v.IsAsset() {
			return fmt.Errorf("expected an asset")
		}
		*v = d.v.AssetValue()
		return nil
	case **resource.Archive:
		if !d.v.IsArchive() {
			return fmt.Errorf("expected an archive")
		}
		*v = d.v.ArchiveValue()
		return nil
	case *resource.ResourceReference:
		if !d.v.IsResourceReference() {
			return fmt.Errorf("expected a resource reference")
		}
		*v = d.v.ResourceReferenceValue()
		return nil
	case Unmarshaler:
		return v.UnmarshalPropertyValue(d.v)
	default:
		return ds.Deserialize(d)
	}
}

type ElemDecoder struct {
	d Decoder
}

func (d ElemDecoder) Element(v any, ds codec.Deserializer) error {
	return d.d.decode(v, ds)
}

type SeqDecoder struct {
	v []resource.PropertyValue
}

func (d *SeqDecoder) Size() (int, bool) {
	return len(d.v), true
}

func (d *SeqDecoder) NextElement(x any, ds codec.Deserializer) (bool, error) {
	if len(d.v) == 0 {
		return false, nil
	}
	v := d.v[0]
	d.v = d.v[1:]
	return true, Decoder{v}.decode(x, ds)
}

type MapDecoder struct {
	size int
	iter *reflect.MapIter
}

func (d *MapDecoder) Size() (int, bool) {
	return d.size, true
}

func (d *MapDecoder) NextKey(k any, ds codec.Deserializer) (bool, error) {
	if !d.iter.Next() {
		return false, nil
	}
	return true, Decoder{resource.NewStringProperty(d.iter.Key().String())}.decode(k, ds)
}

func (d *MapDecoder) NextValue(v any, ds codec.Deserializer) error {
	return Decoder{d.iter.Value().Interface().(resource.PropertyValue)}.decode(v, ds)
}
