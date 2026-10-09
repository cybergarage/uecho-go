// Package mra adapts the MIT-licensed official MRA 1.3.0 snapshot. It deliberately
// supports only state, bounded integer number, raw and their simple oneOf unions.
package mra

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

//go:embed data.json
var database []byte

type text struct {
	EN string `json:"en"`
}
type release struct{ From, To string }
type choice struct {
	EDT          string `json:"edt"`
	Name         string
	Descriptions text
}
type schema struct {
	Type, Ref, Format, Unit string
	Size, MinSize, MaxSize  int
	Minimum, Maximum        *float64
	Multiple                float64
	Enum                    []choice
	OneOf                   []json.RawMessage
}

func (s *schema) UnmarshalJSON(b []byte) error {
	type plain schema
	var p struct {
		plain
		Ref string `json:"$ref"`
	}
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	*s = schema(p.plain)
	s.Ref = p.Ref
	return nil
}

type property struct {
	EPC          string
	PropertyName text
	ValidRelease release
	Data         json.RawMessage
	AccessRule   struct{ Get, Set, Inf string }
}
type class struct {
	EOJ          string
	ClassName    text
	ValidRelease release
	ElProperties []property
}
type snapshot struct {
	Definitions map[string]json.RawMessage
	Classes     []class
}

var db snapshot

func init() {
	if err := json.Unmarshal(database, &db); err != nil {
		panic(err)
	}
}

type Option struct {
	Label string
	Data  []byte
}
type Field struct {
	Kind, Unit       string
	Size             int
	Min, Max, Scale  float64
	Options          []Option
	MinSize, MaxSize int
}
type Property struct {
	EPC          byte
	Name, Reason string
	Fields       []Field
	SetDefined   bool
}

func Name(eoj uint32) string {
	for _, c := range db.Classes {
		v, _ := strconv.ParseUint(strings.TrimPrefix(c.EOJ, "0x"), 16, 16)
		if uint32(v) == eoj>>8 {
			return c.ClassName.EN
		}
	}
	return "Unknown class"
}
func Properties(eoj uint32, rel byte) map[byte]Property {
	grouped := map[byte][]property{}
	// Superclass is followed by the class definition: a class override wins.
	for _, code := range []uint32{0, eoj >> 8} {
		for _, c := range db.Classes {
			v, _ := strconv.ParseUint(strings.TrimPrefix(c.EOJ, "0x"), 16, 16)
			if uint32(v) != code || (rel != 0 && !applies(c.ValidRelease, rel)) {
				continue
			}
			local := map[byte][]property{}
			for _, p := range c.ElProperties {
				ep, err := strconv.ParseUint(strings.TrimPrefix(p.EPC, "0x"), 16, 8)
				if err == nil {
					local[byte(ep)] = append(local[byte(ep)], p)
				}
			}
			for ep, ps := range local {
				grouped[ep] = ps
			}
		}
	}
	out := map[byte]Property{}
	for ep, ps := range grouped {
		p := Property{EPC: ep, Name: ps[0].PropertyName.EN, Reason: "MRA release/schema ambiguous"}
		candidates := []property{}
		for _, v := range ps {
			if rel == 0 || applies(v.ValidRelease, rel) {
				candidates = append(candidates, v)
			}
		}
		if len(candidates) == 0 {
			p.Reason = "MRA definition unavailable for device release"
			out[ep] = p
			continue
		}
		base := candidates[0]
		same := true
		if rel == 0 {
			for r := byte('A'); r <= 'R'; r++ {
				covered := false
				for _, v := range candidates {
					if applies(v.ValidRelease, r) {
						covered = true
					}
				}
				if !covered {
					same = false
				}
			}
		}
		for _, v := range candidates[1:] {
			if !bytes.Equal(v.Data, base.Data) || v.AccessRule.Set != base.AccessRule.Set {
				same = false
			}
		}
		if same {
			p.Name = base.PropertyName.EN
			p.SetDefined = base.AccessRule.Set != "" && base.AccessRule.Set != "notApplicable"
			p.Fields = fields(base.Data, 0)
			p.Reason = ""
			if len(p.Fields) == 0 {
				p.Reason = "MRA type unsupported; raw display only"
			}
		}
		out[ep] = p
	}
	if Name(eoj) == "Unknown class" {
		for ep, p := range out {
			p.Fields = nil
			p.SetDefined = false
			p.Reason = "unknown class; raw display only"
			out[ep] = p
		}
	}
	return out
}
func applies(r release, rel byte) bool {
	return rel >= 'A' && rel <= 'R' && (r.From == "" || rel >= r.From[0]) && (r.To == "" || r.To == "latest" || rel <= r.To[0])
}
func fields(b json.RawMessage, depth int) []Field {
	if depth > 12 {
		return nil
	}
	var s schema
	if json.Unmarshal(b, &s) != nil {
		return nil
	}
	if s.Ref != "" {
		return fields(db.Definitions[strings.TrimPrefix(s.Ref, "#/definitions/")], depth+1)
	}
	if len(s.OneOf) > 0 {
		out := []Field{}
		for _, v := range s.OneOf {
			f := fields(v, depth+1)
			if len(f) == 0 {
				return nil
			}
			out = append(out, f...)
		}
		return out
	}
	f := Field{Kind: s.Type, Unit: s.Unit, Size: s.Size, MinSize: s.MinSize, MaxSize: s.MaxSize, Scale: s.Multiple}
	if f.Scale == 0 {
		f.Scale = 1
	}
	switch s.Type {
	case "state":
		if s.Size < 1 || s.Size > 255 || len(s.Enum) == 0 {
			return nil
		}
		for _, c := range s.Enum {
			b, e := hex.DecodeString(strings.TrimPrefix(c.EDT, "0x"))
			if e != nil || len(b) != s.Size {
				return nil
			}
			label := c.Descriptions.EN
			if label == "" {
				label = c.Name
			}
			f.Options = append(f.Options, Option{fmt.Sprintf("%s (%X)", label, b), b})
		}
	case "number":
		var attributes map[string]json.RawMessage
		_ = json.Unmarshal(b, &attributes)
		for _, key := range []string{"coefficient", "enum", "multipleOf"} {
			if _, ok := attributes[key]; ok {
				return nil
			}
		}
		switch s.Format {
		case "uint8":
			f.Size = 1
		case "int8":
			f.Size = -1
		case "uint16":
			f.Size = 2
		case "int16":
			f.Size = -2
		case "uint32":
			f.Size = 4
		case "int32":
			f.Size = -4
		default:
			return nil
		}
		if s.Minimum == nil || s.Maximum == nil || f.Scale <= 0 {
			return nil
		}
		f.Min = *s.Minimum
		f.Max = *s.Maximum
	case "raw":
		if f.MinSize < 1 || f.MaxSize < f.MinSize || f.MaxSize > 255 {
			return nil
		}
	default:
		return nil
	}
	return []Field{f}
}
func (p Property) Editable() bool { return p.SetDefined && len(p.Fields) > 0 }
func (p Property) Options() []Option {
	o := []Option{}
	for _, f := range p.Fields {
		o = append(o, f.Options...)
	}
	return o
}
func (p Property) InputField() (Field, bool) {
	for _, f := range p.Fields {
		if f.Kind == "number" || f.Kind == "raw" {
			return f, true
		}
	}
	return Field{}, false
}
func (p Property) Validate(b []byte) error {
	for _, f := range p.Fields {
		if _, ok := f.decode(b); ok {
			return nil
		}
	}
	return fmt.Errorf("EDT is outside the supported MRA schema")
}
func (f Field) decode(b []byte) (string, bool) {
	switch f.Kind {
	case "state":
		for _, o := range f.Options {
			if bytes.Equal(b, o.Data) {
				return o.Label, true
			}
		}
	case "raw":
		if len(b) >= f.MinSize && len(b) <= f.MaxSize {
			return fmt.Sprintf("raw %X", b), true
		}
	case "number":
		size := f.Size
		if size < 0 {
			size = -size
		}
		if len(b) != size {
			return "", false
		}
		var n int64
		for _, v := range b {
			n = n<<8 | int64(v)
		}
		if f.Size < 0 && b[0]&0x80 != 0 {
			n -= int64(1) << (size * 8)
		}
		if float64(n) < f.Min || float64(n) > f.Max {
			return "", false
		}
		return fmt.Sprintf("%s %s", strconv.FormatFloat(float64(n)*f.Scale, 'f', -1, 64), f.Unit), true
	}
	return "", false
}
func (p Property) Decode(b []byte) string {
	for _, f := range p.Fields {
		if v, ok := f.decode(b); ok {
			return v
		}
	}
	if len(b) == 0 {
		return "unread"
	}
	return "raw / unknown"
}
func (f Field) Encode(input string) ([]byte, error) {
	if f.Kind == "raw" {
		b, e := hex.DecodeString(strings.Join(strings.Fields(input), ""))
		if e == nil {
			if _, ok := f.decode(b); ok {
				return b, nil
			}
		}
		return nil, fmt.Errorf("enter %d..%d raw bytes", f.MinSize, f.MaxSize)
	}
	if f.Kind != "number" {
		return nil, fmt.Errorf("unsupported editor")
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	n := value / f.Scale
	if err != nil || math.IsNaN(n) || math.IsInf(n, 0) || math.Abs(n-math.Round(n)) > 1e-7 || n < f.Min || n > f.Max {
		return nil, fmt.Errorf("enter %.6g..%.6g %s in steps of %.6g", f.Min*f.Scale, f.Max*f.Scale, f.Unit, f.Scale)
	}
	size := f.Size
	lo, hi := float64(0), math.Pow(2, float64(size*8))-1
	if size < 0 {
		size = -size
		lo = -math.Pow(2, float64(size*8-1))
		hi = -lo - 1
	}
	if n < lo || n > hi {
		return nil, fmt.Errorf("number cannot fit schema format")
	}
	b := make([]byte, size)
	bits := uint64(int64(math.Round(n)))
	for i := size - 1; i >= 0; i-- {
		b[i] = byte(bits)
		bits >>= 8
	}
	return b, nil
}
func Codes(ps map[byte]Property) []byte {
	out := []byte{}
	for ep := range ps {
		out = append(out, ep)
	}
	slices.Sort(out)
	return out
}
