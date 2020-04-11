package poll

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var pollParamsRegex = regexp.MustCompile(`(?m)((?:("[^"]+"|[^\s"]+))?)(?:\s+|$)`)

type Poll struct {
	title         string
	options       []string
	emojiToOption map[string]*Option
}

var numberToEmoji = map[int]string{
	1:  ":one:",
	2:  ":two:",
	3:  ":three:",
	4:  ":four:",
	5:  ":five",
	6:  ":six",
	7:  ":seven:",
	8:  ":eight:",
	9:  ":nine:",
	10: ":ten",
}

func FromStringParams(params string) (p *Poll, err error) {
	var parsedParams []string

	for _, match := range pollParamsRegex.FindAllString(params, -1) {
		res := strings.Trim(strings.TrimSpace(match), "\"")
		parsedParams = append(parsedParams, res)
	}

	if len(parsedParams) <= 2 {
		err = errors.New("aside from a `title`, please provide at least two other options for a vote to make sense")
		return
	}

	p = &Poll{
		title:         parsedParams[0],
		options:       parsedParams[1:],
		emojiToOption: make(map[string]*Option),
	}

	return
}

func (p *Poll) GeneratePollMessage() string {
	messages := []string{p.title}

	for index, option := range p.options {
		emoji := numberToEmoji[index+1]
		p.emojiToOption[emoji] = &Option{name: option}

		formattedOption := fmt.Sprintf("%s. %s", emoji, option)
		messages = append(messages, formattedOption)
	}

	return strings.Join(messages, "\n")
}

func (p *Poll) AddResult(reactionCounts map[string]int) {
	for emoji, count := range reactionCounts {
		p.emojiToOption[emoji].SetCount(count)
	}
}

func (p *Poll) GetVerdict() string {
	options := []Option{}

	for _, option := range p.emojiToOption {
		options = append(options, *option)
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].count < options[j].count
	})

	results := []string{"Results:"}

	for _, option := range options {
		formattedResult := fmt.Sprintf("%s: %v", option.name, option.count)
		results = append(results, formattedResult)
	}

	return strings.Join(results, "\n")
}
