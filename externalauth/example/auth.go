package main

import (
	auth "github.com/herb-go/herbmodules/externalauth"
	"github.com/herb-go/util"
)

var Auth = auth.New()

var qrProvider *auth.Provider
var oauthProvider *auth.Provider
var githubProvider *auth.Provider
var windowsliveProvider *auth.Provider
var mpProvider *auth.Provider
var facebookProvider *auth.Provider

func InitAuth() {
	option := auth.OptionCommon(Config.Host+"/auth", Session)
	qrProvider = Auth.MustRegisterWithCreator("qr", Config.WechatworkQR)
	oauthProvider = Auth.MustRegisterWithCreator("oauth", Config.WechatworkOauth)
	githubProvider = Auth.MustRegisterWithCreator("github", Config.Github)
	windowsliveProvider = Auth.MustRegisterWithCreator("windowslive", Config.Windowslive)
	mpProvider = Auth.MustRegisterWithCreator("wechatmp", Config.Wechatmp)
	facebookProvider = Auth.MustRegisterWithCreator("facebook", Config.Facebook)
	util.Must(Auth.Init(option))
}
