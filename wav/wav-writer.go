// Package wav provides functionality for writing WAV audio files
package wav

import (
	"encoding/binary"
	"os"
)

// Writer handles the creation and writing of WAV audio files
type Writer struct {
	SampleRate    int
	NumChannels   int
	BitsPerSample int
}

// NewWriter creates a new WAV writer with default settings
func NewWriter(sampleRate int) *Writer {
	return &Writer{
		SampleRate:    sampleRate,
		NumChannels:   1,  // Mono audio
		BitsPerSample: 16, // 16-bit audio
	}
}

// WriteFile writes audio samples to a WAV file
func (w *Writer) WriteFile(filename string, data []float64) error {
	// Create or truncate output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Calculate WAV file parameters
	byteRate := w.SampleRate * w.NumChannels * w.BitsPerSample / 8
	blockAlign := w.NumChannels * w.BitsPerSample / 8
	dataSize := len(data) * w.BitsPerSample / 8
	chunkSize := 36 + dataSize

	// Write RIFF header chunk
	if _, err := file.Write([]byte("RIFF")); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(chunkSize)); err != nil {
		return err
	}
	if _, err := file.Write([]byte("WAVE")); err != nil {
		return err
	}

	// Write format chunk
	if _, err := file.Write([]byte("fmt ")); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(16)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(1)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(w.NumChannels)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(w.SampleRate)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(byteRate)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(blockAlign)); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(w.BitsPerSample)); err != nil {
		return err
	}

	// Write data chunk header
	if _, err := file.Write([]byte("data")); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(dataSize)); err != nil {
		return err
	}

	// Write audio samples
	for _, sample := range data {
		// Clamp samples to prevent overflow
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}

		// Convert float64 (-1.0 to 1.0) to int16 (-32768 to 32767)
		intSample := int16(sample * 32767)
		if err := binary.Write(file, binary.LittleEndian, intSample); err != nil {
			return err
		}
	}

	return nil
}
