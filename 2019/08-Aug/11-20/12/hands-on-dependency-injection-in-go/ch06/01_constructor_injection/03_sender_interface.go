package constructor_injection

type Sender interface {
	Send(to string, body string) error
}

func NweWelcomeSenderV2(in Sender) *WelcomeSenderV2 {
	return &WelcomeSenderV2{
		sender: in,
	}
}

type WelcomeSenderV2 struct {
	sender Sender
}

func (w *WelcomeSenderV2) Send(to string) error {
	body := w.buildMessage()

	return w.sender.Send(to, body)
}

func (w *WelcomeSenderV2) buildMessage() string {
	return ""
}
