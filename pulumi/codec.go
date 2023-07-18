package pulumi

import (
	"github.com/pgavlin/codec"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

type Unmarshaler interface {
	UnmarshalPropertyValue(pv resource.PropertyValue) error
}

type Marshaler interface {
	MarshalPropertyValue() (resource.PropertyValue, error)
}

type Value[T any] struct {
	unknown bool
	secret  bool
	t       T
}

func NewUnknown[T any]() Value[T] {
	return Value[T]{unknown: true}
}

func NewSecret[T any](v T) Value[T] {
	return Value[T]{secret: true, t: v}
}

func NewSecretUnknown[T any]() Value[T] {
	return Value[T]{unknown: true, secret: true}
}

func NewValue[T any](v T) Value[T] {
	return Value[T]{t: v}
}

func (v Value[T]) IsUnknown() bool {
	return v.unknown
}

func (v Value[T]) IsSecret() bool {
	return v.secret
}

func (v Value[T]) Value() T {
	return v.t
}

func (v *Value[T]) UnmarshalPropertyValue(pv resource.PropertyValue) error {
	for pv.IsSecret() {
		v.secret = true
		pv = pv.SecretValue().Element
	}
	if pv.IsComputed() || pv.IsOutput() {
		v.unknown = true
		return nil
	}
	return codec.GetDeserializer(&v.t).Deserialize(NewDecoder(pv))
}

func (v Value[T]) MarshalPropertyValue() (pv resource.PropertyValue, err error) {
	if v.unknown {
		pv = resource.MakeComputed(resource.NewStringProperty(""))
	} else if err = codec.GetSerializer(v.t).Serialize(NewEncoder(&pv)); err != nil {
		return
	}
	if v.secret {
		pv = resource.MakeSecret(pv)
	}
	return
}
