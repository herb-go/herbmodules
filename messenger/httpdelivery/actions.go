package httpdelivery

import (
	"net/http"

	"github.com/herb-go/herbmodules/messenger"

	"github.com/herb-go/notification/notificationdelivery"
)

type DeliveryServerOutput struct {
	Delivery     string `json:"delivery"`
	DeliveryType string `json:"delivery-type"`
	Disabled     bool   `json:"disabled"`
	Description  string `json:"description"`
}

func CreateListDeliveryServersAction(c notificationdelivery.DeliveryCenter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		list, err := c.List()
		if err != nil {
			panic(err)
		}
		output := make([]*DeliveryServerOutput, len(list))
		for k, v := range list {
			output[k] = &DeliveryServerOutput{
				Delivery:     v.Delivery,
				DeliveryType: v.DeliveryType(),
				Disabled:     v.Disabled,
				Description:  v.Description,
			}
		}
		messenger.MustRenderJSON(w, output, 200)
	})
}
