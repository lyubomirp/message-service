package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"message-service/health"
	"message-service/helpers"
	"message-service/models"
	"message-service/rabbit"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func setupDbAndTables() *gorm.DB {
	sslMode := "enable"

	if os.Getenv("ENV") != "prod" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"), os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"), os.Getenv(sslMode),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	helpers.CheckForError(err, "Database failed to initialize")

	if (!db.Migrator().HasTable(&models.Log{})) {
		err := db.Migrator().CreateTable(&models.Log{})
		helpers.CheckForError(err, "Table creation failed")
	}

	return db
}

func main() {
	// Set up logs
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	// Load .env variables
	err := godotenv.Load(".env")

	if err != nil {
		panic("ENV variables failed to load")
	}

	go health.InitHealthChecks()

	db := setupDbAndTables()
	openRmqConnection(db)
}

func openRmqConnection(db *gorm.DB) {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s/",
			os.Getenv("RMQ_USER"), os.Getenv("RMQ_PASS"),
			os.Getenv("RMQ_URL"), os.Getenv("RMQ_PORT"),
		),
	)

	if err != nil {
		log.WithFields(log.Fields{
			"message": "RMQ connection failed",
		}).Error(err.Error())
	}

	chErrorChannel := make(chan *amqp.Error)
	channel, err := conn.Channel()

	// reconnect on channel not opening
	if err != nil {
		log.Println("Reconnecting After 10 seconds: " + err.Error())
		time.Sleep(10 * time.Second)
		openRmqConnection(db)
	}

	channel.NotifyClose(chErrorChannel)

	// goroutine for capturing dropped connections
	// this is watching our channel
	// since graceful closes don't throw errors
	// in RMQ connection
	go func() {
		err := <-chErrorChannel
		log.Println("Reconnecting After 10 seconds: " + err.Error())
		time.Sleep(10 * time.Second)
		openRmqConnection(db)
	}()

	receiver := rabbit.NewReceiver(channel)
	receiver.Consume(db)

	func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)
}
