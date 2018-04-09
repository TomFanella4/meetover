package router

import (
	"net/http"

	"meetover/backend/router/handlers"
)

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
		handlers.Index,
	},
	Route{
		"Login",
		"POST",
		"/login/{code}",
		handlers.VerifyUser,
	},
	Route{
		"Refresh",
		"POST",
		"/refreshtoken",
		handlers.RefreshCustomToken,
	},
	Route{
		"Initiate Meetover",
		"POST",
		"/meetover/{otherID}",
		handlers.InitiateMeetover,
	},
	Route{
		"Processes the user's decision to MeetOver",
		"POST",
		"/meetover/decision/{otherID}",
		handlers.ProcessDecision,
	},
	Route{
		"Get other user's profile(s)",
		"POST",
		"/userprofiles",
		handlers.GetUserProfiles,
	},
	Route{
		"Send Push Notification",
		"POST",
		"/sendpush",
		handlers.SendPush,
	},
	Route{
		"Matching",
		"POST",
		"/match/{uid}",
		handlers.Match,
	},
	Route{
		"Debugging Endpoint",
		"POST",
		"/test/{testType}",
		handlers.Test,
	},
}
