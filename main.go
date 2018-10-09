package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/smtp"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/ivpusic/grpool"
)

func main() {

	var errorCount int
	var totalCount int
	var outbound string
	var metrics []map[string]time.Duration
	var outbndcnt map[string]int

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
	outboundip := flag.String("outbound", "", "-outbound=4.2.2.1")
	balance := flag.Bool("balance", false, "-balance")

	flag.Usage = usage

	flag.Parse()

	errorCount = 0
	totalCount = 0
	startTime := time.Now()
	metrics = []map[string]time.Duration{}
	outbndcnt = map[string]int{}

	pool := grpool.NewPool(*workers, *jobs)

	defer pool.Release()
	pool.WaitCount(*count)

	outbounds, err := ipv4list()

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < *count; i++ {

		if *balance {
			outbound = sequental(i, outbounds)

			if _, ok := outbndcnt[outbound]; ok {
				outbndcnt[outbound] = outbndcnt[outbound] + 1
			} else {
				outbndcnt[outbound] = 1
			}
		} else {
			outbound = *outboundip
		}

		pool.JobQueue <- func() {

			metric, err := SendMail(outbound,
				*host,
				*from,
				*to,
				*subject,
				*body,
				*helo)

			if err != nil {
				fmt.Printf("%d: %v\n", totalCount+1, err)
				errorCount++
			}

			metrics = append(metrics, metric)

			defer func() {
				totalCount++
				pool.JobDone()
			}()
		}
	}

	pool.WaitAll()
	endtime := time.Now()

	fmt.Println("Bomberman - SMTP Performans Test Tool")
	fmt.Println("---------------------------------")
	fmt.Printf("Message Count\t: %d\n", totalCount)
	fmt.Printf("Error\t\t: %d\n", errorCount)
	fmt.Printf("Start\t\t: %v\n", startTime)
	fmt.Printf("End\t\t: %v\n", endtime)
	fmt.Printf("Time\t\t: %v\n", endtime.Sub(startTime))

	if *balance {
		fmt.Println("")
		fmt.Println("Outbounds:")
		fmt.Println("")
		for k, v := range outbndcnt {
			fmt.Printf("%s\t: %d\n", k, v)
		}
	}

	fmt.Println("")
	fmt.Println("SMTP Commands:")
	fmt.Println("")

	mkeys := metricKeys(metrics)

	for i := 0; i < len(mkeys); i++ {
		m := mkeys[i]
		min, max, me := getMetric(m, metrics)
		cnt := countMetric(m, metrics)
		fmt.Printf("%s (%d)\t: min. %v, max. %v, med. %v\n", m, cnt, min, max, me)
	}
}

//SendMail mail.
func SendMail(outbound, smtpServer, from, to, subject, body, helo string) (metric map[string]time.Duration, err error) {

	var wc io.WriteCloser
	var msg string

	startTime := time.Now()
	metric = map[string]time.Duration{}
	host, _, _ := net.SplitHostPort(smtpServer)
	conn, err := newDialer(outbound, smtpServer)

	if err != nil {
		err = fmt.Errorf("DIAL: %v (out:%s)", err, outbound)
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

func ipv4list() (iplist []string, err error) {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return
	}

	for i := 0; i < len(addrs); i++ {

		addr := addrs[i].String()

		if strings.Contains(addr, ":") {
			continue
		}

		nt, _, err := net.ParseCIDR(addr)

		if err != nil {
			continue
		}

		if !nt.IsGlobalUnicast() {
			continue
		}

		iplist = append(iplist, nt.String())
	}

	return
}

func newDialer(outboundip, remotehost string) (conn net.Conn, err error) {

	dialer := &net.Dialer{Timeout: time.Second * 6}
	dialer.LocalAddr = &net.TCPAddr{IP: net.ParseIP(outboundip)}

	conn, err = dialer.Dial("tcp", remotehost)

	return
}

func sequental(index int, list []string) string {

	var ob string
	ln := len(list)

	if index < ln {
		ob = list[index]
	} else {
		li := index % ln
		ob = list[li]
	}

	return ob
}
