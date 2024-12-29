# AI Security Alert Generator

A Go-based tool that generates realistic IT security alert voicemails using OpenAI's GPT-4 for content generation and ElevenLabs for text-to-speech conversion.

## Overview

This tool creates convincing IT security alert voicemails by:
1. Generating fictional but legitimate-sounding IT department details
2. Creating an urgent security alert script
3. Converting the script to a natural-sounding voicemail using AI voice synthesis

## Prerequisites

- Go 1.x
- OpenAI API key
- ElevenLabs API key

## Environment Variables

The following environment variables need to be set:
```bash
OPENAI_API_KEY=your_openai_api_key
ELEVENLABS_API_KEY=your_elevenlabs_api_key
```

## Usage

1. Set up your environment variables
2. Run the program:
```bash
go run main.go
```

3. The program will:
   - Generate company IT department details
   - Create a security alert script
   - Allow you to review and regenerate content if needed
   - Generate an MP3 voicemail file

## Output

The program generates a `voicemail.mp3` file in the current directory.

## Note

This tool uses:
- OpenAI's GPT-4 for content generation
- ElevenLabs' text-to-speech API (using the "Brian" voice)
- Interactive prompts for content verification

## Disclaimer

This tool should only be used for legitimate purposes and in accordance with all applicable laws and regulations.
