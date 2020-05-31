package main

const (
	H_SIGNATURE = 0xF7
)

type ProtocolHeader struct {
	Signature byte
	Length    uint16
	Checksum  uint16
	Payload   []byte
}

func UnmarshalProtocolHeader(data []byte) (*ProtocolHeader, error) {
	//
	return nil, nil
}

func (udp *ProtocolHeader) Marshal() ([]byte, error) {
	return nil, nil
}
