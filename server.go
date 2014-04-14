package main

import "github.com/go-martini/martini"

func main() {
	go webserver()
}

func webserver() {
	m := martini.Classic()

	m.Use(martini.Static("public"))

	m.Run()
}
