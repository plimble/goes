package goes

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Event) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "a":
			z.AggregateID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "b":
			z.AggregateType, err = dc.ReadString()
			if err != nil {
				return
			}
		case "c":
			z.EventID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "d":
			z.EventType, err = dc.ReadString()
			if err != nil {
				return
			}
		case "e":
			z.Revision, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "f":
			z.Time, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "g":
			z.Data, err = dc.ReadBytes(z.Data)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Event) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 7
	// write "a"
	err = en.Append(0x87, 0xa1, 0x61)
	if err != nil {
		return err
	}
	err = en.WriteString(z.AggregateID)
	if err != nil {
		return
	}
	// write "b"
	err = en.Append(0xa1, 0x62)
	if err != nil {
		return err
	}
	err = en.WriteString(z.AggregateType)
	if err != nil {
		return
	}
	// write "c"
	err = en.Append(0xa1, 0x63)
	if err != nil {
		return err
	}
	err = en.WriteString(z.EventID)
	if err != nil {
		return
	}
	// write "d"
	err = en.Append(0xa1, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.EventType)
	if err != nil {
		return
	}
	// write "e"
	err = en.Append(0xa1, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Revision)
	if err != nil {
		return
	}
	// write "f"
	err = en.Append(0xa1, 0x66)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.Time)
	if err != nil {
		return
	}
	// write "g"
	err = en.Append(0xa1, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Event) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 7
	// string "a"
	o = append(o, 0x87, 0xa1, 0x61)
	o = msgp.AppendString(o, z.AggregateID)
	// string "b"
	o = append(o, 0xa1, 0x62)
	o = msgp.AppendString(o, z.AggregateType)
	// string "c"
	o = append(o, 0xa1, 0x63)
	o = msgp.AppendString(o, z.EventID)
	// string "d"
	o = append(o, 0xa1, 0x64)
	o = msgp.AppendString(o, z.EventType)
	// string "e"
	o = append(o, 0xa1, 0x65)
	o = msgp.AppendInt(o, z.Revision)
	// string "f"
	o = append(o, 0xa1, 0x66)
	o = msgp.AppendInt64(o, z.Time)
	// string "g"
	o = append(o, 0xa1, 0x67)
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Event) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "a":
			z.AggregateID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "b":
			z.AggregateType, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "c":
			z.EventID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "d":
			z.EventType, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "e":
			z.Revision, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "f":
			z.Time, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "g":
			z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Event) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.AggregateID) + 2 + msgp.StringPrefixSize + len(z.AggregateType) + 2 + msgp.StringPrefixSize + len(z.EventID) + 2 + msgp.StringPrefixSize + len(z.EventType) + 2 + msgp.IntSize + 2 + msgp.Int64Size + 2 + msgp.BytesPrefixSize + len(z.Data)
	return
}
