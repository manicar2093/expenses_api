package json

import (
	"encoding/json"
	"fmt"
)

func MustMarshall(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		message := fmt.Errorf("marshall has failed: %v", err)
		log.Println(message)
		panic(message)
	}

	return string(bytes)
}
