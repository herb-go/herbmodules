package httpsession

import (
	"net/http"
	"sync"
)

type InstallerName string

type Installer interface {
	InstallerMiddleware(*Store) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type InstallerFactory func(func(v interface{}) error) (Installer, error)

var installerFactories sync.Map

func RegisterInstaller(name InstallerName, f InstallerFactory) {
	enginesFactories.Store(name, f)
}
func CreateInstaller(name InstallerName, loader func(v interface{}) error) (Installer, error) {
	v, ok := enginesFactories.Load(name)
	if !ok {
		return nil, NewUnknownSessionInstallerError(name)
	}
	return v.(InstallerFactory)(loader)
}
