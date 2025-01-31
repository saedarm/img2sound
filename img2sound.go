package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image2sound/wav"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	minDuration  = 500
	maxDuration  = 5000
	minBaseFreq  = 220
	maxBaseFreq  = 880
	minDensity   = 1
	maxDensity   = 20
	sampleRate   = 44100
)

type Slider struct {
	X, Y     float64
	Width    int
	Value    int
	Min, Max int
	Dragging bool
}

type Game struct {
	LengthSlider    *Slider
	FrequencySlider *Slider
	DensitySlider   *Slider
	Image           image.Image
	AudioData       []float64
	Converting      bool
	OutputDir       string
}

// NewSlider creates a new slider control
func NewSlider(x, y float64, width, min, max, initialValue int) *Slider {
	return &Slider{
		X:     x,
		Y:     y,
		Width: width,
		Min:   min,
		Max:   max,
		Value: initialValue,
	}
}

// Helper function to keep values within bounds
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Generate unique filenames for output
func generateUniqueFilename(outputDir string) string {
	timestamp := time.Now().Format("20060102_150405")
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)

	filename := fmt.Sprintf("output_%s_%s.wav", timestamp, randomString)
	if outputDir != "" {
		return filepath.Join(outputDir, filename)
	}
	return filename
}

// Initialize the output directory
func initializeOutputDirectory() string {
	outputDir := "wav_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		return ""
	}
	return outputDir
}

// Generate the basic waveform with harmonics
func generateHarmonicWave(frequency, amplitude, t float64) float64 {
	fundamental := amplitude * math.Sin(2*math.Pi*frequency*t)
	firstHarmonic := 0.5 * amplitude * math.Sin(4*math.Pi*frequency*t)
	secondHarmonic := 0.25 * amplitude * math.Sin(6*math.Pi*frequency*t)
	thirdHarmonic := 0.125 * amplitude * math.Sin(8*math.Pi*frequency*t)

	return fundamental + firstHarmonic + secondHarmonic + thirdHarmonic
}

// Apply ADSR envelope to shape the sound
func applyEnvelope(sample float64, position, duration float64) float64 {
	attackTime := 0.1
	decayTime := 0.2
	sustainLevel := 0.7
	releaseTime := 0.2

	normalized := position / duration

	switch {
	case normalized < attackTime:
		return sample * (normalized / attackTime)
	case normalized < attackTime+decayTime:
		decayPosition := (normalized - attackTime) / decayTime
		return sample * (1.0 - (1.0-sustainLevel)*decayPosition)
	case normalized > (1.0 - releaseTime):
		releasePosition := (normalized - (1.0 - releaseTime)) / releaseTime
		return sample * (sustainLevel * (1.0 - releasePosition))
	default:
		return sample * sustainLevel
	}
}

// Convert image data to musical intensities
func imageToNoteIntensities(img image.Image, row int) []float64 {
	bounds := img.Bounds()
	intensities := make([]float64, bounds.Max.X-bounds.Min.X)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		pixel := img.At(x, row)
		r, g, b, _ := pixel.RGBA()
		intensity := (float64(r)*0.3 + float64(g)*0.59 + float64(b)*0.11) / 65535
		intensities[x-bounds.Min.X] = intensity
	}

	return intensities
}

// Generate the final musical samples
func generateMusicalSamples(intensities []float64, durationMs int, baseFreq float64, density int) []float64 {
	durationSec := float64(durationMs) / 1000.0
	totalSamples := int(float64(sampleRate) * durationSec)
	samples := make([]float64, totalSamples)

	noteSpacing := len(intensities) / density
	if noteSpacing < 1 {
		noteSpacing = 1
	}

	for i := 0; i < len(intensities); i += noteSpacing {
		amplitude := intensities[i]
		if amplitude < 0.05 {
			continue
		}

		noteOffset := float64(i) / float64(len(intensities)) * 24
		frequency := baseFreq * math.Pow(2, noteOffset/12.0)

		for j := 0; j < totalSamples; j++ {
			t := float64(j) / float64(sampleRate)
			harmonicWave := generateHarmonicWave(frequency, amplitude, t)
			envelopedWave := applyEnvelope(harmonicWave, float64(j)/float64(totalSamples), durationSec)
			samples[j] += envelopedWave
		}
	}

	maxAmplitude := 0.0
	for i := range samples {
		if math.Abs(samples[i]) > maxAmplitude {
			maxAmplitude = math.Abs(samples[i])
		}
	}

	if maxAmplitude > 0 {
		for i := range samples {
			samples[i] /= maxAmplitude
			samples[i] = math.Tanh(samples[i])
		}
	}

	return samples
}

// Handle the image to audio conversion process
func (g *Game) convertImageToAudio() {
	if g.Image == nil {
		return
	}

	bounds := g.Image.Bounds()
	duration := g.LengthSlider.Value
	baseFrequency := float64(g.FrequencySlider.Value)
	density := g.DensitySlider.Value

	middleRow := bounds.Min.Y + (bounds.Max.Y-bounds.Min.Y)/2
	intensities := imageToNoteIntensities(g.Image, middleRow)

	g.AudioData = generateMusicalSamples(intensities, duration, baseFrequency, density)

	filename := generateUniqueFilename(g.OutputDir)

	writer := wav.NewWriter(sampleRate)
	err := writer.WriteFile(filename, g.AudioData)
	if err != nil {
		log.Printf("Failed to write WAV: %v", err)
	} else {
		log.Printf("Created new audio file: %s", filename)
	}

	g.Converting = false
}

// Clean up resources before exit
func (g *Game) cleanup() {
	if g.Converting {
		log.Println("Waiting for conversion to complete...")
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("Closing application...")
}

// Update game state
func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.cleanup()
		os.Exit(0)
	}

	g.updateSlider(g.LengthSlider)
	g.updateSlider(g.FrequencySlider)
	g.updateSlider(g.DensitySlider)

	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.Converting {
		g.Converting = true
		go g.convertImageToAudio()
	}

	return nil
}

// Handle slider updates
func (g *Game) updateSlider(s *Slider) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if mx >= int(s.X) && mx <= int(s.X)+s.Width &&
			my >= int(s.Y)-10 && my <= int(s.Y)+10 {
			s.Dragging = true
		}
	} else {
		s.Dragging = false
	}

	if s.Dragging {
		mx, _ := ebiten.CursorPosition()
		newValue := (mx-int(s.X))*(s.Max-s.Min)/s.Width + s.Min
		s.Value = clamp(newValue, s.Min, s.Max)
	}
}

// Draw the game
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 40, 255})

	ebitenutil.DebugPrintAt(screen, "Duration (ms):", 10, 95)
	g.drawSlider(screen, g.LengthSlider)

	ebitenutil.DebugPrintAt(screen, "Base Pitch (Hz):", 10, 145)
	g.drawSlider(screen, g.FrequencySlider)

	ebitenutil.DebugPrintAt(screen, "Note Density:", 10, 195)
	g.drawSlider(screen, g.DensitySlider)

	if g.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 250)
		screen.DrawImage(ebiten.NewImageFromImage(g.Image), op)
	}

	ebitenutil.DebugPrintAt(screen, "Press SPACE to convert image to music", 10, 10)
	if g.Converting {
		ebitenutil.DebugPrintAt(screen, "Converting...", 10, 30)
	}
}

// Draw a slider
func (g *Game) drawSlider(screen *ebiten.Image, s *Slider) {
	ebitenutil.DrawRect(screen, s.X, s.Y, float64(s.Width), 5, color.Gray{Y: 200})
	pos := float64((s.Value-s.Min)*s.Width) / float64(s.Max-s.Min)
	ebitenutil.DrawRect(screen, s.X+pos-5, s.Y-5, 10, 15, color.White)
	ebitenutil.DebugPrintAt(screen, strconv.Itoa(s.Value), int(s.X+pos)-10, int(s.Y)-20)
}

// Handle window layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	outputDir := initializeOutputDirectory()

	file, err := os.Open("image.png")
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	game := &Game{
		LengthSlider:    NewSlider(100, 100, 200, minDuration, maxDuration, 2000),
		FrequencySlider: NewSlider(100, 150, 200, minBaseFreq, maxBaseFreq, 440),
		DensitySlider:   NewSlider(100, 200, 200, minDensity, maxDensity, 10),
		Image:           img,
		OutputDir:       outputDir,
	}

	defer game.cleanup()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		game.cleanup()
		os.Exit(0)
	}()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Musical Image to Sound Converter")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
