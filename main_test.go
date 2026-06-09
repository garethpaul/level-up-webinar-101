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

func TestLoadConfigAllowsDryRunWithMalformedCredentials(t *testing.T) {
	config, err := loadConfig(mapLookup(map[string]string{
		"DRY_RUN":             "true",
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  "not-an-account-sid",
		"TWILIO_AUTH_TOKEN":   "not-an-auth-token",
	}))

	if err != nil {
		t.Fatalf("expected dry-run config to ignore Twilio credentials, got %v", err)
	}
	if !config.DryRun {
		t.Fatal("expected dry run to be enabled")
	}
}

func TestLoadConfigRejectsAmbiguousDryRunValue(t *testing.T) {
	_, err := loadConfig(mapLookup(map[string]string{
		"DRY_RUN":             "send-no-message",
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  testAccountSID(),
		"TWILIO_AUTH_TOKEN":   testAuthToken(),
	}))

	if err == nil {
		t.Fatal("expected invalid DRY_RUN error")
	}
	errorText := err.Error()
	if !strings.Contains(errorText, "DRY_RUN") {
		t.Fatalf("expected DRY_RUN name in error, got %q", err)
	}
	if strings.Contains(errorText, "send-no-message") {
		t.Fatalf("error should not echo DRY_RUN value, got %q", err)
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

func TestLoadConfigRejectsMalformedPhoneNumbers(t *testing.T) {
	_, err := loadConfig(mapLookup(map[string]string{
		"DRY_RUN":             "1",
		"TO_PHONE_NUMBER":     "5558675310",
		"TWILIO_PHONE_NUMBER": "+1555INVALID",
	}))

	if err == nil {
		t.Fatal("expected invalid phone number error")
	}
	errorText := err.Error()
	if !strings.Contains(errorText, "TO_PHONE_NUMBER") || !strings.Contains(errorText, "TWILIO_PHONE_NUMBER") {
		t.Fatalf("expected both phone number names in error, got %q", err)
	}
	if strings.Contains(errorText, "5558675310") || strings.Contains(errorText, "+1555INVALID") {
		t.Fatalf("error should not echo phone number values, got %q", err)
	}
}

func TestLoadConfigRejectsMalformedAccountSID(t *testing.T) {
	_, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  "not-an-account-sid",
		"TWILIO_AUTH_TOKEN":   "token",
	}))

	if err == nil {
		t.Fatal("expected invalid account SID error")
	}
	errorText := err.Error()
	if !strings.Contains(errorText, "TWILIO_ACCOUNT_SID") {
		t.Fatalf("expected account SID name in error, got %q", err)
	}
	if strings.Contains(errorText, "not-an-account-sid") {
		t.Fatalf("error should not echo account SID value, got %q", err)
	}
}

func TestLoadConfigRejectsMalformedAuthToken(t *testing.T) {
	_, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  testAccountSID(),
		"TWILIO_AUTH_TOKEN":   "not-an-auth-token",
	}))

	if err == nil {
		t.Fatal("expected invalid auth token error")
	}
	errorText := err.Error()
	if !strings.Contains(errorText, "TWILIO_AUTH_TOKEN") {
		t.Fatalf("expected auth token name in error, got %q", err)
	}
	if strings.Contains(errorText, "not-an-auth-token") {
		t.Fatalf("error should not echo auth token value, got %q", err)
	}
}

func TestLoadConfigRejectsAllZeroCredentialPlaceholders(t *testing.T) {
	accountSID := accountSIDPrefix() + strings.Repeat("0", accountSIDBodyLength())
	authToken := strings.Repeat("0", authTokenLength())

	_, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  accountSID,
		"TWILIO_AUTH_TOKEN":   authToken,
	}))

	if err == nil {
		t.Fatal("expected invalid placeholder credential error")
	}
	errorText := err.Error()
	if !strings.Contains(errorText, "TWILIO_ACCOUNT_SID") || !strings.Contains(errorText, "TWILIO_AUTH_TOKEN") {
		t.Fatalf("expected both credential names in error, got %q", err)
	}
	if strings.Contains(errorText, accountSID) || strings.Contains(errorText, authToken) {
		t.Fatalf("error should not echo credential placeholder values, got %q", err)
	}
}

func TestLoadConfigReadsMessageBodyAndCredentials(t *testing.T) {
	config, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     " +15558675310 ",
		"TWILIO_PHONE_NUMBER": " +15558675309 ",
		"TWILIO_ACCOUNT_SID":  " " + testAccountSID() + " ",
		"TWILIO_AUTH_TOKEN":   " " + testAuthToken() + " ",
		"MESSAGE_BODY":        " Webinar reminder ",
	}))

	if err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
	if config.ToPhoneNumber != "+15558675310" || config.FromPhoneNumber != "+15558675309" {
		t.Fatalf("expected trimmed phone numbers, got %#v", config)
	}
	if config.AccountSID != testAccountSID() || config.AuthToken != testAuthToken() {
		t.Fatalf("expected trimmed credentials, got %#v", config)
	}
	if config.MessageBody != "Webinar reminder" {
		t.Fatalf("expected custom message body, got %q", config.MessageBody)
	}
}

func TestLoadConfigRejectsOversizedMessageBody(t *testing.T) {
	body := strings.Repeat("a", maxMessageBodyCharacters+1)
	_, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  testAccountSID(),
		"TWILIO_AUTH_TOKEN":   testAuthToken(),
		"MESSAGE_BODY":        body,
	}))

	if err == nil {
		t.Fatal("expected invalid message body error")
	}
	errorText := err.Error()
	if !strings.Contains(errorText, "MESSAGE_BODY") {
		t.Fatalf("expected MESSAGE_BODY name in error, got %q", err)
	}
	if strings.Contains(errorText, body) {
		t.Fatalf("error should not echo message body value, got %q", err)
	}
}

func TestLoadConfigAcceptsMaxLengthMessageBody(t *testing.T) {
	body := strings.Repeat("a", maxMessageBodyCharacters)
	config, err := loadConfig(mapLookup(map[string]string{
		"TO_PHONE_NUMBER":     "+15558675310",
		"TWILIO_PHONE_NUMBER": "+15558675309",
		"TWILIO_ACCOUNT_SID":  testAccountSID(),
		"TWILIO_AUTH_TOKEN":   testAuthToken(),
		"MESSAGE_BODY":        body,
	}))

	if err != nil {
		t.Fatalf("expected max-length message body to be accepted, got %v", err)
	}
	if config.MessageBody != body {
		t.Fatalf("expected configured message body to be preserved")
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
		"TWILIO_ACCOUNT_SID":  testAccountSID(),
		"TWILIO_AUTH_TOKEN":   testAuthToken(),
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
		"TWILIO_ACCOUNT_SID":  testAccountSID(),
		"TWILIO_AUTH_TOKEN":   testAuthToken(),
	}), &bytes.Buffer{}, func(smsConfig) error {
		return errors.New("twilio rejected request")
	})

	if err == nil || !strings.Contains(err.Error(), "send SMS: twilio rejected request") {
		t.Fatalf("expected wrapped sender error, got %v", err)
	}
}

func TestParseDryRun(t *testing.T) {
	truthyValues := []string{"1", "true", "TRUE", "t", "yes", "Y", "on"}
	for _, value := range truthyValues {
		parsed, err := parseDryRun(value)
		if err != nil {
			t.Fatalf("expected %q to parse, got %v", value, err)
		}
		if !parsed {
			t.Fatalf("expected %q to be truthy", value)
		}
	}

	falseValues := []string{"", "0", "false", "FALSE", "f", "no", "N", "off", " OFF "}
	for _, value := range falseValues {
		parsed, err := parseDryRun(value)
		if err != nil {
			t.Fatalf("expected %q to parse, got %v", value, err)
		}
		if parsed {
			t.Fatalf("expected %q to be false", value)
		}
	}

	if _, err := parseDryRun("maybe"); err == nil {
		t.Fatal("expected ambiguous dry-run value to be rejected")
	}
}

func TestValidE164PhoneNumber(t *testing.T) {
	valid := []string{"+15558675310", "+442071838750", "+12025550123"}
	for _, value := range valid {
		if !validE164PhoneNumber(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}

	invalid := []string{"", "+", "+0", "15558675310", "+1555INVALID", "+1555 867 5310", "+1234567890123456"}
	for _, value := range invalid {
		if validE164PhoneNumber(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}

func TestValidTwilioAccountSID(t *testing.T) {
	valid := []string{
		testAccountSID(),
		accountSIDPrefix() + strings.Repeat("A", accountSIDBodyLength()),
	}
	for _, value := range valid {
		if !validTwilioAccountSID(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}

	invalid := []string{
		"",
		"AC" + "123",
		string([]byte{83, 75}) + strings.Repeat("1", accountSIDBodyLength()),
		accountSIDPrefix() + strings.Repeat("0", accountSIDBodyLength()),
		"AC" + strings.Repeat("1", 31) + "g",
		testAccountSID() + "0",
	}
	for _, value := range invalid {
		if validTwilioAccountSID(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}

func TestValidTwilioAuthToken(t *testing.T) {
	valid := []string{
		testAuthToken(),
		strings.Repeat("A", authTokenLength()),
	}
	for _, value := range valid {
		if !validTwilioAuthToken(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}

	invalid := []string{
		"",
		strings.Repeat("1", authTokenLength()-1),
		strings.Repeat("1", authTokenLength()) + "0",
		strings.Repeat("0", authTokenLength()),
		strings.Repeat("1", authTokenLength()-1) + "g",
	}
	for _, value := range invalid {
		if validTwilioAuthToken(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}

func testAccountSID() string {
	return accountSIDPrefix() + strings.Repeat("1", accountSIDBodyLength())
}

func testAuthToken() string {
	return strings.Repeat("a", authTokenLength())
}

func accountSIDPrefix() string {
	return string([]byte{65, 67})
}

func accountSIDBodyLength() int {
	return 32
}

func authTokenLength() int {
	return 32
}

func mapLookup(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
