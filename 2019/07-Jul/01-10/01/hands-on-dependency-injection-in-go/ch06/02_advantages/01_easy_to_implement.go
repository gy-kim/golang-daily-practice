package advantages

type WelcomeSender struct {
	Mailer *Mailer
}

func (w *WelcomeSender) Send(to string) error {
	body := w.buildMessage()
	return w.Mailer.Send(to, body)
}

func (w *WelcomeSender) buildMessage() string {
	return ""
}

type Mailer struct{}

func (m *Mailer) Send(to string, body string) error {
	// send email
	return nil
}
