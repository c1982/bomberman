package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/ivpusic/grpool"
)

func main() {

	var errorCount int
	var totalCount int
	var metrics []map[string]time.Duration

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	host := flag.String("host", "localhost:25", "-host=example.org:25")
	from := flag.String("from", "me@example.org", "-from=me@example.org")
	to := flag.String("to", "to@example.net", "-to=me@example.net")
	subject := flag.String("subject", "Test Email", "-subject=Test Email")
	body := flag.String("body", "Load Test Generator", "-body=Load Test Generator")
	helo := flag.String("helo", "mail.example.org", "-helo=mail.example.org")
	count := flag.Int("count", 10, "-count=10")
	workers := flag.Int("workers", 100, "-workers=100")
	jobs := flag.Int("jobs", 50, "-jobs=50")
	outboundip := flag.String("outbound", "127.0.0.1", "-outbound=127.0.0.1")
	flag.Usage = usage

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	errorCount = 0
	totalCount = 0
	metrics = []map[string]time.Duration{}

	pool := grpool.NewPool(*workers, *jobs)

	defer pool.Release()

	pool.WaitCount(*count)

	startTime := time.Now()

	for i := 0; i < *count; i++ {

		pool.JobQueue <- func() {
			totalCount++

			metric, err := SendMail(*outboundip,
				*host,
				*from,
				*to,
				*subject,
				*body,
				*helo)

			if err != nil {
				fmt.Printf("%d: %v", i, err)
				errorCount++
			}

			metrics = append(metrics, metric)

			defer pool.JobDone()
		}
	}

	pool.WaitAll()
	endtime := time.Now()

	fmt.Printf("Count\t: %d\n", totalCount)
	fmt.Printf("Error\t: %d\n", errorCount)
	fmt.Printf("Start\t: %v\n", startTime)
	fmt.Printf("End\t: %v\n", endtime)
	fmt.Printf("Time\t: %v\n", endtime.Sub(startTime))

	mkeys := metricKeys(metrics)

	for i := 0; i < len(mkeys); i++ {
		m := mkeys[i]
		min, max, me := getMetric(m, metrics)
		cnt := countMetric(m, metrics)
		fmt.Printf("%s (%d)\t: min. %v, max. %v, med. %v\n", m, cnt, min, max, me)
	}
}

//SendMail mail.
func SendMail(outboundIPaddr, smtpServer, from, to, subject, body, helo string) (metric map[string]time.Duration, err error) {
	metric = map[string]time.Duration{}

	var wc io.WriteCloser
	var msg string

	startTime := time.Now()

	dialer := &net.Dialer{Timeout: time.Second * 6}
	dialer.LocalAddr = &net.TCPAddr{IP: net.ParseIP(outboundIPaddr)}

	host, _, _ := net.SplitHostPort(smtpServer)
	conn, err := dialer.Dial("tcp", smtpServer)

	if err != nil {
		err = fmt.Errorf("DIAL: %v", err)
		metric["DIAL"] = time.Now().Sub(startTime)

		return
	}

	metric["DIAL"] = time.Now().Sub(startTime)

	c, err := smtp.NewClient(conn, host)

	if err != nil {

		err = fmt.Errorf("TOUCH: %v", err)
		metric["TOUCH"] = time.Now().Sub(startTime)

		return
	}

	metric["TOUCH"] = time.Now().Sub(startTime)
	defer c.Close()

	err = c.Hello(helo)

	if err != nil {
		err = fmt.Errorf("HELO: %v", err)
		metric["HELO"] = time.Now().Sub(startTime)

		return
	}

	metric["HELO"] = time.Now().Sub(startTime)

	err = c.Mail(from)

	if err != nil {
		err = fmt.Errorf("MAIL: %v", err)
		metric["MAIL"] = time.Now().Sub(startTime)

		return
	}

	metric["MAIL"] = time.Now().Sub(startTime)

	err = c.Rcpt(to)

	if err != nil {
		err = fmt.Errorf("RCPT: %v", err)
		metric["RCPT"] = time.Now().Sub(startTime)
		return
	}

	metric["RCPT"] = time.Now().Sub(startTime)

	msg = ""
	msg += fmt.Sprintf("from: <%s>\r\n", from)
	msg += fmt.Sprintf("to: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("\r\n%s", body)

	wc, err = c.Data()

	if err != nil {
		err = fmt.Errorf("DATA: %v", err)
		metric["DATA"] = time.Now().Sub(startTime)
		return
	}

	_, err = fmt.Fprintf(wc, msg)

	err = wc.Close()

	if err != nil {
		err = fmt.Errorf("DATA: %v", err)
		metric["DATA"] = time.Now().Sub(startTime)
		return
	}

	metric["DATA"] = time.Now().Sub(startTime)

	err = c.Quit()

	if err != nil {
		err = fmt.Errorf("QUIT: %v", err)
		metric["QUIT"] = time.Now().Sub(startTime)
		return
	}

	metric["SUCCESS"] = time.Now().Sub(startTime)

	return
}

func getMetric(name string, metrics []map[string]time.Duration) (max, min, med time.Duration) {

	totaldur, _ := time.ParseDuration("0ms")
	list := []time.Duration{}

	for i := 0; i < len(metrics); i++ {
		m := metrics[i]

		if t, ok := m[name]; ok {
			totaldur += t
			list = append(list, t)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] > list[j]
	})

	min = list[0]
	max = list[len(list)-1]
	med = totaldur / time.Duration(len(list))

	return
}

func countMetric(name string, metrics []map[string]time.Duration) (cnt int) {

	for i := 0; i < len(metrics); i++ {
		m := metrics[i]

		for mkey := range m {

			if mkey == name {
				cnt++
			}
		}
	}

	return
}

func metricKeys(metrics []map[string]time.Duration) (keys []string) {

	for i := 0; i < len(metrics); i++ {
		m := metrics[i]

		for mkey := range m {
			contain := isContain(mkey, keys)

			if !contain {
				keys = append(keys, mkey)
			}
		}
	}

	return
}

func isContain(key string, keys []string) bool {

	exists := false

	for z := 0; z < len(keys); z++ {
		if keys[z] == key {
			exists = true
			break
		}
	}

	return exists
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "USAGE:")
	fmt.Fprintln(os.Stderr, "./bomberman -host=mail.server.com:25 -from=test@mydomain.com -to=user@remotedomain.com -workers=100 -jobs=50 -count=100 -outbound=YOUR_PUBLIC_IP -helo=mydomain.com -subject=\"Test Email\"")
	fmt.Fprintln(os.Stderr, "")
}
