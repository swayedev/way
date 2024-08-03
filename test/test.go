package main

import "github.com/swayedev/way"

func main() {
	w := way.New(nil)
	w.GET("/", helloHandler)

	w.Start(":8081")
}

func helloHandler(c *way.Context) {
	// c.Response.Header().Set("Content-Type", "application/json")
	// c.Response.Write([]byte("Hello World"))
	c.JSON(200, "Hello World")
}
