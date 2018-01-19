package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

func main() {

	host := flag.String("host", "localhost:25", "-host=example.org:25")

	//TODO: Autogenerate
	from := flag.String("from", "me@example.org", "-from=me@example.org")

	//TODO: Autogenerate
	to := flag.String("to", "to@example.net", "-to=me@example.net")

	//TODO: Autogenerate
	subject := flag.String("subject", "Test Email", "-subject=Test Email")

	//TODO: Autogenerate
	//TODO: Import from file
	body := flag.String("body", "Load Test Generator", "-body=Load Test Generator")
	helo := flag.String("helo", "mail.example.org", "-helo=mail.example.org")
	concurrent := flag.Int("thread", 10, "-thread=10")

	flag.Parse()

	ch := make(chan time.Duration)

	for i := 0; i < *concurrent; i++ {
		go SendMail(*host,
			*from,
			*to,
			*subject,
			*body,
			*helo, ch)
	}

	for i := 0; i < *concurrent; i++ {
		fmt.Println(<-ch)
	}

}

func SendMail(host, from, to, subject, body, helo string, ch chan<- time.Duration) {

	startTime := time.Now()

	c, err := smtp.Dial(host)

	if err != nil {
		log.Println(err)
		ch <- time.Now().Sub(startTime)
		return
	}

	defer c.Close()

	if err := c.Hello(helo); err != nil {
		log.Println(err)
		ch <- time.Now().Sub(startTime)
		return
	}

	if err := c.Mail(from); err != nil {
		log.Println(err)
		ch <- time.Now().Sub(startTime)
		return
	}

	if err := c.Rcpt(to); err != nil {
		log.Println(err)
		ch <- time.Now().Sub(startTime)
		return
	}

	msg := ""
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("\r\n%s", body)

	wc, err := c.Data()
	_, err = fmt.Fprintf(wc, msg)
	_ = wc.Close()

	err = c.Quit()

	if err != nil {
		log.Println(err)
		ch <- time.Now().Sub(startTime)
		return
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	ch <- elapsed
}
