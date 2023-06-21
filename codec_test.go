package codec

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	v reflect.Value
}

func data(v any) testData {
	return testData{reflect.ValueOf(v)}
}

func (d testData) DecodeAny(v Visitor) error {
	switch d.v.Kind() {
	case reflect.Bool:
		return v.VisitBool(d.v.Bool())
	case reflect.Int:
		return v.VisitInt(int(d.v.Int()))
	case reflect.Int8:
		return v.VisitInt8(int8(d.v.Int()))
	case reflect.Int16:
		return v.VisitInt16(int16(d.v.Int()))
	case reflect.Int32:
		return v.VisitInt32(int32(d.v.Int()))
	case reflect.Int64:
		return v.VisitInt64(d.v.Int())
	case reflect.Uint:
		return v.VisitUint(uint(d.v.Uint()))
	case reflect.Uint8:
		return v.VisitUint8(uint8(d.v.Uint()))
	case reflect.Uint16:
		return v.VisitUint16(uint16(d.v.Uint()))
	case reflect.Uint32:
		return v.VisitUint32(uint32(d.v.Uint()))
	case reflect.Uint64:
		return v.VisitUint64(d.v.Uint())
	case reflect.Uintptr:
		return v.VisitUintptr(uintptr(d.v.Uint()))
	case reflect.Float32:
		return v.VisitFloat32(float32(d.v.Float()))
	case reflect.Float64:
		return v.VisitFloat64(d.v.Float())
	case reflect.Complex64:
		return v.VisitComplex64(complex64(d.v.Complex()))
	case reflect.Complex128:
		return v.VisitComplex128(d.v.Complex())
	case reflect.String:
		return v.VisitString(d.v.String())
	case reflect.Pointer:
		if d.v.IsNil() {
			return v.VisitNil()
		}
		return testData{d.v.Elem()}.DecodeAny(v)
	case reflect.Array, reflect.Slice:
		return v.VisitSeq(&testSeqDecoder{d.v, d.v.Len(), 0})
	case reflect.Map:
		return v.VisitMap(&testMapDecoder{d.v.MapRange()})
	case reflect.Struct:
		return v.VisitMap(&testStructDecoder{d.v, reflect.VisibleFields(d.v.Type()), 0})
	default:
		panic("unsupported")
	}
}

func (d testData) DecodeNil(v Visitor) error                 { return d.DecodeAny(v) }
func (d testData) DecodeBool(v Visitor) error                { return d.DecodeAny(v) }
func (d testData) DecodeInt(v Visitor) error                 { return d.DecodeAny(v) }
func (d testData) DecodeInt8(v Visitor) error                { return d.DecodeAny(v) }
func (d testData) DecodeInt16(v Visitor) error               { return d.DecodeAny(v) }
func (d testData) DecodeInt32(v Visitor) error               { return d.DecodeAny(v) }
func (d testData) DecodeInt64(v Visitor) error               { return d.DecodeAny(v) }
func (d testData) DecodeUint(v Visitor) error                { return d.DecodeAny(v) }
func (d testData) DecodeUint8(v Visitor) error               { return d.DecodeAny(v) }
func (d testData) DecodeUint16(v Visitor) error              { return d.DecodeAny(v) }
func (d testData) DecodeUint32(v Visitor) error              { return d.DecodeAny(v) }
func (d testData) DecodeUint64(v Visitor) error              { return d.DecodeAny(v) }
func (d testData) DecodeUintptr(v Visitor) error             { return d.DecodeAny(v) }
func (d testData) DecodeFloat32(v Visitor) error             { return d.DecodeAny(v) }
func (d testData) DecodeFloat64(v Visitor) error             { return d.DecodeAny(v) }
func (d testData) DecodeComplex64(v Visitor) error           { return d.DecodeAny(v) }
func (d testData) DecodeComplex128(v Visitor) error          { return d.DecodeAny(v) }
func (d testData) DecodeString(v Visitor) error              { return d.DecodeAny(v) }
func (d testData) DecodeBytes(v Visitor) error               { return d.DecodeAny(v) }
func (d testData) DecodeOption(v Visitor) error              { return d.DecodeAny(v) }
func (d testData) DecodeSeq(v Visitor) error                 { return d.DecodeAny(v) }
func (d testData) DecodeMap(v Visitor) error                 { return d.DecodeAny(v) }
func (d testData) DecodeStruct(name string, v Visitor) error { return d.DecodeAny(v) }

type testSeqDecoder struct {
	seq reflect.Value
	len int
	i   int
}

func (d *testSeqDecoder) Size() (int, bool) {
	return d.len, true
}

func (d *testSeqDecoder) NextElement(ds Deserializer) (bool, error) {
	if d.i == d.len {
		return false, nil
	}
	v := d.seq.Index(d.i)
	d.i++
	return true, ds.Deserialize(testData{v})
}

type testMapDecoder struct {
	iter *reflect.MapIter
}

func (d *testMapDecoder) Size() (int, bool) {
	return 0, false
}

func (d *testMapDecoder) NextKey(ds Deserializer) (bool, error) {
	if !d.iter.Next() {
		return false, nil
	}
	return true, ds.Deserialize(testData{d.iter.Key()})
}

func (d *testMapDecoder) NextValue(ds Deserializer) error {
	return ds.Deserialize(testData{d.iter.Value()})
}

type testStructDecoder struct {
	struct_ reflect.Value
	fields  []reflect.StructField
	i       int
}

func (d *testStructDecoder) Size() (int, bool) {
	return 0, false
}

func (d *testStructDecoder) NextKey(ds Deserializer) (bool, error) {
	for {
		if d.i == len(d.fields) {
			return false, nil
		}
		f := &d.fields[d.i]
		d.i++
		if !f.Anonymous {
			return true, ds.Deserialize(testData{reflect.ValueOf(f.Name)})
		}
	}
}

func (d *testStructDecoder) NextValue(ds Deserializer) error {
	return ds.Deserialize(testData{d.struct_.FieldByIndex(d.fields[d.i-1].Index)})
}

type boolStruct struct {
	Field bool
}

type boolStructVisitor struct {
	DefaultVisitor
	struct_ *boolStruct
}

func (v boolStructVisitor) VisitMap(map_ MapDecoder) error {
	for {
		var k string
		ok, err := map_.NextKey(NewString(&k))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		var d Deserializer
		switch k {
		case "Field":
			d = NewBool(&v.struct_.Field)
		default:
			d = SkipCodec{}
		}

		if err := map_.NextValue(d); err != nil {
			return err
		}
	}
}

func (v boolStructVisitor) Deserialize(d Decoder) error {
	return d.DecodeStruct("boolStruct", v)
}

func TestCodec(t *testing.T) {
	var b bool
	err := NewBool(&b).Deserialize(data(true))
	require.NoError(t, err)
	assert.Equal(t, true, b)

	var bp *bool
	err = NewPtr[BoolCodec[bool]](&bp).Deserialize(data(true))
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	var bools []bool
	err = NewSeq[BoolCodec[bool]](&bools).Deserialize(data([]bool{true, false}))
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, bools)

	var boolMap map[int]bool
	err = NewMap[IntCodec[int], BoolCodec[bool]](&boolMap).Deserialize(data(map[int]bool{42: true}))
	require.NoError(t, err)
	require.Equal(t, map[int]bool{42: true}, boolMap)

	var struct_ boolStruct
	err = boolStructVisitor{struct_: &struct_}.Deserialize(data(boolStruct{Field: true}))
	require.NoError(t, err)
	require.Equal(t, boolStruct{Field: true}, struct_)

	var any_ any
	err = NewAny(&any_).Deserialize(data(boolStruct{Field: true}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"Field": true}, any_)
}

func TestCodecReflect(t *testing.T) {
	var b bool
	err := GetDeserializer(&b).Deserialize(data(true))
	require.NoError(t, err)
	assert.Equal(t, true, b)

	var bp *bool
	err = GetDeserializer(&bp).Deserialize(data(true))
	require.NoError(t, err)
	require.NotNil(t, bp)
	assert.Equal(t, true, *bp)

	var bools []bool
	err = GetDeserializer(&bools).Deserialize(data([]bool{true, false}))
	require.NoError(t, err)
	require.Equal(t, []bool{true, false}, bools)

	var boolMap map[int]bool
	err = GetDeserializer(&boolMap).Deserialize(data(map[int]bool{42: true}))
	require.NoError(t, err)
	require.Equal(t, map[int]bool{42: true}, boolMap)

	var struct_ boolStruct
	err = GetDeserializer(&struct_).Deserialize(data(boolStruct{Field: true}))
	require.NoError(t, err)
	require.Equal(t, boolStruct{Field: true}, struct_)

	var any_ any
	err = GetDeserializer(&any_).Deserialize(data(boolStruct{Field: true}))
	require.NoError(t, err)
	require.Equal(t, map[string]any{"Field": true}, any_)
}
