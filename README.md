# Bomberman
SMTP Performance Test Tool

## Installation

bomberman requires Go 1.11 or later.

```
$ go get github.com/c1982/bomberman
```

or


[download](https://github.com/c1982/bomberman/releases)

## Flags

| Flag        | Desc           | 
| ------------- |-------------| 
| host | Remote SMTP server with Port. Default: mail.server.com:25 | 
| from | From email address | 
| to | To email address| 
| subject | Email subject text | 
| body | Email body text | 
| helo | SMTP client helo name. Default: mail.server.com | 
| count | Email message count. Default: 10|
| workers | Thread workers for SMTP client. Default: 100 |
| jobs | Job queue lenght in workers. Default: 50 |
| outbound | Outbound IP address for SMTP client |
| balance | Tool is use all IP address for outbound ip with sequental balance. Defalut: false |


## Server Configuration Checklist

* Check SPF value in domain dns zone
* Check PTR record your outbound IP address

## Usage

Send 10 email to mail.server.com:25 50 workers

```
./bomberman -host=mail.server.com:25 -from=test@mydomain.com -to=user@remotedomain.com -workers=50 -jobs=25 -count=10 -balance
```

## Output

```
Bomberman - SMTP Performance Test Tool
---------------------------------
Message Count	: 1000
Error		: 0
Start		: 2018-10-09 14:52:28.770383156 +0300 EEST m=+0.000754807
End		: 2018-10-09 14:53:40.00788398 +0300 EEST m=+71.238255580
Time		: 1m11.237500773s

Outbounds:

10.10.10.216	: 250
10.10.10.222	: 250
10.10.10.238	: 250
10.10.10.239	: 250

SMTP Commands:

SUCCESS (999)	: min. 508.368509ms, max. 28.674771244s, med. 6.8859535s
DIAL (999)	: min. 8.095947ms, max. 53.888706ms, med. 9.773941ms
TOUCH (999)	: min. 112.503563ms, max. 20.010673679s, med. 912.603978ms
HELO (999)	: min. 159.324966ms, max. 20.057296582s, med. 1.516430795s
MAIL (999)	: min. 206.220242ms, max. 20.104045938s, med. 2.004478873s
RCPT (999)	: min. 260.078873ms, max. 20.151196116s, med. 2.401031078s
DATA (999)	: min. 461.493221ms, max. 23.624083138s, med. 6.311582036s
```

## Features

* Linux/BSD/Windows supported.
* SMTP RFC 5321 support
* Outbount IP selection
* SMTP Command duration min, max, mean metrics
* Multi-thread support
* Workers and Job Queue support
* Balancing outbound ip address automatically

## Built With

* [grpool](https://github.com/ivpusic/grpool) - Lightweight Goroutine pool
* [net/smtp](https://golang.org/pkg/net/smtp/) - Golang SMTP Package

## Author

* **Oğuzhan** - *MaestroPanel Tech Lead* - [c1982](https://github.com/c1982)

## ARO

Cemil and Rıza ^_^