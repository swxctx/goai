// Code generated by 'malatd gen' command.
// DO NOT EDIT!

package api

import (
	td "github.com/swxctx/malatd"
)

func Route(srv *td.Server, rootGroup string) {
	//自定义路由处理
	routeLogic(srv, rootGroup)

	// APIs...
	{

		v1_chat := srv.Group(rootGroup + "/v1/chat")
		v1_chat.Post("/do", DoHandle)
		v1_chat.Get("/do", DoHandle)
	}

}
