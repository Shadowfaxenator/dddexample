package schema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type EventType uint

const (
	Update EventType = iota
	Insert
)

type Event struct {
	FieldName string
	Type      EventType
	Value     any
}

type Kind uint

const (
	String Kind = iota + 1
	Number
	Date
	Array
)

func (t Kind) String() string {
	switch t {
	case Array:
		return "Array"
	case Date:
		return "Date"
	case Number:
		return "Number"
	case String:
		return "String"
	}
	return ""
}

func (t Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Kind) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	switch str {
	case "Array":
		*t = Array
	case "Date":
		*t = Date
	case "Number":
		*t = Number
	case "String":
		*t = String
	default:
		return fmt.Errorf("invalid status string: %s", str)
	}
	return nil
}

type Field struct {
	Type   Kind
	MaxL   *uint
	MinL   *int
	Patern regexp.Regexp
}
type Schema map[string]Field
type Object map[string]any

type CustomType struct {
	Object
	index []any `json:"-"`
	*Schema
}

func (t *CustomType) UnmarshalJSON(b []byte) error {
	t.newObject()
	if err := json.Unmarshal(b, &t.Object); err != nil {
		return err
	}

	return nil

}

func NewCustomType(s *Schema) *CustomType {
	return &CustomType{Schema: s}
}

func (s *CustomType) newObject() {

	s.Object = Object{}
	for k, v := range *s.Schema {
		switch v.Type {
		case Number:
			s.Object[k] = 0
		case String:
			s.Object[k] = ""
		case Array:
			s.Object[k] = []Object{}
		case Date:
			s.Object[k] = time.Time{}
		}
	}
}

// func (e Entity) MarshalJSON() ([]byte, error) {
// 	m := make(map[string]any)

// 	for _, v := range e {
// 		m[v.Name] = v.Value
// 	}
// 	return json.Marshal(m)
// 	//json.Indent()
// }

// func (f Field) MarshalJSON() ([]byte, error) {
// 	// if f.kind == Date {
// 	// 	return json.Marshal(f.Value.(time.Time))
// 	// }
// 	return json.Marshal(f.Value)
// 	//json.Indent()
// }

// func (f *Field) UnmarshalJSON(b []byte) error {
// 	fmt.Printf("f.kind: %v\n", *f)
// 	switch f.kind {
// 	case Number:

// 		var value float64
// 		if err := json.Unmarshal(b, &value); err != nil {
// 			return err
// 		}
// 		f.Value = value
// 		return nil
// 	case Date:
// 		var value time.Time
// 		if err := json.Unmarshal(b, &value); err != nil {
// 			return err
// 		}
// 		f.Value = value
// 		return nil
// 	// case EntitySlice:
// 	// 	fmt.Println("Slice")
// 	// 	f.kind = Array
// 	// 	return nil
// 	default:
// 	}
// 	return nil
// 	//json.Indent()
// }
