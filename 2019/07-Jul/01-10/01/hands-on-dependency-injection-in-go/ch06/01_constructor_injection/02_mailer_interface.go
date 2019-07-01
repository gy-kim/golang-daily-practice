package constructor_injection

type MailerInterface interface {
	Send(to string, body string) error
	Receive(address string) (string, error)
}
