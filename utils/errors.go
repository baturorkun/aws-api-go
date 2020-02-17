package utils

import "encoding/json"

func JsonError(err error) string {

	data := map[string]string{"Error": err.Error()}

	jsonData, _ := json.Marshal(data)

	return string(jsonData)
}
