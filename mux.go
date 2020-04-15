package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

// Mux is a multiplexer that routes requests to the right handler
type Mux struct {
	handlers map[string]func(event Event, param string)
	matchers map[string]Matcher
}

// Matcher implements the ability to match a command
type Matcher interface {
	Match(message string) (matched bool, arg string, err error)
	GetName() string
}

// NewMux creates a new instance of a Mux
func NewMux() *Mux {
	m := Mux{}
	m.handlers = make(map[string]func(event Event, param string))
	m.matchers = make(map[string]Matcher)

	return &m
}

// RegisterHandler will register a handler along with its matcher in a Mux
func (mux *Mux) RegisterHandler(matcher Matcher, handlerFunc func(event Event, param string)) {
	name := matcher.GetName()
	mux.handlers[name] = handlerFunc
	mux.matchers[name] = matcher
}

// Route a request by finding a right handler and routing the
// request along with its parameters to the handler
func (mux *Mux) Route(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	for handlerName, matcher := range mux.matchers {
		matched, arg, err := matcher.Match(m.Content)
		if !matched {
			continue
		}

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		handlerFunc, exists := mux.handlers[handlerName]

		if !exists {
			log.Fatal("Handler matched but no handler func registered.")
		}

		handlerFunc(DiscordEvent{s, m}, arg)
		return
	}
}
