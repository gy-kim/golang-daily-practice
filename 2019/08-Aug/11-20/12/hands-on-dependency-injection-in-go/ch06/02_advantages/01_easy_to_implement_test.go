package advantages

func ExampleWelcomeSender_Send() {
	welcomeSender := &WelcomeSender{
		Mailer: &Mailer{},
	}
	welcomeSender.Send("me@home.com")
}
