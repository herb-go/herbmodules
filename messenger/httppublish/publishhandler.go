package httppublish

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/herb-go/notification/notificationdelivery/notificationqueue"

	"github.com/herb-go/herbmodules/messenger"

	"github.com/herb-go/notification"

	"github.com/herb-go/notification/notificationdelivery"
)

type PublisherResult struct {
	NotificationID string `json:"notification-id"`
	Published      bool   `json:"published"`
}

type PublisherHandler struct {
	Publisher *notificationqueue.Publisher
	Builder   messenger.NotificationBuilder
}

func LoadNotificationHeader(h http.Header) notification.Header {
	result := notification.NewHeader()
	for k := range h {
		name := strings.ToLower(k)
		if strings.HasPrefix(name, messenger.NotificationHeaderPrefix) {
			result.Set(strings.TrimPrefix(name, messenger.NotificationHeaderPrefix), h.Get(k))
		}
	}
	return result
}
func (h *PublisherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	d, err := h.Publisher.DeliveryCenter.Get(p)
	if err != nil {
		if notificationdelivery.IsErrDeliveryNotFound(err) {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		panic(err)
	}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}
	content := notification.NewContent()
	err = json.Unmarshal(bs, &content)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	invalids, err := d.CheckInvalidContent(content)
	if err != nil {
		panic(err)
	}
	if len(invalids) > 0 {
		messenger.MustRenderInvalidFields(w, invalids...)
		return
	}
	n := notification.New()
	n.Delivery = p
	n.Content = content
	n.Header = LoadNotificationHeader(r.Header)
	if h.Builder != nil {
		h.Builder(r, n)
	}
	ttlheader := r.Header.Get(messenger.HeaderTTL)
	if ttlheader != "" {
		i, err := strconv.Atoi(ttlheader)
		if err == nil {
			n.ExpiredTime = time.Now().Add(time.Duration(i) * time.Second).Unix()
		}
	}
	if n.ExpiredTime <= 0 {
		n.ExpiredTime = time.Now().Add(notification.SuggestedNotificationTTL).Unix()
	}
	nid, published, err := h.Publisher.PublishNotification(n)
	if err != nil {
		panic(err)
	}

	messenger.MustRenderJSON(w, &PublisherResult{NotificationID: nid, Published: published}, 200)
}

func CreatePublishHandler(publisher *notificationqueue.Publisher, builder messenger.NotificationBuilder) *PublisherHandler {
	return &PublisherHandler{
		Publisher: publisher,
		Builder:   builder,
	}
}
