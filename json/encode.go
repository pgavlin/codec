package json

import (
	"encoding"
	"math"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/pgavlin/codec"
	"github.com/segmentio/asm/base64"
)

const hex = "0123456789abcdef"

type Encoder struct {
	enc encoder
	out []byte
}

func (e *Encoder) encode(v any, s codec.Serializer) (err error) {
	codec := getCodec(v)
	if codec.encode != nil {
		e.out, err = codec.encode(e.enc, e.out, v)
		return
	}
	return s.Serialize(e)
}

func (e *Encoder) EncodeNil() (err error) {
	e.out, err = e.enc.encodeNull(e.out)
	return
}

func (e *Encoder) EncodeBool(v bool) (err error) {
	e.out, err = e.enc.encodeBool(e.out, v)
	return
}

func (e *Encoder) EncodeInt(v int) (err error) {
	e.out, err = e.enc.encodeInt(e.out, v)
	return
}

func (e *Encoder) EncodeInt8(v int8) (err error) {
	e.out, err = e.enc.encodeInt8(e.out, v)
	return
}

func (e *Encoder) EncodeInt16(v int16) (err error) {
	e.out, err = e.enc.encodeInt16(e.out, v)
	return
}

func (e *Encoder) EncodeInt32(v int32) (err error) {
	e.out, err = e.enc.encodeInt32(e.out, v)
	return
}

func (e *Encoder) EncodeInt64(v int64) (err error) {
	e.out, err = e.enc.encodeInt64(e.out, v)
	return
}

func (e *Encoder) EncodeUint(v uint) (err error) {
	e.out, err = e.enc.encodeUint(e.out, v)
	return
}

func (e *Encoder) EncodeUintptr(v uintptr) (err error) {
	e.out, err = e.enc.encodeUintptr(e.out, v)
	return
}

func (e *Encoder) EncodeUint8(v uint8) (err error) {
	e.out, err = e.enc.encodeUint8(e.out, v)
	return
}

func (e *Encoder) EncodeUint16(v uint16) (err error) {
	e.out, err = e.enc.encodeUint16(e.out, v)
	return
}

func (e *Encoder) EncodeUint32(v uint32) (err error) {
	e.out, err = e.enc.encodeUint32(e.out, v)
	return
}

func (e *Encoder) EncodeUint64(v uint64) (err error) {
	e.out, err = e.enc.encodeUint64(e.out, v)
	return
}

func (e *Encoder) EncodeFloat32(v float32) (err error) {
	e.out, err = e.enc.encodeFloat32(e.out, v)
	return
}

func (e *Encoder) EncodeFloat64(v float64) (err error) {
	e.out, err = e.enc.encodeFloat64(e.out, v)
	return
}

func (e *Encoder) EncodeComplex64(v complex64) error {
	return &UnsupportedTypeError{Type: complex64Type}
}

func (e *Encoder) EncodeComplex128(v complex128) error {
	return &UnsupportedTypeError{Type: complex128Type}
}

func (e *Encoder) EncodeString(v string) (err error) {
	e.out, err = e.enc.encodeString(e.out, v)
	return
}

func (e *Encoder) EncodeBytes(v []byte) (err error) {
	e.out, err = e.enc.encodeBytes(e.out, v)
	return
}

func (e *Encoder) EncodeElem(v any, s codec.Serializer) error {
	return e.encode(v, s)
}

func (e *Encoder) EncodeSeq(count int) (codec.SeqEncoder, error) {
	e.out = append(e.out, '[')
	return &SeqEncoder{enc: e, first: true}, nil
}

func (e *Encoder) EncodeMap(len int) (codec.MapEncoder, error) {
	e.out = append(e.out, '{')
	return &MapEncoder{enc: e, first: true}, nil
}

func (e *Encoder) EncodeStruct(name string) (codec.StructEncoder, error) {
	e.out = append(e.out, '{')
	return &StructEncoder{enc: e, first: true}, nil
}

type SeqEncoder struct {
	enc   *Encoder
	first bool
}

func (e *SeqEncoder) Close() error {
	e.enc.out = append(e.enc.out, ']')
	return nil
}

func (e *SeqEncoder) EncodeElement(v any, s codec.Serializer) error {
	if !e.first {
		e.enc.out = append(e.enc.out, ',')
	} else {
		e.first = false
	}
	return e.enc.encode(v, s)
}

type MapEncoder struct {
	enc   *Encoder
	first bool
}

func (e *MapEncoder) Close() error {
	e.enc.out = append(e.enc.out, '}')
	return nil
}

func (e *MapEncoder) EncodeKey(k any, s codec.Serializer) error {
	if !e.first {
		e.enc.out = append(e.enc.out, ',')
	} else {
		e.first = false
	}
	return mapKeyEncoder{enc: e.enc}.encode(k, s)
}

func (e *MapEncoder) EncodeValue(v any, s codec.Serializer) error {
	e.enc.out = append(e.enc.out, ':')
	return e.enc.encode(v, s)
}

type StructEncoder struct {
	enc   *Encoder
	first bool
}

func (e *StructEncoder) Close() error {
	e.enc.out = append(e.enc.out, '}')
	return nil
}

func (e *StructEncoder) EncodeField(key string, v any, s codec.Serializer) error {
	if !e.first {
		e.enc.out = append(e.enc.out, ',')
	} else {
		e.first = false
	}

	if err := e.enc.EncodeString(key); err != nil {
		return err
	}
	e.enc.out = append(e.enc.out, ':')
	return e.enc.encode(v, s)
}

func (e encoder) encodeNull(b []byte) ([]byte, error) {
	return append(b, "null"...), nil
}

func (e encoder) encodeBool(b []byte, v bool) ([]byte, error) {
	if v {
		return append(b, "true"...), nil
	}
	return append(b, "false"...), nil
}

func (e encoder) encodeInt(b []byte, v int) ([]byte, error) {
	return appendInt(b, int64(v)), nil
}

func (e encoder) encodeInt8(b []byte, v int8) ([]byte, error) {
	return appendInt(b, int64(v)), nil
}

func (e encoder) encodeInt16(b []byte, v int16) ([]byte, error) {
	return appendInt(b, int64(v)), nil
}

func (e encoder) encodeInt32(b []byte, v int32) ([]byte, error) {
	return appendInt(b, int64(v)), nil
}

func (e encoder) encodeInt64(b []byte, v int64) ([]byte, error) {
	return appendInt(b, v), nil
}

func (e encoder) encodeUint(b []byte, v uint) ([]byte, error) {
	return appendUint(b, uint64(v)), nil
}

func (e encoder) encodeUintptr(b []byte, v uintptr) ([]byte, error) {
	return appendUint(b, uint64(v)), nil
}

func (e encoder) encodeUint8(b []byte, v uint8) ([]byte, error) {
	return appendUint(b, uint64(v)), nil
}

func (e encoder) encodeUint16(b []byte, v uint16) ([]byte, error) {
	return appendUint(b, uint64(v)), nil
}

func (e encoder) encodeUint32(b []byte, v uint32) ([]byte, error) {
	return appendUint(b, uint64(v)), nil
}

func (e encoder) encodeUint64(b []byte, v uint64) ([]byte, error) {
	return appendUint(b, v), nil
}

func (e encoder) encodeFloat32(b []byte, v float32) ([]byte, error) {
	return e.encodeFloat(b, float64(v), 32)
}

func (e encoder) encodeFloat64(b []byte, v float64) ([]byte, error) {
	return e.encodeFloat(b, v, 64)
}

func (e encoder) encodeFloat(b []byte, f float64, bits int) ([]byte, error) {
	switch {
	case math.IsNaN(f):
		return b, &UnsupportedValueError{Value: reflect.ValueOf(f), Str: "NaN"}
	case math.IsInf(f, 0):
		return b, &UnsupportedValueError{Value: reflect.ValueOf(f), Str: "inf"}
	}

	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}

	b = strconv.AppendFloat(b, f, fmt, -1, int(bits))

	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(b)
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}

	return b, nil
}

func (e encoder) encodeString(b []byte, v string) ([]byte, error) {
	s := v
	if len(s) == 0 {
		return append(b, `""`...), nil
	}
	i := 0
	j := 0
	escapeHTML := (e.flags & EscapeHTML) != 0

	b = append(b, '"')

	if len(s) >= 8 {
		if j = escapeIndex(s, escapeHTML); j < 0 {
			return append(append(b, s...), '"'), nil
		}
	}

	for j < len(s) {
		c := s[j]

		if c >= 0x20 && c <= 0x7f && c != '\\' && c != '"' && (!escapeHTML || (c != '<' && c != '>' && c != '&')) {
			// fast path: most of the time, printable ascii characters are used
			j++
			continue
		}

		switch c {
		case '\\', '"':
			b = append(b, s[i:j]...)
			b = append(b, '\\', c)
			i = j + 1
			j = j + 1
			continue

		case '\n':
			b = append(b, s[i:j]...)
			b = append(b, '\\', 'n')
			i = j + 1
			j = j + 1
			continue

		case '\r':
			b = append(b, s[i:j]...)
			b = append(b, '\\', 'r')
			i = j + 1
			j = j + 1
			continue

		case '\t':
			b = append(b, s[i:j]...)
			b = append(b, '\\', 't')
			i = j + 1
			j = j + 1
			continue

		case '<', '>', '&':
			b = append(b, s[i:j]...)
			b = append(b, `\u00`...)
			b = append(b, hex[c>>4], hex[c&0xF])
			i = j + 1
			j = j + 1
			continue
		}

		// This encodes bytes < 0x20 except for \t, \n and \r.
		if c < 0x20 {
			b = append(b, s[i:j]...)
			b = append(b, `\u00`...)
			b = append(b, hex[c>>4], hex[c&0xF])
			i = j + 1
			j = j + 1
			continue
		}

		r, size := utf8.DecodeRuneInString(s[j:])

		if r == utf8.RuneError && size == 1 {
			b = append(b, s[i:j]...)
			b = append(b, `\ufffd`...)
			i = j + size
			j = j + size
			continue
		}

		switch r {
		case '\u2028', '\u2029':
			// U+2028 is LINE SEPARATOR.
			// U+2029 is PARAGRAPH SEPARATOR.
			// They are both technically valid characters in JSON strings,
			// but don't work in JSONP, which has to be evaluated as JavaScript,
			// and can lead to security holes there. It is valid JSON to
			// escape them, so we do so unconditionally.
			// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
			b = append(b, s[i:j]...)
			b = append(b, `\u202`...)
			b = append(b, hex[r&0xF])
			i = j + size
			j = j + size
			continue
		}

		j += size
	}

	b = append(b, s[i:]...)
	b = append(b, '"')
	return b, nil
}

func (e encoder) encodeBytes(b, v []byte) ([]byte, error) {
	if v == nil {
		return append(b, "null"...), nil
	}

	n := base64.StdEncoding.EncodedLen(len(v)) + 2

	if avail := cap(b) - len(b); avail < n {
		newB := make([]byte, cap(b)+(n-avail))
		copy(newB, b)
		b = newB[:len(b)]
	}

	i := len(b)
	j := len(b) + n

	b = b[:j]
	b[i] = '"'
	base64.StdEncoding.Encode(b[i+1:j-1], v)
	b[j-1] = '"'
	return b, nil
}

func (e encoder) encodeJSONMarshaler(b []byte, v any) ([]byte, error) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		if reflect.ValueOf(v).IsNil() {
			return append(b, "null"...), nil
		}
	}

	j, err := v.(Marshaler).MarshalJSON()
	if err != nil {
		return nil, err
	}

	d := decoder{}
	s, _, _, err := d.parseValue(j)
	if err != nil {
		return b, &MarshalerError{Type: t, Err: err}
	}

	// TODO: escape HTML

	return append(b, s...), nil
}

func (e encoder) encodeTextMarshaler(b []byte, v any) ([]byte, error) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		if reflect.ValueOf(v).IsNil() {
			return append(b, "null"...), nil
		}
	}

	s, err := v.(encoding.TextMarshaler).MarshalText()
	if err != nil {
		return nil, err
	}

	return e.encodeString(b, string(s))
}
