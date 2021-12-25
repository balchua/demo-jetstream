package monitoring

type ConsumerInfo struct {
	Stream        string `json:"stream_name"`
	Name          string `json:"name"`
	NumAckPending uint64 `json:"num_ack_pending"`
	NumPending    uint64 `json:"num_pending"`
}

type StreamDetail struct {
	Name      string          `json:"name"`
	Consumers []*ConsumerInfo `json:"consumer_detail,omitempty"`
}

type AccountDetail struct {
	Name    string         `json:"name"`
	Id      string         `json:"id"`
	Streams []StreamDetail `json:"stream_detail,omitempty"`
}

// JSInfo has detailed information on JetStream.
type JSInfo struct {
	ID string `json:"server_id"`

	// aggregate raft info
	AccountDetails []*AccountDetail `json:"account_details,omitempty"`
}
