package midi

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// MIDI event type constants
const (
	EventNoteOff       = 0x80
	EventNoteOn        = 0x90
	EventProgramChange = 0xC0

	DefaultNoteDuration   = 500 * time.Millisecond
	DefaultNoteVelocity   = 64
	DefaultNoteChannel    = 0
	DefaultNoteInstrument = 0
	DefaultNoteProgram    = 0
)

// MIDI represents a MIDI file or a sequence of MIDI events
type MIDI struct {
	Format         uint16
	Division       uint16
	BPM            float64
	Tracks         []*Track
	ChannelProgram map[uint8]uint8
}

// Track represents a track in a MIDI file or a sequence of MIDI events
type Track struct {
	NoteMap map[time.Duration][]*Note
	Events  []*Event
}

// Event represents a MIDI event
type Event struct {
	DeltaTime uint32
	Type      uint8
	Channel   uint8
	Program   uint8
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
	StartDelay time.Duration
	Program    uint8
}

// NewMIDI creates a new MIDI file or sequence of MIDI events
func NewMIDI(format, division uint16, bpm float64) *MIDI {
	return &MIDI{
		Format:         format,
		Division:       division,
		BPM:            bpm,
		Tracks:         make([]*Track, 0),
		ChannelProgram: make(map[uint8]uint8),
	}
}

// NewTrack creates a new Track
func NewTrack() *Track {
	return &Track{
		NoteMap: NewNoteMap(),
		Events:  make([]*Event, 0),
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

func (m *MIDI) AddNote(t *Track, note *Note) {
	// Convert frequency to MIDI note
	midiNote, _ := FrequencyToMidi(note.Frequency)

	// Convert the note start pause and duration to ticks
	startPauseTicks := m.DurationToTicks(note.StartDelay)
	durationTicks := m.DurationToTicks(note.Duration)

	// Check if program change is needed
	currentProgram := m.GetProgram(note.Channel)
	if note.Program != currentProgram {
		// Create program change event
		programChange := &Event{
			DeltaTime: startPauseTicks,
			Type:      EventProgramChange,
			Channel:   note.Channel,
			Program:   note.Program,
			Data:      []byte{note.Program},
		}
		t.AddEvent(programChange)
		m.SetProgram(note.Channel, note.Program)
	}

	// Create "note on" event
	noteOn := &Event{
		DeltaTime: startPauseTicks,
		Type:      EventNoteOn,
		Channel:   note.Channel,
		Program:   note.Program,
		Data:      []byte{midiNote, note.Velocity},
	}

	// Create "note off" event
	noteOff := &Event{
		DeltaTime: durationTicks,
		Type:      EventNoteOff,
		Channel:   note.Channel,
		Program:   note.Program,
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

// NewNoteMap creates a new map for storing notes by their start time
func NewNoteMap() map[time.Duration][]*Note {
	return make(map[time.Duration][]*Note)
}

// AddNoteToMap adds a note to a track's note map by its start time
func (t *Track) AddNoteToMap(startTime time.Duration, note *Note) {
	t.NoteMap[startTime] = append(t.NoteMap[startTime], note)
}

// AddNotesFromMap adds notes to a track from its note map
func (t *Track) AddNotesFromMap(m *MIDI) {
	m.AddNotesFromMap(t, t.NoteMap)
}

// AddNotesFromMap adds notes to a track from a map by their start time
func (m *MIDI) AddNotesFromMap(t *Track, noteMap map[time.Duration][]*Note) {
	// Convert map to a list of note start times and sort it
	var startTimes []time.Duration
	for startTime := range noteMap {
		startTimes = append(startTimes, startTime)
	}
	sort.Slice(startTimes, func(i, j int) bool {
		return startTimes[i] < startTimes[j]
	})

	// Add notes to track in order of start time
	var lastStartTime time.Duration
	for _, startTime := range startTimes {
		notes := noteMap[startTime]
		for _, note := range notes {
			// Adjust the note's start delay to be relative to the last note
			note.StartDelay = startTime - lastStartTime
			m.AddNote(t, note)
			lastStartTime = startTime
		}
	}
}

func (m *MIDI) AddNoteFromNoteString(t *Track, noteString string, startDelay time.Duration) error {
	note := &Note{
		StartDelay: startDelay,
		Duration:   DefaultNoteDuration,
		Velocity:   DefaultNoteVelocity,
		Channel:    DefaultNoteChannel,
		Instrument: DefaultNoteInstrument,
		Program:    DefaultNoteProgram,
	}

	parts := strings.SplitN(noteString, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid note string format")
	}

	noteName := parts[0]
	note.Frequency = NoteNameToFrequency(noteName)

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return err
	}
	note.Duration = duration

	t.AddNoteToMap(startDelay, note)
	return nil
}
