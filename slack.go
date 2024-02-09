package main

import (
	"log"

	"github.com/ashwanthkumar/slack-go-webhook"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
)

func slacklog(session *mgo.Session, project, logString string) error {
	p, err := getProject(session, project)
	if err != nil {
		return err
	}
	if p.SlackWebhookURL != "" {
		payload := slack.Payload{
			Text:    logString,
			Channel: "#" + project,
		}
		err := slack.Send(p.SlackWebhookURL, "", payload)
		if len(err) > 0 {
			for _, e := range err {
				if *flagDebug {
					log.Println(e) // Slack에 로깅하는 것은 중요한 기능이 아니기 때문에 에러 발생시 로그만 출력한다.
				}
			}
		}
	}
	return nil
}

func slacklogV2(client *mongo.Client, project, logString string) error {
	p, err := getProjectV2(client, project)
	if err != nil {
		return err
	}
	if p.SlackWebhookURL != "" {
		payload := slack.Payload{
			Text:    logString,
			Channel: "#" + project,
		}
		err := slack.Send(p.SlackWebhookURL, "", payload)
		if len(err) > 0 {
			for _, e := range err {
				if *flagDebug {
					log.Println(e) // Slack에 로깅하는 것은 중요한 기능이 아니기 때문에 에러 발생시 로그만 출력한다.
				}
			}
		}
	}
	return nil
}
