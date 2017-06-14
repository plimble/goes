package gogoprotobuf

import (
	"errors"
)

type Marshaler interface {
	Marshal() (dAtA []byte, err error)
	Unmarshal(dAtA []byte) error
}

// ProtobufEncoder is a protobuf implementation for EncodedConn
// This encoder will use the builtin protobuf lib to Marshal
// and Unmarshal structs.
type ProtobufEncoder struct {
	// Empty
}

func New() *ProtobufEncoder {
	return &ProtobufEncoder{}
}

var (
	ErrInvalidProtoMsgEncode = errors.New("Invalid gogoprotobuf proto.Message object passed to encode")
	ErrInvalidProtoMsgDecode = errors.New("Invalid gogoprotobuf proto.Message object passed to decode")
)

// Encode
func (pb *ProtobufEncoder) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	i, found := v.(Marshaler)
	if !found {
		return nil, ErrInvalidProtoMsgEncode
	}

	b, err := i.Marshal()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Decode
func (pb *ProtobufEncoder) Decode(data []byte, vPtr interface{}) error {
	if data == nil || len(data) == 0 {
		return nil
	}
	if _, ok := vPtr.(*interface{}); ok {
		return nil
	}
	i, found := vPtr.(Marshaler)
	if !found {
		return ErrInvalidProtoMsgDecode
	}

	err := i.Unmarshal(data)
	if err != nil {
		return err
	}
	return nil
}
