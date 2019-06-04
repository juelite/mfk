package logic

import (
	"log"
	"os"
	"runtime"
	"strings"

	"mfk/common"
)

func init() {
	if runtime.GOOS == "windows" {
		DirSep = "\\"
	} else {
		DirSep = "/"
	}
}

func Run() {
	err, _, errMsg := common.ExecCommand("go build")
	if err != nil {
		log.Fatal(errMsg)
	}
	path, _ := os.Getwd()
	path_sep := strings.Split(path, DirSep)
	app := path_sep[len(path_sep)-1]
	common.ExecLiveCommand("./" + app)
}
