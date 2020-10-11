package main

import "github.com/herb-go/herbmodules/externalauth/testsession"

var Session = &testsession.TestCookieSession{
	Name: "herbgoexternalauthlogin",
	Path: "/",
}
