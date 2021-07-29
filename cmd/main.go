package main

import (
	"time"

	HTTPServer "shopee.com/zeliang-entry-task/HTTPServer/main"

	"shopee.com/zeliang-entry-task/TCPServer"
)

func main() {
	go TCPServer.TCPServerMain()
	time.Sleep(time.Second * 5)
	HTTPServer.HttpServerMain()
}
