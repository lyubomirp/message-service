package health

import (
	"fmt"
	"github.com/hellofresh/health-go/v5"
	log "github.com/sirupsen/logrus"
	"message-service/helpers"
	"net/http"
	"os"
	"time"

	healthHttp "github.com/hellofresh/health-go/v5/checks/http"
	healthPg "github.com/hellofresh/health-go/v5/checks/postgres"
)

func InitHealthChecks() {
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    os.Getenv("APP_NAME"),
		Version: os.Getenv("APP_VER"),
	}))

	rmqConnString := fmt.Sprintf(
		`http://%s:%s@%s:%s/api/aliveness-test/%s`,
		os.Getenv("RMQ_USER"), os.Getenv("RMQ_PASS"),
		os.Getenv("RMQ_URL"), os.Getenv("RMQ_WEB_PORT"),
		"%2f",
	)

	err := h.Register(health.Config{
		Name:      "rabbitmq",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthHttp.New(healthHttp.Config{
			URL: rmqConnString,
		}),
	})
	helpers.CheckForError(err, "RMQ healthcheck registration failed")

	sslMode := "enable"

	if os.Getenv("ENV") != "prod" {
		sslMode = "disable"
	}

	pgConnString := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=%s`,
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"), sslMode,
	)

	err = h.Register(health.Config{
		Name: "postgres",
		Check: healthPg.New(healthPg.Config{
			DSN: pgConnString,
		}),
	})
	helpers.CheckForError(err, "PSQL healthcheck registration failed")

	httpf := http.NewServeMux()
	httpf.Handle("/status", h.Handler())

	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_URL"), os.Getenv("APP_PORT"))

	srv := &http.Server{
		Addr:    addr,
		Handler: httpf,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(fmt.Sprintf(
			"Failed to initialize health check endpoint, Error: %s",
			err.Error(),
		))
	}

	log.Info("Health-checks started")

}
