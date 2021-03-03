package httptemplate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/herbmodules/messenger/httpview"
	"github.com/herb-go/herbtext"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/notification/notificationview"
)

func validateRequest(w http.ResponseWriter, r *http.Request) *Template {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return nil
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
		return nil
	}
	invalid, err := t.Validate()
	if err != nil {
		panic(err)
	}
	if invalid != "" {
		messenger.MustRenderInvalidFields(w, invalid)
		return nil
	}
	return t
}
func TOMLAction(w http.ResponseWriter, r *http.Request) {
	t := validateRequest(w, r)
	if t == nil {
		return
	}
	_, err := w.Write(t.MustTOML())
	if err != nil {
		panic(err)
	}
}

func renderRequest(t *Template, b messenger.NotificationBuilder, r *http.Request) (*notification.Notification, error) {

	view, err := t.Parse()
	if err != nil {
		return nil, err
	}
	msg := notificationview.NewMessage()
	herbtext.MergeSet(msg, herbtext.Map(t.Data))
	n, err := view.Render(msg)
	if err != nil {
		return nil, err
	}
	b(r, n)
	return n, nil
}

func CreateRenderAction(b messenger.NotificationBuilder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		t := validateRequest(w, r)
		if t == nil {
			return
		}
		n, err := renderRequest(t, b, r)
		if !httpview.HTTPError(err, w, r) {
			return
		}
		messenger.MustRenderNotification(w, n)
	})
}

func CreateSendAction(b messenger.NotificationBuilder, sender notification.Sender) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		t := validateRequest(w, r)
		if t == nil {
			return
		}
		n, err := renderRequest(t, b, r)
		if !httpview.HTTPError(err, w, r) {
			return
		}
		err = sender.Send(n)
		if err != nil {
			panic(err)
		}
		messenger.MustRenderJSON(w, map[string]interface{}{"status": notificationdelivery.DeliveryStatusSuccess}, 200)
	})
}
