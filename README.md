Image to Sound Converter
This program transforms images into musical sounds, with a particular focus on spectrogram-like images. It provides an interactive interface for adjusting sound parameters and generates unique WAV files for each conversion.
Installation Requirements
Installing Go
First, you'll need to install Go on your system. Follow these steps:

Visit the official Go downloads page: https://go.dev/dl/
Download the appropriate version for your operating system:

For Windows: Download the .msi installer
For macOS: Download the .pkg installer
For Linux: Download the .tar.gz archive


Install Go following your system's standard installation process
Verify the installation by opening a terminal and running:
bashCopygo version


Setting Up the Project
Follow these steps to set up the project structure:

Create a new directory for your project:
bashCopymkdir image2sound
cd image2sound

Initialize the Go module:
bashCopygo mod init image2sound

Create the necessary subdirectories:
bashCopymkdir wav_output

Install required dependencies:
bashCopygo get github.com/hajimehoshi/ebiten/v2
go mod tidy


Project Structure
Your project should look like this:
Copyimage2sound/
├── main.go
├── go.mod
├── go.sum
├── wav/
│   └── writer.go
├── wav_output/
└── image.png
Creating Optimal Spectrograms
For the best results, your input images should follow these specifications:
Image Format and Size

File format: PNG
Recommended width: 400-600 pixels
Recommended height: 200-300 pixels
Color depth: 24-bit RGB or 32-bit RGBA

Design Guidelines
The program reads the middle row of pixels to generate sound. Create your images with these principles in mind:

Background

Use pure black (RGB: 0,0,0) for silence
Higher brightness values create louder sounds


Frequency Mapping

Left side of the image = lower frequencies
Right side of the image = higher frequencies
Vertical position doesn't affect the sound


Pattern Suggestions

Horizontal lines create steady tones
Diagonal lines create frequency sweeps
Dots create brief tones
Vertical spacing affects note density



Example Patterns
For testing, try creating these basic patterns:

Simple Tone Test
CopyWidth: 400px
Height: 200px
Pattern: Single horizontal white line in the middle
Expected Result: Clean, steady tone

Frequency Sweep
CopyWidth: 400px
Height: 200px
Pattern: Diagonal white line from bottom-left to top-right
Expected Result: Rising tone


Running the Program

Place your image file in the project directory as "image.png"
Run the program:
bashCopygo run main.go

Use the interface controls:

Duration: Controls sound length (500ms to 5000ms)
Base Pitch: Sets the fundamental frequency (220Hz to 880Hz)
Note Density: Adjusts polyphony (1 to 20 notes)


Press SPACE to convert the image to sound
Find your generated WAV files in the wav_output directory

Understanding the Output
The program generates WAV files with these characteristics:

Sample Rate: 44100 Hz (CD quality)
Bit Depth: 16-bit
Channels: Mono
Format: Uncompressed PCM

Each conversion creates a unique file named:
Copyoutput_YYYYMMDD_HHMMSS_XXXX.wav
Where:

YYYYMMDD is the date
HHMMSS is the time
XXXX is a random identifier
