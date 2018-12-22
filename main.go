package main

import (
	"io"
	"log"
	"os"
)

func main() {

	//**********************//
	// make sure log.txt exists first
	// use touch command to create if log.txt does not exist
	logFile, err := os.OpenFile("log.txt", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	// direct all log messages to log.txt
	//log.SetOutput(logFile)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	//**********************//

	a := App{}
	// You need to set your Username and Password here
	//a.Initialize("syncyours", "1Q@w3e4r5t6yaz", "166.62.27.55", "3306", "syncyours")
	//database new host
	a.Initialize("leenhal1_sync", "1Q@w3e4r5t6yaz", "192.254.181.110", "3306", "leenhal1_syncyours")
	a.Run(":8080")
}
