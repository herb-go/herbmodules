package healthcheck

import (
	"encoding/json"
	"net/http"
)

var Hanlder = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := Check()
	w.WriteHeader(result.StatusCode())
	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(data)
	if err != nil {
		panic(err)
	}
})
