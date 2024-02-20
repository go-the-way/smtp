// Copyright 2024 smtp Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package smtp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

const (
	defaultPort         = "25"
	defaultStartTLSPort = "587"
	defaultTLSPort      = "465"
	ctHtml              = "Content-Type: text/html"
)

type mail struct {
	username, password string
	auth               smtp.Auth
	from               string
	to, cc, bcc        []string
	subject            string
	message            string
	headers            []string
	host, port, addr   string
	starttls, tls      bool
	tlc                *tls.Config
}

func Mail() *mail {
	f := func() []string { return make([]string, 0) }
	return &mail{
		to:      f(),
		cc:      f(),
		bcc:     f(),
		headers: f(),
	}
}

func (m *mail) Username(username string) *mail   { m.username = username; return m }
func (m *mail) Password(password string) *mail   { m.password = password; return m }
func (m *mail) Auth(auth smtp.Auth) *mail        { m.auth = auth; return m }
func (m *mail) From(from string) *mail           { m.from = from; return m }
func (m *mail) To(to ...string) *mail            { m.to = to; return m }
func (m *mail) Cc(cc ...string) *mail            { m.cc = cc; return m }
func (m *mail) Bcc(bcc ...string) *mail          { m.bcc = bcc; return m }
func (m *mail) Subject(subject string) *mail     { m.subject = subject; return m }
func (m *mail) Message(message string) *mail     { m.message = message; return m }
func (m *mail) Html(html string) *mail           { m.AddHeader(ctHtml); return m.Message(html) }
func (m *mail) Header(header ...string) *mail    { m.headers = header; return m }
func (m *mail) AddHeader(header ...string) *mail { m.headers = append(m.headers, header...); return m }
func (m *mail) Host(host string) *mail           { m.host = host; return m }
func (m *mail) Port(port string) *mail           { m.port = port; return m }
func (m *mail) TLS(tlc *tls.Config) *mail        { m.tls = true; m.tlc = tlc; return m }
func (m *mail) StartTLS() *mail                  { m.starttls = true; return m }

func (m *mail) Send() (err error) {
	m.setPort()
	m.setAuth()
	if m.tls {
		return m.sendTLS()
	}
	return m.send()
}

func (m *mail) setPort() {
	defer func() { m.addr = fmt.Sprintf("%s:%s", m.host, m.port) }()
	if m.port != "" {
		return
	}
	if m.tls {
		m.port = defaultTLSPort
	} else if m.starttls {
		m.port = defaultStartTLSPort
	} else {
		m.port = defaultPort
	}
}

func (m *mail) setAuth() {
	if m.auth == nil {
		m.auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}
}

func (m *mail) buildMsg() []byte {
	buf := &bytes.Buffer{}
	if from := m.from; from != "" {
		buf.WriteString(fmt.Sprintf("From: %s\r\n", m.from))
	}
	if to := m.to; len(to) > 0 {
		buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ";")))
	}
	if cc := m.cc; len(cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ";")))
	}
	if bcc := m.bcc; len(bcc) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(bcc, ";")))
	}
	if subject := m.subject; subject != "" {
		buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	}
	for _, header := range m.headers {
		buf.WriteString(fmt.Sprintf("%s\r\n", header))
	}
	buf.WriteString("\r\n")
	buf.WriteString(m.message)
	return buf.Bytes()
}

func (m *mail) send() (err error) {
	return smtp.SendMail(m.addr, m.auth, m.from, m.to, m.buildMsg())
}

func (m *mail) sendTLS() (err error) {
	tlsConn, tlsErr := tls.Dial("tcp", m.addr, m.tlc)
	if tlsErr != nil {
		return tlsErr
	}
	client, clientErr := smtp.NewClient(tlsConn, m.host)
	if clientErr != nil {
		return clientErr
	}
	if err = client.Hello(m.host); err != nil {
		return
	}
	if err = client.Auth(m.auth); err != nil {
		return
	}
	if err = client.Mail(m.from); err != nil {
		return
	}
	for _, to := range m.to {
		if err = client.Rcpt(to); err != nil {
			return
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(m.buildMsg())
	if err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	_ = client.Quit()
	_ = client.Close()
	return
}
