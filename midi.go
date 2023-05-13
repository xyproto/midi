package midi

import (
	"io"
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
	BPM      float64
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
	StartPause time.Duration
}

// NewMIDI creates a new MIDI file or sequence of MIDI events
func NewMIDI(format, division uint16, bpm float64) *MIDI {
	return &MIDI{
		Format:   format,
		Division: division,
		BPM:      bpm,
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
func (m *MIDI) AddNote(t *Track, note *Note) {
	// Convert frequency to MIDI note
	midiNote, _ := FrequencyToMidi(note.Frequency)

	// Convert the note start pause and duration to ticks
	startPauseTicks := m.DurationToTicks(note.StartPause)
	durationTicks := m.DurationToTicks(note.Duration)

	// Create "note on" event
	noteOn := &Event{
		DeltaTime: startPauseTicks,
		Type:      EventNoteOn,
		Channel:   note.Channel,
		Data:      []byte{midiNote, note.Velocity},
	}

	// Create "note off" event
	noteOff := &Event{
		DeltaTime: durationTicks,
		Type:      EventNoteOff,
		Channel:   note.Channel,
		Data:      []byte{midiNote, 0}, // Velocity is 0
	}

	// Add the events to the track
	t.AddEvent(noteOn)
	t.AddEvent(noteOff)
}

// Size returns the byte size of an Event
func (e *Event) Size() int {
	return 1 + len(e.Data) // 1 byte for the event type, plus the size of the data
}

// Write writes the MIDI data to an io.Writer.
func (m *MIDI) Write(w io.Writer) error {
	// Write MIDI header
	if err := writeMIDIHeader(w, m); err != nil {
		return err
	}

	// Write each track
	for _, track := range m.Tracks {
		if err := writeTrack(w, track); err != nil {
			return err
		}
	}

	return nil
}

// DurationToTicks converts a time duration to the number of ticks
func (m *MIDI) DurationToTicks(d time.Duration) uint32 {
	ticksPerBeat := float64(m.Division)
	ticksPerSecond := ticksPerBeat * m.BPM / 60.0
	return uint32(d.Seconds() * ticksPerSecond)
}

// TicksToDuration converts the number of ticks to a time duration
func (m *MIDI) TicksToDuration(ticks uint32) time.Duration {
	ticksPerBeat := float64(m.Division)
	ticksPerSecond := ticksPerBeat * m.BPM / 60.0
	return time.Duration(float64(ticks) / ticksPerSecond * float64(time.Second))
}
