package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetPeople",
		"GET",
		"/people/{id}",
		GetPeople,
	},
	Route{
		"TodoShow",
		"POST",
		"/login/{code}",
		VerifyUser,
	},
	Route{
		"TodoShow",
		"POST",
		"/userinfo/{code}",
		GetUserProfile,
	},
	Route{
		"TodoShow",
		"POST",
		"/match/{ouser}",
		Match,
	},
	Route{
		"TodoShow",
		"GET",
		"/test/{testType}",
		Test,
	},
}
