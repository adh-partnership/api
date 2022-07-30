package discord

type Message struct {
	Username  *string  `json:"username,omitempty"`
	AvatarURL *string  `json:"avatar_URL,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Embeds    *[]Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       *string    `json:"title,omitempty"`
	URL         *string    `json:"URL,omitempty"`
	Description *string    `json:"description,omitempty"`
	Color       *string    `json:"color,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      *[]Field   `json:"fields,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
}

type Author struct {
	Name    *string `json:"name,omitempty"`
	URL     *string `json:"URL,omitempty"`
	IconURL *string `json:"icon_URL,omitempty"`
}

type Field struct {
	Name   *string `json:"name,omitempty"`
	Value  *string `json:"value,omitempty"`
	Inline *bool   `json:"inline,omitempty"`
}

type Thumbnail struct {
	URL *string `json:"URL,omitempty"`
}

type Image struct {
	URL *string `json:"URL,omitempty"`
}

type Footer struct {
	Text    *string `json:"text,omitempty"`
	IconURL *string `json:"icon_URL,omitempty"`
}
