/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package discord

type Message struct {
	Username  *string  `json:"username,omitempty"`
	AvatarURL *string  `json:"avatar_URL,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Embeds    []*Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       *string    `json:"title,omitempty"`
	URL         *string    `json:"URL,omitempty"`
	Description *string    `json:"description,omitempty"`
	Color       *int       `json:"color,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      []*Field   `json:"fields,omitempty"`
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
