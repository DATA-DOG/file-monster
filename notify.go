package main

import "gopkg.in/jmcvetta/napping.v1"

type slackNotification struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Channel   string `json:"channel"`
	Text      string `json:"text"`
}

func notifyUser(user string, message string) {
	notif := slackNotification{}

	notif.Channel = "@" + user
	notif.Username = "File Monster"
	notif.IconEmoji = ":japanese_ogre:"
	notif.Text = message

	napping.Post(config.SlackWebookURL, notif, nil, nil)
}
