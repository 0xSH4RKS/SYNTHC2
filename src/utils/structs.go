package utils

type Command struct {
	Name        string
	Description string
	Help        string
	SubArgs     map[string][]string
}

type Listener struct {
	Name     string
	Address  string
	Port     string
	Uri      string
	Password string
}
