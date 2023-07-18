package json

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/pgavlin/codec/typecache"
)

const (
	// 1000 is the value used by the standard encoding/json package.
	//
	// https://cs.opensource.google/go/go/+/refs/tags/go1.17.3:src/encoding/json/encode.go;drc=refs%2Ftags%2Fgo1.17.3;l=300
	startDetectingCyclesAfter = 1000
)

type encoder struct {
	flags AppendFlags
	// ptrDepth tracks the depth of pointer cycles, when it reaches the value
	// of startDetectingCyclesAfter, the ptrSeen map is allocated and the
	// encoder starts tracking pointers it has seen as an attempt to detect
	// whether it has entered a pointer cycle and needs to error before the
	// goroutine runs out of stack space.
	ptrDepth uint32
	ptrSeen  map[unsafe.Pointer]struct{}
}

type decoder struct {
	flags ParseFlags
}

type jsonCodec struct {
	encode func(e encoder, b []byte, v any) ([]byte, error)
	decode func(d decoder, b []byte, v any) ([]byte, error)
}

var codecs typecache.Cache[jsonCodec]

func getCodec(v any) jsonCodec {
	return codecs.GetOrCreate(reflect.TypeOf(v), func(t reflect.Type) jsonCodec {
		if t == nil {
			return jsonCodec{}
		}

		// TOOD:
		// - json.Number
		// - json.RawMessage

		// TODO: pointer cases?

		var c jsonCodec
		switch {
		case t.Implements(jsonMarshalerType):
			c.encode = encoder.encodeJSONMarshaler
		case t.Implements(textMarshalerType):
			c.encode = encoder.encodeTextMarshaler
		}

		p := t
		if p.Kind() != reflect.Pointer {
			p = reflect.PtrTo(t)
		}
		switch {
		case p.Implements(jsonUnmarshalerType):
			c.decode = decoder.decodeJSONUnmarshaler
		case p.Implements(textUnmarshalerType):
			c.decode = decoder.decodeTextUnmarshaler
		}

		return c
	})
}

func unmarshalTypeError(b []byte, t reflect.Type) error {
	return &UnmarshalTypeError{Value: strconv.Quote(prefix(b)), Type: t}
}

func unmarshalOverflow(b []byte, t reflect.Type) error {
	return &UnmarshalTypeError{Value: "number " + prefix(b) + " overflows", Type: t}
}

func unexpectedEOF(b []byte) error {
	return syntaxError(b, "unexpected end of JSON input")
}

var syntaxErrorMsgOffset = ^uintptr(0)

func init() {
	t := reflect.TypeOf(SyntaxError{})
	for i, n := 0, t.NumField(); i < n; i++ {
		if f := t.Field(i); f.Type.Kind() == reflect.String {
			syntaxErrorMsgOffset = f.Offset
		}
	}
}

func syntaxError(b []byte, msg string, args ...interface{}) error {
	e := new(SyntaxError)
	i := syntaxErrorMsgOffset
	if i != ^uintptr(0) {
		s := "json: " + fmt.Sprintf(msg, args...) + ": " + prefix(b)
		p := unsafe.Pointer(e)
		// Hack to set the unexported `msg` field.
		*(*string)(unsafe.Pointer(uintptr(p) + i)) = s
	}
	return e
}

func objectKeyError(b []byte, err error) ([]byte, error) {
	if len(b) == 0 {
		return nil, unexpectedEOF(b)
	}
	switch err.(type) {
	case *UnmarshalTypeError:
		err = syntaxError(b, "invalid character '%c' looking for beginning of object key", b[0])
	}
	return b, err
}

func prefix(b []byte) string {
	if len(b) < 32 {
		return string(b)
	}
	return string(b[:32]) + "..."
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
	nilType  = reflect.TypeOf(nil)
	boolType = reflect.TypeOf((*bool)(nil)).Elem()

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

	float32Type    = reflect.TypeOf((*float32)(nil)).Elem()
	float64Type    = reflect.TypeOf((*float64)(nil)).Elem()
	complex64Type  = reflect.TypeOf((*complex64)(nil)).Elem()
	complex128Type = reflect.TypeOf((*complex64)(nil)).Elem()

	ptrType    = reflect.TypeOf((*any)(nil))
	sliceType  = reflect.TypeOf((*[]any)(nil)).Elem()
	mapType    = reflect.TypeOf((*map[string]any)(nil)).Elem()
	structType = reflect.TypeOf((*struct{})(nil)).Elem()

	bytesType  = reflect.TypeOf((*[]byte)(nil)).Elem()
	numberType = reflect.TypeOf((*Number)(nil)).Elem()

	jsonMarshalerType   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	jsonUnmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	textMarshalerType   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

// =============================================================================
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// appendDuration appends a human-readable representation of d to b.
//
// The function copies the implementation of time.Duration.String but prevents
// Go from making a dynamic memory allocation on the returned value.
func appendDuration(b []byte, d time.Duration) []byte {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte
	w := len(buf)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			return append(b, '0', 's')
		case u < uint64(time.Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		w--
		buf[w] = 's'

		w, u = fmtFrac(buf[:w], u, 9)

		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'm'
			w = fmtInt(buf[:w], u%60)
			u /= 60

			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 {
				w--
				buf[w] = 'h'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return append(b, buf[w:]...)
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros.  it omits the decimal
// point too when the fraction is 0.  It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

// =============================================================================
