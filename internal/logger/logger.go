package logger

import (
  "log"
  "os"
)

type AppLogger struct {
  InfoLogger *log.Logger
  ErrorLogger *log.Logger
  DebugLogger *log.Logger
}

/*
//Will need to run in function that constructs the Logger (Main) to create the *File and defer closing it until main ends
if len(l.logFilePath) == 0 {
    return errors.New("Error: Empty Log File Path provided for App Logger. Logger cannot be created.")
  }
  //Creates log file if it doesn't exist and append if it does.
  file, err := os.OpenFile(l.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    return err
  }

//Construct Logger Here
  
defer file.Close()
*/
func (l *AppLogger) Logger(file *os.File) {
  l.InfoLogger = log.New(file, "Info: " , log.Ldate|log.Ltime)
  l.ErrorLogger = log.New(file, "Error: " , log.Ldate|log.Ltime)
  l.DebugLogger = log.New(file, "Debug: " , log.Ldate|log.Ltime)
}


//Function to Log different level/prefixes
func (l *AppLogger) LogInfo(info string) {
  l.InfoLogger.Println(info)
}

func (l *AppLogger) LogError(info string) {
  l.ErrorLogger.Println(info)
}

func (l *AppLogger) LogDebug(info string) {
  l.DebugLogger.Println(info)
}
