package config

import (
	"archerystar/component/jsonloader"
	"fmt"
)

var gameCfg *GameConfig

func init() {
	gameCfg = &GameConfig{}

	jsonloader.Load("./config/conf/server.json", &gameCfg.Server)

	fmt.Printf("gameconfig:%+v", *gameCfg)
}

func Gameconfig() *GameConfig {
	return gameCfg
}
