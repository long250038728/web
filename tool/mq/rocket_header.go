package mq

//
//const (
//	RocketTypeNORMAL = iota
//	RocketTypeFIFO
//	RocketTypeDELAY
//	RocketTypeTRANSACTION
//)
//
//type MsgHeader struct {
//	MsgType        int       `json:"msg_type"`
//	IsAsync        bool      `json:"is_async"`
//	Tag            string    `json:"tag"`             //设置过滤tag
//	MessageGroup   string    `json:"message_group"`   //fifo
//	DelayTimestamp time.Time `json:"delay_timestamp"` //delay
//}
//
//const (
//	MsgType        = "msg_type"
//	IsAsync        = "is_async"
//	Tag            = "tag"
//	MessageGroup   = "message_group"
//	DelayTimestamp = "delay_timestamp"
//)
//
//func NewMessageHeaderNORMAL() *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeNORMAL}
//}
//
//func NewMessageHeaderFIFO(messageGroup string) *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeFIFO, MessageGroup: messageGroup}
//}
//
//func NewMessageHeaderDELAY(delayTimestamp time.Time) *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeDELAY, DelayTimestamp: delayTimestamp}
//}
//
//func NewMessageHeaderTransaction() *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeTRANSACTION}
//}
//
//func NewAsyncMessageHeaderNORMAL() *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeNORMAL, IsAsync: true}
//}
//
//func NewAsyncMessageHeaderFIFO(messageGroup string) *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeFIFO, MessageGroup: messageGroup, IsAsync: true}
//}
//
//func NewAsyncMessageHeaderDELAY(delayTimestamp time.Time) *MsgHeader {
//	return &MsgHeader{MsgType: RocketTypeDELAY, DelayTimestamp: delayTimestamp, IsAsync: true}
//}
//
//func (header *MsgHeader) SetAsync(isAsync bool) *MsgHeader {
//	header.IsAsync = isAsync
//	return header
//}
//
//func (header *MsgHeader) SetTag(tag string) *MsgHeader {
//	header.Tag = tag
//	return header
//}
//
//func parseHeader(header []Header) (h *MsgHeader, err error) {
//	h = &MsgHeader{}
//	for _, head := range header {
//		switch head.Key {
//		case MsgType:
//			h.MsgType, err = strconv.Atoi(string(head.Value))
//		case IsAsync:
//			h.IsAsync, err = strconv.ParseBool(string(head.Value))
//		case Tag:
//			h.Tag = string(head.Value)
//		case MessageGroup:
//			h.MessageGroup = string(head.Value)
//		case DelayTimestamp:
//			h.DelayTimestamp, err = time.Parse(time.DateTime, string(head.Value))
//		}
//		if err != nil {
//			return nil, err
//		}
//	}
//	return h, nil
//}
//
//func NewRocketHeader(h *MsgHeader) []Header {
//	return []Header{
//		{Key: "msg_type", Value: []byte(strconv.Itoa(h.MsgType))},
//		{Key: "is_async", Value: []byte(strconv.FormatBool(h.IsAsync))},
//		{Key: "tag", Value: []byte(h.Tag)},
//		{Key: "message_group", Value: []byte(h.MessageGroup)},
//		{Key: "delay_timestamp", Value: []byte(h.DelayTimestamp.Format(time.DateTime))},
//	}
//}
