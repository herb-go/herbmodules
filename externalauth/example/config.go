package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	facebookauth "github.com/herb-go/herbmodule-drivers/externalauth-drivers/facebookauth"
	githubauth "github.com/herb-go/herbmodule-drivers/externalauth-drivers/githubauth"
	wechatmpauth "github.com/herb-go/herbmodule-drivers/externalauth-drivers/wechatmpauth"
	wechatworkauth "github.com/herb-go/herbmodule-drivers/externalauth-drivers/wechatworkauth"
	windowsliveauth "github.com/herb-go/herbmodule-drivers/externalauth-drivers/windowsliveauth"
)

type APPConfig struct {
	Host            string
	Addr            string
	Github          *githubauth.OauthAuthConfig
	WechatworkOauth *wechatworkauth.OauthAuthConfig
	WechatworkQR    *wechatworkauth.QRAuthConfig
	Windowslive     *windowsliveauth.OauthAuthConfig
	Wechatmp        *wechatmpauth.OauthAuthConfig
	Facebook        *facebookauth.OauthAuthConfig
}

func (c *APPConfig) Server() *http.Server {
	server := &http.Server{
		Addr: c.Addr,
	}
	return server
}

func (c *APPConfig) MustLoadJSON(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
}

var Config = &APPConfig{}

func init() {
	Config.MustLoadJSON("config.json")
}
