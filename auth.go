package main

import "github.com/cagnosolutions/web"

var USER = web.Auth{
	Roles:    []string{"ADMIN", "USER"},
	Redirect: "/",
	Msg:      "You need to be logged in to do that",
}
