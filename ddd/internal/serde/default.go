package serde

import "encoding/json"

type Serder interface {
	Serialize(v any) ([]byte, error)
	Deserialize(b []byte, out any) error
}

type Serializer interface {
	Serialize() ([]byte, error)
}

type Deserializer interface {
	Deserialize(in []byte) error
}

type DefaultSerder struct{}

func (DefaultSerder) Serialize(v any) ([]byte, error) {

	return json.Marshal(v)
}

func (DefaultSerder) Deserialize(b []byte, out any) error {

	return json.Unmarshal(b, out)
}
