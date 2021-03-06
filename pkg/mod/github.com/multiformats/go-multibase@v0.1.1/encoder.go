package multibase

import (
	"fmt"
	"unicode/utf8"
)

// Encoder is a multibase encoding that is verified to be supported and
// supports an Encode method that does not return an error
type Encoder struct {
	enc Encoding
}

// NewEncoder create a new Encoder from an Encoding
func NewEncoder(base Encoding) (Encoder, error) {
	_, ok := EncodingToStr[base]
	if !ok {
		return Encoder{-1}, fmt.Errorf("unsupported multibase encoding: %d", base)
	}
	return Encoder{base}, nil
}

// MustNewEncoder is like NewEncoder but will panic if the encoding is
// invalid.
func MustNewEncoder(base Encoding) Encoder {
	_, ok := EncodingToStr[base]
	if !ok {
		panic("Unsupported multibase encoding")
	}
	return Encoder{base}
}

// EncoderByName creates an encoder from a string, the string can
// either be the multibase name or single character multibase prefix
func EncoderByName(str string) (Encoder, error) {
	var base Encoding
	var ok bool
	if len(str) == 0 {
		return Encoder{-1}, fmt.Errorf("empty multibase encoding")
	} else if utf8.RuneCountInString(str) == 1 {
		r, _ := utf8.DecodeRuneInString(str)
		base = Encoding(r)
		_, ok = EncodingToStr[base]
	} else {
		base, ok = Encodings[str]
	}
	if !ok {
		return Encoder{-1}, fmt.Errorf("unsupported multibase encoding: %s", str)
	}
	return Encoder{base}, nil
}

func (p Encoder) Encoding() Encoding {
	return p.enc
}

// Encode encodes the multibase using the given Encoder.
func (p Encoder) Encode(data []byte) string {
	str, err := Encode(p.enc, data)
	if err != nil {
		// should not happen
		panic(err)
	}
	return str
}
