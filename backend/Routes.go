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
		"Login",
		"POST",
		"/login/{code}",
		VerifyUser,
	},
	Route{
		"LinkedIn Profile",
		"GET",
		"/userinfo/{accessToken}",
		GetUserProfile,
	},
	Route{
		"Matching",
		"POST",
		"/match/{ouser}",
		Match,
	},
	Route{
		"Debugging Endpoint",
		"POST",
		"/test/{testType}",
		Test,
	},
}
