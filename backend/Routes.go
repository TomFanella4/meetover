package main

import "net/http"

// Route - REST endpoint route used by REST router
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes for REST router
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Login",
		"POST",
		"/login/{code}",
		VerifyUser,
	},
	Route{
		"Refresh",
		"POST",
		"/refreshtoken",
		RefreshCustomToken,
	},
	Route{
		"Initiate Meetover",
		"POST",
		"/meetover/{otherId}",
		InitiateMeetover,
	},
	Route{
		"Get another user's profile",
		"GET",
		"/userprofile/{id}",
		GetUserProfile,
	},
	Route{
		"Send Push Notification",
		"POST",
		"/sendpush",
		SendPush,
	},
	Route{
		"Matching",
		"POST",
		"/match/{uid}",
		Match,
	},
	Route{
		"Debugging Endpoint",
		"POST",
		"/test/{testType}",
		Test,
	},
}
