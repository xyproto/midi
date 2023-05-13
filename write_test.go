// write_test.go

package midi

import (
	"os"
	"testing"
)

func TestWriteMIDI(t *testing.T) {
	m := &MIDI{
		Tracks: []*Track{
			{
				Events: []*Event{
					{
						DeltaTime: 0,
						Type:      EventNoteOn,
						Data:      []byte{0x40, 0x60}, // Note 64, velocity 96
					},
					{
						DeltaTime: 480, // 1 quarter note later
						Type:      EventNoteOff,
						Data:      []byte{0x40, 0x00}, // Note 64, velocity 0
					},
				},
			},
		},
	}

	f, err := os.Create("test.mid")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	err = WriteMIDI(f, m)
	if err != nil {
		t.Fatalf("Failed to write MIDI: %v", err)
	}
}
