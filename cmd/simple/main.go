package main

import (
	"io/ioutil"
	"log"

	"github.com/xyproto/midi"
)

func main() {
	// Create a slice of MIDI notes
	notes := []midi.MidiNote{
		{440, 100, 127, 1, 0, false},
		{523.25, 100, 127, 1, 0, false},
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
