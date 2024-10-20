package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/TejasGupta-27/dfs/config"
)

type File struct {
	Name   string
	Size   int64
	Chunks []Chunk
}

type Chunk struct {
	ID   string
	Data []byte
}

func SplitFile(filePath string, cfg *config.Config) (*File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(filePath)
	fileSize := fileInfo.Size()
	chunkSize := cfg.ChunkSize

	splitFile := &File{
		Name: fileName,
		Size: fileSize,
	}

	buffer := make([]byte, chunkSize)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		chunk := Chunk{
			Data: buffer[:bytesRead],
		}
		chunk.ID = generateChunkID(chunk.Data)
		splitFile.Chunks = append(splitFile.Chunks, chunk)
	}

	return splitFile, nil
}

func ReassembleFile(file *File, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, chunk := range file.Chunks {
		_, err := outFile.Write(chunk.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateChunkID(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}