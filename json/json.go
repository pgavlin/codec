package json

import (
	"encoding/json"
	"math/bits"
	"sync"

	"github.com/pgavlin/codec"
)

// Delim is documented at https://golang.org/pkg/encoding/json/#Delim
type Delim = json.Delim

// InvalidUTF8Error is documented at https://golang.org/pkg/encoding/json/#InvalidUTF8Error
type InvalidUTF8Error = json.InvalidUTF8Error

// InvalidUnmarshalError is documented at https://golang.org/pkg/encoding/json/#InvalidUnmarshalError
type InvalidUnmarshalError = json.InvalidUnmarshalError

// Marshaler is documented at https://golang.org/pkg/encoding/json/#Marshaler
type Marshaler = json.Marshaler

// MarshalerError is documented at https://golang.org/pkg/encoding/json/#MarshalerError
type MarshalerError = json.MarshalerError

// Number is documented at https://golang.org/pkg/encoding/json/#Number
type Number = json.Number

// RawMessage is documented at https://golang.org/pkg/encoding/json/#RawMessage
type RawMessage = json.RawMessage

// A SyntaxError is a description of a JSON syntax error.
type SyntaxError = json.SyntaxError

// Token is documented at https://golang.org/pkg/encoding/json/#Token
type Token = json.Token

// UnmarshalFieldError is documented at https://golang.org/pkg/encoding/json/#UnmarshalFieldError
type UnmarshalFieldError = json.UnmarshalFieldError

// UnmarshalTypeError is documented at https://golang.org/pkg/encoding/json/#UnmarshalTypeError
type UnmarshalTypeError = json.UnmarshalTypeError

// Unmarshaler is documented at https://golang.org/pkg/encoding/json/#Unmarshaler
type Unmarshaler = json.Unmarshaler

// UnsupportedTypeError is documented at https://golang.org/pkg/encoding/json/#UnsupportedTypeError
type UnsupportedTypeError = json.UnsupportedTypeError

// UnsupportedValueError is documented at https://golang.org/pkg/encoding/json/#UnsupportedValueError
type UnsupportedValueError = json.UnsupportedValueError

// AppendFlags is a type used to represent configuration options that can be
// applied when formatting json output.
type AppendFlags uint32

const (
	// EscapeHTML is a formatting flag used to to escape HTML in json strings.
	EscapeHTML AppendFlags = 1 << iota

	// SortMapKeys is formatting flag used to enable sorting of map keys when
	// encoding JSON (this matches the behavior of the standard encoding/json
	// package).
	SortMapKeys

	// TrustRawMessage is a performance optimization flag to skip value
	// checking of raw messages. It should only be used if the values are
	// known to be valid json (e.g., they were created by json.Unmarshal).
	TrustRawMessage

	// appendNewline is a formatting flag to enable the addition of a newline
	// in Encode (this matches the behavior of the standard encoding/json
	// package).
	appendNewline
)

// ParseFlags is a type used to represent configuration options that can be
// applied when parsing json input.
type ParseFlags uint32

func (flags ParseFlags) has(f ParseFlags) bool {
	return (flags & f) != 0
}

func (f ParseFlags) kind() Kind {
	return Kind((f >> kindOffset) & 0xFF)
}

func (f ParseFlags) withKind(kind Kind) ParseFlags {
	return (f & ^(ParseFlags(0xFF) << kindOffset)) | (ParseFlags(kind) << kindOffset)
}

const (
	// DisallowUnknownFields is a parsing flag used to prevent decoding of
	// objects to Go struct values when a field of the input does not match
	// with any of the struct fields.
	DisallowUnknownFields ParseFlags = 1 << iota

	// UseNumber is a parsing flag used to load numeric values as Number
	// instead of float64.
	UseNumber

	// DontCopyString is a parsing flag used to provide zero-copy support when
	// loading string values from a json payload. It is not always possible to
	// avoid dynamic memory allocations, for example when a string is escaped in
	// the json data a new buffer has to be allocated, but when the `wire` value
	// can be used as content of a Go value the decoder will simply point into
	// the input buffer.
	DontCopyString

	// DontCopyNumber is a parsing flag used to provide zero-copy support when
	// loading Number values (see DontCopyString and DontCopyRawMessage).
	DontCopyNumber

	// DontCopyRawMessage is a parsing flag used to provide zero-copy support
	// when loading RawMessage values from a json payload. When used, the
	// RawMessage values will not be allocated into new memory buffers and
	// will instead point directly to the area of the input buffer where the
	// value was found.
	DontCopyRawMessage

	// DontMatchCaseInsensitiveStructFields is a parsing flag used to prevent
	// matching fields in a case-insensitive way. This can prevent degrading
	// performance on case conversions, and can also act as a stricter decoding
	// mode.
	DontMatchCaseInsensitiveStructFields

	// ZeroCopy is a parsing flag that combines all the copy optimizations
	// available in the package.
	//
	// The zero-copy optimizations are better used in request-handler style
	// code where none of the values are retained after the handler returns.
	ZeroCopy = DontCopyString | DontCopyNumber | DontCopyRawMessage

	// validAsciiPrint is an internal flag indicating that the input contains
	// only valid ASCII print chars (0x20 <= c <= 0x7E). If the flag is unset,
	// it's unknown whether the input is valid ASCII print.
	validAsciiPrint ParseFlags = 1 << 28

	// noBackslach is an internal flag indicating that the input does not
	// contain a backslash. If the flag is unset, it's unknown whether the
	// input contains a backslash.
	noBackslash ParseFlags = 1 << 29

	// Bit offset where the kind of the json value is stored.
	//
	// See Kind in token.go for the enum.
	kindOffset ParseFlags = 16
)

// Kind represents the different kinds of value that exist in JSON.
type Kind uint

const (
	Undefined Kind = 0

	Null Kind = 1 // Null is not zero, so we keep zero for "undefined".

	Bool  Kind = 2 // Bit two is set to 1, means it's a boolean.
	False Kind = 2 // Bool + 0
	True  Kind = 3 // Bool + 1

	Num   Kind = 4 // Bit three is set to 1, means it's a number.
	Uint  Kind = 5 // Num + 1
	Int   Kind = 6 // Num + 2
	Float Kind = 7 // Num + 3

	String    Kind = 8 // Bit four is set to 1, means it's a string.
	Unescaped Kind = 9 // String + 1

	Array  Kind = 16 // Equivalent to Delim == '['
	Object Kind = 32 // Equivalent to Delim == '{'
)

// Class returns the class of k.
func (k Kind) Class() Kind { return Kind(1 << uint(bits.Len(uint(k))-1)) }

// Append acts like Marshal but appends the json representation to b instead of
// always reallocating a new slice.
func Append(b []byte, x any, s codec.Serializer, flags AppendFlags) ([]byte, error) {
	e := Encoder{
		enc: encoder{flags: flags},
		out: b,
	}
	err := e.encode(x, s)
	return e.out, err
}

// Marshal is documented at https://golang.org/pkg/encoding/json/#Marshal
func Marshal(x any) ([]byte, error) {
	var err error
	var buf = encoderBufferPool.Get().(*encoderBuffer)

	if buf.data, err = Append(buf.data[:0], x, codec.GetSerializer(x), EscapeHTML|SortMapKeys); err != nil {
		return nil, err
	}

	b := make([]byte, len(buf.data))
	copy(b, buf.data)
	encoderBufferPool.Put(buf)
	return b, nil
}

// Unmarshal is documented at https://golang.org/pkg/encoding/json/#Unmarshal
func Unmarshal(b []byte, x any) error {
	r, err := Parse(b, x, codec.GetDeserializer(x), 0)
	if len(r) != 0 {
		if _, ok := err.(*SyntaxError); !ok {
			// The encoding/json package prioritizes reporting errors caused by
			// unexpected trailing bytes over other issues; here we emulate this
			// behavior by overriding the error.
			err = syntaxError(r, "invalid character '%c' after top-level value", r[0])
		}
	}
	return err
}

// Parse behaves like Unmarshal but the caller can pass a set of flags to
// configure the parsing behavior.
func Parse(b []byte, x any, ds codec.Deserializer, flags ParseFlags) ([]byte, error) {
	b = skipSpaces(b)
	d := &Decoder{rest: b, flags: flags | internalParseFlags(b)}
	err := d.decode(x, ds)
	return d.rest, err
}

var encoderBufferPool = sync.Pool{
	New: func() interface{} { return &encoderBuffer{data: make([]byte, 0, 4096)} },
}

type encoderBuffer struct{ data []byte }
