package healthcheck

type Code int

type Info struct {
	Msg  string
	Code *Code
	Data interface{}
}

func (i *Info) WithMsg(msg string) *Info {
	i.Msg = msg
	return i
}

func (i *Info) WithCode(code *Code, data interface{}) *Info {
	if code != nil {
		i.Code = code
		i.Data = data
	}
	return i
}

func NewInfo() *Info {
	return &Info{}
}

type Infos []*Info

func (i *Infos) Data() *Infos {
	if i == nil || len(*i) == 0 {
		return nil
	}
	return i
}
func (i *Infos) Append(info *Info) {
	if info != nil {
		*i = append(*i, info)
	}
}
func NewInfos() *Infos {
	return &Infos{}
}
