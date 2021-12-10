package httpsession

import (
	"errors"
	"fmt"
)

var ErrUnknownSessionEngine = errors.New("unknown session engine")

func NewUnknownSessionEngineError(name EngineName) error {
	return fmt.Errorf("%w [%s]", ErrUnknownSessionEngine, name)
}

var ErrUnknownSessionInstaller = errors.New("unknown session installer")

func NewUnknownSessionInstallerError(name InstallerName) error {
	return fmt.Errorf("%w [%s]", ErrUnknownSessionInstaller, name)
}
