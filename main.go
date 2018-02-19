package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/gh0o0st/uploader/COS"
)

var HOME = os.Getenv("HOME")
var configFile = HOME + "/.uploader.json"

func main() {
	config, err := loadConfig(configFile)
	if err != nil {
		panic(err)
	}
	if config.APIURL == "" || config.PrefixURL == "" || config.SecretKey == "" || config.SecretID == "" {
		panic("some field in config file is missing")
	}
	client, err := COS.NewClient(config.SecretID, config.SecretKey, config.APIURL)
	if err != nil {
		panic(err)
	}
	if len(os.Args) < 2 {
		panic("need an argument")
	}
	fileToUpload := os.Args[1]
	dir := withDate()
	err = client.PutObject(fileToUpload, dir)
	if err != nil {
		panic(err)
	}
	info, _ := os.Stat(fileToUpload)
	fmt.Println(config.PrefixURL + dir + info.Name())
}

func withDate() string {
	year, month, _ := time.Now().Date()
	return strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/"
}

type Config struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	APIURL    string `json:"api_url"`
	PrefixURL string `json:"prefix_url"`
}

func loadConfig(file string) (*Config, error) {
	config := &Config{}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		//if file not exist, create an empty file
		if os.IsNotExist(err) {

			f, err := os.Create(file)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			b, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return nil, err
			}
			f.Write(b)
		}
		return nil, err
	}

	if err := json.Unmarshal(b, config); err != nil {
		return nil, err
	}
	return config, nil
}
