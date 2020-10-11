package dchest

import (
	"crypto/md5"
	"encoding/base64"
	"time"

	"github.com/herb-go/uniqueid"
)

var DefaultSigner = func(data string) (string, error) {
	sign := md5.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(sign[:]), nil
}

var DefaultTokenGenerator = func() (string, error) {
	return uniqueid.DefaultGenerator.GenerateID()
}

var DefaultTTL = 10 * time.Minute

var DefaultSize = 8

var DefaultWidth = 200
var DefaultHeight = 50

var DefaultOptionalBytes = []byte("0123456789")
