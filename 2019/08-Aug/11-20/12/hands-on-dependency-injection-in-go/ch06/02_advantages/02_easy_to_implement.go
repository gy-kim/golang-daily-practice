package advantages

func NewWelcomeSenderV2(mailer *Mailer) *WelcomeSenderV2 {
	return &WelcomeSenderV2{
		mailer: mailer,
	}
}

type WelcomeSenderV2 struct {
	mailer *Mailer
}

func (w *WelcomeSenderV2) Send(to string) error {
	body := w.buildMessage()

	return w.mailer.Send(to, body)
}

func (w *WelcomeSenderV2) buildMessage() string {
	return ""
}
