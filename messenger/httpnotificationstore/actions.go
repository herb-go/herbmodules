package httpnotificationstore

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/notification"
)

func CreateCountStoreAction(storeloader func() notification.Store) http.Handler {
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
		store := storeloader()
		unsupported := options.MustCheckUnsupported(store)
		if len(unsupported) > 0 {
			messenger.MustRenderUnsupportedConditions(w, unsupported)
			return
		}
		messenger.MustRenderJSON(w, options.MustCount(store), 200)
	})
}

func CreateListStoreAction(storeloader func() notification.Store) http.Handler {
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
		store := storeloader()
		unsupported := options.MustCheckUnsupported(store)
		if len(unsupported) > 0 {
			messenger.MustRenderUnsupportedConditions(w, unsupported)
			return
		}
		messenger.MustRenderJSON(w, messenger.MustList(options, store), 200)
	})
}

func CreateRemoveAction(storeloader func() notification.Store) http.Handler {
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
		store := storeloader()
		_, err = store.Remove(nid.NotificationID)
		if err != nil {
			if !notification.IsErrNotificationIDNotFound(err) {
				panic(err)
			}
		}
		messenger.MustRenderOK(w)
	})
}

func CreateStoreSupportedConditionsAction(storeloader func() notification.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supported, err := storeloader().SupportedConditions()
		if err != nil {
			panic(err)
		}
		messenger.MustRenderJSON(w, supported, 200)
	})
}
