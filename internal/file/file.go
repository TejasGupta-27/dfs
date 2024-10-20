package file

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"dfs/config"
)

type FileSystem struct {
	cfg *config.Config
}

type File struct {
	Name   string
	Size   int64
	Chunks []Chunk
}

type Chunk struct {
	ID   string
	Data []byte
}

func NewFileSystem(cfg *config.Config) *FileSystem {
	return &FileSystem{cfg: cfg}
}

func (fs *FileSystem) SplitFile(filePath string) (*File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	splitFile := &File{
		Name: filepath.Base(filePath),
		Size: fileInfo.Size(),
	}

	buffer := make([]byte, fs.cfg.ChunkSize)
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

func (fs *FileSystem) ReassembleFile(file *File, outputPath string) error {
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