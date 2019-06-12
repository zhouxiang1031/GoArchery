package route

type RouteTable map[uint]Route

type Route []RouteStep

type RouteStep struct {
	Svc int //to which svc this step
	Sn  int //to which svc	next step when success
	Fn  int //to which svc	next step when fail
}

type RouteFile struct {
	Msgid uint
	Route Route
}

type RouteFileTable map[string]RouteFile

type Router struct {
	routeTable RouteTable
}
