package pulumi

import (
	"github.com/pgavlin/codec"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const unknownRepr = "04da6b54-80e4-46f7-96ec-b56ff0331ba9"

type Serializer struct {
	v resource.PropertyValue
}

func Serialize(v resource.PropertyValue, enc codec.Encoder) error {
	return NewSerializer(v).Serialize(enc)
}

func NewSerializer(v resource.PropertyValue) Serializer {
	return Serializer{v: v}
}

func (s Serializer) Serialize(enc codec.Encoder) error {
	switch {
	case s.v.IsNull():
		return enc.EncodeNil()
	case s.v.IsComputed(), s.v.IsOutput():
		return enc.EncodeString(unknownRepr)
	case s.v.IsBool():
		return enc.EncodeBool(s.v.BoolValue())
	case s.v.IsNumber():
		return enc.EncodeFloat64(s.v.NumberValue())
	case s.v.IsString():
		return enc.EncodeString(s.v.StringValue())
	case s.v.IsArchive():
		return codec.GetSerializer(s.v.ArchiveValue().Serialize()).Serialize(enc)
	case s.v.IsAsset():
		return codec.GetSerializer(s.v.AssetValue().Serialize()).Serialize(enc)
	case s.v.IsResourceReference():
		panic("todo")
	case s.v.IsSecret():
		panic("todo")
	case s.v.IsArray():
		vals := s.v.ArrayValue()
		seq, err := enc.EncodeSeq(len(vals))
		if err != nil {
			return err
		}
		for _, v := range vals {
			if err := seq.EncodeElement(v, NewSerializer(v)); err != nil {
				return err
			}
		}
		return seq.Close()
	case s.v.IsObject():
		obj := s.v.ObjectValue()

		map_, err := enc.EncodeMap(len(obj))
		if err != nil {
			return err
		}
		keys := maps.Keys(obj)
		slices.Sort(keys)
		for _, k := range keys {
			if err := map_.EncodeKey(k, codec.NewString(&k)); err != nil {
				return err
			}
			if err := map_.EncodeValue(obj[k], NewSerializer(obj[k])); err != nil {
				return err
			}
		}
		return map_.Close()
	default:
		panic("unsupported")
	}

}
