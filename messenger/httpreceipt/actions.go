package httpreceipt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue/receiptstore"
)

func CreateCountStoreAction(storeloader func() receiptstore.ReceiptStore) http.Handler {
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

func CreateListStoreAction(storeloader func() receiptstore.ReceiptStore) http.Handler {
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
		messenger.MustRenderJSON(w, MustList(options, store), 200)
	})
}

type Output struct {
	ID          string               `json:"id"`
	Delivery    string               `json:"delivery"`
	CreatedTime int64                `json:"createdtime"`
	ExpiredTime int64                `json:"expiredtime"`
	Message     string               `json:"message"`
	Status      int64                `json:"status"`
	Header      notification.Header  `json:"header"`
	Content     notification.Content `json:"content"`
}

func ConvertOutput(r *notificationqueue.Receipt) *Output {
	return &Output{
		Status:      int64(r.Status),
		Message:     r.Message,
		ID:          r.Notification.ID,
		Delivery:    r.Notification.Delivery,
		CreatedTime: r.Notification.CreatedTime,
		ExpiredTime: r.Notification.ExpiredTime,
		Header:      r.Notification.Header,
		Content:     r.Notification.Content,
	}
}

type ListResult struct {
	Iter   string    `json:"iter"`
	Result []*Output `json:"result"`
}

func MustList(o *messenger.ListOptions, s receiptstore.ReceiptStore) *ListResult {
	c := o.Count
	if c == 0 {
		c = notification.DefaultStoreListLimit
	}
	list, iter, err := s.List(o.ConvertConditions(), o.Iter, o.Asc, c+1)
	if err != nil {
		panic(err)
	}

	if len(list) <= c {
		iter = ""
		c = len(list)
	} else {
		iter = list[c-1].Notification.ID
	}
	return CreateListResult(list[:c], iter)
}

func CreateListResult(receipts []*notificationqueue.Receipt, iter string) *ListResult {
	result := &ListResult{}
	result.Iter = iter
	result.Result = make([]*Output, len(receipts))
	for k := range receipts {
		result.Result[k] = ConvertOutput(receipts[k])
	}
	return result
}

func CreateStoreSupportedConditionsAction(storeloader func() receiptstore.ReceiptStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supported, err := storeloader().SupportedConditions()
		if err != nil {
			panic(err)
		}
		messenger.MustRenderJSON(w, supported, 200)
	})
}

func CreateFlushAction(storeloader func() receiptstore.ReceiptStore, recover func()) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := storeloader()
		d, err := s.RetentionDays()
		if err != nil {
			panic(err)
		}

		if d > 0 {
			go func() {
				defer recover()
				retention := time.Now().Add(time.Duration(d) * 24 * time.Hour)
				c := []*notification.Condition{
					{
						Keyword: notification.ConditionBeforeTimestamp,
						Value:   strconv.FormatInt(retention.Unix(), 10),
					},
				}
				var iter string
				var results []*notificationqueue.Receipt
				for {
					results, iter, err = s.List(c, iter, true, notification.DefaultStoreListLimit)
					if err != nil {
						panic(err)
					}
					for _, v := range results {
						result := v
						go func() {
							defer recover()
							s.Remove(result.Notification.ID)
						}()
					}
					if iter == "" {
						return
					}
				}
			}()
		}
		messenger.MustRenderOK(w)
	})
}
