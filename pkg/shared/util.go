package shared

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
)

func CalculateUUIDv5(data []byte) string {
	return calculateUUIDv5(ProductId, data)
}

func calculateUUIDv5(namespace string, data []byte) string {
	u := uuid.NewSHA1(uuid.Must(uuid.Parse(namespace)), data)
	return u.String()
}

func CalculateUUIDv5FromMap(m map[string]interface{}) string {
	return calculateUUIDv5FromMap(ProductId, m)
}

func calculateUUIDv5FromMap(namespace string, m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	h, err := jcs.Transform(b)
	if err != nil {
		panic(err)
	}
	return calculateUUIDv5(namespace, h)
}
