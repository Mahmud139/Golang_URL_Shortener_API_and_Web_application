package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	db   struct {
		db_address  string
		db_password string
		db_no       int
		ctx         context.Context
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	rdb      *redis.Client
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.db_address, "db-addr", os.Getenv(""), "Redis DB address")
	flag.StringVar(&cfg.db.db_password, "db-passwd", os.Getenv(""), "Redis DB password")
	flag.IntVar(&cfg.db.db_no, "db-no", 0, "Redis DB number")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	cfg.db.ctx = context.Background()

	if *displayVersion {
		fmt.Printf("Version: \t%s\n", version)
		fmt.Printf("Build time: \t%s\n", buildTime)
		os.Exit(0)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	rdb, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("database connection is established!")

	app := application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		rdb:      rdb,
	}

	fiberApp := fiber.New()

	app.routes(fiberApp)

	infoLog.Printf("starting %s server on port :%d\n", cfg.env, cfg.port)

	err = fiberApp.Listen(fmt.Sprintf(":%d", cfg.port))
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(cfg config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.db.db_address,
		Password: cfg.db.db_password,
		DB:       cfg.db.db_no,
	})

	res, err := rdb.Ping(cfg.db.ctx).Result()
	if err != nil {
		return nil, err
	}

	if res != "PONG" {
		return nil, fmt.Errorf("couldn't connect to Redis DB")
	}

	return rdb, nil
}
