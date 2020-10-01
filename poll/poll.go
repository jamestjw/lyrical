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

// FromStringParams is a factory method to generate a poll
// based on a string of params
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

func (p *Poll) GeneratePollMessage() (string, []string) {
	messages := []string{
		"A poll has been started!",
		utils.Bold(p.title),
	}

	emojis := []string{}

	for index, option := range p.options {
		emoji := numberToEmoji[index+1]
		emojis = append(emojis, emoji)
		p.emojiToOption[emoji] = &Option{name: option}

		formattedOption := fmt.Sprintf("%s. %s", emoji, option)
		messages = append(messages, formattedOption)
	}

	finalMessage := fmt.Sprintf("Exercise your right to vote by reacting accordingly! The poll will close in %s.", p.GetDuration())

	messages = append(messages, finalMessage)

	return strings.Join(messages, "\n"), emojis
}

// AddResult accepts a map of emojis to array of user IDs and updates the results
// of the poll.
// reactions: Map of emoji to array of user IDs
// excludedUserID: User ID to exlclude from array of user IDs passed in reactions
func (p *Poll) AddResult(reactions map[string][]string, excludedUserID string) {
	for emoji, userIDs := range reactions {
		option, exists := p.emojiToOption[emoji]
		if exists {
			sanitisedUserIDs := make([]string, 0)
			for _, userID := range userIDs {
				if userID == excludedUserID {
					continue
				}
				sanitisedUserIDs = append(sanitisedUserIDs, userID)
			}
			option.AddResult(sanitisedUserIDs)
		}
	}
}

// GetVerdict produces a string that contains the verdict of the poll
// given the current scores of the poll. This will only make sense if
// AddResult was called prior to this.
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
		// If the option with most votes has 0 votes
		verdictMessage = "Unfortunately no votes were received, hence a decision was unable to be made."
	} else if options[0].count == options[1].count {
		// If option with most votes has same vote count has the runner-up
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

// GetParticipants returns a list of IDs of users that participated in this poll.
func (p *Poll) GetParticipants() []string {
	userIDExists := make(map[string]bool)
	uniqueUserIDs := make([]string, 0)
	for _, opt := range p.emojiToOption {
		for _, userID := range opt.userIDs {
			if _, value := userIDExists[userID]; !value {
				userIDExists[userID] = true
				uniqueUserIDs = append(uniqueUserIDs, userID)
			}
		}
	}
	return uniqueUserIDs
}
