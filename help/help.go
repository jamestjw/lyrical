package help

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// HelpJSON is a struct containing the help text to be displayed to users
var HelpJSON helpText

type helpText struct {
	Commands []command `json:"commands"`
}

type command struct {
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Arguments   []string `json:"arguments"`
}

func init() {
	jsonFile, err := os.Open("help.json")
	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(jsonBytes, &HelpJSON)

	fmt.Println(HelpJSON.Commands[0].Description)
}
