package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Ignore          []string
	ExternalUrlBase string
}

func Load(filename string) *Config {
	var result Config
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &result)
		if err != nil {
			usage(filename, err)
		}
	}
	if result.Ignore == nil {
		result.Ignore = []string{}
	}
	return &result
}

func usage(filename string, err error) {
	fmt.Println("Error reading config file ", filename)
	fmt.Println(err)
	fmt.Println("The config file must be well formed json.")
	fmt.Println("Known keys and types are:")
	fmt.Println("\tignore : []string.  .gitignore format lines to ignore files")
	fmt.Println("\texternalUrlBase : string.  Add links to external content - eg the github location for the project")
	log.Fatal()
}
