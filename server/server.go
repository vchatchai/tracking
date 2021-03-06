package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"tracking/db"
	"tracking/web"

	env "github.com/Netflix/go-env"
)

// var (
// 	debug         = flag.Bool("debug", true, "enable debugging")
// 	password      = flag.String("password", "Manager1", "the database password")
// 	port     *int = flag.Int("port", 1433, "the database port")
// 	server        = flag.String("server", "192.168.1.115", "the database server")
// 	user          = flag.String("user", "tracking", "the database user")
// 	database      = flag.String("database", "tracking", "the database name")
// )

func main() {

	es, err := env.UnmarshalFromEnviron(&web.Environment)
	if err != nil {
		log.Fatal(err)
	}
	// Remaining web.environment variables.
	web.Environment.Extras = es

	d, err := sql.Open("sqlserver", dataSource())
	if err != nil {
		log.Fatal(err)
	}
	d.SetMaxOpenConns(0)
	d.SetMaxIdleConns(10)

	defer d.Close()
	// CORS is enabled only in prod profile
	cors := os.Getenv("profile") == "prod"
	app := web.NewApp(db.NewDB(d), cors)
	err = app.Serve()
	log.Println("Error", err)
}

func dataSource() string {
	// host := "localhost"
	// pass := "manager1"
	// if os.Getenv("profile") == "prod" {
	// 	host = "db"
	// 	pass = os.Getenv("db_pass")
	// }

	if web.Environment.Debug {
		// fmt.Printf(" password:%s\n", web.Environment.Password)
		fmt.Printf(" port:%d\n", web.Environment.DatabasePort)
		fmt.Printf(" server:%s\n", web.Environment.Server)
		fmt.Printf(" user:%s\n", web.Environment.User)
		fmt.Printf(" database:%s\n", web.Environment.Database)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", web.Environment.Server, web.Environment.User, web.Environment.Password, web.Environment.DatabasePort, web.Environment.Database)
	if web.Environment.Debug {
		fmt.Printf(" connString:%s\n", connString)
	}

	return connString

	// return "goxygen:" + pass + "@tcp(" + host + ":3306)/goxygen"
}
