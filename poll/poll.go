package poll

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jamestjw/lyrical/utils"
)

var pollParamsRegex = regexp.MustCompile(`(?m)((?:("[^"]+"|[^\s"]+))?)(?:\s+|$)`)

type Poll struct {
	title             string
	options           []string
	emojiToOption     map[string]*Option
	durationInSeconds int
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

	if len(parsedParams) <= 3 {
		err = errors.New("aside from a `title` and `poll-duration`, please provide at least two other options for a vote to make sense")
		return
	} else if len(parsedParams) > 12 {
		err = errors.New("we currently only support a maximum number of **10** options in a poll")
		return
	}

	duration, err := strconv.Atoi(parsedParams[1])
	if err != nil {
		err = errors.New("`poll-duration` needs to be an integer")
		return
	}

	p = &Poll{
		title:             parsedParams[0],
		options:           parsedParams[2:],
		emojiToOption:     make(map[string]*Option),
		durationInSeconds: duration,
	}

	return
}

func (p *Poll) GeneratePollMessage() string {
	messages := []string{
		"A poll has been started!",
		utils.Bold(p.title),
	}

	for index, option := range p.options {
		emoji := numberToEmoji[index+1]
		p.emojiToOption[emoji] = &Option{name: option}

		formattedOption := fmt.Sprintf("%s. %s", emoji, option)
		messages = append(messages, formattedOption)
	}

	finalMessage := fmt.Sprintf("Exercise your right to vote by reacting accordingly! The poll will close in %s.", p.GetDuration())

	messages = append(messages, finalMessage)

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

	var verdictMessage string

	if options[0].count == 0 {
		verdictMessage = "Unfortunately no votes were received, hence a decision was unable to be made."
	} else if options[0].count == options[1].count {
		verdictMessage = fmt.Sprintf("Looks like we have a tie between %s", getTiedOptions(options))
	} else {
		verdictMessage = fmt.Sprintf("The people have spoken, **%s** it shall be.", options[0].name)
	}

	results = append(results, verdictMessage)

	return strings.Join(results, "\n")
}

// GetDuration returns a time.Duration corresponding to how long
// the poll should last.
func (p *Poll) GetDuration() time.Duration {
	return time.Duration(p.durationInSeconds) * time.Second
}

// GetTiedOptions returns a string with the options that are tied on a particular
// score. tiedScore should be a sorted array of Options in descending order
// of count.
func getTiedOptions(options []Option) string {
	tiedScore := options[0].count
	var tiedNames []string

	for _, o := range options {
		if o.count == tiedScore {
			tiedNames = append(tiedNames, utils.Bold(o.name))
		} else {
			break
		}
	}

	return strings.Join(tiedNames, ", ")
}
