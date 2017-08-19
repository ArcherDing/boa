package main

import "github.com/ArcherDing/boa"

func hello(ctx *boa.Context)  {
	ctx.String(200, "Hello boa")
}

func main() {
	println(boa.Version)
	b := boa.New()
	b.GET("/hello", hello)
	b.Run(":3000")
}
