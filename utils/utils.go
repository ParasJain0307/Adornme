package utils

import (
	"context"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"
)

func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		log.Printf("Attempt %d failed: %v. Retrying in %s...", i+1, err, delay)
		time.Sleep(delay)
	}
	return err
}

func GetPodID() string {
	if podID := os.Getenv("POD_ID"); podID != "" {
		return podID
	}
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return "pod-unknown"
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "requestID", requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if v := ctx.Value("requestID"); v != nil {
		return v.(string)
	}
	return ""
}

func IsEmail(input string) bool {
	_, err := mail.ParseAddress(input)
	return err == nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func IsValidIndianPhone(phone string) bool {
	if len(phone) != 10 {
		return false
	}
	return phone[0] >= '6' && phone[0] <= '9'
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

func CleanInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\u00A0", "")
	input = strings.ReplaceAll(input, "＋", "+")
	return input
}

func IsPhone(input string) bool {
	input = CleanInput(input)
	return phoneRegex.MatchString(input)
}

func NormalizePhone(phone string) string {
	phone = CleanInput(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.TrimPrefix(phone, "+91")
	phone = strings.TrimPrefix(phone, "91")
	return phone
}
