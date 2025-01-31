# Image to Sound Converter

This program transforms images into musical sounds, with a particular focus on spectrogram-like images. It provides an interactive interface for adjusting sound parameters and generates unique WAV files for each conversion.

## Installation Requirements

### Installing Go

First, you'll need to install Go on your system. Follow these steps:

1. Visit the official Go downloads page: https://go.dev/dl/
2. Download the appropriate version for your operating system:
   - For Windows: Download the .msi installer
   - For macOS: Download the .pkg installer
   - For Linux: Download the .tar.gz archive
3. Install Go following your system's standard installation process
4. Verify the installation by opening a terminal and running:
   ```bash
   go version
   ```
### Recommended Development Environment
https://code.visualstudio.com/download

### Recommended VS Code Extensions
**For Core Development:
"Go" by Go Team at Google**

- This provides our foundational Go support
- Includes the language server protocol integration
- Handles automatic imports and formatting


**For Better Dependency Management:
"Dependi" by Fill Labs**

- Visualizes dependency relationships
- Shows direct and indirect dependencies
- Helps identify potential dependency conflicts
- Creates interactive dependency graphs



**For the .PNG File We're Using:
"Image Preview" by Kiss Tamás**

- Essential for working with our spectrogram images
- Shows image dimensions and format details
- Provides quick visual feedback




### Setting Up the Project

Follow these steps to set up the project structure:

1. Create a new directory for your project:
   ```bash
   mkdir image2sound
   cd image2sound
   ```

2. Initialize the Go module:
   ```bash
   go mod init image2sound
   ```

3. Create the necessary subdirectories:
   ```bash
   mkdir wav_output
   ```

4. Install required dependencies:
   ```bash
   go get github.com/hajimehoshi/ebiten/v2
   go mod tidy
   ```

## Project Structure

Your project should look like this:
```
image2sound/
├── main.go
├── go.mod
├── go.sum
├── wav/
│   └── writer.go
├── wav_output/
└── image.png
```

## Creating Optimal Spectrograms

For the best results, your input images should follow these specifications:

### Image Format and Size
- File format: PNG
- Recommended width: 400-600 pixels
- Recommended height: 200-300 pixels
- Color depth: 24-bit RGB or 32-bit RGBA

### Design Guidelines
The program reads the middle row of pixels to generate sound. Create your images with these principles in mind:

1. Background
   - Use pure black (RGB: 0,0,0) for silence
   - Higher brightness values create louder sounds

2. Frequency Mapping
   - Left side of the image = lower frequencies
   - Right side of the image = higher frequencies
   - Vertical position doesn't affect the sound

3. Pattern Suggestions
   - Horizontal lines create steady tones
   - Diagonal lines create frequency sweeps
   - Dots create brief tones
   - Vertical spacing affects note density

### Example Patterns
For testing, try creating these basic patterns:

1. Simple Tone Test
   ```
   Width: 400px
   Height: 200px
   Pattern: Single horizontal white line in the middle
   Expected Result: Clean, steady tone
   ```



## Running the Program

1. Place your image file in the project directory as "image.png"

2. Run the program:
   ```bash
   go run main.go
   ```

3. Use the interface controls:
   - Duration: Controls sound length (500ms to 5000ms)
   - Base Pitch: Sets the fundamental frequency (220Hz to 880Hz)
   - Note Density: Adjusts polyphony (1 to 20 notes)

4. Press SPACE to convert the image to sound

5. Find your generated WAV files in the wav_output directory

## Understanding the Output

The program generates WAV files with these characteristics:
- Sample Rate: 44100 Hz (CD quality)
- Bit Depth: 16-bit
- Channels: Mono
- Format: Uncompressed PCM

Each conversion creates a unique file named:
```
output_YYYYMMDD_HHMMSS_XXXX.wav
```
Where:
- YYYYMMDD is the date
- HHMMSS is the time
- XXXX is a random identifier

## Troubleshooting

Common issues and solutions:

1. "Image file not found"
   - Ensure image.png exists in the project root directory
   - Check file permissions

2. "Failed to create WAV file"
   - Verify wav_output directory exists
   - Check write permissions
