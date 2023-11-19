#!/bin/bash

# Windows
GOOS=windows GOARCH=amd64 go build -o structera-windows-amd64.exe  # 64-bit Windows
GOOS=windows GOARCH=386 go build -o structera-windows-386.exe      # 32-bit Windows

# macOS (Darwin)
GOOS=darwin GOARCH=amd64 go build -o structera-macos-amd64         # Intel Macs (64-bit)
GOOS=darwin GOARCH=arm64 go build -o structera-macos-arm64         # Apple Silicon Macs (M1/M2)

# Linux
GOOS=linux GOARCH=amd64 go build -o structera-linux-amd64          # 64-bit Linux
GOOS=linux GOARCH=386 go build -o structera-linux-386              # 32-bit Linux
GOOS=linux GOARCH=arm64 go build -o structera-linux-arm64          # ARM64 Linux (e.g., Raspberry Pi 64-bit)

# Additional Linux Architectures (Optional)
GOOS=linux GOARCH=arm go build -o structera-linux-arm              # ARM Linux (32-bit, e.g., older Raspberry Pi)
GOOS=linux GOARCH=mips go build -o structera-linux-mips            # MIPS Linux (for MIPS architecture devices)
GOOS=linux GOARCH=mips64 go build -o structera-linux-mips64        # MIPS64 Linux (64-bit MIPS architecture)
