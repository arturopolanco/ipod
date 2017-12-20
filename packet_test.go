package ipod_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/oandrew/ipod"
)

type ShortWriter struct {
	max     int
	written int
}

func (sw *ShortWriter) Write(p []byte) (int, error) {
	if sw.written+len(p) > sw.max {
		return 0, io.ErrShortWrite
	}
	sw.written += len(p)
	return len(p), nil
}

func NewShortWriter(max int) io.Writer {
	return &ShortWriter{
		max: max,
	}
}

// func TestBinWritePanic(t *testing.T) {
// 	buf := NewShortWriter(1)
// 	err := ipod.MarshalSmallPacket(buf, &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{0x03, 0x04}})
// 	if err == nil {
// 		t.Error("error == nil")
// 	}
// 	t.Logf("binWrite err = %v", err)
// }

func TestMarshalRawPacket(t *testing.T) {
	largeData := bytes.Repeat([]byte{0xee}, 255)

	tests := []struct {
		name    string
		pkt     *ipod.Packet
		want    []byte
		wantErr bool
	}{
		//{"nil-data", &ipod.Packet{ipod.NewLingoCmdID(0x01, 0x02), nil, nil},
		//	[]byte{0x01, 0x02}, false},
		{"empty-data", &ipod.Packet{ipod.NewLingoCmdID(0x01, 0x02), nil, []byte{}},
			[]byte{0x01, 0x02}, false},
		{"small-data", &ipod.Packet{ipod.NewLingoCmdID(0x01, 0x02), nil, []byte{0xfd}},
			[]byte{0x01, 0x02, 0xfd}, false},
		{"large-data", &ipod.Packet{ipod.NewLingoCmdID(0x01, 0x02), nil, largeData},
			append([]byte{0x1, 0x02}, largeData...), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipod.MarshalPacket(tt.pkt)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalSmallPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual([]byte(got), tt.want) {
				t.Errorf("MarshalSmallPacket() = [% 02x], want [% 02x]", got, tt.want)
			}
		})
	}
}

// func TestUnmarshalRawPacket(t *testing.T) {
// 	largeData := bytes.Repeat([]byte{0xee}, 255)

// 	tests := []struct {
// 		name    string
// 		raw     []byte
// 		want    *ipod.RawPacket
// 		wantErr bool
// 	}{
// 		{"no-data", []byte{0x55, 0x02, 0x01, 0x02, 256 - 0x05}, &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{}}, false},
// 		{"with-data", []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0xfd}, &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), []byte{0xfd}}, false},
// 		{"bad-crc", []byte{0x55, 0x03, 0x01, 0x02, 0xfd, 0x22}, nil, true},
// 		{"wrong-start-byte", []byte{0xff, 0x03, 0x01, 0x02, 0xfd, 0xfd}, nil, true},
// 		{"too-small-payload-length", []byte{0xff, 0x01, 0x01, 0x02, 0xfd, 0xfd}, nil, true},
// 		{"too-large-payload-length", []byte{0xff, 0x33, 0x01, 0x02, 0xfd, 0xfd}, nil, true},
// 		{"no-len-marker", []byte{0x55, 0xdd, 0x01, 0x01, 0x1, 0x02}, nil, true},
// 		{"bad-crc", append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0x22)...), nil, true},
// 		{"with-data", append([]byte{0x55, 0x00, 0x01, 0x01, 0x1, 0x02}, append(largeData, 0xe9)...), &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), largeData}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ipod.UnmarshalRawPacket(tt.raw)
// 			t.Logf("%+v", got)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UnmarshalSmallPacket() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("UnmarshalSmallPacket() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestMarshalPacket(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		pkt         *ipod.Packet
// 		wantLenByte byte
// 		wantErr     bool
// 	}{
// 		{"small", &ipod.Packet{ipod.NewLingoCmdID(0x00, 0x02), nil, general.ACK{0xfd, 0x00}}, 0x04, false},
// 		//{"large", &ipod.RawPacket{ipod.NewLingoCmdID(0x01, 0x02), make([]byte, 254)}, 0x00, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := bytes.Buffer{}
// 			err := ipod.MarshalPacket(&got, tt.pkt)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TestMarshalPacket() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			lenByte := got.Bytes()[1]
// 			if lenByte != tt.wantLenByte {
// 				t.Errorf("TestMarshalPacket() = %02x, want %02x", lenByte, tt.wantLenByte)
// 			}
// 		})
// 	}
// }

// func TestUnmarshalPacket(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		raw     []byte
// 		want    *ipod.Packet
// 		wantErr bool
// 	}{

// 		{"small", []byte{0x55, 0x04, 0x00, 0x02, 0xfd, 0x00, 0xfd}, &ipod.Packet{ipod.NewLingoCmdID(0x00, 0x02), nil, general.ACK{0xfd, 0x00}}, false},
// 		{"large", []byte{0x55, 0x00, 0x00, 0x04, 0x00, 0x02, 0xfd, 0x00, 0xfd}, &ipod.Packet{ipod.NewLingoCmdID(0x00, 0x02), nil, general.ACK{0xfd, 0x00}}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := ipod.Packet{}
// 			err := ipod.UnmarshalPacket(bytes.NewReader(tt.raw), &got)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UnmarshalPacket() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(&got, tt.want) {
// 				t.Errorf("UnmarshalPacket() = %v, want %v", &got, tt.want)
// 			}
// 		})
// 	}
// }

var shortFrameB = []byte{0x14, 0x00, 0x55, 0x88, 0x00, 0x15, 0x00, 0x06, 0x02, 0x00, 0x00, 0x07,
	0x30, 0x82, 0x03, 0xad, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d,
	0x01, 0x07, 0x02, 0xa0, 0x82, 0x03, 0x9e, 0x30, 0x82, 0x03, 0x9a, 0x02,
	0x01, 0x01, 0x31, 0x00, 0x30, 0x0b, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86,
	0xf7, 0x0d, 0x01, 0x07, 0x01, 0xa0, 0x82, 0x03, 0x80, 0x30, 0x82, 0x03,
	0x7c, 0x30, 0x82, 0x02, 0x64, 0xa0, 0x03, 0x02, 0x01, 0x02, 0x02, 0x0f,
	0x34, 0x34, 0xaa, 0x11, 0x10, 0x09, 0xaa, 0x06, 0xaa, 0x00, 0x02, 0xaa,
	0x05, 0x70, 0x41, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7,
	0x0d, 0x01, 0x01, 0x05, 0x05, 0x00, 0x30, 0x81, 0x92, 0x31, 0x0b, 0x30,
	0x09, 0x06, 0x03, 0x55, 0x04, 0x06, 0x13, 0x02, 0x55, 0x53, 0x31, 0x1d,
	0x30, 0x1b, 0x06, 0x03, 0x55, 0x04, 0x0a, 0x13, 0x14, 0x41, 0x70, 0x70,
	0x6c, 0x65, 0x20, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0xa3, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00}

func BenchmarkUnmarshalRawPacket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := ipod.NewPacketReader(shortFrameB)
		r.ReadPacket()
	}
}

// func BenchmarkMarshalRawPacket(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		ipod.MarshalRawPacket(&ipod.RawPacket{ID: ipod.NewLingoCmdID(0x00, 0x15), Data: shortFrameB[6 : 6+134]})
// 	}
// }

// func TestPacketPrint(t *testing.T) {
// 	p := ipod.Packet{ID: ipod.NewLingoCmdID(0x03, 0x0002), Payload: []byte{0xff}}

// 	t.Logf("val: string (%%s): %s", p)
// 	t.Logf("ptr: string (%%s): %s", &p)
// 	t.Logf("val: value (%%v): %v", p)
// 	t.Logf("ptr: value (%%v): %v", &p)
// 	t.Logf("val: value+fields (%%+v): %+v", p)
// 	t.Logf("ptr: value+fields (%%+v): %+v", &p)
// 	t.Logf("val: go-syntax (%%#v): %#v", p)
// 	t.Logf("ptr: go-syntax (%%#v): %#v", &p)

// }
