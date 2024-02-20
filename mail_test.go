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
	"testing"
)

var (
	from     = "noreply@example.com"
	to       = []string{"ex1@example.com", "ex2@example.com"}
	bcc      = []string{"ex1@example.com", "ex2@example.com"}
	cc       = []string{"ex1@example.com", "ex2@example.com"}
	host     = "smtp.example.com"
	port     = "587" // 25 for smtp, 587 for STARTTLS, 465 for TLS
	portTLS  = "465" // 25 for smtp, 587 for STARTTLS, 465 for TLS
	username = "mailuser"
	password = "mailpasswd"
	subject  = "Test message"
	message  = `This is a test message by Go rwscode/smtp`
)

func TestMailSend(t *testing.T) {
	err := Mail().Username(username).Password(password).
		From(from).To(to...).Cc(cc...).Bcc(bcc...).Subject(subject).Message(message).
		Host(host).Port(port).
		Send()
	if err != nil {
		fmt.Println("send mail error:", err)
		return
	}
	fmt.Println("send mail successful")
}

func TestMailSendTLS(t *testing.T) {
	err := Mail().Username(username).Password(password).
		From(from).To(to...).Cc(cc...).Bcc(bcc...).Subject(subject).Message(message).
		Host(host).Port(portTLS).
		TLS(&tls.Config{InsecureSkipVerify: true, ServerName: host}).Send()
	if err != nil {
		fmt.Println("send mail error:", err)
		return
	}
	fmt.Println("send mail successful")
}
