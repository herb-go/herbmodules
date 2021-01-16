package messenger

import (
	"net/http"

	"github.com/herb-go/notification"
)

//NotificationBuilder build notification with given requret.
//Panic if any error raised
type NotificationBuilder func(r *http.Request, n *notification.Notification)
