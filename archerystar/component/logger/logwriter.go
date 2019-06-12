package logger

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const MAX_LOGFILE_LIMIT int64 = 1024 * 1024 * 20

type logWiriter struct {
	file *os.File
	size int64
	day  int
}

func newWirter() *logWiriter {
	now := time.Now()
	file, err := os.OpenFile("./mylog"+strconv.FormatInt(now.Unix(), 10), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)

	if err != nil {
		fmt.Println("log writer init failed")
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	return &logWiriter{file, info.Size(), now.Day()}
}

func (w *logWiriter) Write(data []byte) (int, error) {
	if w == nil {
		return 0, errors.New("logWrite is nil!")
	}

	if w.file == nil {
		return 0, errors.New("log file has not be openned!")
	}

	w.CheckCreateNewFile()

	n, err := w.file.Write(data)

	return n, err
}

func (w *logWiriter) CheckCreateNewFile() {
	now := time.Now()

	if now.Day() != w.day || w.size >= MAX_LOGFILE_LIMIT {
		w.file.Close()
		fmt.Println("log file is full")
		w.file, _ = os.OpenFile("./mylog"+strconv.FormatInt(time.Now().Unix(), 10), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		w.size = 0
		w.day = now.Day()
	}

}

func (w *logWiriter) initWriter() *logWiriter {
	now := time.Now()
	file, err := os.OpenFile("./mylog"+strconv.FormatInt(now.Unix(), 10), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)

	if err != nil {
		fmt.Println("log writer init failed")
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	return &logWiriter{file, info.Size(), now.Day()}
}
