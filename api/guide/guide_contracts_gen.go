package guide

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *RequestMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "previous_hash":
			err = dc.ReadExtension(&z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		case "tour_number":
			z.TourNumber, err = dc.ReadByte()
			if err != nil {
				err = msgp.WrapError(err, "TourNumber")
				return
			}
		case "tour_length":
			z.TourLength, err = dc.ReadByte()
			if err != nil {
				err = msgp.WrapError(err, "TourLength")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z RequestMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "previous_hash"
	err = en.Append(0x83, 0xad, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteExtension(&z.PreviousHash)
	if err != nil {
		err = msgp.WrapError(err, "PreviousHash")
		return
	}
	// write "tour_number"
	err = en.Append(0xab, 0x74, 0x6f, 0x75, 0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	if err != nil {
		return
	}
	err = en.WriteByte(z.TourNumber)
	if err != nil {
		err = msgp.WrapError(err, "TourNumber")
		return
	}
	// write "tour_length"
	err = en.Append(0xab, 0x74, 0x6f, 0x75, 0x72, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68)
	if err != nil {
		return
	}
	err = en.WriteByte(z.TourLength)
	if err != nil {
		err = msgp.WrapError(err, "TourLength")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z RequestMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "previous_hash"
	o = append(o, 0x83, 0xad, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68)
	o, err = msgp.AppendExtension(o, &z.PreviousHash)
	if err != nil {
		err = msgp.WrapError(err, "PreviousHash")
		return
	}
	// string "tour_number"
	o = append(o, 0xab, 0x74, 0x6f, 0x75, 0x72, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72)
	o = msgp.AppendByte(o, z.TourNumber)
	// string "tour_length"
	o = append(o, 0xab, 0x74, 0x6f, 0x75, 0x72, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68)
	o = msgp.AppendByte(o, z.TourLength)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RequestMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "previous_hash":
			bts, err = msgp.ReadExtensionBytes(bts, &z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		case "tour_number":
			z.TourNumber, bts, err = msgp.ReadByteBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "TourNumber")
				return
			}
		case "tour_length":
			z.TourLength, bts, err = msgp.ReadByteBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "TourLength")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z RequestMsg) Msgsize() (s int) {
	s = 1 + 14 + msgp.ExtensionPrefixSize + z.PreviousHash.Len() + 12 + msgp.ByteSize + 12 + msgp.ByteSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ResponseMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "previous_hash":
			err = dc.ReadExtension(&z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		case "hash":
			err = dc.ReadExtension(&z.Hash)
			if err != nil {
				err = msgp.WrapError(err, "Hash")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ResponseMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "previous_hash"
	err = en.Append(0x82, 0xad, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteExtension(&z.PreviousHash)
	if err != nil {
		err = msgp.WrapError(err, "PreviousHash")
		return
	}
	// write "hash"
	err = en.Append(0xa4, 0x68, 0x61, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteExtension(&z.Hash)
	if err != nil {
		err = msgp.WrapError(err, "Hash")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ResponseMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "previous_hash"
	o = append(o, 0x82, 0xad, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68)
	o, err = msgp.AppendExtension(o, &z.PreviousHash)
	if err != nil {
		err = msgp.WrapError(err, "PreviousHash")
		return
	}
	// string "hash"
	o = append(o, 0xa4, 0x68, 0x61, 0x73, 0x68)
	o, err = msgp.AppendExtension(o, &z.Hash)
	if err != nil {
		err = msgp.WrapError(err, "Hash")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ResponseMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "previous_hash":
			bts, err = msgp.ReadExtensionBytes(bts, &z.PreviousHash)
			if err != nil {
				err = msgp.WrapError(err, "PreviousHash")
				return
			}
		case "hash":
			bts, err = msgp.ReadExtensionBytes(bts, &z.Hash)
			if err != nil {
				err = msgp.WrapError(err, "Hash")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ResponseMsg) Msgsize() (s int) {
	s = 1 + 14 + msgp.ExtensionPrefixSize + z.PreviousHash.Len() + 5 + msgp.ExtensionPrefixSize + z.Hash.Len()
	return
}
