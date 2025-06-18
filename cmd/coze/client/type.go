package client

type Conversation struct {
	ID            string            `json:"id"`
	CreatedAt     int               `json:"created_at"`
	MetaData      map[string]string `json:"meta_data,omitempty"`
	LastSectionID string            `json:"last_section_id"`
}

type Chat struct {
	// The ID of the chat.
	ID string `json:"id"`
	// The ID of the conversation.
	ConversationID string `json:"conversation_id"`
	// The ID of the bot.
	BotID string `json:"bot_id"`
	// Indicates the create time of the chat. The value format is Unix timestamp in seconds.
	CreatedAt int `json:"created_at"`
	// Indicates the end time of the chat. The value format is Unix timestamp in seconds.
	CompletedAt int `json:"completed_at,omitempty"`
	// Indicates the failure time of the chat. The value format is Unix timestamp in seconds.
	FailedAt int `json:"failed_at,omitempty"`
	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages.
	MetaData map[string]string `json:"meta_data,omitempty"`
	// When the chat encounters an auth_error, this field returns detailed error information.
	LastError string `json:"last_error,omitempty"`
	// The running status of the session.
	Status string `json:"status"`
}

type StreamChat struct {
	Content string
	Err     error
}

type Message struct {
	// The entity that sent this message.
	Role string `json:"role"`

	// The type of message.
	Type string `json:"type"`

	// The content of the message. It supports various types of content, including plain text,
	// multimodal (a mix of text, images, and files), message cards, and more.
	Content string `json:"content"`

	// The reasoning_content of the thought process message
	ReasoningContent string `json:"reasoning_content"`

	// The type of message content.
	ContentType string `json:"content_type"`

	// Additional information when creating a message, and this additional information will also be
	// returned when retrieving messages. Custom key-value pairs should be specified in Map object
	// format, with a length of 16 key-value pairs. The length of the key should be between 1 and 64
	// characters, and the length of the value should be between 1 and 512 characters.
	MetaData map[string]string `json:"meta_data,omitempty"`

	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`

	// section_id is used to distinguish the context sections of the session history. The same section
	// is one context.
	SectionID string `json:"section_id"`
	BotID     string `json:"bot_id"`
	ChatID    string `json:"chat_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
