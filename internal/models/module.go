package models

import "html"

type ModuleCreate struct {
	Title string `json:"title,omitempty"`
	ID    int    `json:"id"`
}

type ModuleList struct {
	Modules []ModuleCreate `json:"modules"`
}

func (mc *ModuleCreate) Sanitize() {
	mc.Title = html.EscapeString(mc.Title)
}
