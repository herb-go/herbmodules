package healthcheck

type Result struct {
	Status   Status
	Msgs     *Infos `json:",omitempty"`
	Errors   *Infos `json:",omitempty"`
	Warnings *Infos `json:",omitempty"`
}

func (r *Result) StatusCode() int {
	switch r.Status {
	case StatusHealthy:
		return StatusCodeHealthy
	case StatusWarning:
		return StatusCodeWarning
	}
	return StatusCodeError
}
func NewResult() *Result {
	return &Result{
		Status: StatusHealthy,
	}
}

type Status int

var StatusHealthy = Status(0)
var StatusWarning = Status(1)
var StatusError = Status(-1)

var StatusCodeHealthy = int(200)
var StatusCodeWarning = int(200)
var StatusCodeError = int(503)
