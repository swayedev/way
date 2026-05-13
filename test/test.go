package main

import "github.com/swayedev/way"

// "github.com/swayedev/way/database/pgx"

func main() {
	w := way.New()
	// w.GET("/", helloHandler)
	// uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.GetDbUser(), config.GetDbPassword(), config.GetDbHost(), config.GetDbPort(), config.GetDbName())
	// conn, err := pgx.Connect(context.Background(), uri)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }

	w.Start(":8081")
}

// func helloHandler(c *way.Context) {
// 	c.Response.Header().Set("Content-Type", "application/json")
// 	c.Response.Write([]byte("Hello World"))
// }
