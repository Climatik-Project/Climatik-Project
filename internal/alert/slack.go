package alert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

type SlackClient interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
	UpdateMessage(channelID, timestamp string, options ...slack.MsgOption) (string, string, string, error)
}

type SlackAlertManager struct {
	webhookURL string
	channel    string
	client     SlackClient
}

type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "INFO"
	AlertLevelWarning  AlertLevel = "WARNING"
	AlertLevelCritical AlertLevel = "CRITICAL"
)

type SlackAlert struct {
	PodName       string
	PowerCapValue int
	CurrentPower  float64
	Devices       map[string]string
	Level         AlertLevel
	Timestamp     time.Time
}

func NewSlackAlertManager(webhookURL, token, channel string) (*SlackAlertManager, error) {
	return &SlackAlertManager{
		webhookURL: webhookURL,
		channel:    channel,
		client:     slack.New(token),
	}, nil
}

func (s *SlackAlertManager) CreateAlert(podName string, powerCapValue int, devices map[string]string) error {
	alert := SlackAlert{
		PodName:       podName,
		PowerCapValue: powerCapValue,
		CurrentPower:  float64(powerCapValue), // You might want to get the actual current power from somewhere
		Devices:       devices,
		Level:         AlertLevelWarning,
		Timestamp:     time.Now(),
	}
	return s.sendAPIAlert(alert)
}

func (s *SlackAlertManager) sendAPIAlert(alert SlackAlert) error {
	attachment := s.createAttachment(alert)

	_, _, err := s.client.PostMessage(
		s.channel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)

	return err
}

func (s *SlackAlertManager) createAttachment(alert SlackAlert) slack.Attachment {
	return slack.Attachment{
		Color:  getColorForAlertLevel(alert.Level),
		Title:  fmt.Sprintf("Power Cap Alert for Pod %s", alert.PodName),
		Text:   fmt.Sprintf("Power cap value of %d exceeded. Current power: %.2f", alert.PowerCapValue, alert.CurrentPower),
		Fields: getFieldsForDevices(alert.Devices),
		Footer: "Climatik Power Management",
		Ts:     json.Number(fmt.Sprintf("%d", alert.Timestamp.Unix())),
	}
}

func getColorForAlertLevel(level AlertLevel) string {
	switch level {
	case AlertLevelInfo:
		return "#36a64f" // Green
	case AlertLevelWarning:
		return "#ffcc00" // Yellow
	case AlertLevelCritical:
		return "#ff0000" // Red
	default:
		return "#808080" // Gray
	}
}

func getFieldsForDevices(devices map[string]string) []slack.AttachmentField {
	var fields []slack.AttachmentField
	for device, value := range devices {
		fields = append(fields, slack.AttachmentField{
			Title: device,
			Value: value,
			Short: true,
		})
	}
	return fields
}

func (s *SlackAlertManager) UpdateAlert(alertTimestamp string, updatedAlert SlackAlert) error {
	attachment := s.createAttachment(updatedAlert)
	attachment.Title = fmt.Sprintf("Updated: %s", attachment.Title)
	attachment.Footer = "Climatik Power Management (Updated)"

	_, _, _, err := s.client.UpdateMessage(
		s.channel,
		alertTimestamp,
		slack.MsgOptionAttachments(attachment),
	)

	return err
}

func (s *SlackAlertManager) ClearAlert(alertTimestamp string, podName string) error {
	attachment := slack.Attachment{
		Color:  "#36a64f", // Green
		Title:  fmt.Sprintf("Cleared: Power Cap Alert for Pod %s", podName),
		Text:   "The power consumption has returned to normal levels.",
		Footer: "Climatik Power Management (Cleared)",
		Ts:     json.Number(alertTimestamp),
	}

	_, _, _, err := s.client.UpdateMessage(
		s.channel,
		alertTimestamp,
		slack.MsgOptionAttachments(attachment),
	)

	return err
}
