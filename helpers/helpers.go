package helpers

import (
	"gorm.io/gorm"
	"message-service/models"

	log "github.com/sirupsen/logrus"
)

func CheckForError(err error, msg string) {
	if err != nil {
		log.WithFields(log.Fields{
			"message": msg,
		}).Error(err.Error())
	}
}

func LogMessageError(err error, db *gorm.DB, body []byte) {
	log.WithFields(log.Fields{
		"message": "Message sending failed irreparably",
	}).Error(err)

	// Log the body of the failed message
	failedMessage := models.Log{
		Body: string(body),
	}

	result := db.Create(&failedMessage)

	// If DB creation fails, we should panic the instance away
	if result.Error != nil {
		log.WithFields(log.Fields{
			"message": "DB saving failed",
		}).Panic(result.Error)
	}
}
