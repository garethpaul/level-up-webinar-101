package main

import (
	"fmt"
	"io"
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
	if err := run(os.Getenv, os.Stdout, sendSMS); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(lookup func(string) string, output io.Writer, send func(smsConfig) error) error {
	config, err := loadConfig(lookup)
	if err != nil {
		return err
	}

	if config.DryRun {
		fmt.Fprintln(output, "Dry run: SMS configuration validated; no message sent.")
		return nil
	}

	if err := send(config); err != nil {
		return fmt.Errorf("send SMS: %w", err)
	}

	fmt.Fprintln(output, "SMS sent successfully!")
	return nil
}

func loadConfig(lookup func(string) string) (smsConfig, error) {
	dryRun, err := parseDryRun(lookup("DRY_RUN"))
	if err != nil {
		return smsConfig{}, err
	}

	config := smsConfig{
		ToPhoneNumber:   strings.TrimSpace(lookup("TO_PHONE_NUMBER")),
		FromPhoneNumber: strings.TrimSpace(lookup("TWILIO_PHONE_NUMBER")),
		AccountSID:      strings.TrimSpace(lookup("TWILIO_ACCOUNT_SID")),
		AuthToken:       strings.TrimSpace(lookup("TWILIO_AUTH_TOKEN")),
		MessageBody:     strings.TrimSpace(lookup("MESSAGE_BODY")),
		DryRun:          dryRun,
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

	var invalid []string
	if !validE164PhoneNumber(config.ToPhoneNumber) {
		invalid = append(invalid, "TO_PHONE_NUMBER")
	}
	if !validE164PhoneNumber(config.FromPhoneNumber) {
		invalid = append(invalid, "TWILIO_PHONE_NUMBER")
	}
	if !config.DryRun && !validTwilioAccountSID(config.AccountSID) {
		invalid = append(invalid, "TWILIO_ACCOUNT_SID")
	}
	if !config.DryRun && !validTwilioAuthToken(config.AuthToken) {
		invalid = append(invalid, "TWILIO_AUTH_TOKEN")
	}
	if len(invalid) > 0 {
		return smsConfig{}, fmt.Errorf("invalid environment variables: %s", strings.Join(invalid, ", "))
	}

	return config, nil
}

func validE164PhoneNumber(value string) bool {
	if len(value) < 3 || len(value) > 16 || value[0] != '+' {
		return false
	}
	if value[1] < '1' || value[1] > '9' {
		return false
	}
	for _, digit := range value[2:] {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return true
}

func validTwilioAccountSID(value string) bool {
	if len(value) != 34 || !strings.HasPrefix(value, "AC") {
		return false
	}
	for _, char := range value[2:] {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}
	return true
}

func validTwilioAuthToken(value string) bool {
	if len(value) != 32 {
		return false
	}
	for _, char := range value {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}
	return true
}

func parseDryRun(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "yes", "y", "on":
		return true, nil
	case "", "0", "false", "f", "no", "n", "off":
		return false, nil
	}
	return false, fmt.Errorf("invalid environment variables: DRY_RUN")
}

func sendSMS(config smsConfig) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:   config.AccountSID,
		Password:   config.AuthToken,
		AccountSid: config.AccountSID,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(config.ToPhoneNumber)
	params.SetFrom(config.FromPhoneNumber)
	params.SetBody(config.MessageBody)

	_, err := client.Api.CreateMessage(params)
	return err
}
