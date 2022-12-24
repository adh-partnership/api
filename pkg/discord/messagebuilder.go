package discord

import "strconv"

var (
	DefaultUsername  = "Web API"
	DefaultAvatarURL = "https://avatars.githubusercontent.com/u/115051386?s=200&v=4"
)

func NewMessage() *Message {
	return &Message{
		Username:  &DefaultUsername,
		AvatarURL: &DefaultAvatarURL,
	}
}

func (m *Message) SetUsername(username string) *Message {
	m.Username = &username

	return m
}

func (m *Message) SetAvatarURL(url string) *Message {
	m.AvatarURL = &url

	return m
}

func (m *Message) SetContent(content string) *Message {
	m.Content = &content

	return m
}

func (m *Message) AddEmbed(e *Embed) *Message {
	if m.Embeds == nil {
		m.Embeds = []*Embed{}
	}
	m.Embeds = append(m.Embeds, e)

	return m
}

func (m *Message) Send(webhook string) error {
	return SendWebhookMessageObj(webhook, *m)
}

func NewEmbed() *Embed {
	c := GetColor("00", "ff", "00")

	return &Embed{
		Color: &c,
	}
}

func (e *Embed) SetTitle(title string) *Embed {
	e.Title = &title

	return e
}

func (e *Embed) SetURL(url string) *Embed {
	e.URL = &url

	return e
}

func (e *Embed) SetDescription(desc string) *Embed {
	e.Description = &desc

	return e
}

func (e *Embed) SetColor(color int) *Embed {
	e.Color = &color

	return e
}

func (e *Embed) SetAuthor(a *Author) *Embed {
	e.Author = a

	return e
}

func (e *Embed) AddField(f *Field) *Embed {
	if e.Fields == nil {
		e.Fields = []*Field{}
	}
	e.Fields = append(e.Fields, f)

	return e
}

func (e *Embed) SetThumbnail(t *Thumbnail) *Embed {
	e.Thumbnail = t

	return e
}

func (e *Embed) SetImage(i *Image) *Embed {
	e.Image = i

	return e
}

func (e *Embed) SetFooter(f *Footer) *Embed {
	e.Footer = f

	return e
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) SetName(name string) *Field {
	f.Name = &name

	return f
}

func (f *Field) SetValue(value string) *Field {
	f.Value = &value

	return f
}

func (f *Field) SetInline(inline bool) *Field {
	f.Inline = &inline

	return f
}

// Convert RGB hex values to integer
func GetColor(red, green, blue string) int {
	r, _ := strconv.ParseInt(red, 16, 64)
	g, _ := strconv.ParseInt(green, 16, 64)
	b, _ := strconv.ParseInt(blue, 16, 64)

	return int(r<<16 | g<<8 | b)
}
