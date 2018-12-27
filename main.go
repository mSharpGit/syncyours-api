package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var config tomlConfig
var mail mailer

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//log.Println()
	//**********toml configuration initialization************//
	//log.Println(randSeq(10))
	if _, err := toml.DecodeFile("config/config.toml", &config); err != nil {
		log.Println(err)
		return
	}
	//log.Println(os.O_WRONLY)
	//**********toml configuration initialization************//

	//**********mail initialization************//

	mail.Initialize(config.Smtp.Port, config.Smtp.Host, config.Smtp.Username, config.Smtp.Password, config.Smtp.From, config.Smtp.CC)

	//**********************//

	//**********************//
	// make sure log.txt exists first
	// use touch command to create if log.txt does not exist
	//perm, _ := fmt.Printf("%o", config.Logs.Permission)

	perm := os.FileMode(config.Logs.Permission)
	logFile, err := os.OpenFile(config.Logs.Location, config.Logs.Flag, perm)
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
	a.Initialize(config.DB.User, config.DB.Password, config.DB.Server, config.DB.Port, config.DB.Dbname)
	a.Run(":" + config.Api.Listener)
}
