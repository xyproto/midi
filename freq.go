package midi

import "math"

func FrequencyToMidi(frequency float64) (note uint8, bend int) {
	const A4 = 440.0
	const A4MidiNote = 69

	midi := math.Log2(frequency/A4)*12 + A4MidiNote

	// Calculate the pitch bend value based on the fraction part of the midi number
	bend = int(math.Round((midi - math.Floor(midi)) * 8192)) // 8192 is the range for pitch bend (-8192 to 8191)

	// Round up if the fractional part is 0.5 or more
	midi = math.Floor(midi + 0.5)

	// Limit the midi note to valid range
	if midi < 0 {
		midi = 0
	} else if midi > 127 {
		midi = 127
	}

	return uint8(midi), bend
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
