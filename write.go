package midi

import (
	"bytes"
	"io"
)

func WriteMIDI(w io.Writer, m *MIDI) error {
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

func writeMIDIHeader(w io.Writer, m *MIDI) error {
	// Chunk type: "MThd"
	if _, err := w.Write([]byte("MThd")); err != nil {
		return err
	}

	// Chunk size: 6 bytes
	if err := writeMIDIUint32(w, 6); err != nil {
		return err
	}

	// Format type
	if err := writeMIDIUint16(w, m.Format); err != nil {
		return err
	}

	// Number of tracks
	numTracks := uint16(len(m.Tracks))
	if err := writeMIDIUint16(w, numTracks); err != nil {
		return err
	}

	// Time division
	if err := writeMIDIUint16(w, m.Division); err != nil {
		return err
	}

	return nil
}

func writeTrack(w io.Writer, t *Track) error {
	// Buffer the track data
	buf := new(bytes.Buffer)

	// Write each event to the buffer
	for _, event := range t.Events {
		if err := writeEvent(buf, event); err != nil {
			return err
		}
	}

	// Chunk type: "MTrk"
	if _, err := w.Write([]byte("MTrk")); err != nil {
		return err
	}

	// Calculate track size and write it
	trackSize := uint32(buf.Len())
	if err := writeMIDIUint32(w, trackSize); err != nil {
		return err
	}

	// Write the buffered track data
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func writeEvent(w io.Writer, e *Event) error {
	// Write delta time
	if err := writeVariableLengthQuantity(w, e.DeltaTime); err != nil {
		return err
	}

	// Write event type
	if err := writeMIDIUint8(w, e.Type); err != nil {
		return err
	}

	// Write event data
	if _, err := w.Write(e.Data); err != nil {
		return err
	}

	return nil
}
