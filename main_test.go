package main

import (
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

func TestTruthy(t *testing.T) {
	truthyValues := []string{"1", "true", "TRUE", "t", "yes", "Y"}
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
