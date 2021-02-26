package httptemplate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbmodules/messenger"
)

func ValidateAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}
	t := &Template{}
	err = json.Unmarshal(bs, t)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	invalid, err := t.Validate()
	if err != nil {
		panic(err)
	}
	if invalid != "" {
		messenger.MustRenderInvalidFields(w, invalid)
		return
	}
	_, err = w.Write(t.MustTOML())
	if err != nil {
		panic(err)
	}
}
