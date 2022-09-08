package muggie

type Message interface {
	ID() string
	Content() string
}

type MsgHandlerFunc func(msg Message) error

type MessageProvider interface {
	OnMsg(MsgHandlerFunc) error
}

type MessageReplier interface {
	ReplyTo(t Message, content string) error
}
