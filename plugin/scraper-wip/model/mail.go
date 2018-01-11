package model

// Mail is the container of a single e-mail
type Mail struct {
	Title   string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Link    string `json:"link,omitempty" yaml:"link,omitempty" toml:"link,omitempty"`
	Author  string `json:"author,omitempty" yaml:"author,omitempty" toml:"author,omitempty"`
	Date    string `json:"date,omitempty" yaml:"date,omitempty" toml:"date,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty" toml:"message,omitempty"`
}
