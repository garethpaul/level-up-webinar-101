package main

import (
	"fmt"
	"os"
	"strings"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const defaultMessageBody = "Hello from Golang!"

type smsConfig struct {
	ToPhoneNumber   string
	FromPhoneNumber string
	AccountSID      string
	AuthToken       string
	MessageBody     string
	DryRun          bool
}

func main() {
	config, err := loadConfig(os.Getenv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if config.DryRun {
		fmt.Println("Dry run: SMS configuration validated; no message sent.")
		return
	}

	if err := sendSMS(config); err != nil {
		fmt.Fprintf(os.Stderr, "send SMS: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SMS sent successfully!")
}

func loadConfig(lookup func(string) string) (smsConfig, error) {
	config := smsConfig{
		ToPhoneNumber:   strings.TrimSpace(lookup("TO_PHONE_NUMBER")),
		FromPhoneNumber: strings.TrimSpace(lookup("TWILIO_PHONE_NUMBER")),
		AccountSID:      strings.TrimSpace(lookup("TWILIO_ACCOUNT_SID")),
		AuthToken:       strings.TrimSpace(lookup("TWILIO_AUTH_TOKEN")),
		MessageBody:     strings.TrimSpace(lookup("MESSAGE_BODY")),
		DryRun:          truthy(lookup("DRY_RUN")),
	}

	if config.MessageBody == "" {
		config.MessageBody = defaultMessageBody
	}

	var missing []string
	if config.ToPhoneNumber == "" {
		missing = append(missing, "TO_PHONE_NUMBER")
	}
	if config.FromPhoneNumber == "" {
		missing = append(missing, "TWILIO_PHONE_NUMBER")
	}
	if !config.DryRun {
		if config.AccountSID == "" {
			missing = append(missing, "TWILIO_ACCOUNT_SID")
		}
		if config.AuthToken == "" {
			missing = append(missing, "TWILIO_AUTH_TOKEN")
		}
	}
	if len(missing) > 0 {
		return smsConfig{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return config, nil
}

func truthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "yes", "y":
		return true
	default:
		return false
	}
}

func sendSMS(config smsConfig) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(config.ToPhoneNumber)
	params.SetFrom(config.FromPhoneNumber)
	params.SetBody(config.MessageBody)

	_, err := client.Api.CreateMessage(params)
	return err
}
