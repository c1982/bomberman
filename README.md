# Bomberman
SMTP Performance Test Tool

[![Build Status](https://travis-ci.org/c1982/bomberman.svg?branch=master)](https://travis-ci.org/c1982/bomberman) [![Go Report Card](https://goreportcard.com/badge/github.com/c1982/bomberman)](https://goreportcard.com/report/github.com/c1982/bomberman)

![Bomberman Logo](https://github.com/c1982/bomberman/blob/master/logo.jpg?raw=true)

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
| size | Set email size Kilobytes (1024 Kilobyte = 1Mbyte). Default: 5Kb |
| helo | SMTP client helo name. Default: mail.server.com | 
| count | Email message count. Default: 10|
| workers | Thread workers for SMTP client. Default: 10 |
| jobs | Job queue lenght in workers. Default: 10 |
| outbound | Outbound IP address for SMTP client |
| showerror | Print SMTP errors |
| balance | Tool is use all IP address for outbound ip with sequental balance. Defalut: false |

## Server Configuration Checklist

* Set SPF value in from email address domain.
* Set PTR record your outbound IP addresses
* Increase ulimit on your server (ulimit -n 10000)

## Usage

Send 50 email to mail.server.com:25 50 workers

```
./bomberman -host=mail.server.com:25 -from=test@mydomain.com -to=user@remotedomain.com -workers=50 -jobs=50 -count=50 -size=75 -balance
```

## Output

```
Bomberman - SMTP Performance Test Tool
--------------------------------------
Message Count		: 1022
Message Size		: 75K
Error			: 168
Start			: 2018-10-12 06:42:56.808098931 +0300 EEST m=+0.000932257
End			: 2018-10-12 06:43:34.049561955 +0300 EEST m=+37.242392313
Time			: 37.241460056s

Source IP Stats:

10.0.5.216		: 256
10.0.5.222		: 256
10.0.5.238		: 256
10.0.5.239		: 256

Destination IP Stats:

5.4.0.248:25            : 856

SMTP Commands:

DATA (854)	        : min. 775.377638ms, max. 20.662139316s, med. 8.870254307s
DIAL (1022)	        : min. 27.323µs, max. 6.000565014s, med. 1.511920428s
HELO (854)	        : min. 34.061919ms, max. 3.80865823s, med. 343.306129ms
MAIL (854)	        : min. 42.455906ms, max. 6.150506182s, med. 943.313477ms
RCPT (854)	        : min. 34.972014ms, max. 3.151397545s, med. 497.683671ms
SUCCESS (854)           : min. 1.480909163s, max. 37.223728269s, med. 15.673002296s
TOUCH (854)	        : min. 112.109537ms, max. 17.759899662s, med. 3.985871341s
```

## Features

* Linux/BSD/Windows supported.
* SMTP RFC 5321 support
* Outbount IP selection
* SMTP Command duration min, max, mean metrics
* Multi-thread support
* Workers and Job Queue support
* Balancing outbound ip address automatically
* Set email body size of Kilobyte
* Count and Report destination IP changes 
* Count and Report source IP changes

## Built With

* [grpool](https://github.com/ivpusic/grpool) - Lightweight Goroutine pool
* [net/smtp](https://golang.org/pkg/net/smtp/) - Golang SMTP Package

## Author

* **Oğuzhan** - *MaestroPanel Tech Lead* - [c1982](https://github.com/c1982)

## ARO

Cemil and Rıza and Osman ^_^
