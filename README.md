# smtp

The extended go src/net/smtp package

## Example
```go

package main

import (
  "fmt"
  
  "github.com/rwscode/smtp"
)

var (
  from     = "noreply@example.com"
  to       = []string{"ex1@example.com", "ex2@example.com"}
  host     = "smtp.example.com"
  port     = "587" // 25 for smtp, 587 for STARTTLS, 465 for TLS
  portTLS  = "465" // 25 for smtp, 587 for STARTTLS, 465 for TLS
  username = "mailuser"
  password = "mailpasswd"
  subject  = "Test message"
  message  = `This is a test message by Go rwscode/smtp`
)

func main() {
  err := smtp.Mail().Message(&smtp.Message{
    From:    &smtp.Email{Addr: from},
    To:      to,
    Subject: subject,
    Content: smtp.Content{ContentType: smtp.Plain, Body: message},
  }).PlainAuth(username, password, host).Send(host, port, false)
  if err != nil {
    fmt.Println("send mail error:", err)
    return
  }
  fmt.Println("send mail successful")
}

func mainTLS() {
  err := smtp.Mail().Message(&smtp.Message{
    From:    &smtp.Email{Addr: from},
    To:      to,
    Subject: subject,
    Content: smtp.Content{ContentType: smtp.Plain, Body: message},
  }).PlainAuth(username, password, host).Send(host, portTLS, true)
  if err != nil {
    fmt.Println("send mail error:", err)
    return
  }
  fmt.Println("send mail successful")
}
```

## SMTP transport example

```
S: 220 smtp.example.com ESMTP Postfix
C: HELO relay.example.org
S: 250 Hello relay.example.org, I am glad to meet you
C: MAIL FROM:<bob@example.org>
S: 250 Ok
C: RCPT TO:<alice@example.com>
S: 250 Ok
C: RCPT TO:<theboss@example.com>
S: 250 Ok
C: DATA
S: 354 End data with <CR><LF>.<CR><LF>
C: From: "Bob Example" <bob@example.org>
C: To: "Alice Example" <alice@example.com>
C: Cc: theboss@example.com
C: Date: Tue, 15 Jan 2008 16:02:43 -0500
C: Subject: Test message
C:
C: Hello Alice.
C: This is a test message with 5 header fields and 4 lines in the message body.
C: Your friend,
C: Bob
C: .
S: 250 Ok: queued as 12345
C: QUIT
S: 221 Bye
```

## SMTP Extensions

- 8BITMIME – 8 bit data transmission, [RFC 6152](https://datatracker.ietf.org/doc/html/rfc5152)
- ATRN – Authenticated TURN for On-Demand Mail Relay, [RFC 2645](https://datatracker.ietf.org/doc/html/rfc2645)
- AUTH – Authenticated SMTP, [RFC 4954](https://datatracker.ietf.org/doc/html/rfc4954)
- CHUNKING – Chunking, [RFC 4954](https://datatracker.ietf.org/doc/html/rfc3030)
- DSN – Delivery status notification, [RFC 3461](https://datatracker.ietf.org/doc/html/rfc3461) (See Variable envelope
  return path)
- ETRN – Extended version of remote message queue starting command
  TURN, [RFC 1985](https://datatracker.ietf.org/doc/html/rfc1985)
- HELP – Supply helpful information, [RFC 821](https://datatracker.ietf.org/doc/html/rfc821)
- PIPELINING – Command pipelining, [RFC 2920](https://datatracker.ietf.org/doc/html/rfc2920)
- SIZE – Message size declaration, [RFC 1870](https://datatracker.ietf.org/doc/html/rfc1870)
- STARTTLS – Transport Layer Security, [RFC 3207](https://datatracker.ietf.org/doc/html/rfc3207)
- SMTPUTF8 – Allow UTF-8 encoding in mailbox names and header
  fields, [RFC 6531](https://datatracker.ietf.org/doc/html/rfc6531)
- UTF8SMTP – Allow UTF-8 encoding in mailbox names and header
  fields, [RFC 5336](https://datatracker.ietf.org/doc/html/rfc5336) (deprecated[28])
