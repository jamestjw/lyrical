package help

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
)

// HelpJSON is a struct containing the help text to be displayed to users
var HelpJSON helpJSON

type helpJSON struct {
	Commands []command `json:"commands"`
	HelpText string    `json:"helpText"`
}

type command struct {
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Arguments   []string `json:"arguments"`
}

func (c command) prettyArguments() string {
	args := make([]string, 0)
	for _, arg := range c.Arguments {
		args = append(args, fmt.Sprintf("<%s>", arg))
	}

	if len(args) > 0 {
		return strings.Join(args, ",")
	}
	return "nil"
}

func (c command) prettyName() string {
	return fmt.Sprintf("!%s", c.Name)
}

// InitialiseHelpText should be called by application on startup
func InitialiseHelpText() {
	jsonFile, err := os.Open("help/help.json")
	jsonBytes, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(jsonBytes, &HelpJSON)

	log.Println("Loaded help text.")
}

// Message is the message sent to channels with list
// of possible commands that users can execute
func Message() string {
	var result bytes.Buffer
	table := tablewriter.NewWriter(&result)
	table.SetHeader([]string{"Description", "Name", "Parameters"})
	table.SetCaption(true, HelpJSON.HelpText)

	for _, command := range HelpJSON.Commands {
		v := []string{command.Description, command.prettyName(), command.prettyArguments()}
		table.Append(v)
	}

	table.Render()
	tableString := string(result.Bytes())
	return codeBlockify(tableString)
}

func codeBlockify(msg string) string {
	return fmt.Sprintf("```\n%s\n```", msg)
}
