package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func WriteLog(fileName string, content string) {
	err := os.Mkdir("logs", 0755)
	if err != nil && !os.IsExist(err) {
		log.Println("failed to write logs: " + err.Error())
		return
	}

	logFilePath := filepath.Join("logs", fileName+".log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("failed to write logs: " + err.Error())
		return
	}
	defer file.Close()

	s := fmt.Sprintf("---\n%s\n\n%s", time.Now(), content)

	_, err = file.WriteString(s)
	if err != nil {
		log.Println("failed to write logs: " + err.Error())
	}
}
