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
	1:  "1️⃣",
	2:  "2️⃣",
	3:  "3️⃣",
	4:  "4️⃣",
	5:  "5️⃣",
	6:  "6️⃣",
	7:  "7️⃣",
	8:  "8️⃣",
	9:  "9️⃣",
	10: "🔟",
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
		option, exists := p.emojiToOption[emoji]
		if exists {
			option.SetCount(count)
		}
	}
}

func (p *Poll) GetVerdict() string {
	options := []Option{}

	for _, option := range p.emojiToOption {
		options = append(options, *option)
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].count > options[j].count
	})

	results := []string{"**Results:**"}

	for _, option := range options {
		formattedResult := fmt.Sprintf("%s: %v", option.name, option.count)
		results = append(results, formattedResult)
	}

	var verdictMessage string

	if options[0].count == 0 {
		verdictMessage = "Unfortunately no votes were received and a decision was unable to be made."
	} else if options[0].count == options[1].count {
		verdictMessage = fmt.Sprintf("Looks like we have a tie between **%s** and **%s**", options[0].name, options[1].name)
	} else {
		verdictMessage = fmt.Sprintf("The people have spoken, **%s** it shall be.", options[0].name)
	}

	results = append(results, verdictMessage)

	return strings.Join(results, "\n")
}
