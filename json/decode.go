package json

import (
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/pgavlin/codec"
)

func (d decoder) decodeNull(b []byte, cv codec.Visitor) ([]byte, error) {
	if hasNullPrefix(b) {
		return b[4:], cv.VisitNil()
	}
	if len(b) < 4 {
		return b[len(b):], unexpectedEOF(b)
	}
	return b, syntaxError(b, "expected 'null' but found invalid token")
}

func (d decoder) decodeTrue(b []byte, cv codec.Visitor) ([]byte, error) {
	if hasTruePrefix(b) {
		return b[4:], cv.VisitBool(true)
	}
	if len(b) < 4 {
		return b[len(b):], unexpectedEOF(b)
	}
	return b, syntaxError(b, "expected 'true' but found invalid token")
}

func (d decoder) decodeFalse(b []byte, cv codec.Visitor) ([]byte, error) {
	if hasFalsePrefix(b) {
		return b[5:], cv.VisitBool(false)
	}
	if len(b) < 5 {
		return b[len(b):], unexpectedEOF(b)
	}
	return b, syntaxError(b, "expected 'false' but found invalid token")
}

var uint64Type = reflect.TypeOf((*uint64)(nil)).Elem()

// convUint converts an already-parsed number into a uint64.
func convUint(b []byte) (uint64, error) {
	if len(b) == 1 && b[0] == '0' {
		return 0, nil
	}

	var value uint64
	for _, c := range b {
		x := uint64(c - '0')
		next := value*10 + x
		if next < value {
			return 0, unmarshalOverflow(b, uint64Type)
		}
		value = next
	}
	return value, nil
}

var int64Type = reflect.TypeOf((*int64)(nil)).Elem()

// convInt converts an already-parsed number into an int64.
func convInt(b []byte) (int64, error) {
	var value int64

	if b[0] == '-' {
		const max = math.MinInt64
		const lim = max / 10

		if len(b) == 2 && b[1] == '0' {
			return 0, nil
		}

		for _, c := range b[1:] {
			if value < lim {
				return 0, unmarshalOverflow(b, int64Type)
			}

			value *= 10
			x := int64(c - '0')

			if value < (max + x) {
				return 0, unmarshalOverflow(b, int64Type)
			}

			value -= x
		}
	} else {
		if len(b) == 1 && b[1] == '0' {
			return 0, nil
		}

		for _, c := range b {
			x := int64(c - '0')
			next := value*10 + x
			if next < value {
				return 0, unmarshalOverflow(b, int64Type)
			}
			value = next
		}
	}

	return value, nil
}

func (d decoder) decodeNumber(b []byte, cv codec.Visitor) ([]byte, error) {
	v, r, kind, err := d.parseNumber(b)
	if err != nil {
		return r, err
	}

	switch kind {
	case Uint:
		u, err := convUint(v)
		if err != nil {
			return r, err
		}
		return r, cv.VisitUint64(u)
	case Int:
		i, err := convInt(v)
		if err != nil {
			return r, err
		}
		return r, cv.VisitInt64(i)
	case Float:
		f, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&v)), 64)
		if err != nil {
			return r, err
		}
		return r, cv.VisitFloat64(f)
	default:
		panic("unexpected number kind")
	}
}

func (d decoder) decodeString(b []byte, cv codec.Visitor) ([]byte, error) {
	s, r, new, err := d.parseStringUnquote(b, nil)
	if err != nil {
		return r, err
	}

	var str string
	if new || (d.flags&DontCopyString) != 0 {
		str = *(*string)(unsafe.Pointer(&s))
	} else {
		str = string(s)
	}

	return r, cv.VisitString(str)
}

func (d decoder) decodeObject(b []byte, cv codec.Visitor) ([]byte, error) {
	if len(b) < 2 {
		return b[len(b):], unexpectedEOF(b)
	}

	if b[0] != '{' {
		return b, syntaxError(b, "expected '{' at the beginning of an object value")
	}

	dec := MapDecoder{first: true, rest: b[1:], flags: d.flags}
	err := cv.VisitMap(&dec)
	return dec.rest, err
}

func (d decoder) decodeArray(b []byte, cv codec.Visitor) ([]byte, error) {
	if len(b) < 2 {
		return b[len(b):], unexpectedEOF(b)
	}

	if b[0] != '[' {
		return b, syntaxError(b, "expected '[' at the beginning of array value")
	}

	dec := SeqDecoder{first: true, rest: b[1:], flags: d.flags}
	err := cv.VisitSeq(&dec)
	return dec.rest, err
}

func (d decoder) decodeValue(b []byte, cv codec.Visitor) ([]byte, error) {
	if len(b) == 0 {
		return b, syntaxError(b, "unexpected end of JSON input")
	}

	var err error
	switch b[0] {
	case '{':
		b, err = d.decodeObject(b, cv)
	case '[':
		b, err = d.decodeArray(b, cv)
	case '"':
		b, err = d.decodeString(b, cv)
	case 'n':
		b, err = d.decodeNull(b, cv)
	case 't':
		b, err = d.decodeTrue(b, cv)
	case 'f':
		b, err = d.decodeFalse(b, cv)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		b, err = d.decodeNumber(b, cv)
	default:
		err = syntaxError(b, "invalid character '%c' looking for beginning of value", b[0])
	}
	return b, err
}

type SeqDecoder struct {
	first bool
	rest  []byte
	flags ParseFlags
}

func (d *SeqDecoder) Size() (int, bool) {
	return 0, false
}

func (d *SeqDecoder) NextElement(ds codec.Deserializer) (bool, error) {
	b := skipSpaces(d.rest)

	if len(b) == 0 {
		d.rest = b
		return false, syntaxError(b, "missing closing ']' after array value")
	}

	if b[0] == ']' {
		d.rest = b[1:]
		return false, nil
	}

	if !d.first {
		if len(b) == 0 {
			d.rest = b
			return false, syntaxError(b, "unexpected EOF after array element")
		}
		if b[0] != ',' {
			d.rest = b
			return false, syntaxError(b, "expected ',' after array element but found '%c'", b[0])
		}
		b = skipSpaces(b[1:])
		if len(b) == 0 {
			d.rest = b
			return false, unexpectedEOF(b)
		}
		if b[0] == ']' {
			d.rest = b
			return false, syntaxError(b, "unexpected trailing comma after object field")
		}
	} else {
		d.first = false
	}

	dec := Decoder{rest: b, flags: d.flags}
	err := ds.Deserialize(&dec)
	d.rest = dec.rest
	return err == nil, err
}

type MapDecoder struct {
	first bool
	rest  []byte
	flags ParseFlags
}

func (d *MapDecoder) Size() (int, bool) {
	return 0, false
}

func (d *MapDecoder) NextKey(ds codec.Deserializer) (bool, error) {
	b := skipSpaces(d.rest)

	if len(b) == 0 {
		d.rest = b
		return false, syntaxError(b, "cannot decode object from empty input")
	}

	if b[0] == '}' {
		d.rest = b[1:]
		return false, nil
	}

	if !d.first {
		if len(b) == 0 {
			d.rest = b
			return false, syntaxError(b, "unexpected EOF after object field value")
		}
		if b[0] != ',' {
			d.rest = b
			return false, syntaxError(b, "expected ',' after object field value but found '%c'", b[0])
		}
		b = skipSpaces(b[1:])
		if len(b) == 0 {
			d.rest = b
			return false, unexpectedEOF(b)
		}
		if b[0] == '}' {
			d.rest = b
			return false, syntaxError(b, "unexpected trailing comma after object field")
		}
	} else {
		d.first = false
	}

	dec := Decoder{rest: b, flags: d.flags}
	err := ds.Deserialize(&dec)
	d.rest = dec.rest
	return err == nil, err
}

func (d *MapDecoder) NextValue(ds codec.Deserializer) error {
	b := d.rest
	if len(b) == 0 {
		d.rest = b
		return syntaxError(b, "unexpected EOF after object field key")
	}
	if b[0] != ':' {
		d.rest = b
		return syntaxError(b, "expected ':' after object field key but found '%c'", b[0])
	}
	b = skipSpaces(b[1:])

	dec := Decoder{rest: b, flags: d.flags}
	err := ds.Deserialize(&dec)
	d.rest = dec.rest
	return err
}

type Decoder struct {
	rest  []byte
	flags ParseFlags
}

func (d *Decoder) DecodeAny(v codec.Visitor) (err error) {
	dec := decoder{flags: d.flags}
	d.rest, err = dec.decodeValue(d.rest, v)
	return
}

func (d *Decoder) DecodeNil(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeBool(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeInt(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeInt8(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeInt16(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeInt32(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeInt64(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeUint(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeUint8(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeUint16(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeUint32(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeUint64(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeUintptr(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeFloat32(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeFloat64(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeComplex64(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeComplex128(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeString(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeBytes(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeOption(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeSeq(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeMap(v codec.Visitor) error {
	return d.DecodeAny(v)
}

func (d *Decoder) DecodeStruct(name string, v codec.Visitor) error {
	return d.DecodeAny(v)
}
