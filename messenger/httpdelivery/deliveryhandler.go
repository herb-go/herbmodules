package httpdelivery

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/herb-go/herbmodules/messenger"

	"github.com/herb-go/notification"

	"github.com/herb-go/notification/notificationdelivery"
)

type DeliveryResult struct {
	Status notificationdelivery.DeliveryStatus `json:"status"`
	Msg    string                              `json:"msg"`
}
type DeliveryHandler struct {
	Center notificationdelivery.DeliveryCenter
}

func (h *DeliveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	p := r.URL.Path
	if p[0] == '/' {
		p = p[1:]
	}
	d, err := h.Center.Get(p)
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
		messenger.MustRenderInvalidContents(w, invalids)
		return
	}
	status, receipt, err := d.Deliver(content)
	if err != nil {
		panic(err)
	}
	messenger.MustRenderJSON(w, &DeliveryResult{Status: status, Msg: receipt}, 200)
}

func CreateDeliveryHandler(center notificationdelivery.DeliveryCenter) *DeliveryHandler {
	return &DeliveryHandler{
		Center: center,
	}
}
