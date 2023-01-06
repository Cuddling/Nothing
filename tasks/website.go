package tasks

import "strings"

type Website struct {
	Name string
	Url  string
}

// GetHostName Returns the "host" name of the website's URL
func (w *Website) GetHostName() string {
	name := strings.Replace(w.Url, "http://", "", -1)
	name = strings.Replace(name, "https://", "", -1)
	return name
}
