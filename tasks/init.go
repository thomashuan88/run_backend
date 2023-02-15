package tasks

import (
	"fmt"
	"run-backend/conf"
)

func init() {
	fmt.Println("this is tasks init")
	conf.Init()
}
