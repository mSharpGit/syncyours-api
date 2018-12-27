package main

import gomail "gopkg.in/gomail.v2"

type mailer struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
	CC       string
	Subject  string
	Message  string
}

func (m *mailer) Initialize(Port int, Host, Username, Password, From, CC string) {
	m.Host = Host
	m.Port = Port
	m.Username = Username
	m.Password = Password
	m.From = From
	m.CC = CC
}

func (m *mailer) Send(To, Subject, Message string) {
	m.To = To
	m.Subject = Subject
	m.Message = Message

	body := gomail.NewMessage()
	body.SetHeader("From", m.From)
	body.SetHeader("To", m.To)
	if m.CC != "" {
		body.SetAddressHeader("Cc", m.CC, "sync")
	}
	body.SetHeader("Subject", m.Subject)
	body.SetBody("text/html", m.Message)

	d := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(body); err != nil {
		panic(err)
	}
}
