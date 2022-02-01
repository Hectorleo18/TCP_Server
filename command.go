package main

type commandID int

const (
	CMD_NICK commandID = iota //<----Autoincremental
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_FILE
	CMD_QUIT
)

type command struct {
	id     commandID
	client *client
	args   []string
}
