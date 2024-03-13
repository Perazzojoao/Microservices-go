package main

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Passwd      string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}
