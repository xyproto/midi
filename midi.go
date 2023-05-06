package midi

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math"
	"time"
)

type Note struct {
	Frequency  float64
	Duration   time.Duration
	Velocity   byte
	Channel    int
	Instrument int
	Slur       bool
	Start      int
}

type MIDIFile struct {
	formatType uint16
	tracks     uint16
	division   uint16
}

func NewMIDIFile(formatType, tracks, division uint16) *MIDIFile {
	return &MIDIFile{
		formatType: formatType,
		tracks:     tracks,
		division:   division,
	}
}

func FrequencyToMidi(freq float64) (int, int) {
	midiNote := 69 + 12*math.Log2(freq/440.0)
	midiNoteRounded := int(math.Round(midiNote))

	pitchBend := int(math.Round(8192 * (midiNote - float64(midiNoteRounded))))
	pitchBend = pitchBend + 8192 // Center value for pitch bend

	return midiNoteRounded, pitchBend
}

func NoteToFrequency(note string) float64 {
	noteNameToPitch := map[string]float64{
		"C":  0,
		"C#": 1,
		"Db": 1,
		"D":  2,
		"D#": 3,
		"Eb": 3,
		"E":  4,
		"F":  5,
		"F#": 6,
		"Gb": 6,
		"G":  7,
		"G#": 8,
		"Ab": 8,
		"A":  9,
		"A#": 10,
		"Bb": 10,
		"B":  11,
		"H":  11,
	}

	if len(note) < 2 || len(note) > 3 {
		return 0
	}

	noteName := note[:len(note)-1]
	pitch, valid := noteNameToPitch[noteName]
	if !valid {
		return 0
	}

	octave := float64(note[len(note)-1] - '0')
	midiNumber := 12*(octave+1) + pitch
	return 440 * math.Pow(2, (midiNumber-69)/12)
}

func ConvertToMIDI(tracks [][]Note) ([]byte, error) {
	buf := new(bytes.Buffer)

	midi := NewMIDIFile(1, uint16(len(tracks)), 96)
	index, err := midi.WriteHeader(buf)
	if err != nil {
		return nil, err
	}

	for _, track := range tracks {
		trackData, err := ConvertToMIDITrack(track)
		if err != nil {
			return nil, err
		}
		index = midi.WriteTrack(buf, index, trackData)
	}

	return buf.Bytes(), nil
}

func (m *MIDIFile) WriteHeader(w io.Writer) (int, error) {
	index := 0
	n, err := w.Write([]byte("MThd"))
	if err != nil {
		return index, err
	}
	index += n

	err = binary.Write(w, binary.BigEndian, uint32(6))
	if err != nil {
		return index, err
	}
	index += 4

	err = binary.Write(w, binary.BigEndian, m.formatType)
	if err != nil {
		return index, err
	}
	index += 2

	err = binary.Write(w, binary.BigEndian, m.tracks)
	if err != nil {
		return index, err
	}
	index += 2

	err = binary.Write(w, binary.BigEndian, m.division)
	if err != nil {
		return index, err
	}
	index += 2

	return index, nil
}

func (m *MIDIFile) WriteTrack(w io.Writer, index int, events []byte) int {
	n, err := w.Write([]byte("MTrk"))
	if err != nil {
		log.Fatalln(err)
	}
	index += n

	err = binary.Write(w, binary.BigEndian, uint32(len(events)))
	if err != nil {
		log.Fatalln(err)
	}
	index += 4

	n, err = w.Write(events)
	if err != nil {
		log.Fatalln(err)
	}
	index += n

	return index
}

func ConvertToMIDITrack(notes []Note) ([]byte, error) {
	buf := new(bytes.Buffer)
	lastNoteOffTime := 0
	for _, note := range notes {
		midiNote, pitchBend := FrequencyToMidi(note.Frequency)

		// Write program change
		WriteProgramChange(buf, note.Channel, note.Instrument)

		if note.Slur {
			WritePitchBend(buf, note.Channel, pitchBend)
		}

		deltaTime := note.Start - lastNoteOffTime
		writeDeltaTime(buf, deltaTime)
		WriteNoteOn(buf, note.Channel, midiNote, note.Velocity)

		noteOffTicks := int(note.Duration.Seconds() * 96)
		writeDeltaTime(buf, noteOffTicks)
		WriteNoteOff(buf, note.Channel, midiNote, 0)

		lastNoteOffTime = note.Start + noteOffTicks
	}
	return buf.Bytes(), nil
}

func WritePitchBend(w io.Writer, channel int, pitchBend int) {
	pitchBendLSB := pitchBend & 0x7F
	pitchBendMSB := (pitchBend >> 7) & 0x7F
	w.Write([]byte{0x00, byte(0xE0 + channel - 1), byte(pitchBendLSB), byte(pitchBendMSB)})
}

func WriteNoteOn(w io.Writer, channel int, midiNote int, velocity byte) {
	w.Write([]byte{0x00, byte(0x90 + channel - 1), byte(midiNote), velocity})
}

func WriteNoteOff(w io.Writer, channel int, midiNote int, ticks int) {
	w.Write([]byte{byte(ticks), byte(0x80 + channel - 1), byte(midiNote), 0x00})
}

func WriteProgramChange(w io.Writer, channel int, program int) {
	w.Write([]byte{0x00, byte(0xC0 + channel - 1), byte(program)})
}

func writeDeltaTime(w io.Writer, deltaTime int) error {
	vlq := deltaTimeToVLQ(deltaTime)
	return writeBytes(w, vlq)
}

func deltaTimeToVLQ(deltaTime int) []byte {
	var buf []byte
	current := deltaTime & 0x7F

	for deltaTime >>= 7; deltaTime > 0; deltaTime >>= 7 {
		buf = append([]byte{byte(current | 0x80)}, buf...)
		current = deltaTime & 0x7F
	}
	buf = append([]byte{byte(current)}, buf...)
	return buf
}

func writeBytes(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
