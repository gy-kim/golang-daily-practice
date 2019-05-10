package constructor_injection

// Mailer sends and receives emails
type MailerInterface interface {
	Send(to string, bodpy string) error
	Receive(address string) (string, error)
}
