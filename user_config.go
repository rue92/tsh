package main

import "log"
import "fmt"
import "os"
import "io/ioutil"
import "path/filepath"
import "encoding/xml"

type UserConfig struct {
	XMLName  xml.Name `xml:"Config"`
	Username string   `xml:"Username"`
	Oauth    string   `xml:"Oauth"`
}

var User UserConfig

func ReadConfig() {
	userHome := os.Getenv("HOME")
	fileAsBytes, err := ioutil.ReadFile(filepath.Join(userHome, ".config/tsh.xml"))
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}
	err = xml.Unmarshal([]byte(fileAsBytes), &User)
	if err != nil {
		log.Printf("Error reading config: %v", err)
		return
	}
	log.Printf("Read Config %s\n", User.String())
}

func (config *UserConfig) String() string {
	return fmt.Sprintf("Username: %s, "+
		"Oauth: %s",
		config.Username, config.Oauth)
}
