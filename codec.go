package codec

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"github.com/pgavlin/codec/typecache"
	"github.com/segmentio/asm/keyset"
)

func GetDeserializer(v any, format Format) Deserializer {
	if d, ok := v.(Deserializer); ok {
		return d
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return unsupportedTypeCodec{t: rv.Type()}
	}
	return getCodec(rv.Type().Elem(), format).new(rv.UnsafePointer())
}

func GetSerializer(v any, format Format) Serializer {
	if s, ok := v.(Serializer); ok {
		return s
	}
	if v == nil {
		return NilCodec{}
	}
	return getCodec(reflect.TypeOf(v), format).new((*iface)(unsafe.Pointer(&v)).ptr)
}

type format struct {
	codecs typecache.Cache[codec]
}

var codecs typecache.Cache[codec]

func getCodec(t reflect.Type, format Format) codec {
	return codecs.GetOrCreate(t, func(t reflect.Type) codec {
		// TODO: inlined...?
		//	if inlined(t) {
		//		c.encode = constructInlineValueEncodeFunc(c.encode)
		//	}
		return constructCodec(t, map[reflect.Type]*structType{}, t.Kind() == reflect.Ptr)
	})
}

func IsAny[T any]() bool {
	var v T
	return reflect.TypeOf(v) == anyType
}

const (
	// 1000 is the value used by the standard encoding/json package.
	//
	// https://cs.opensource.google/go/go/+/refs/tags/go1.17.3:src/encoding/json/encode.go;drc=refs%2Ftags%2Fgo1.17.3;l=300
	startDetectingCyclesAfter = 1000
)

type encoder struct {
	// ptrDepth tracks the depth of pointer cycles, when it reaches the value
	// of startDetectingCyclesAfter, the ptrSeen map is allocated and the
	// encoder starts tracking pointers it has seen as an attempt to detect
	// whether it has entered a pointer cycle and needs to error before the
	// goroutine runs out of stack space.
	ptrDepth uint32
	ptrSeen  map[unsafe.Pointer]struct{}
}

type decoder struct{}

type emptyFunc func(unsafe.Pointer) bool
type sortFunc func([]reflect.Value)

func constructCodec(t reflect.Type, seen map[reflect.Type]*structType, canAddr bool) (c codec) {
	switch t {
	case nullType:
		return nilCodec{}
	case boolType:
		return boolCodec{}
	case intType:
		return intCodec{}
	case int8Type:
		return int8Codec{}
	case int16Type:
		return int16Codec{}
	case int32Type:
		return int32Codec{}
	case int64Type:
		return int64Codec{}
	case uintType:
		return uintCodec{}
	case uint8Type:
		return uint8Codec{}
	case uint16Type:
		return uint16Codec{}
	case uint32Type:
		return uint32Codec{}
	case uint64Type:
		return uint64Codec{}
	case uintptrType:
		return uintptrCodec{}
	case float32Type:
		return float32Codec{}
	case float64Type:
		return float64Codec{}
	case stringType:
		return stringCodec{}
	case bytesType:
		return bytesCodec{}
	}

	switch t.Kind() {
	case reflect.Bool:
		c = boolCodec{}
	case reflect.Int:
		c = intCodec{}
	case reflect.Int8:
		c = int8Codec{}
	case reflect.Int16:
		c = int16Codec{}
	case reflect.Int32:
		c = int32Codec{}
	case reflect.Int64:
		c = int64Codec{}
	case reflect.Uint:
		c = uintCodec{}
	case reflect.Uintptr:
		c = uintptrCodec{}
	case reflect.Uint8:
		c = uint8Codec{}
	case reflect.Uint16:
		c = uint16Codec{}
	case reflect.Uint32:
		c = uint32Codec{}
	case reflect.Uint64:
		c = uint64Codec{}
	case reflect.Float32:
		c = float32Codec{}
	case reflect.Float64:
		c = float64Codec{}
	case reflect.String:
		c = stringCodec{}
	case reflect.Interface:
		c = AnyCodec{}
	case reflect.Array:
		c = constructArrayCodec(t, seen, canAddr)
	case reflect.Slice:
		c = constructSliceCodec(t, seen)
	case reflect.Map:
		c = constructMapCodec(t, seen)
	case reflect.Struct:
		c = constructStructCodec(t, seen, canAddr)
	case reflect.Ptr:
		c = constructPointerCodec(t, seen)
	default:
		c = constructUnsupportedTypeCodec(t)
	}

	if t.Implements(codecSerializerType) {
		c = constructSerializerCodec(t, c)
	}
	if reflect.PtrTo(t).Implements(codecDeserializerType) {
		c = constructDeserializerCodec(t, c)
	}

	return
}

func constructArrayCodec(t reflect.Type, seen map[reflect.Type]*structType, canAddr bool) codec {
	e := t.Elem()
	c := constructCodec(e, seen, canAddr)
	n := t.Len()
	return arrayCodec{
		arrayType: &arrayType{
			n:    n,
			t:    t,
			elem: c,
		},
	}
}

func constructSliceCodec(t reflect.Type, seen map[reflect.Type]*structType) codec {
	e := t.Elem()
	s := alignedSize(e)
	c := constructCodec(e, seen, true)

	if e == uint8Type {
		return bytesCodec{}
	}

	return sliceCodec{
		sliceType: &sliceType{
			size: s,
			t:    t,
			elem: c,
		},
	}
}

func constructMapCodec(t reflect.Type, seen map[reflect.Type]*structType) codec {
	var sortKeys sortFunc
	kt := t.Key()
	vt := t.Elem()

	kc := constructCodec(kt, seen, false)
	vc := constructCodec(vt, seen, false)

	switch kt.Kind() {
	case reflect.String:
		sortKeys = func(keys []reflect.Value) {
			sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
		}

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		sortKeys = func(keys []reflect.Value) {
			sort.Slice(keys, func(i, j int) bool { return keys[i].Int() < keys[j].Int() })
		}

	case reflect.Uint,
		reflect.Uintptr,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		sortKeys = func(keys []reflect.Value) {
			sort.Slice(keys, func(i, j int) bool { return keys[i].Uint() < keys[j].Uint() })
		}

	default:
		return constructUnsupportedTypeCodec(t)
	}

	// TODO: inlined...?
	//	if inlined(v) {
	//		vc.encode = constructInlineValueEncodeFunc(vc.encode)
	//	}

	kz := reflect.Zero(kt)
	vz := reflect.Zero(vt)

	return mapCodec{
		mapType: &mapType{
			t:        t,
			kt:       kt,
			vt:       vt,
			kz:       kz,
			vz:       vz,
			kc:       kc,
			vc:       vc,
			sortKeys: sortKeys,
		},
	}
}

func constructStructCodec(t reflect.Type, seen map[reflect.Type]*structType, canAddr bool) codec {
	st := constructStructType(t, seen, canAddr)
	return structCodec{
		structType: st,
	}
}

func constructStructType(t reflect.Type, seen map[reflect.Type]*structType, canAddr bool) *structType {
	// Used for preventing infinite recursion on types that have pointers to
	// themselves.
	st := seen[t]

	if st == nil {
		st = &structType{
			fields:      make([]structField, 0, t.NumField()),
			fieldsIndex: make(map[string]*structField),
			ficaseIndex: make(map[string]*structField),
			typ:         t,
		}

		seen[t] = st
		st.fields = appendStructFields(st.fields, t, 0, seen, canAddr)

		for i := range st.fields {
			f := &st.fields[i]
			s := strings.ToLower(f.name)
			st.fieldsIndex[f.name] = f
			// When there is ambiguity because multiple fields have the same
			// case-insensitive representation, the first field must win.
			if _, exists := st.ficaseIndex[s]; !exists {
				st.ficaseIndex[s] = f
			}
		}

		// At a certain point the linear scan provided by keyset is less
		// efficient than a map. The 32 was chosen based on benchmarks in the
		// segmentio/asm repo run with an Intel Kaby Lake processor and go1.17.
		if len(st.fields) <= 32 {
			keys := make([][]byte, len(st.fields))
			for i, f := range st.fields {
				keys[i] = []byte(f.name)
			}
			st.keyset = keyset.New(keys)
		}
	}

	return st
}

func appendStructFields(fields []structField, t reflect.Type, offset uintptr, seen map[reflect.Type]*structType, canAddr bool) []structField {
	type embeddedField struct {
		index      int
		offset     uintptr
		pointer    bool
		unexported bool
		subtype    *structType
		subfield   *structField
	}

	names := make(map[string]struct{})
	embedded := make([]embeddedField, 0, 10)

	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)

		var (
			name       = f.Name
			anonymous  = f.Anonymous
			tag        = false
			omitempty  = false
			unexported = len(f.PkgPath) != 0
		)

		if unexported && !anonymous { // unexported
			continue
		}

		if parts := strings.Split(f.Tag.Get("codec"), ","); len(parts) != 0 {
			if len(parts[0]) != 0 {
				name, tag = parts[0], true
			}

			if name == "-" && len(parts) == 1 { // ignored
				continue
			}

			if !isValidTag(name) {
				name = f.Name
			}

			for _, tag := range parts[1:] {
				switch tag {
				case "omitempty":
					omitempty = true
				}
			}
		}

		if anonymous && !tag { // embedded
			typ := f.Type
			ptr := f.Type.Kind() == reflect.Ptr

			if ptr {
				typ = f.Type.Elem()
			}

			if typ.Kind() == reflect.Struct {
				// When the embedded fields is inlined the fields can be looked
				// up by offset from the address of the wrapping object, so we
				// simply add the embedded struct fields to the list of fields
				// of the current struct type.
				subtype := constructStructType(typ, seen, canAddr)

				for j := range subtype.fields {
					embedded = append(embedded, embeddedField{
						index:      i<<32 | j,
						offset:     offset + f.Offset,
						pointer:    ptr,
						unexported: unexported,
						subtype:    subtype,
						subfield:   &subtype.fields[j],
					})
				}

				continue
			}

			if unexported { // ignore unexported non-struct types
				continue
			}
		}

		codec := constructCodec(f.Type, seen, canAddr)

		fields = append(fields, structField{
			codec:     codec,
			offset:    offset + f.Offset,
			empty:     emptyFuncOf(f.Type),
			tag:       tag,
			omitempty: omitempty,
			name:      name,
			index:     i << 32,
			typ:       f.Type,
			zero:      reflect.Zero(f.Type),
		})

		names[name] = struct{}{}
	}

	// Only unambiguous embedded fields must be serialized.
	ambiguousNames := make(map[string]int)
	ambiguousTags := make(map[string]int)

	// Embedded types can never override a field that was already present at
	// the top-level.
	for name := range names {
		ambiguousNames[name]++
		ambiguousTags[name]++
	}

	for _, embfield := range embedded {
		ambiguousNames[embfield.subfield.name]++
		if embfield.subfield.tag {
			ambiguousTags[embfield.subfield.name]++
		}
	}

	for _, embfield := range embedded {
		subfield := *embfield.subfield

		if ambiguousNames[subfield.name] > 1 && !(subfield.tag && ambiguousTags[subfield.name] == 1) {
			continue // ambiguous embedded field
		}

		if embfield.pointer {
			subfield.embedded = &embeddedStructField{
				unexported: embfield.unexported,
				offset:     embfield.offset,
			}
		} else {
			subfield.offset += embfield.offset
		}

		// To prevent dominant flags more than one level below the embedded one.
		subfield.tag = false

		// To ensure the order of the fields in the output is the same is in the
		// struct type.
		subfield.index = embfield.index

		fields = append(fields, subfield)
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].index < fields[j].index })
	return fields
}

func constructPointerCodec(t reflect.Type, seen map[reflect.Type]*structType) codec {
	e := t.Elem()
	c := constructCodec(e, seen, true)
	return ptrCodec{
		ptrType: &ptrType{
			t:    t,
			elem: c,
		},
	}
}

func constructUnsupportedTypeCodec(t reflect.Type) codec {
	return unsupportedTypeCodec{t: t}
}

func constructSerializerCodec(t reflect.Type, next codec) codec {
	return serializerCodec{
		t:    t,
		next: next,
	}
}

func constructDeserializerCodec(t reflect.Type, next codec) codec {
	return deserializerCodec{
		t:    t,
		next: next,
	}
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input. noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func alignedSize(t reflect.Type) uintptr {
	a := t.Align()
	s := t.Size()
	return align(uintptr(a), uintptr(s))
}

func align(align, size uintptr) uintptr {
	if align != 0 && (size%align) != 0 {
		size = ((size / align) + 1) * align
	}
	return size
}

func inlined(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr:
		return true
	case reflect.Map:
		return true
	case reflect.Struct:
		return t.NumField() == 1 && inlined(t.Field(0).Type)
	default:
		return false
	}
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

func emptyFuncOf(t reflect.Type) emptyFunc {
	switch t {
	case bytesType:
		return func(p unsafe.Pointer) bool { return (*slice)(p).len == 0 }
	}

	//	if t.Implements(codecSerializerType) {
	//		pointer := t.Kind() == reflect.Pointer
	//		return func(p unsafe.Pointer) bool {
	//			v := reflect.NewAt(t, p)
	//			return v.Interface().(Serializer).IsEmpty()
	//		}
	//	}

	switch t.Kind() {
	case reflect.Array:
		if t.Len() == 0 {
			return func(unsafe.Pointer) bool { return true }
		}

	case reflect.Map:
		return func(p unsafe.Pointer) bool { return reflect.NewAt(t, p).Elem().Len() == 0 }

	case reflect.Slice:
		return func(p unsafe.Pointer) bool { return (*slice)(p).len == 0 }

	case reflect.String:
		return func(p unsafe.Pointer) bool { return len(*(*string)(p)) == 0 }

	case reflect.Bool:
		return func(p unsafe.Pointer) bool { return !*(*bool)(p) }

	case reflect.Int, reflect.Uint:
		return func(p unsafe.Pointer) bool { return *(*uint)(p) == 0 }

	case reflect.Uintptr:
		return func(p unsafe.Pointer) bool { return *(*uintptr)(p) == 0 }

	case reflect.Int8, reflect.Uint8:
		return func(p unsafe.Pointer) bool { return *(*uint8)(p) == 0 }

	case reflect.Int16, reflect.Uint16:
		return func(p unsafe.Pointer) bool { return *(*uint16)(p) == 0 }

	case reflect.Int32, reflect.Uint32:
		return func(p unsafe.Pointer) bool { return *(*uint32)(p) == 0 }

	case reflect.Int64, reflect.Uint64:
		return func(p unsafe.Pointer) bool { return *(*uint64)(p) == 0 }

	case reflect.Float32:
		return func(p unsafe.Pointer) bool { return *(*float32)(p) == 0 }

	case reflect.Float64:
		return func(p unsafe.Pointer) bool { return *(*float64)(p) == 0 }

	case reflect.Ptr:
		return func(p unsafe.Pointer) bool { return *(*unsafe.Pointer)(p) == nil }

	case reflect.Interface:
		return func(p unsafe.Pointer) bool { return (*iface)(p).ptr == nil }
	}

	return func(unsafe.Pointer) bool { return false }
}

type iface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

type slice struct {
	data unsafe.Pointer
	len  int
	cap  int
}

func unmarshalTypeError(b []byte, t reflect.Type) error {
	return &UnmarshalTypeError{Value: strconv.Quote(prefix(b)), Type: t}
}

func unmarshalOverflow(b []byte, t reflect.Type) error {
	return &UnmarshalTypeError{Value: "number " + prefix(b) + " overflows", Type: t}
}

func prefix(b []byte) string {
	if len(b) < 32 {
		return string(b)
	}
	return string(b[:32]) + "..."
}

func intStringsAreSorted(i0, i1 int64) bool {
	var b0, b1 [32]byte
	return string(strconv.AppendInt(b0[:0], i0, 10)) < string(strconv.AppendInt(b1[:0], i1, 10))
}

func uintStringsAreSorted(u0, u1 uint64) bool {
	var b0, b1 [32]byte
	return string(strconv.AppendUint(b0[:0], u0, 10)) < string(strconv.AppendUint(b1[:0], u1, 10))
}

func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: *(*unsafe.Pointer)(unsafe.Pointer(&s)),
		Len:  len(s),
		Cap:  len(s),
	}))
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

var (
	nullType = reflect.TypeOf(nil)
	boolType = reflect.TypeOf(false)

	intType   = reflect.TypeOf(int(0))
	int8Type  = reflect.TypeOf(int8(0))
	int16Type = reflect.TypeOf(int16(0))
	int32Type = reflect.TypeOf(int32(0))
	int64Type = reflect.TypeOf(int64(0))

	uintType    = reflect.TypeOf(uint(0))
	uint8Type   = reflect.TypeOf(uint8(0))
	uint16Type  = reflect.TypeOf(uint16(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	uint64Type  = reflect.TypeOf(uint64(0))
	uintptrType = reflect.TypeOf(uintptr(0))

	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))

	stringType = reflect.TypeOf("")
	bytesType  = reflect.TypeOf(([]byte)(nil))

	anyType               = reflect.TypeOf((*any)(nil)).Elem()
	codecSerializerType   = reflect.TypeOf((*Serializer)(nil)).Elem()
	codecDeserializerType = reflect.TypeOf((*Deserializer)(nil)).Elem()
)

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "codec: unsupported type: " + e.Type.String()
}

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "codec: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "codec: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
