package healthcheck

var registered = []Checker{}

func Reset() {
	registered = []Checker{}
}

type Checker func() (Status, *Info)

func OnCheck(c Checker) {
	registered = append(registered, c)
}

func Check() *Result {
	result := NewResult()
	msgs := NewInfos()
	warnings := NewInfos()
	errmsgs := NewInfos()
	for k := range registered {
		status, info := registered[k]()
		switch status {
		case StatusHealthy:
			msgs.Append(info)
		case StatusWarning:
			if result.Status == StatusHealthy {
				result.Status = StatusWarning
			}
			warnings.Append(info)
		default:
			if result.Status == StatusHealthy || result.Status == StatusWarning {
				result.Status = StatusError
			}
			errmsgs.Append(info)
		}
	}
	result.Msgs = msgs.Data()
	result.Warnings = warnings.Data()
	result.Errors = errmsgs.Data()
	return result
}
