package herbsession

import (
	"net/http"
)

type InstallerID string

type Installer interface {
	InstallerMiddleware(*Store) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	InstallerID() InstallerID
}
