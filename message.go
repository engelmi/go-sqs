package gosqs

type OutgoingMessage struct {
	DeduplicationId *string
	GroupId         *string
	Payload         []byte
	Attributes      map[string]string
}

type IncomingMessage struct {
	MessageId              *string
	ReceiptHandle          *string
	Body                   *string
	MD5OfBody              *string
	MessageAttributes      map[string]*IncomingMessageAttributeValue
	MD5OfMessageAttributes *string
	Attributes             map[string]*string
}

type IncomingMessageAttributeValue struct {
	BinaryValue []byte
	DataType    *string
	StringValue *string
}
