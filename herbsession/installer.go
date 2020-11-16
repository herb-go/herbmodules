package herbsession

import (
	"net/http"
)

type RequestInstaller interface {
	InstallToRequest(Session) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
