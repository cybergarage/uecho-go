package mra

import (
	"bytes"
	"testing"
)

func TestLightingSchema(t *testing.T) {
	ps := Properties(0x029001, 'R')
	if Name(0x029001) != "General lighting" {
		t.Fatal("class")
	}
	state := ps[0x80]
	if !state.Editable() || state.Decode([]byte{0x30}) != "On (30)" || state.Validate([]byte{0x32}) == nil {
		t.Fatalf("state %+v", state)
	}
	number := ps[0xb0]
	f, ok := number.InputField()
	if !ok || !number.Editable() {
		t.Fatalf("number %+v", number)
	}
	for _, v := range []string{"-1", "101", "50.5", "NaN", "Inf"} {
		if _, err := f.Encode(v); err == nil {
			t.Fatal(v)
		}
	}
	b, err := f.Encode("75")
	if err != nil || !bytes.Equal(b, []byte{75}) {
		t.Fatal(b, err)
	}
	if _, ok := ps[0xff]; ok {
		t.Fatal("unknown property inferred")
	}
}
func TestSignedScaledNumber(t *testing.T) {
	b := db.Definitions["number_-273.2-3276.6Celsius"]
	fs := fields(b, 0)
	if len(fs) != 1 {
		t.Fatal(fs)
	}
	f := fs[0]
	v, e := f.Encode("-12.3")
	if e != nil || !bytes.Equal(v, []byte{0xff, 0x85}) {
		t.Fatal(v, e)
	}
	if text, ok := f.decode(v); !ok || text != "-12.3 Celsius" {
		t.Fatal(text, ok)
	}
	if _, e = f.Encode("-12.35"); e == nil {
		t.Fatal("fractional raw integer")
	}
}
func TestConservativeSchemas(t *testing.T) {
	for _, b := range []string{`{"type":"object"}`, `{"$ref":"#/definitions/missing"}`, `{"type":"number","format":"uint8","minimum":0,"maximum":100,"coefficient":["0xE1"]}`, `{"oneOf":[{"type":"raw","minSize":1,"maxSize":1},{"type":"time","size":3}]}`} {
		if fields([]byte(b), 0) != nil {
			t.Fatal(b)
		}
	}
	ps := Properties(0x029001, 'Z')
	if ps[0xb0].Editable() {
		t.Fatal("future release incorrectly inferred")
	}
	if Properties(0x029001, 0)[0xb1].Editable() {
		t.Fatal("historical schema guessed without release")
	}
	unknown := Properties(0xeeee01, 0)
	if unknown[0xb0].Editable() || unknown[0x80].Editable() {
		t.Fatal("unknown class-specific property inferred")
	}
}
