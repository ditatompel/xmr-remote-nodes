package monero

// RuckniumNodeData represents a single remote node data from Rucknium's MRL ban list
type RuckniumNodeData struct {
	ID                uint   `json:"id,omitempty" db:"id"`
	Date              string `json:"date" db:"scan_date"`
	ConnectedNodeIP   string `json:"connected_node_ip" db:"connected_node_ip"`
	IsSpyNode         int    `json:"is_spy_node" db:"is_spy_node"`                   // 0 = no, 1 = yes, 2 = not applied
	MRLBanListEnabled int    `json:"mrl_ban_list_enabled" db:"mrl_ban_list_enabled"` // 0 = no, 1 = yes, 2 = not applied
	DNSBanListEnabled int    `json:"dns_ban_list_enabled" db:"dns_ban_list_enabled"` // 0 = no, 1 = yes, 2 = not applied
}
