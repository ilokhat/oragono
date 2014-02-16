package irc

import (
	"errors"
	"fmt"
)

//
// simple types
//

// a string with wildcards
type Mask string

// add, remove, list modes
type ModeOp rune

func (op ModeOp) String() string {
	return string(op)
}

// user mode flags
type UserMode rune

func (mode UserMode) String() string {
	return fmt.Sprintf("%c", mode)
}

type Phase uint

type Numeric uint

func (code Numeric) String() string {
	return fmt.Sprintf("%03d", code)
}

// channel mode flags
type ChannelMode rune

func (mode ChannelMode) String() string {
	return fmt.Sprintf("%c", mode)
}

type ChannelNameMap map[string]*Channel

func (channels ChannelNameMap) Add(channel *Channel) error {
	if channels[channel.name] != nil {
		return fmt.Errorf("%s: already set", channel.name)
	}
	channels[channel.name] = channel
	return nil
}

func (channels ChannelNameMap) Remove(channel *Channel) error {
	if channel != channels[channel.name] {
		return fmt.Errorf("%s: mismatch", channel.name)
	}
	delete(channels, channel.name)
	return nil
}

type ClientNameMap map[string]*Client

var (
	ErrNickMissing   = errors.New("nick missing")
	ErrNicknameInUse = errors.New("nickname in use")
)

func (clients ClientNameMap) Add(client *Client) error {
	if !client.HasNick() {
		return ErrNickMissing
	}
	if clients[client.nick] != nil {
		return ErrNicknameInUse
	}
	clients[client.nick] = client
	return nil
}

func (clients ClientNameMap) Remove(client *Client) error {
	if clients[client.nick] != client {
		return fmt.Errorf("%s: mismatch", client.nick)
	}
	delete(clients, client.nick)
	return nil
}

type ChannelModeSet map[ChannelMode]bool

type ClientSet map[*Client]ChannelModeSet

func (clients ClientSet) Add(client *Client) {
	clients[client] = make(ChannelModeSet)
}

func (clients ClientSet) Remove(client *Client) {
	delete(clients, client)
}

func (clients ClientSet) HasMode(client *Client, mode ChannelMode) bool {
	modes, ok := clients[client]
	if !ok {
		return false
	}
	return modes[mode]
}

func (clients ClientSet) Has(client *Client) bool {
	_, ok := clients[client]
	return ok
}

type ChannelSet map[*Channel]bool

func (channels ChannelSet) Add(channel *Channel) {
	channels[channel] = true
}

func (channels ChannelSet) Remove(channel *Channel) {
	delete(channels, channel)
}

func (channels ChannelSet) First() *Channel {
	for channel := range channels {
		return channel
	}
	return nil
}

//
// interfaces
//

type Identifier interface {
	Id() string
	Nick() string
}

type Replier interface {
	Reply(Reply)
}

type Reply interface {
	Format(*Client) []string
	Source() Identifier
}

type Command interface {
	Name() string
	Client() *Client
	Source() Identifier
	Reply(Reply)
}

type ServerCommand interface {
	Command
	HandleServer(*Server)
}

type AuthServerCommand interface {
	Command
	HandleAuthServer(*Server)
}

type RegServerCommand interface {
	Command
	HandleRegServer(*Server)
}

type ChannelCommand interface {
	Command
	HandleChannel(channel *Channel)
}

//
// structs
//

type UserMask struct {
	nickname Mask
	username Mask
	hostname Mask
}

func (mask *UserMask) String() string {
	return fmt.Sprintf("%s!%s@%s", mask.nickname, mask.username, mask.hostname)
}