package main

func (app *Config) sendEMail(msg Message) {
	app.Wait.Add(1)
	app.Mailer.MailerChan <- msg
}
