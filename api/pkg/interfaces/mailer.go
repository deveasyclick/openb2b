package interfaces

type Mailer interface {
	SendWithAttachment(to, subject, body, filename string, pdfBytes []byte) error
	Send(to, subject, body string) error
}
