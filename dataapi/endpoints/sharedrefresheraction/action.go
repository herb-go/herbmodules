package sharedrefresheraction

import (
	"io"
	"net/http"

	"github.com/herb-go/herbdata/datautil/sharedrefresher"
)

//NewSharedRefresherAction create action for giver shared refresher.
func NewSharedRefresherAction(refresher sharedrefresher.SharedRefresher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		new, err := refresher.RefreshShared(data)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(new)
		if err != nil {
			panic(err)
		}
	}
}
