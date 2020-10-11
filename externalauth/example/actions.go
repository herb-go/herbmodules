package main

import (
	"fmt"
	"net/http"

	"github.com/herb-go/herb/ui/render"
)

func ActionIndex(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.NotFound(w, r)
		return
	}
	outputtemplate := `<html>
	<body>
	<p><a href='%s'>qr</a></p>
	<p><a href='%s'>oauth</a></p>
	<p><a href='%s'>github</a></p>
	<p><a href='%s'>windows</a></p>
	<p><a href='%s'>wechatmp</a></p>
	<p><a href='%s'>facebook</a></p>
	</body></html>`
	output := []byte(fmt.Sprintf(outputtemplate,
		qrProvider.LoginURL(),
		oauthProvider.LoginURL(),
		githubProvider.LoginURL(),
		windowsliveProvider.LoginURL(),
		mpProvider.LoginURL(),
		facebookProvider.LoginURL(),
	))
	render.MustWriteHTML(w, output, 200)
}

func ActionSuccess(w http.ResponseWriter, r *http.Request) {
	result := Auth.MustGetResult(r)
	render.MustJSON(w, result, 200)

}
