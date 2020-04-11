package poll

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var pollParamsRegex = regexp.MustCompile(`(?m)((?:("[^"]+"|[^\s"]+))?)(?:\s+|$)`)

type Poll struct {
	title   string
	options []string
}

func FromStringParams(params string) (p *Poll, err error) {
	var parsedParams []string

	for _, match := range pollParamsRegex.FindAllString(params, -1) {
		res := strings.Trim(strings.TrimSpace(match), "\"")
		parsedParams = append(parsedParams, res)
	}

	if len(parsedParams) <= 2 {
		err = errors.New("Aside from a `title`, please provide at least two other options for a vote to make sense!")
		return
	}

	p = &Poll{
		title:   parsedParams[0],
		options: parsedParams[1:],
	}

	return
}

func (p *Poll) GeneratePollMessage() string {
	messages := []string{p.title}

	for index, option := range p.options {
		formattedOption := fmt.Sprintf("%v. %s", index+1, option)
		messages = append(messages, formattedOption)
	}

	return strings.Join(messages, "\n")
}
