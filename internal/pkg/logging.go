package pkg

import (
	"log"
	"os"
	"time"
)

var (
	LoggerInfo  *log.Logger
	LoggerWarn  *log.Logger
	LoggerError *log.Logger
	logFile     *os.File
	Info        = "\033[34m" // Blue
	Error       = "\033[31m" // Red
	Warn        = "\033[33m" // Yellow
	Reset       = "\033[0m"
)

const maxFileSize = 24 * 1024 // 24 Ko

func InitLogging() {
	var err error
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Println("Failed to create logs directory:", err)
			return
		}
	}

	logFile, err = os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file:", err)
		return
	}

	LoggerInfo = log.New(logFile, Info+"[INFO]: "+Reset, log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile)
	LoggerError = log.New(logFile, Error+"[ERROR]: "+Reset, log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile)
	LoggerWarn = log.New(logFile, Warn+"[WARN]: "+Reset, log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile)
}

func RotateLogFile() {
	if stat, err := logFile.Stat(); err != nil {
		if stat.Size() >= maxFileSize {
			err := logFile.Close()
			if err != nil {
				log.Println("Failed to close log file:", err)
				return
			}

			newName := "logs/app-" + time.Now().Format("20060102150405") + ".log"
			err = os.Rename("logs/app.log", newName)
			if err != nil {
				log.Println("Failed to rename log file:", err)
				LoggerError.Println("Failed to rename log file:", err)
				return
			}
			LoggerInfo.Println("Rotate log file successfully")
			InitLogging()
		}
	}
}
