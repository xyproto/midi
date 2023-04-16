package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

func PrintMIDI(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading MIDI file:", err)
		return
	}

	pos := 0

	// Read MIDI header
	header := string(data[pos : pos+4])
	pos += 4
	if header != "MThd" {
		fmt.Println("Invalid MIDI header")
		return
	}

	headerLength := binary.BigEndian.Uint32(data[pos : pos+4])
	pos += 4
	formatType := binary.BigEndian.Uint16(data[pos : pos+2])
	pos += 2
	numTracks := binary.BigEndian.Uint16(data[pos : pos+2])
	pos += 2
	division := binary.BigEndian.Uint16(data[pos : pos+2])
	pos += 2

	fmt.Printf("Header length: %d\nFormat type: %d\nNumber of tracks: %d\nDivision: %d\n", headerLength, formatType, numTracks, division)

	// Read tracks
	for trackIndex := 0; trackIndex < int(numTracks); trackIndex++ {
		trackHeader := string(data[pos : pos+4])
		pos += 4
		if trackHeader != "MTrk" {
			fmt.Println("Invalid track header")
			return
		}

		trackLength := binary.BigEndian.Uint32(data[pos : pos+4])
		pos += 4

		fmt.Printf("\nTrack %d length: %d\n", trackIndex+1, trackLength)

		trackEnd := pos + int(trackLength)
		for pos < trackEnd {
			// Read delta-time
			deltaTime := 0
			for {
				byteValue := data[pos]
				pos++
				deltaTime = (deltaTime << 7) | int(byteValue&0x7F)
				if byteValue&0x80 == 0 {
					break
				}
			}

			eventType := data[pos]
			pos++

			// Meta event
			if eventType == 0xFF {
				metaType := data[pos]
				pos++
				metaLength := int(data[pos])
				pos++

				fmt.Printf("DeltaTime: %d, Meta Event: %X, Length: %d\n", deltaTime, metaType, metaLength)
				pos += metaLength
			} else {
				// MIDI event
				data1 := data[pos]
				pos++

				if (eventType&0xF0)>>4 != 12 && (eventType&0xF0)>>4 != 13 {
					data2 := data[pos]
					pos++
					fmt.Printf("DeltaTime: %d, Event: %X, Data1: %X, Data2: %X\n", deltaTime, eventType, data1, data2)
				} else {
					fmt.Printf("DeltaTime: %d, Event: %X, Data1: %X\n", deltaTime, eventType, data1)
				}
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pm <MIDI file>")
		return
	}

	filename := os.Args[1]
	PrintMIDI(filename)
}
