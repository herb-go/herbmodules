package messenger

import (
	"net/http"
	"strconv"

	"github.com/herb-go/notification"
)

type NotificationID struct {
	NotificationID string `json:"notification-id"`
}

func MustRenderNotification(w http.ResponseWriter, n *notification.Notification) {
	header := w.Header()
	if n.ExpiredTime != 0 {
		header.Set(HeaderTTL, strconv.FormatInt(n.ExpiredTime, 10))
	}
	for k, v := range n.Header {
		header.Set(NotificationHeaderPrefix+k, v)
	}
	if n.ID != "" {
		header.Set(HeaderID, n.ID)
	}
	header.Set(HeaderDelivery, n.Delivery)
	w.WriteHeader(200)
	MustRenderJSON(w, n.Content, 200)
}
