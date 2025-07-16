package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bwebb-hx/figma-to-html/cli/figmatic/internal/utils"
)

func createLogString(s string) string {
	return fmt.Sprintf("---\n%s\n\n%s\n", time.Now(), s)
}

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

	_, err = file.WriteString(createLogString(content))
	if err != nil {
		log.Println("failed to write logs: " + err.Error())
	}
}

func Error(errorMsg any) {
	s := fmt.Sprintf("ERROR: %s", errorMsg)
	utils.Colors.ErrorPrint(s)
	WriteLog("error", s)
}

func Fatal(errorMsg any) {
	Error(errorMsg)
	os.Exit(1)
}
