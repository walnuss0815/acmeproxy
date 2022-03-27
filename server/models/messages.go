package models

type MessageDefault struct {
	FQDN  string `json:"fqdn"`
	Value string `json:"value"`
}

type MessageRaw struct {
	Domain  string `json:"domain"`
	Token   string `json:"token"`
	KeyAuth string `json:"keyAuth"`
}

type MessageIncoming struct {
	MessageDefault
	MessageRaw
}
