package httpview

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbtext"
	"github.com/herb-go/notification"

	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/notification/notificationview"
)

var ErrInvalidJSON = errors.New("invalid json")

func renderRequest(c notificationview.ViewCenter, r *http.Request) (*notification.Notification, error) {
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	view, err := c.Get(p)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	data := map[string]string{}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return nil, ErrInvalidJSON
	}
	msg := notificationview.NewMessage()
	herbtext.MergeSet(msg, herbtext.Map(data))
	return view.Render(msg)
}

func httpError(err error, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		if notificationview.IsErrViewNotFound(err) {
			http.NotFound(w, r)
			return false
		} else if err == ErrInvalidJSON {
			http.Error(w, http.StatusText(400), 400)
			return false

		} else if notification.IsErrInvalidContent(err) {
			ce := err.(*notification.ErrInvalidContent)
			messenger.MustRenderInvalidFields(w, ce.Fields...)
			return false
		}
		panic(err)
	}
	return true
}
func CreateRenderAction(c notificationview.ViewCenter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		n, err := renderRequest(c, r)
		if !httpError(err, w, r) {
			return
		}
		messenger.MustRenderNotification(w, n)
	})
}

func CreateSendAction(c notificationview.ViewCenter, sender notification.Sender) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		n, err := renderRequest(c, r)
		if !httpError(err, w, r) {
			return
		}
		err = sender.Send(n)
		if err != nil {
			panic(err)
		}
		messenger.MustRenderOK(w)
	})
}
