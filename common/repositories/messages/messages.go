package messages

// MessageType indicates the type to which a Message belongs to.
type MessageType string

// String returns the string type of MessageType
func (m MessageType) String() string {
	return string(m)
}

// Constants for our message types
const (
	MessageTypeEventIntro   MessageType = "@event/intro"
	MessageTypeEventOutro   MessageType = "@event/outro"
	MessageTypeActionSingle MessageType = "@action/single"
)
