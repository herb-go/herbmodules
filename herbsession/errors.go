package herbsession

import "errors"

var ErrSessionNotStarted = errors.New("session not started")
var ErrInstallerNotFound = errors.New("installer not found")
