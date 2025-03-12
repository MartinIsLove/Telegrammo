package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/user", rt.doLogin)
	rt.router.PUT("/user/name", rt.setMyUserName)
	rt.router.PUT("/user/photo", rt.setMyPhoto)

	rt.router.GET("/users/:id", rt.getMyUser)

	rt.router.GET("/user/:username", rt.checknames)

	rt.router.POST("/conversation", rt.createChat)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
