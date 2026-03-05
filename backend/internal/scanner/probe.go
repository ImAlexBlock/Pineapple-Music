package scanner

import (
	"encoding/binary"
	"io"
	"math"
	"os"
)

// ProbeAudioDuration attempts to determine audio file duration in seconds.
// Supports FLAC, MP3, M4A/MP4, OGG Vorbis, and WAV.
func ProbeAudioDuration(filePath string, format string, fileSize int64) float64 {
	f, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer f.Close()

	switch format {
	case "flac":
		return probeFlac(f)
	case "mp3":
		return probeMp3(f, fileSize)
	case "m4a", "mp4", "aac":
		return probeM4A(f)
	case "ogg":
		return probeOgg(f)
	case "wav":
		return probeWav(f, fileSize)
	}
	return 0
}

// probeFlac reads the STREAMINFO metadata block to get total samples and sample rate.
func probeFlac(f *os.File) float64 {
	// FLAC starts with "fLaC" magic, then metadata blocks
	magic := make([]byte, 4)
	if _, err := io.ReadFull(f, magic); err != nil {
		return 0
	}
	if string(magic) != "fLaC" {
		return 0
	}

	// Read metadata block header (4 bytes)
	header := make([]byte, 4)
	if _, err := io.ReadFull(f, header); err != nil {
		return 0
	}

	// Block type (lower 7 bits of first byte), STREAMINFO = 0
	blockType := header[0] & 0x7f
	blockLen := int(header[1])<<16 | int(header[2])<<8 | int(header[3])

	if blockType != 0 || blockLen < 34 {
		return 0
	}

	data := make([]byte, blockLen)
	if _, err := io.ReadFull(f, data); err != nil {
		return 0
	}

	// Bytes 10-12: sample rate (20 bits), starting at bit 80
	// Bits layout at byte offset 10: SSSS SSSS | SSSS SSSS | SSSS CCCC
	// S = sample rate bits (20 bits), C = channel bits
	sampleRate := uint32(data[10])<<12 | uint32(data[11])<<4 | uint32(data[12])>>4

	// Total samples: 36 bits starting at bit 116 (byte 14, bit 4)
	// Bytes 13-17: ....TTTT TTTTTTTT TTTTTTTT TTTTTTTT TTTTTTTT
	totalSamples := uint64(data[13]&0x0f)<<32 |
		uint64(data[14])<<24 |
		uint64(data[15])<<16 |
		uint64(data[16])<<8 |
		uint64(data[17])

	if sampleRate == 0 {
		return 0
	}

	return float64(totalSamples) / float64(sampleRate)
}

// probeMp3 reads MP3 to find duration.
// Tries Xing/VBRI header first, falls back to CBR estimate.
func probeMp3(f *os.File, fileSize int64) float64 {
	buf := make([]byte, 16384)
	n, err := f.Read(buf)
	if err != nil || n < 256 {
		return 0
	}
	buf = buf[:n]

	// Find first valid frame sync
	frameOff := -1
	for i := 0; i < len(buf)-4; i++ {
		if buf[i] == 0xff && (buf[i+1]&0xe0) == 0xe0 {
			frameOff = i
			break
		}
	}
	if frameOff < 0 {
		return 0
	}

	header := binary.BigEndian.Uint32(buf[frameOff:])
	mpegVer := (header >> 19) & 3    // 0=2.5, 2=2, 3=1
	layer := (header >> 17) & 3       // 1=III, 2=II, 3=I
	brIdx := (header >> 12) & 0x0f
	srIdx := (header >> 10) & 3

	if brIdx == 0 || brIdx == 15 || srIdx == 3 || layer == 0 || mpegVer == 1 {
		return 0
	}

	// Bitrate table for MPEG1 Layer III
	bitrateTable := [16]int{0, 32, 40, 48, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0, 0}
	sampleRateTable := [4]int{44100, 48000, 32000, 0}

	bitrate := 0
	sampleRate := 0

	if mpegVer == 3 && layer == 1 { // MPEG1 Layer III
		bitrate = bitrateTable[brIdx] * 1000
		sampleRate = sampleRateTable[srIdx]
	} else if mpegVer == 3 && layer == 2 { // MPEG1 Layer II
		br2 := [16]int{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0}
		bitrate = br2[brIdx] * 1000
		sampleRate = sampleRateTable[srIdx]
	} else if mpegVer == 3 && layer == 3 { // MPEG1 Layer I
		br1 := [16]int{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0}
		bitrate = br1[brIdx] * 1000
		sampleRate = sampleRateTable[srIdx]
	} else if (mpegVer == 2 || mpegVer == 0) && layer == 1 { // MPEG2/2.5 Layer III
		br23 := [16]int{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}
		bitrate = br23[brIdx] * 1000
		if mpegVer == 2 {
			sr2 := [4]int{22050, 24000, 16000, 0}
			sampleRate = sr2[srIdx]
		} else {
			sr25 := [4]int{11025, 12000, 8000, 0}
			sampleRate = sr25[srIdx]
		}
	}

	if bitrate == 0 || sampleRate == 0 {
		return 0
	}

	// Check for Xing/VBRI header within the first frame
	// Xing header is at offset 36 (for MPEG1 stereo) or 21 (mono) bytes from frame start
	xingOffsets := []int{36, 21, 13, 9}
	for _, off := range xingOffsets {
		pos := frameOff + off
		if pos+8 > len(buf) {
			continue
		}
		tag := string(buf[pos : pos+4])
		if tag == "Xing" || tag == "Info" {
			flags := binary.BigEndian.Uint32(buf[pos+4:])
			if flags&1 != 0 { // frames field present
				frames := binary.BigEndian.Uint32(buf[pos+8:])
				samplesPerFrame := 1152 // MPEG1 Layer III
				if mpegVer != 3 {
					samplesPerFrame = 576
				}
				return float64(frames) * float64(samplesPerFrame) / float64(sampleRate)
			}
		}
	}

	// Fallback: CBR estimate
	// Subtract estimated tag overhead (ID3v2 at start if any)
	audioSize := fileSize
	if frameOff > 0 {
		audioSize -= int64(frameOff)
	}
	if bitrate > 0 {
		return float64(audioSize) * 8 / float64(bitrate)
	}
	return 0
}

// probeM4A reads MP4/M4A container to find duration from mvhd atom.
func probeM4A(f *os.File) float64 {
	// Read up to 64KB to find moov/mvhd
	buf := make([]byte, 65536)
	n, _ := f.Read(buf)
	if n < 8 {
		return 0
	}
	buf = buf[:n]

	return findMvhd(buf)
}

func findMvhd(data []byte) float64 {
	for i := 0; i+8 <= len(data); {
		if i+8 > len(data) {
			break
		}
		size := int(binary.BigEndian.Uint32(data[i:]))
		name := string(data[i+4 : i+8])

		if size < 8 || i+size > len(data)+1 {
			break
		}

		if name == "moov" || name == "trak" || name == "mdia" {
			// Recurse into container atoms
			result := findMvhd(data[i+8 : i+size])
			if result > 0 {
				return result
			}
		}

		if name == "mvhd" {
			// Version 0: timescale at offset 20, duration at offset 24 (4 bytes each)
			// Version 1: timescale at offset 28, duration at offset 32 (8 bytes)
			off := i + 8
			if off >= len(data) {
				break
			}
			version := data[off]
			if version == 0 && off+28 <= len(data) {
				timeScale := binary.BigEndian.Uint32(data[off+12:])
				dur := binary.BigEndian.Uint32(data[off+16:])
				if timeScale > 0 {
					return float64(dur) / float64(timeScale)
				}
			} else if version == 1 && off+36 <= len(data) {
				timeScale := binary.BigEndian.Uint32(data[off+20:])
				dur := binary.BigEndian.Uint64(data[off+24:])
				if timeScale > 0 {
					return float64(dur) / float64(timeScale)
				}
			}
		}

		if name == "mdhd" {
			off := i + 8
			if off >= len(data) {
				break
			}
			version := data[off]
			if version == 0 && off+24 <= len(data) {
				timeScale := binary.BigEndian.Uint32(data[off+12:])
				dur := binary.BigEndian.Uint32(data[off+16:])
				if timeScale > 0 {
					return float64(dur) / float64(timeScale)
				}
			} else if version == 1 && off+32 <= len(data) {
				timeScale := binary.BigEndian.Uint32(data[off+20:])
				dur := binary.BigEndian.Uint64(data[off+24:])
				if timeScale > 0 {
					return float64(dur) / float64(timeScale)
				}
			}
		}

		i += size
	}
	return 0
}

// probeOgg reads OGG Vorbis to determine duration from last page granule position.
func probeOgg(f *os.File) float64 {
	// Read first page to get sample rate from Vorbis identification header
	header := make([]byte, 4096)
	n, err := f.Read(header)
	if err != nil || n < 58 {
		return 0
	}

	// Verify OGG magic
	if string(header[:4]) != "OggS" {
		return 0
	}

	// Find Vorbis identification header
	sampleRate := 0
	for i := 0; i < n-15; i++ {
		if header[i] == 1 && string(header[i+1:i+7]) == "vorbis" {
			sr := binary.LittleEndian.Uint32(header[i+11:])
			sampleRate = int(sr)
			break
		}
	}
	if sampleRate == 0 {
		return 0
	}

	// Read last 16KB to find the last OGG page with granule position
	fi, err := f.Stat()
	if err != nil {
		return 0
	}
	tailSize := int64(16384)
	if fi.Size() < tailSize {
		tailSize = fi.Size()
	}
	f.Seek(fi.Size()-tailSize, io.SeekStart)
	tail := make([]byte, tailSize)
	n, _ = f.Read(tail)
	tail = tail[:n]

	// Search backwards for last OggS sync
	var granule int64
	for i := len(tail) - 14; i >= 0; i-- {
		if string(tail[i:i+4]) == "OggS" {
			granule = int64(binary.LittleEndian.Uint64(tail[i+6:]))
			if granule > 0 {
				break
			}
		}
	}

	if granule > 0 && sampleRate > 0 {
		return float64(granule) / float64(sampleRate)
	}
	return 0
}

// probeWav reads WAV/RIFF header to get duration.
func probeWav(f *os.File, fileSize int64) float64 {
	header := make([]byte, 44)
	if _, err := io.ReadFull(f, header); err != nil {
		return 0
	}

	if string(header[:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return 0
	}

	// Find "fmt " chunk
	sampleRate := 0
	byteRate := 0
	pos := 12
	buf := make([]byte, 8192)
	copy(buf, header[12:44])
	remaining := 32

	for pos < int(fileSize) && remaining >= 8 {
		chunkID := string(buf[:4])
		chunkSize := int(binary.LittleEndian.Uint32(buf[4:]))

		if chunkID == "fmt " && chunkSize >= 16 && remaining >= 24 {
			sampleRate = int(binary.LittleEndian.Uint32(buf[12:]))
			byteRate = int(binary.LittleEndian.Uint32(buf[16:]))
			break
		}

		skip := 8 + chunkSize
		if skip%2 != 0 {
			skip++ // padding byte
		}
		pos += skip
		f.Seek(int64(pos), io.SeekStart)
		n, err := f.Read(buf)
		if err != nil || n < 8 {
			break
		}
		remaining = n
	}

	_ = sampleRate
	if byteRate > 0 {
		// Estimate: total file size / byte rate (approximate, includes headers)
		dataSize := float64(fileSize) - 44 // subtract header estimate
		if dataSize < 0 {
			dataSize = float64(fileSize)
		}
		return math.Max(0, dataSize/float64(byteRate))
	}
	return 0
}
