package poll

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactoryHappyFlow(t *testing.T) {
	params := "title 30 option1 option2 option3"
	p, err := FromStringParams(params)

	assert.Nil(t, err)
	assert.Equal(t, "title", p.title)
	assert.Equal(t, 30, p.durationInSeconds)
	assert.ElementsMatch(t, []string{"option1", "option2", "option3"}, p.options)
}

func TestFactoryNotEnoughOptions(t *testing.T) {
	params := "title 30 option1"
	p, err := FromStringParams(params)

	assert.Nil(t, p)
	assert.Equal(t, "aside from a `title` and `poll-duration`, please provide at least two other options for a vote to make sense", err.Error())
}

func TestFactoryNonIntegerDuration(t *testing.T) {
	params := "title 30.0 option1 option2"
	p, err := FromStringParams(params)

	assert.Nil(t, p)
	assert.Equal(t, "`poll-duration` needs to be an integer", err.Error())
}

func TestFactoryTooManyOptions(t *testing.T) {
	params := "title 30 o o o o o o o o o o o"
	p, err := FromStringParams(params)

	assert.Nil(t, p)
	assert.Equal(t, "we currently only support a maximum number of **10** options in a poll", err.Error())
}

func TestGeneratePollMessage(t *testing.T) {
	p := &Poll{
		title:             "title",
		options:           []string{"option1", "option2"},
		durationInSeconds: 5,
		emojiToOption:     make(map[string]*Option),
	}

	expectedMessage := "A poll has been started!\n**title**\n1️⃣. option1\n2️⃣. option2\nExercise your right to vote by reacting accordingly! The poll will close in 5s."
	assert.Equal(t, expectedMessage, p.GeneratePollMessage(), "should have right message")
}

func TestAddResult(t *testing.T) {
	var emojiToOption = map[string]*Option{
		"emojiOne": &Option{},
		"emojiTwo": &Option{},
	}

	var reactionCounts = map[string]int{
		"emojiOne": 10,
		"emojiTwo": 20,
	}

	p := &Poll{emojiToOption: emojiToOption}
	p.AddResult(reactionCounts)

	assert.Equal(t, 10, p.emojiToOption["emojiOne"].count, "should have updated count")
	assert.Equal(t, 20, p.emojiToOption["emojiTwo"].count, "should have updated count")
}

func TestTiedOptionsMessage(t *testing.T) {
	options := []Option{
		Option{name: "option1", count: 5},
		Option{name: "option2", count: 5},
		Option{name: "option3", count: 5},
		Option{name: "option4", count: 4},
		Option{name: "option5", count: 3},
	}

	res := getTiedOptions(options)
	expectedText := "**option1**, **option2**, **option3**"
	assert.Equal(t, expectedText, res, "returns bolded options")
}

func TestGetVerdictClearWinner(t *testing.T) {
	var emojiToOption = map[string]*Option{
		"emojiOne": &Option{name: "option1", count: 5},
		"emojiTwo": &Option{name: "option2", count: 4},
	}

	p := &Poll{emojiToOption: emojiToOption}
	res := p.GetVerdict()
	expectedResult := "**Results:**\nThe people have spoken, **option1** it shall be."
	assert.Equal(t, expectedResult, res)
}

func TestGetVerdictNoVotes(t *testing.T) {
	var emojiToOption = map[string]*Option{
		"emojiOne": &Option{name: "option1", count: 0},
		"emojiTwo": &Option{name: "option2", count: 0},
	}

	p := &Poll{emojiToOption: emojiToOption}
	res := p.GetVerdict()
	expectedResult := "**Results:**\nUnfortunately no votes were received, hence a decision was unable to be made."
	assert.Equal(t, expectedResult, res)
}

func TestGetVerdictTies(t *testing.T) {
	var emojiToOption = map[string]*Option{
		"emojiOne":   &Option{name: "option1", count: 2},
		"emojiTwo":   &Option{name: "option2", count: 2},
		"emojiThree": &Option{name: "option3", count: 2},
	}

	p := &Poll{emojiToOption: emojiToOption}
	res := p.GetVerdict()
	assert.Contains(t, res, "Looks like we have a tie between")
}
