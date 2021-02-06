package httpview

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbtext"
	"github.com/herb-go/notification"

	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/notification/notificationview"
)

func CreateRenderAction(c notificationview.ViewCenter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		p := r.URL.Path
		if p[0] == '/' {
			p = p[1:]
		}
		view, err := c.Get(p)
		if err != nil {
			if notificationview.IsErrViewNotFound(err) {
				http.NotFound(w, r)
				return
			}
			panic(err)
		}
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		data := map[string]string{}
		err = json.Unmarshal(bs, &data)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		msg := notificationview.NewMessage()
		herbtext.MergeSet(msg, herbtext.Map(data))
		n, err := view.Render(msg)
		if err != nil {
			panic(err)
		}
		messenger.MustRenderJSON(w, messenger.ConvertNotification(n), 200)
	})
}

func CreateSendAction(c notificationview.ViewCenter, sender notification.Sender) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		p := r.URL.Path
		if p[0] == '/' {
			p = p[1:]
		}
		view, err := c.Get(p)
		if err != nil {
			if notificationview.IsErrViewNotFound(err) {
				http.NotFound(w, r)
				return
			}
			panic(err)
		}
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		data := map[string]string{}
		err = json.Unmarshal(bs, &data)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		msg := notificationview.NewMessage()
		herbtext.MergeSet(msg, herbtext.Map(data))
		n, err := view.Render(msg)
		if err != nil {
			if notification.IsErrInvalidContent(err) {
				ce := err.(*notification.ErrInvalidContent)
				messenger.MustRenderInvalidFields(w, ce.Fields...)
				return
			}
			panic(err)
		}
		err = sender.Send(n)
		if err != nil {
			panic(err)
		}
		messenger.MustRenderOK(w)
	})
}
