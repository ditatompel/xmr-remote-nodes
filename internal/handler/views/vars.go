package views

type link struct {
	Text string
	URI  string
}

var communityLinks = []link{
	{Text: "moneroworld.com", URI: "https://moneroworld.com"},
	{Text: "monero.how", URI: "https://www.monero.how"},
	{Text: "monero.observer", URI: "https://www.monero.observer"},
	{Text: "revuo-xmr.com", URI: "https://revuo-xmr.com"},
	{Text: "themonoeromoon.com", URI: "https://www.themoneromoon.com"},
	{Text: "monerotopia.com", URI: "https://monerotopia.com"},
	{Text: "sethforprivacy.com", URI: "https://sethforprivacy.com"},
}

var refreshIntevals = []string{"5s", "10s", "30s", "1m"}
