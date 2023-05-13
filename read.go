package midi

import (
	"fmt"
	"io"
)

func ReadMIDI(r io.Reader) (*MIDI, error) {
	m := &MIDI{}

	// Read MIDI header
	numTracks, err := readMIDIHeader(r, m)
	if err != nil {
		return nil, err
	}

	// Read each track
	for i := 0; i < numTracks; i++ {
		track, err := readTrack(r)
		if err != nil {
			return nil, err
		}
		m.Tracks = append(m.Tracks, track)
	}

	return m, nil
}

func readMIDIHeader(r io.Reader, m *MIDI) (int, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return 0, err
	}

	if string(header[:]) != "MThd" {
		return 0, fmt.Errorf("invalid MIDI header: %v", header)
	}

	chunkSize, err := readMIDIUint32(r)
	if err != nil {
		return 0, err
	}

	if chunkSize != 6 {
		return 0, fmt.Errorf("invalid MIDI header chunk size: %v", chunkSize)
	}

	m.Format, err = readMIDIUint16(r)
	if err != nil {
		return 0, err
	}

	numTracks, err := readMIDIUint16(r)
	if err != nil {
		return 0, err
	}

	m.Division, err = readMIDIUint16(r)
	if err != nil {
		return 0, err
	}

	return int(numTracks), nil
}

func readTrack(r io.Reader) (*Track, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}

	if string(header[:]) != "MTrk" {
		return nil, fmt.Errorf("invalid track header: %v", header)
	}

	trackSize, err := readMIDIUint32(r)
	if err != nil {
		return nil, err
	}

	track := &Track{}
	for i := uint32(0); i < trackSize; {
		event, bytesRead, err := readEvent(r)
		if err != nil {
			return nil, err
		}
		track.Events = append(track.Events, event)
		i += bytesRead
	}

	return track, nil
}

func readEvent(r io.Reader) (*Event, uint32, error) {
	deltaTime, err := readVariableLengthQuantity(r)
	if err != nil {
		return nil, 0, err
	}

	typeAndChannel, err := readMIDIUint8(r)
	if err != nil {
		return nil, 0, err
	}

	event := &Event{
		DeltaTime: deltaTime,
		Type:      typeAndChannel & 0xF0,
		Channel:   typeAndChannel & 0x0F,
	}

	if event.Type == SystemExclusive {
		length, err := readVariableLengthQuantity(r)
		if err != nil {
			return nil, 0, err
		}
		event.Data = make([]byte, length)
		if _, err := io.ReadFull(r, event.Data); err != nil {
			return nil, 0, err
		}
	} else {
		event.Data = make([]byte, 2)
		if _, err := io.ReadFull(r, event.Data); err != nil {
			return nil, 0, err
		}
	}

	bytesRead := uint32(1 + len(event.Data))
	if event.Type == SystemExclusive {
		bytesRead += uint32(len(event.Data))
	} else {
		bytesRead += 2
	}

	return event, bytesRead, nil
}
