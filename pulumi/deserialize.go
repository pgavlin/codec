package pulumi

import (
	"fmt"

	"github.com/pgavlin/codec"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Deserializer struct {
	v *resource.PropertyValue
}

func Deserialize(dec codec.Decoder) (v resource.PropertyValue, err error) {
	err = NewDeserializer(&v).Deserialize(dec)
	return
}

func NewDeserializer(v *resource.PropertyValue) Deserializer {
	return Deserializer{v: v}
}

func (d Deserializer) VisitNil() error                    { return NewEncoder(d.v).EncodeNil() }
func (d Deserializer) VisitBool(v bool) error             { return NewEncoder(d.v).EncodeBool(v) }
func (d Deserializer) VisitInt(v int) error               { return NewEncoder(d.v).EncodeInt(v) }
func (d Deserializer) VisitInt8(v int8) error             { return NewEncoder(d.v).EncodeInt8(v) }
func (d Deserializer) VisitInt16(v int16) error           { return NewEncoder(d.v).EncodeInt16(v) }
func (d Deserializer) VisitInt32(v int32) error           { return NewEncoder(d.v).EncodeInt32(v) }
func (d Deserializer) VisitInt64(v int64) error           { return NewEncoder(d.v).EncodeInt64(v) }
func (d Deserializer) VisitUint(v uint) error             { return NewEncoder(d.v).EncodeUint(v) }
func (d Deserializer) VisitUint8(v uint8) error           { return NewEncoder(d.v).EncodeUint8(v) }
func (d Deserializer) VisitUint16(v uint16) error         { return NewEncoder(d.v).EncodeUint16(v) }
func (d Deserializer) VisitUint32(v uint32) error         { return NewEncoder(d.v).EncodeUint32(v) }
func (d Deserializer) VisitUint64(v uint64) error         { return NewEncoder(d.v).EncodeUint64(v) }
func (d Deserializer) VisitUintptr(v uintptr) error       { return NewEncoder(d.v).EncodeUintptr(v) }
func (d Deserializer) VisitFloat32(v float32) error       { return NewEncoder(d.v).EncodeFloat32(v) }
func (d Deserializer) VisitFloat64(v float64) error       { return NewEncoder(d.v).EncodeFloat64(v) }
func (d Deserializer) VisitComplex64(v complex64) error   { return NewEncoder(d.v).EncodeComplex64(v) }
func (d Deserializer) VisitComplex128(v complex128) error { return NewEncoder(d.v).EncodeComplex128(v) }
func (d Deserializer) VisitBytes(v []byte) error          { return NewEncoder(d.v).EncodeBytes(v) }

func (d Deserializer) VisitString(v string) error {
	if v == unknownRepr {
		*d.v = resource.MakeComputed(resource.NewStringProperty(""))
		return nil
	}
	*d.v = resource.NewStringProperty(v)
	return nil
}

func (d Deserializer) VisitElem(elem codec.ElemDecoder) error {
	return elem.Element(nil, d)
}

func (d Deserializer) VisitSeq(seq codec.SeqDecoder) error {
	var vals []resource.PropertyValue
	if len, ok := seq.Size(); ok {
		vals = make([]resource.PropertyValue, 0, len)
	}
	for {
		var v resource.PropertyValue
		ok, err := seq.NextElement(&v, NewDeserializer(&v))
		if err != nil {
			return err
		}
		if !ok {
			*d.v = resource.NewArrayProperty(vals)
			return nil
		}
		vals = append(vals, v)
	}
}

func (d Deserializer) VisitMap(map_ codec.MapDecoder) error {
	var m resource.PropertyMap
	if len, ok := map_.Size(); ok {
		m = make(resource.PropertyMap, len)
	} else {
		m = make(resource.PropertyMap)
	}
	for {
		var k resource.PropertyKey
		ok, err := map_.NextKey(&k, codec.NewString(&k))
		if err != nil {
			return err
		}
		if !ok {
			return d.visitSig(m)
		}

		var v resource.PropertyValue
		if err := map_.NextValue(&v, NewDeserializer(&v)); err != nil {
			return err
		}
		m[k] = v
	}
}

func (d Deserializer) visitSig(m resource.PropertyMap) error {
	sig, ok := m[resource.PropertyKey(resource.SigKey)]
	if !ok {
		*d.v = resource.NewObjectProperty(m)
		return nil
	}

	if !sig.IsString() {
		return fmt.Errorf("unrecognized signature of type %v", sig.TypeString())
	}

	switch sig.StringValue() {
	case resource.ArchiveSig:
		var obj map[string]any
		if err := codec.GetDeserializer(&obj).Deserialize(NewDecoder(resource.NewObjectProperty(m))); err != nil {
			return err
		}
		archive, _, err := resource.DeserializeArchive(obj)
		if err != nil {
			return err
		}
		*d.v = resource.NewArchiveProperty(archive)
		return nil
	case resource.AssetSig:
		var obj map[string]any
		if err := codec.GetDeserializer(&obj).Deserialize(NewDecoder(resource.NewObjectProperty(m))); err != nil {
			return err
		}
		asset, _, err := resource.DeserializeAsset(obj)
		if err != nil {
			return err
		}
		*d.v = resource.NewAssetProperty(asset)
		return nil
	case resource.ResourceReferenceSig:
		panic("todo")
	case resource.SecretSig:
		panic("todo")
	default:
		return fmt.Errorf("unrecognized signature %q", sig.StringValue())
	}
}

func (d Deserializer) Deserialize(dec codec.Decoder) error {
	return dec.DecodeAny(d)
}
