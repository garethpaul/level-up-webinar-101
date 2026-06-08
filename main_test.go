package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestLoadConfigRequiresPhoneNumbers(t *testing.T) {
	_, err := loadConfig(mapLookup(map[string]string{
		"DRY_RUN": "1",
	}))

	if err == nil {
		t.Fatal("expected missing phone number error")
	}
	if !strings.Contains(err.Error(), "TO_PHONE_NUMBER") || !strings.Contains(err.Error(), "TWILIO_PHONE_NUMBER") {
		t.Fatalf("expected both phone numbers in error, got %q", err)
	}
}

func TestLoadConfigAllowsDryRunWithoutCredentials(t *testing.T) {
	config, err := loadConfig(mapLookup(map[string]string{
		"DRY_RUN":             "true",
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
	}))

	if err != nil {
		t.Fatalf("expected dry-run config without credentials, got %v", err)
	}
	if !config.DryRun {
		t.Fatal("expected dry run to be enabled")
	}
	if config.MessageBody != defaultMessageBody {
		t.Fatalf("expected default message body, got %q", config.MessageBody)
	}
}

func TestLoadConfigRequiresCredentialsWhenSending(t *testing.T) {
	_, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
	}))

	if err == nil {
		t.Fatal("expected missing credential error")
	}
	if !strings.Contains(err.Error(), "TWILIO_ACCOUNT_SID") || !strings.Contains(err.Error(), "TWILIO_AUTH_TOKEN") {
		t.Fatalf("expected both credentials in error, got %q", err)
	}
}

func TestLoadConfigReadsMessageBodyAndCredentials(t *testing.T) {
	config, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     " +15558675310 ",
		"TWILIO_PHONE_NUMBER": " +15558675309 ",
		"TWILIO_ACCOUNT_SID":  " AC123 ",
		"TWILIO_AUTH_TOKEN":   " token ",
		"MESSAGE_BODY":        " Webinar reminder ",
	}))

	if err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
	if config.ToPhoneNumber != "+15558675310" || config.FromPhoneNumber != "+15558675309" {
		t.Fatalf("expected trimmed phone numbers, got %#v", config)
	}
	if config.AccountSID != "AC123" || config.AuthToken != "token" {
		t.Fatalf("expected trimmed credentials, got %#v", config)
	}
	if config.MessageBody != "Webinar reminder" {
		t.Fatalf("expected custom message body, got %q", config.MessageBody)
	}
}

func TestRunDryRunSkipsSender(t *testing.T) {
	var output bytes.Buffer
	called := false

	err := run(mapLookup(map[string]string{
		"DRY_RUN":             "1",
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
	}), &output, func(smsConfig) error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("expected dry run to succeed, got %v", err)
	}
	if called {
		t.Fatal("dry run should not call sender")
	}
	if !strings.Contains(output.String(), "Dry run") {
		t.Fatalf("expected dry-run output, got %q", output.String())
	}
}

func TestRunSendsConfiguredMessage(t *testing.T) {
	var output bytes.Buffer
	var sent smsConfig

	err := run(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  "AC123",
		"TWILIO_AUTH_TOKEN":   "token",
		"MESSAGE_BODY":        "Webinar reminder",
	}), &output, func(config smsConfig) error {
		sent = config
		return nil
	})

	if err != nil {
		t.Fatalf("expected send to succeed, got %v", err)
	}
	if sent.MessageBody != "Webinar reminder" {
		t.Fatalf("expected configured message body, got %q", sent.MessageBody)
	}
	if !strings.Contains(output.String(), "SMS sent successfully") {
		t.Fatalf("expected success output, got %q", output.String())
	}
}

func TestRunWrapsSenderError(t *testing.T) {
	err := run(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  "AC123",
		"TWILIO_AUTH_TOKEN":   "token",
	}), &bytes.Buffer{}, func(smsConfig) error {
		return errors.New("twilio rejected request")
	})

	if err == nil || !strings.Contains(err.Error(), "send SMS: twilio rejected request") {
		t.Fatalf("expected wrapped sender error, got %v", err)
	}
}

func TestTruthy(t *testing.T) {
	truthyValues := []string{"1", "true", "TRUE", "t", "yes", "Y", "on"}
	for _, value := range truthyValues {
		if !truthy(value) {
			t.Fatalf("expected %q to be truthy", value)
		}
	}

	falseValues := []string{"", "0", "false", "no"}
	for _, value := range falseValues {
		if truthy(value) {
			t.Fatalf("expected %q to be false", value)
		}
	}
}

func mapLookup(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
