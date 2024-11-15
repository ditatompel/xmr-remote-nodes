package views

type link struct {
	Text string
	URI  string
}

// Community reference links that are displayed on the home page
var communityLinks = []link{
	{Text: "moneroworld.com", URI: "https://moneroworld.com"},
	{Text: "monero.how", URI: "https://www.monero.how"},
	{Text: "monero.observer", URI: "https://www.monero.observer"},
	{Text: "revuo-xmr.com", URI: "https://revuo-xmr.com"},
	{Text: "themonoeromoon.com", URI: "https://www.themoneromoon.com"},
	{Text: "monerotopia.com", URI: "https://monerotopia.com"},
	{Text: "sethforprivacy.com", URI: "https://sethforprivacy.com"},
}

type nodeStatus struct {
	Code int
	Text string
}

// nodeStatuses is a list of status and their text representation in the UI
//
// The "Status" filter select option in the UI is populated from this list.
var nodeStatuses = []nodeStatus{
	{-1, "ANY"},
	{1, "Online"},
	{0, "Offline"},
}

// refreshIntevals, nettypes, and protocols are used to populate the refresh
// interval, Monero network types, and protocols filter select options in the
// UI
var (
	refreshIntevals = []string{"5s", "10s", "30s", "1m"}
	nettypes        = []string{"mainnet", "stagenet", "testnet"}
	protocols       = []string{"tor", "i2p", "http", "https"}
)
