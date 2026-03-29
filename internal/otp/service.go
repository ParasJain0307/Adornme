package otp

import "errors"

type Service struct {
	store       *Store
	emailSender *EmailSender
	smsSender   *SMSSender
}

func NewService(store *Store, email *EmailSender, sms *SMSSender) *Service {
	return &Service{
		store:       store,
		emailSender: email,
		smsSender:   sms,
	}
}

// 🔥 Send OTP
func (s *Service) SendOTP(userID, target, otpType string) error {

	otp := GenerateOTP()

	err := s.store.Save(userID, otpType, otp)
	if err != nil {
		return err
	}

	if otpType == "email" {
		return s.emailSender.Send(target, otp)
	}

	if otpType == "phone" {
		return s.smsSender.Send(target, otp)
	}

	return errors.New("invalid otp type")
}

// 🔥 Verify OTP
func (s *Service) VerifyOTP(userID, otpType, inputOTP string) (bool, error) {

	storedOTP, err := s.store.Get(userID, otpType)
	if err != nil {
		return false, err
	}

	if storedOTP == inputOTP {
		s.store.Delete(userID, otpType)
		return true, nil
	}

	return false, nil
}
