package auth

//Result provider auth result
type Result struct {
	//Account user account
	Account string
	//Keyword keyword of provider which auth the request.
	Keyword string
	//Data user data from auth provider
	Data Profile
}

//NewResult create new auth result.
func NewResult() *Result {
	return &Result{
		Data: map[ProfileIndex][]string{},
	}
}
