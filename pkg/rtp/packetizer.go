package rtp

type Payloader interface {
	Payload(mtu int, payload []byte) [][]byte
}

type Packetizer interface {
	Packetize(payload []byte) []*Packet
}

type packetizer struct {
	MTU int
	PayloadType uint8
	SSRC int32
	Payloader Payloader
	Sequencer Sequencer
}

func NewPacketizer(mtu int, pt uint8, ssrc int32, payloader Payloader, sequencer Sequencer) Packetizer {
	return &packetizer {
		mtu,
		pt,
		ssrc,
		payloader,
		sequencer,
	}
}

func (p *packetizer) Packetize(payload []byte) []*Packet {
	payloads := p.Payloader.Payload(p.MTU - 12, payload)
	packets := make([]*Packet, len(payloads))

	for i, pp := range payloads {
		packets[i] = &Packet{
			Version: 2,
			Padding: false,
			Extension: false,
			Marker: i == len(payloads) - 1,
			PayloadType: p.PayloadType,
			SequenceNumber: p.Sequencer.NextSequenceNumber(),
			Timestamp: 1, // Figure out how to do timestamps
			Payload: pp,
		}
	}

	return packets
}