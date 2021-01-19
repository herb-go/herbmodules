package httppublish

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

func CreatePublishDraftAction(p *notificationqueue.Publisher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = r.Body.Close()
		if err != nil {
			panic(err)
		}
		nid := &messenger.NotificationID{}
		err = json.Unmarshal(bs, nid)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if nid.NotificationID == "" {
			messenger.MustRenderInvalidFields(w, "notification-id")
			return
		}
		_, err = p.PublishDraft(nid.NotificationID)
		if err != nil {
			if notification.IsErrNotificationIDNotFound(err) {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			panic(err)
		}
		messenger.MustRenderOK(w)
	})
}

func CreatePublishAllAction(p *notificationqueue.Publisher) http.Handler {
	return CreateIterAction(p, func(p *notificationqueue.Publisher) func(*notification.Notification) {
		return func(n *notification.Notification) {
			p.PublishDraft(n.ID)
		}
	})
}

func CreateDiscardAllAction(p *notificationqueue.Publisher) http.Handler {
	return CreateIterAction(p, func(p *notificationqueue.Publisher) func(*notification.Notification) {
		return func(n *notification.Notification) {
			p.Draftbox.Remove(n.ID)
		}
	})
}

func CreateIterAction(
	p *notificationqueue.Publisher,
	builder func(p *notificationqueue.Publisher) func(*notification.Notification),
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = r.Body.Close()
		if err != nil {
			panic(err)
		}
		options := messenger.NewListOptions()
		err = json.Unmarshal(bs, options)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if options.Conditions == nil {
			messenger.MustRenderInvalidFields(w, "conditions")
			return
		}
		store := p.Draftbox
		unsupported := options.MustCheckUnsupported(store)
		if len(unsupported) > 0 {
			messenger.MustRenderUnsupportedConditions(w, unsupported)
			return
		}

		iter := NewIter()
		iter.Options = options
		iter.Store = store
		iter.Recover = p.Recover
		iter.Handler = builder(p)
		go iter.Next()
		messenger.MustRenderOK(w)
	})
}
