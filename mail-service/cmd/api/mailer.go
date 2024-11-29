package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Defines a struct for the Mailer object itself
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// Defines a single email
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// SendSMTPMessage sends an email using a provided Message struct and the SMTP
// server connection details specified in the Mailer struct.
func (m *Mail) SendSMTPMessage(msg Message) error {

	// If From is not specified, use the default From address for the Mailer
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	// Likewise for From Name
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// Create a DataMap based on the messages data
	data := map[string]any{
		"message": msg.Data,
	}
	msg.DataMap = data

	// Create an HTML-formatted version of the email
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// Create a plain-formatted version of the email
	plainText, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// Set up the SMTP client
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Connect to the SMTP server
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Set up the email message
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainText).             // Default body is plain message
		AddAlternative(mail.TextHTML, formattedMessage) // Alternative body is HTML message

	// If there are any attachments, add them to the message
	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	// Send the email
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

// getEncryption returns the appropriate mail.Encryption type based on the input string.
// Supported values are "tls", "ssl", and "none". Any other value defaults to "tls".
func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

// buildPlainTextMessage builds a plain text email message from a template and a Message object.
// It parses the template, executes it with the Message data, and then returns the formatted plain text string.
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	// Create an email template based on the HTML file
	templateToRender := "./templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// Try to execute the template, checking for errors in template file
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// Convert the template to a string
	plainMessage := tpl.String()

	return plainMessage, nil
}

// buildHTMLMessage builds an HTML email message from a template and a Message object.
// It parses the template, executes it with the Message data, and then inlines the CSS.
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	// Create an email template based on the HTML file
	templateToRender := "./templates/mail.html.gohtml"
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// Try to execute the template, checking for errors in template file
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// Convert the template to a string
	formattedMessage := tpl.String()

	// Inline the CSS
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// inlineCSS takes an HTML string, inlines the CSS, and returns the new HTML string.
// It uses the premailer package to do the actual inlining.
func (m *Mail) inlineCSS(s string) (string, error) {

	// Set some options for premailer
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	// Create the premailer object with the above options
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	// Finally, transform the CSS in <style> tags into inline CSS
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
