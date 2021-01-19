package messenger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var ResultOK = []byte(`{"status":"ok"}`)

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

var InvalidFieldMessage = "invalid content"

type InvalidField struct {
	Field   string `json:"Field"`
	Message string `json:"Msg"`
}

func MustRenderInvalidFields(w http.ResponseWriter, invalids ...string) {
	result := make([]*InvalidField, len(invalids))
	for k, v := range invalids {
		result[k] = &InvalidField{
			Field:   v,
			Message: InvalidFieldMessage,
		}
	}
	MustRenderJSON(w, result, 422)
}

var UnsupportedConditionsMessage = "unsupported conditions[ %s ]"

func MustRenderUnsupportedConditions(w http.ResponseWriter, unsupported []string) {
	result := []*InvalidField{
		&InvalidField{
			Field:   "conditions",
			Message: fmt.Sprintf(UnsupportedConditionsMessage, strings.Join(unsupported, " , ")),
		},
	}
	MustRenderJSON(w, result, 422)
}

func MustRenderOK(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	_, err := w.Write(ResultOK)
	if err != nil {
		panic(err)
	}
}
