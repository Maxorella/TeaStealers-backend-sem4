package models

type ModuleCreate struct {
	Title string `json:"title,omitempty"`
	ID    int    `json:"id"`
}

type ModuleList struct {
	Modules []ModuleCreate `json:"modules"`
}
