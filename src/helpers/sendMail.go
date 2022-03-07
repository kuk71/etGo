package helpers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

func SendMail(toEmail, subj, body, fromEmail, fromName, login, password, host string) (err error) {

	from := mail.Address{Name: fromName, Address: fromEmail}
	to := mail.Address{Name: toEmail, Address: toEmail}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := "smtp.yandex.ru:465"

	host, _, _ = net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", login, password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		fmt.Println(1)
		log.Panic(err)
	}

	a, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println(2)
		log.Panic(err)
	}

	// Auth
	if err = a.Auth(auth); err != nil {
		fmt.Println(3)
		log.Panic(err)
	}

	// To && From
	if err = a.Mail(from.Address); err != nil {
		fmt.Println(4)
		log.Panic(err)
	}

	if err = a.Rcpt(to.Address); err != nil {
		fmt.Println(5)
		log.Panic(err)
	}

	// Data
	w, err := a.Data()
	if err != nil {
		fmt.Println(6)
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Println(7)
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		fmt.Println(8)
		log.Panic(err)
	}

	a.Quit()

	return
}
