package teller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

const webhookTolerance = 3 * time.Minute

type WebhookEventType = string

const (
	WebhookEventTypeEnrollmentDisconnected             WebhookEventType = "enrollment.disconnected"
	WebhookEventTypeTransactionsProcessed              WebhookEventType = "transactions.processed"
	WebhookEventTypeAccountNumberVerificationProcessed WebhookEventType = "account.number_verification.processed"
	WebhookEventTypeWebhookTest                        WebhookEventType = "webhook.test"
)

type EnrollmentDisconnectedReasonType = string

const (
	EnrollmentDisconnectedReasonTypeDisconnected                         EnrollmentDisconnectedReasonType = "disconnected"
	EnrollmentDisconnectedReasonTypeAccountLocked                        EnrollmentDisconnectedReasonType = "disconnected.account_locked"
	EnrollmentDisconnectedReasonTypeCredentialsInvalid                   EnrollmentDisconnectedReasonType = "disconnected.credentials_invalid"
	EnrollmentDisconnectedReasonTypeEnrollmentInactive                   EnrollmentDisconnectedReasonType = "disconnected.enrollment_inactive"
	EnrollmentDisconnectedReasonTypeUserActionCaptchaRequired            EnrollmentDisconnectedReasonType = "disconnected.user_action.captcha_required"
	EnrollmentDisconnectedReasonTypeUserActionContactInformationRequired EnrollmentDisconnectedReasonType = "disconnected.user_action.contact_information_required"
	EnrollmentDisconnectedReasonTypeUserActionInsufficientPermissions    EnrollmentDisconnectedReasonType = "disconnected.user_action.insufficient_permissions"
	EnrollmentDisconnectedReasonTypeUserActionMFARequired                EnrollmentDisconnectedReasonType = "disconnected.user_action.mfa_required"
	EnrollmentDisconnectedReasonTypeUserActionWebLoginRequired           EnrollmentDisconnectedReasonType = "disconnected.user_action.web_login_required"
)

type WebhookVerificationStatus = string

const (
	WebhookVerificationStatusCompleted WebhookVerificationStatus = "completed"
	WebhookVerificationStatusExpired   WebhookVerificationStatus = "expired"
)

// WebhookPayload captures the optional fields returned in a Teller webhook event.
type WebhookPayload struct {
	EnrollmentID string                           `json:"enrollment_id,omitempty"`
	Reason       EnrollmentDisconnectedReasonType `json:"reason,omitempty"`
	Transactions []TellerTransaction              `json:"transactions,omitempty"`
	AccountID    string                           `json:"account_id,omitempty"`
	Status       WebhookVerificationStatus        `json:"status,omitempty"`
}

// WebhookEvent represents a Teller webhook event.
type WebhookEvent struct {
	ID        string           `json:"id"`
	Payload   WebhookPayload   `json:"payload"`
	Timestamp time.Time        `json:"timestamp"`
	Type      WebhookEventType `json:"type"`
}

// ConstructWebhook verifies the Teller-Signature header and unmarshal the webhook body.
//
// Parameters:
//   - body: raw request body of the webhook. This should be the exact bytes received.
//   - signatureHeader: value of the Teller-Signature header from the webhook request.
//   - signingSecrets: list of valid signing secrets to verify against.
func ConstructWebhook(body []byte, signatureHeader string, signingSecrets []string) (*WebhookEvent, error) {
	if len(signingSecrets) == 0 {
		return nil, errors.New("signing secrets are required")
	}

	tsValue, signatures, err := parseSignatureHeader(signatureHeader)
	if err != nil {
		return nil, err
	}

	tsUnix, err := strconv.ParseInt(tsValue, 10, 64)
	if err != nil {
		return nil, errors.New("invalid signature timestamp")
	}

	timestamp := time.Unix(tsUnix, 0)
	if time.Since(timestamp) > webhookTolerance {
		return nil, errors.New("signature timestamp too old")
	}

	message := []byte(tsValue + "." + string(body))

	verified := false
	for _, secret := range signingSecrets {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(message)
		expected := mac.Sum(nil)

		for _, sig := range signatures {
			provided, err := hex.DecodeString(sig)
			if err != nil {
				continue
			}
			if hmac.Equal(expected, provided) {
				verified = true
				break
			}
		}

		if verified {
			break
		}
	}

	if !verified {
		return nil, errors.New("signature verification failed")
	}

	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func parseSignatureHeader(header string) (string, []string, error) {
	if header == "" {
		return "", nil, errors.New("missing Teller-Signature header")
	}

	var ts string
	var signatures []string

	parts := strings.SplitSeq(header, ",")
	for part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		val := kv[1]

		switch key {
		case "t":
			ts = val
		case "v1":
			signatures = append(signatures, val)
		}
	}

	if ts == "" {
		return "", nil, errors.New("missing signature timestamp")
	}
	if len(signatures) == 0 {
		return "", nil, errors.New("no signatures found")
	}

	return ts, signatures, nil
}
