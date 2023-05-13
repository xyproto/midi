package midi

import (
	"time"
)

// MIDI event type constants
const (
	EventNoteOff = 0x80
	EventNoteOn  = 0x90
)

// MIDI represents a MIDI file or a sequence of MIDI events
type MIDI struct {
	Format   uint16
	Division uint16
	Tracks   []*Track
}

// Track represents a track in a MIDI file or a sequence of MIDI events
type Track struct {
	Events []*Event
}

// Event represents a MIDI event
type Event struct {
	DeltaTime uint32
	Type      uint8
	Channel   uint8
	Data      []byte
}

// Note represents a musical note in a MIDI track
type Note struct {
	Frequency  float64
	Duration   time.Duration
	Velocity   uint8
	Channel    uint8
	Instrument uint8
	Slur       bool
	Start      uint32
}

// NewMIDI creates a new MIDI file or sequence of MIDI events
func NewMIDI(format, division uint16) *MIDI {
	return &MIDI{
		Format:   format,
		Division: division,
		Tracks:   make([]*Track, 0),
	}
}

// NewTrack creates a new Track
func NewTrack() *Track {
	return &Track{
		Events: make([]*Event, 0),
	}
}

// AddTrack adds a track to a MIDI file
func (m *MIDI) AddTrack(track *Track) {
	m.Tracks = append(m.Tracks, track)
}

// AddEvent adds an event to a Track
func (t *Track) AddEvent(event *Event) {
	t.Events = append(t.Events, event)
}

// AddNote adds a note to a Track as an Event
func (t *Track) AddNote(note *Note) {
	// Convert frequency to MIDI note
	midiNote, _ := FrequencyToMidi(note.Frequency)

	// Create "note on" event
	noteOn := &Event{
		DeltaTime: note.Start,
		Type:      0x90, // Note-on event
		Channel:   note.Channel,
		Data:      []byte{midiNote, note.Velocity},
	}

	// Create "note off" event
	noteOff := &Event{
		DeltaTime: uint32(note.Duration.Milliseconds()),
		Type:      0x90, // Note-on event with velocity 0 = note-off
		Channel:   note.Channel,
		Data:      []byte{midiNote, 0}, // Velocity is 0
	}

	// Add the events to the track
	t.AddEvent(noteOn)
	t.AddEvent(noteOff)
}

func (e *Event) Size() int {
	return 1 + len(e.Data) // 1 byte for the event type, plus the size of the data
}
