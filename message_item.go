package smtp

const (
	Html  ContentType = "html"
	Plain ContentType = "plain"
)

type (
	ContentType string
	Content     struct {
		ContentType ContentType
		Body        string
	}
)

func (c *Content) toString() string {
	switch c.ContentType {
	case Html:
		return "Content-Type: text/html; charset=UTF-8" + crlf + crlf + c.Body + crlf + crlf
	default:
		return "Content-Type: text/plain; charset=UTF-8" + crlf + crlf + c.Body + crlf + crlf
	}
}
