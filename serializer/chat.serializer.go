package serializer

type Message struct {
	Sender    uint   `json:"sender,omitempty"`
	Recipient uint   `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}
