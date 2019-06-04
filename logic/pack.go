package logic

import (
	"fmt"
	"log"
	"os"
	"strings"

	"mfk/common"
)

func Pack() {
	err, _, errMsg := common.ExecCommand("go build")
	if err != nil {
		log.Fatal(errMsg)
	}
	path, _ := os.Getwd()
	path_sep := strings.Split(path, DirSep)
	app := path_sep[len(path_sep)-1]

	err, res, errMsg := common.ExecCommand("tar -czf " + app + ".tar.gz ./" + app)
	if err != nil {
		log.Fatal(errMsg)
	}
	fmt.Println(res)
}
