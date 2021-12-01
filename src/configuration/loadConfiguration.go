package configuration

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"structures"
)

func GetConfiguration() structures.Configuration {

	var c structures.Configuration

	logrus.Infoln("Acquiring configuration from configuration.json file...")

	file, err := os.Open("configuration.json")
	if err != nil {

		logrus.Fatal("Unable to load configuration.json file, now exiting !")
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal(err)
		logrus.Fatal("Unable to load configuration.json file because of invalid permissions, now exiting !")
	}

	err = json.Unmarshal(fileContent, &c)
	if err != nil {
		logrus.Fatal("Malformed configuration.json file. Please check documentation. Program is now exiting.")
	}
	_ = file.Close()

	return c
}
