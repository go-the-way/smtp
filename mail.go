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
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type mail struct {
	auth             smtp.Auth
	host, port, addr string
	message          *Message
}

func Mail() *mail { return &mail{} }

func (m *mail) f(u, p, h string) smtp.Auth              { return smtp.PlainAuth("", u, p, h) }
func (m *mail) PlainAuth(user, pass, host string) *mail { m.Auth(m.f(user, pass, host)); return m }
func (m *mail) Auth(auth smtp.Auth) *mail               { m.auth = auth; return m }
func (m *mail) Message(message *Message) *mail          { m.message = message; return m }

func (m *mail) setAddr(host, port string) {
	m.host = host
	m.port = port
	m.addr = fmt.Sprintf("%s:%s", host, port)
}

func (m *mail) Send(host, port string, TLS bool, opt ...func(tlc *tls.Config)) (err error) {
	m.setAddr(host, port)
	if TLS {
		tlc := &tls.Config{ServerName: host}
		if len(opt) > 0 {
			for _, fn := range opt {
				fn(tlc)
			}
		}
		return m.sendTLS(tlc)
	}
	return m.send()
}

func (m *mail) send() (err error) {
	return smtp.SendMail(m.addr, m.auth, m.message.From.Addr, m.message.To, m.message.toBytes())
}

func (m *mail) sendTLS(tlc *tls.Config) (err error) {
	tlsConn, tlsErr := tls.Dial("tcp", m.addr, tlc)
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
	if err = client.Mail(m.message.From.Addr); err != nil {
		return
	}
	for _, to := range m.message.To {
		if err = client.Rcpt(to); err != nil {
			return
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(m.message.toBytes())
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
