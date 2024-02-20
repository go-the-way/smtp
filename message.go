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
	"encoding/base64"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type (
	Message struct {
		From, ReplyTo *Email
		To, Cc, Bcc   []string
		Subject       string
		Content       Content
		Attachments   []Attachment
	}
	Attachment struct {
		Filename    string
		ContentType string
		Data        []byte
	}
)

func (m *Message) toBytes() []byte {
	buf := &bytes.Buffer{}
	if from := m.From; from != nil {
		buf.WriteString("From: " + from.toString() + crlf)
	}
	if to := m.To; len(to) > 0 {
		buf.WriteString("To: " + strings.Join(m.To, ",") + crlf)
	}
	if cc := m.Cc; len(cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.Cc, ",") + crlf)
	}
	if bcc := m.Bcc; len(bcc) > 0 {
		buf.WriteString("Bcc: " + strings.Join(m.Cc, ",") + crlf)
	}
	if rto := m.ReplyTo; rto != nil {
		buf.WriteString("Reply-To: " + rto.toString() + crlf)
	}
	if m.Subject != "" {
		buf.WriteString("Subject: " + m.Subject + crlf)
	}
	buf.WriteString("MIME-Version: 1.0" + crlf)
	m.writeMixed(buf)
	return buf.Bytes()
}

func (m *Message) writeMixed(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: multipart/mixed; boundary=\"MixedBoundaryString\"" + crlf + crlf)
	buf.WriteString("--MixedBoundaryString" + crlf)
	m.writeRelated(buf)
	for _, a := range m.Attachments {
		m.writeAttachment(buf, a.Filename, a.ContentType, a.Data)
	}
	buf.WriteString("--MixedBoundaryString--")
}

func (m *Message) writeRelated(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: multipart/related; boundary=\"RelatedBoundaryString\"" + crlf + crlf)
	buf.WriteString("--RelatedBoundaryString" + crlf)
	m.writeAlternative(buf)
	buf.WriteString("--RelatedBoundaryString--" + crlf + crlf)
}

func (m *Message) writeAlternative(buf *bytes.Buffer) {
	buf.WriteString("Content-Type: multipart/alternative; boundary=\"AlternativeBoundaryString\"" + crlf + crlf)
	buf.WriteString("--AlternativeBoundaryString" + crlf)
	buf.WriteString(m.Content.toString())
	buf.WriteString("--AlternativeBoundaryString--" + crlf + crlf)
}

func (m *Message) writeAttachment(buf *bytes.Buffer, filename, contentType string, data []byte) {
	buf.WriteString("--MixedBoundaryString" + crlf)
	buf.WriteString(fmt.Sprintf("Content-Type: %s;name=\"%s\""+crlf, contentType, filename))
	buf.WriteString("Content-Transfer-Encoding: base64" + crlf)
	buf.WriteString(fmt.Sprintf("Content-Disposition: attachment;filename=\"%s\""+crlf+crlf, filename))
	encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encodedData, data)
	buf.Write(encodedData)
	buf.WriteString(crlf)
}
