package main

import "github.com/swayedev/way"

func main() {
	w := way.New()
	w.GET("/", helloHandler)

	w.Start(":8080")
}

func helloHandler(c *way.Context) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.Write([]byte("Hello World"))
}
