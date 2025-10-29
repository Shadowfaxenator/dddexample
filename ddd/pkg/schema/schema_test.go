package schema

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func NewCustomerSchema() *Schema {
	return &Schema{
		"ID":       {Type: String},
		"Name":     {Type: String},
		"Age":      {Type: Number},
		"Birthday": {Type: Date},
		"Orders":   {Type: Array},
	}
}

func TestSchema(t *testing.T) {
	order1 := Object{
		"ID": uuid.New().String(),
	}
	order2 := Object{
		"ID": uuid.New().String(),
	}

	customer := Object{
		"ID":       uuid.New().String(),
		"Test":     444,
		"Name":     "Bob",
		"Age":      22,
		"Birthday": time.Now(),
		"Orders": []Object{
			order1,
			order2,
		},
	}

	b, _ := json.MarshalIndent(customer, "", "   ")
	fmt.Printf("b: %s\n", b)

	sch := NewCustomerSchema()
	cust := NewCustomType(sch)
	err := json.Unmarshal(b, &cust)
	if err != nil {
		panic(err)
	}

	fmt.Printf("cust: %+v\n", cust.Object)

}
