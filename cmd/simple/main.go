package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/xyproto/midi"
)

func frequencyToNoteName(freq float64) string {
	noteNames := [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	midiNumber := int(69 + 12*math.Log2(freq/440))
	return noteNames[midiNumber%12] + strconv.Itoa(midiNumber/12-1)
}

func main() {
	// Create a slice of MIDI notes
	notes := [][]midi.Note{
		{
			{Frequency: midi.NoteToFrequency("C4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("D4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("E4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("F4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("G4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("A4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("B4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("C5"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("B4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("A4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("G4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("F4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("E4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("D4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
			{Frequency: midi.NoteToFrequency("C4"), Duration: time.Millisecond * 300, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
		},
	}

	// Convert the notes to MIDI bytes
	data, err := midi.ConvertToMIDI(notes)
	if err != nil {
		log.Fatal(err)
	}

	for i, track := range notes {
		fmt.Printf("Track %d:\n", i+1)
		for _, note := range track {
			fmt.Printf("Note: %s, Duration: %v, Velocity: %d, Channel: %d, Instrument: %d\n", frequencyToNoteName(note.Frequency), note.Duration, note.Velocity, note.Channel, note.Instrument)
		}
		fmt.Println()
	}

	// Write the MIDI bytes to a file
	err = os.WriteFile("output.mid", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
