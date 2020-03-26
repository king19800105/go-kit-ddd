package main

import (
	"github.com/king19800105/go-kit-ddd/cmd/usr-base/app"
	"log"
)

func main() {
	if err := app.Run(); nil != err {
		log.Fatalf("base operating 服务被终止，原因：%v", err)
	}
}
