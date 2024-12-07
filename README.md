# Media Organizer

A command-line tool that automatically organizes your photos and videos by reading EXIF/metadata information and copying them into a structured directory based on the capture date.

NOTICE: ALL CODE IS WRITTEN BY AI, I JUST MODIFIED IT.

[![Go Version](https://img.shields.io/github/go-mod/go-version/yourusername/media-organizer)](https://golang.org/dl/)
[![License](https://img.shields.io/github/license/yourusername/media-organizer)](LICENSE)
[![Build Status](https://github.com/chiyiangel/media-organizer/workflows/Build%20and%20Test/badge.svg)](https://github.com/chiyiangel/media-organizer/actions)

## Features

- 📸 Supports various media formats:
  - Photos: JPG, JPEG, PNG, RAW, CR2, NEF, ARW
  - Videos: MP4, MOV, AVI, MKV, WMV, FLV, M4V, 3GP, WEBM
- 📅 Organizes files by capture date (EXIF) or creation date
- 📁 Creates a structured directory tree (year/month)
- 🚀 Multi-threaded processing for better performance
- 📊 Real-time progress display with TUI
- 📝 Detailed logging with different log levels
- 🖥️ Cross-platform support (Windows, macOS, Linux, Synology NAS)

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/chiyiangel/media-organizer/releases).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/chiyiangel/media-organizer
cd media-organizer

# Build for your current platform
make build

# Or build for a specific platform
make linux-amd64    # For Linux
make darwin-arm64   # For macOS (Apple Silicon)
make windows-amd64  # For Windows
make synology-amd64 # For Synology NAS
```

## Usage

Basic usage:
```bash
./media-organizer -src /path/to/photos -dest /path/to/output
```

Available options:
```bash
  -src string
        Source directory containing media files
  -dest string
        Destination directory for organized files
  -log string
        Log file path (default: logs/media-organizer-{timestamp}.log)
  -quiet
        Quiet mode, only output to log file
```

The tool will:
1. Scan the source directory for media files
2. Read EXIF/metadata information
3. Create a structured directory in the destination:
```
output/
├── photos/
│   ├── 2024/
│   │   ├── 01/
│   │   ├── 02/
│   │   └── ...
│   └── ...
└── videos/
    ├── 2024/
    │   ├── 01/
    │   ├── 02/
    │   └── ...
    └── ...
```

## Building

The project includes a Makefile with various build targets:

```bash
make help     # Show all available build targets
make all      # Build for all platforms
make release  # Build and package all platforms
make test     # Run tests
```

### Supported Platforms

- Linux (AMD64, ARM64, ARMv7)
- macOS (Intel, Apple Silicon)
- Windows (AMD64, ARM64)
- Synology NAS (Intel/AMD, ARM64, ARMv7)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [goexif](https://github.com/rwcarlsen/goexif) - For EXIF metadata extraction
- [bubbletea](https://github.com/charmbracelet/bubbletea) - For the TUI interface
- [bubbles](https://github.com/charmbracelet/bubbles) - For TUI components