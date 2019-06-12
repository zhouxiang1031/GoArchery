package route

import (
	"fmt"

	gmerrs "archerystar/common/errors"
	"archerystar/component/jsonloader"
	"archerystar/component/logger"
	"archerystar/message"
)

const routeTitle = "route"

var router *Router

func init() {
	router = &Router{
		routeTable: RouteTable{},
	}

	rFiles := RouteFileTable{}

	jsonloader.Load("./config/conf/route.json", &rFiles)

	for _, ritem := range rFiles {
		router.routeTable[ritem.Msgid] = ritem.Route
	}

	fmt.Printf("route table is %+v\n", router.routeTable)

}

func NewRoute(msgid uint) *message.MsgRoute {
	msgRoute := &message.MsgRoute{
		MsgId:  msgid,
		Step:   0,
		Result: true,
	}

	return msgRoute
}

func FindNextRoute(msgRoute *message.MsgRoute) error {
	rt, ok := router.routeTable[msgRoute.MsgId]
	if !ok {
		return gmerrs.NewErr("Not found route for msg:", msgRoute.MsgId)
	}

	if msgRoute.Result {
		msgRoute.To = rt[msgRoute.Step].Sn
	} else {
		msgRoute.To = rt[msgRoute.Step].Fn
	}

	msgRoute.Step++

	return nil
}

func Next(entiry RouteEntity, msg *message.NtolMessage) {
	if msg.MsgRoute == nil {
		msg.MsgRoute = NewRoute(msg.Msg.ID)
		if err := FindNextRoute(msg.MsgRoute); err != nil {
			logger.Error(routeTitle, err.Error())
		}
	} else {
		if err := FindNextRoute(msg.MsgRoute); err != nil {
			logger.Error(routeTitle, err.Error())
		}
	}

	entiry.RouteNext(msg)
}
