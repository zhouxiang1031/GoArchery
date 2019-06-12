package config

type Server struct {
	Addr          string
	NetType       int //0:tcp; 1:wc; 2:udp
	UseTls        bool
	Certfile      string
	Keyfile       string
	Release       bool
	SvcCount      int  //how many go routines for service-worker
	LogFile       bool //out log to file?
	LogLev        int
	Heartbeat     Heartbeat
	NtoGBufMax    int
	BuffLimit     BuffLimit
	IsCompression bool
	Room          RoomConf
}

type BuffLimit struct {
	NtogMax int
	GtonMax int
}

type Heartbeat struct {
	Interval int64 //interval of hearbeat,default 20s
	Ttl      int64 //time to live,default 60s
}

type GameConfig struct {
	Server Server
}

type RoomConf struct {
	Update int32
	Buff   int32
}
