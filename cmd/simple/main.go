package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/xyproto/midi"
)

func main() {
	// Create a slice of MIDI notes
	notes := []midi.MidiNote{
		{midi.NoteToFrequency("C4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("D4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("E4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("F4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("G4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("A4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("B4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("C5"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("B4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("A4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("G4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("F4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("E4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("D4"), time.Millisecond * 300, 127, 1, 0, false},
		{midi.NoteToFrequency("C4"), time.Millisecond * 300, 127, 1, 0, false},
	}

	// Convert the notes to MIDI bytes
	tracks := midi.ConvertToMIDITracks(notes)
	data, err := midi.ConvertToMIDI(tracks)
	if err != nil {
		log.Fatal(err)
	}

	// Write the MIDI bytes to a file
	err = ioutil.WriteFile("output.mid", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
