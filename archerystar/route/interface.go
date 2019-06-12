package route

import (
	"archerystar/message"
)

type RouteEntity interface {
	RouteNext(msg *message.NtolMessage)
}
