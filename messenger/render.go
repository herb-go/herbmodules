package messenger

import (
	"encoding/json"
	"net/http"
)

func MustRenderJSON(w http.ResponseWriter, data interface{}, code int) {
	bs, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(bs)
	if err != nil {
		panic(err)
	}
}

var InvalidContentMessage = "invalid content"

type InvalidContentField struct {
	Field   string `json:"field"`
	Message string `json:"msg"`
}

func MustRenderInvalidContents(w http.ResponseWriter, invalids []string) {
	result := make([]*InvalidContentField, len(invalids))
	for k, v := range invalids {
		result[k] = &InvalidContentField{
			Field:   v,
			Message: InvalidContentMessage,
		}
	}
	MustRenderJSON(w, result, 422)
}
