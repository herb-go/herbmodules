package httppublish

import (
	"github.com/herb-go/herbmodules/messenger"
	"github.com/herb-go/notification"
)

type Iter struct {
	Options *messenger.ListOptions
	Store   notification.Store
	Handler func(*notification.Notification)
	Recover func()
}

func (i *Iter) Handle(n *notification.Notification) {
	defer i.Recover()
	i.Handler(n)
}
func (i *Iter) Next() {
	defer i.Recover()
	result, newiter, err := i.Store.List(i.Options.ConvertConditions(), i.Options.Iter, i.Options.Asc, i.Options.Count)
	if err != nil {
		panic(err)
	}
	for k := range result {
		n := result[k]
		i.Handle(n)
	}
	if newiter == "" {
		return
	}
	i.Options.Iter = newiter
	go i.Next()
}

func NewIter() *Iter {
	return &Iter{}
}
