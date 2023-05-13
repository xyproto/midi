package midi

func (m *MIDI) GetProgram(channel uint8) uint8 {
	if program, ok := m.ChannelProgram[channel]; ok {
		return program
	}
	return 0 // Default program number
}

func (m *MIDI) SetProgram(channel, program uint8) {
	m.ChannelProgram[channel] = program
}
