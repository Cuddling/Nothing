package profiles

import "strings"

type Address struct {
	Name     string
	Email    string
	Phone    string
	Line1    string
	Line2    string
	PostCode string
	City     string
	Country  string
	State    string
}

// GetFirstName Returns the first name of the address receiver
func (a *Address) GetFirstName() string {
	return strings.Split(a.Name, " ")[0]
}

// GetLastName Returns the last name of the address receiver
func (a *Address) GetLastName() string {
	ss := strings.Split(a.Name, " ")
	return ss[len(ss)-1]
}
