package messenger

import (
	"github.com/herb-go/notification"
)

type ConditionOption struct {
	Keyword string `json:"keyword"`
	Value   string `json:"Value"`
}

type ListOptions struct {
	Conditions *[]*ConditionOption `json:"conditions"`
	Asc        bool                `json:"asc"`
	Iter       string              `json:"iter"`
	Count      int                 `json:"count"`
}

func (o *ListOptions) ConvertConditions() []*notification.Condition {
	conds := make([]*notification.Condition, len(*o.Conditions))
	for k, v := range *o.Conditions {
		conds[k] = &notification.Condition{
			Keyword: v.Keyword,
			Value:   v.Value,
		}
	}
	return conds
}
func (o *ListOptions) MustCheckUnsupported(store notification.Searchable) []string {
	var result []string
	supported, err := store.SupportedConditions()
	if err != nil {
		panic(err)
	}
	supportedmap := make(map[string]bool, len(supported))
	for _, v := range supported {
		supportedmap[v] = true
	}
	for _, v := range *o.Conditions {
		if !supportedmap[v.Keyword] {
			result = append(result, v.Keyword)
		}
	}
	return result
}
func (o *ListOptions) MustCount(s notification.Searchable) *CountResult {
	conds := make([]*notification.Condition, len(*o.Conditions))
	for k, v := range *o.Conditions {
		conds[k] = &notification.Condition{
			Keyword: v.Keyword,
			Value:   v.Value,
		}
	}
	count, err := s.Count(conds)
	if err != nil {
		panic(err)
	}
	return &CountResult{
		Count: count,
	}
}
func MustList(o *ListOptions, s notification.Store) *ListResult {

	list, iter, err := s.List(o.ConvertConditions(), o.Iter, o.Asc, o.Count)
	if err != nil {
		panic(err)
	}
	return CreateListResult(list, iter)
}
func NewListOptions() *ListOptions {
	return &ListOptions{}
}

type NotificationOutput struct {
	ID          string               `json:"id"`
	Delivery    string               `json:"delivery"`
	CreatedTime int64                `json:"createdtime"`
	ExpiredTime int64                `json:"expiredtime"`
	Header      notification.Header  `json:"header"`
	Content     notification.Content `json:"content"`
}

func ConvertNotification(n *notification.Notification) *NotificationOutput {
	return &NotificationOutput{
		ID:          n.ID,
		Delivery:    n.Delivery,
		CreatedTime: n.CreatedTime,
		ExpiredTime: n.ExpiredTime,
		Header:      n.Header,
		Content:     n.Content,
	}
}

type ListResult struct {
	Iter   string                `json:"iter"`
	Result []*NotificationOutput `json:"result"`
}

func CreateListResult(notifications []*notification.Notification, iter string) *ListResult {
	result := &ListResult{}
	result.Iter = iter
	result.Result = make([]*NotificationOutput, len(notifications))
	for k := range notifications {
		result.Result[k] = ConvertNotification(notifications[k])
	}
	return result
}

type CountResult struct {
	Count int `json:"count"`
}
