package request

type (
	// doc: https://docs.wati.io/reference/webhooks
	WebhookRequest struct {
		ID             string    `json:"id"`
		Created        string    `json:"created"`
		ConversationID string    `json:"conversationId"`
		TicketID       string    `json:"ticketId"`
		Text           string    `json:"text"`
		Type           string    `json:"type"`
		Data           string    `json:"data,omitempty"`
		Timestamp      string    `json:"timestamp"`
		Owner          bool      `json:"owner"`
		EventType      string    `json:"eventType"`
		StatusString   string    `json:"statusString,omitempty"`
		AvatarURL      string    `json:"avatarUrl,omitempty"`
		AssignedID     string    `json:"assignedId,omitempty"`
		OperatorName   string    `json:"operatorName,omitempty"`
		OperatorEmail  string    `json:"operatorEmail,omitempty"`
		WaID           string    `json:"waId"`
		SenderName     string    `json:"senderName"`
		ListReply      ListReply `json:"listReply"`
		ReplyContextID string    `json:"replyContextId"`
	}

	ListReply struct {
		Title      string `json:"title"`
		Desription string `json:"desription"`
		ID         string `json:"id"`
	}
)
