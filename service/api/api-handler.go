package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// --- rotte generiche e di ricerca ---
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/user", rt.doLogin)
	rt.router.PUT("/user/name", rt.setMyUserName)
	rt.router.PUT("/user/photo", rt.setMyPhoto)
	rt.router.GET("/users/:id", rt.getMyUser)
	rt.router.GET("/search/users/:username", rt.checknames)

	// --- conversazioni e gruppi ---
	rt.router.GET("/conversations/forward", rt.getForwardChat)
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.DELETE("/conversation/leave/:idChat", rt.leaveGroup)
	rt.router.POST("/conversation", rt.createChat)
	rt.router.POST("/conversation/group", rt.createGroup)
	rt.router.PUT("/conversation/group/name", rt.setGroupName)
	rt.router.PUT("/conversation/group/photo", rt.setGroupPhoto)
	rt.router.PUT("/conversation/group/user", rt.addToGroup)

	rt.router.GET("/conversation/:idChat/group/users", rt.getGroupUsers)
	rt.router.GET("/conversation/:idChat", rt.getConversation)

	// --- messaggi ---
	rt.router.POST("/message", rt.sendMessage)
	rt.router.DELETE("/message/delete/:idMes", rt.deleteMessage)
	rt.router.POST("/message/forward", rt.forwardMessage)
	rt.router.POST("/message/comment", rt.commentMessage)
	rt.router.POST("/message/uncomment", rt.uncommentMessage)

	// --- special routes ---
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
