package otp

import "fmt"

type SMSSender struct{}

func (s *SMSSender) Send(phone, otp string) error {
	// TODO: integrate Twilio later
	fmt.Printf("Sending OTP %s to phone %s\n", otp, phone)
	return nil
}
