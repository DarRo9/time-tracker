package database 


import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDatabase() *pgxpool.Pool {
	host 	 := os.Getenv("POSTGRES_HOST")
	port 	 := os.Getenv("POSTGRES_PORT")
	user 	 := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname   := os.Getenv("POSTGRES_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil{
		log.Fatalf("oshibka podkl k bd: %v\n", err)
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := dbpool.Ping(context); err != nil{
		log.Fatalf("Ping ne proshel, bd lezhit yopta: %v\n", err)
	}

	return dbpool
}