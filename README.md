# Bomberman
SMTP Load Test tool 

## Installation

bomberman requires Go 1.11 or later.

```
$ go get github.com/c1982/bomberman
```

or

download

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

## Server Configuration Checklist

* Check SPF value in domain dns zone
* Check PTR record your outbound IP address

## Usage

Send 100 email to mail.server.com:25 100 workers

```
./bomberman -host=mail.server.com:25 -from=test@mydomain.com -to=user@remotedomain.com -workers=100 -jobs=50 -count=100 -outbound=YOUR_PUBLIC_IP -helo=mydomain.com -subject="Test Email"
```

## Output


```
Count	: 10
Error	: 0
Start	: 2018-10-08 23:07:17.202394064 +0300 EEST m=+0.000830547
End	: 2018-10-08 23:07:17.792740992 +0300 EEST m=+0.591177492
Time	: 590.346945ms
HELO (10)	: min. 183.457938ms, max. 292.501073ms, med. 197.830802ms
MAIL (10)	: min. 230.36918ms, max. 339.589489ms, med. 246.102343ms
RCPT (10)	: min. 276.966907ms, max. 386.162844ms, med. 292.913193ms
DATA (10)	: min. 417.385041ms, max. 542.583022ms, med. 431.915842ms
SUCCESS (10)	: min. 464.569102ms, max. 589.784965ms, med. 478.931928ms
DIAL (10)	: min. 8.480266ms, max. 16.660143ms, med. 9.456583ms
TOUCH (10)	: min. 136.853636ms, max. 245.775096ms, med. 147.872106ms
```

## Features

* Linux/BSD/Windows supported.
* SMTP RFC 5321 support
* Outbount IP selection
* SMTP Command duration min, max, mean metrics
* Multi-thread support
* Workers and Job Queue support

## Built With

* [grpool](github.com/ivpusic/grpool) - Lightweight Goroutine pool
* [net/smtp](https://golang.org/pkg/net/smtp/) - Golang SMTP Package

## Author

* **Oğuzhan** - *MaestroPanel Tech Lead* - [c1982](https://github.com/c1982)