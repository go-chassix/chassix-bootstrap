package main

import "c5x.io/chassix"

func main() {
	chassix.Init()
	chassix.StartHttpServer(nil)
}
