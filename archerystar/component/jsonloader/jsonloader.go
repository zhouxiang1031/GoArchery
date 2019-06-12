package jsonloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Load(filename string, data interface{}) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Load file: %s failed:%s\n", filename, err.Error())
		return
	}

	err = json.Unmarshal(buf, data)
	if err != nil {
		return
	}

}
