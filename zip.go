package gozip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

type Reader struct {
	reader *zip.ReadCloser
}

type Writer struct {
	writer *zip.Writer
	file   *os.File
}

// OpenReader Opens an existing zip file for reading and returns a reader object
func OpenReader(fileName string) *Reader {
	zipReader, err := zip.OpenReader(fileName)
	if err != nil {
		panic(fmt.Errorf("Error opening source epub: %v", err))
	}
	return &Reader{zipReader}
}

// OpenWriter Open a new zip file for writing and returns a writer object.
// If the zip file already exists, it will be overwritten.
//
// NOTE: A zip file is not flushed automatically. The writer must be closed to ensure
// all content is completely written and archive metadata is written consistently.
func OpenWriter(filename string) *Writer {
	zipFile, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("Error opening target epub: %v", err))
	}
	return &Writer{
		file:   zipFile,
		writer: zip.NewWriter(zipFile),
	}
}

// Close closes the reader
func (r *Reader) Close() {
	r.reader.Close()
}

// Close closes the writer.
func (w *Writer) Close() {
	w.writer.Close()
	w.file.Close()
}

// GetContents returns the contents of a zip file as slice of filenames.
func (r *Reader) GetContents() []string {
	resultFileNames := make([]string, 0)
	for _, fileInfo := range r.reader.File {
		resultFileNames = append(resultFileNames, fileInfo.Name)
	}
	return resultFileNames
}

// ReadFile returns the content of a file as slice of bytes.
func (r *Reader) ReadFile(fileName string) []byte {
	cBlockSize := 1024
	f, err := r.reader.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("Error opening information on encrypted files: %v", err))
	}
	defer f.Close()

	if err != nil {
		panic(fmt.Errorf("Error opening information on encrypted files (cannot read stats): %v", err))
	}
	var fileContent []byte
	buf := make([]byte, cBlockSize)
	for {
		n, _ := f.Read(buf)
		fileContent = append(fileContent, buf[:n]...)
		if n < cBlockSize {
			break
		}
	}

	return fileContent
}

// WriteFile writes a byte array into a new file in the zip archive.
//
// NOTE: the zip file is not flushed. The writer must be closed to ensure the
// content is completely written and the zip archive is properly completed.
func (w *Writer) WriteFile(content []byte, fileName string) {
	f, err := w.writer.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("Error opening file '%s' in target epub: %v", fileName, err))
	}
	n, err := f.Write(content)
	if err != nil || n != len(content) {
		panic(fmt.Errorf("Error writing file '%s' to target epub: %v", fileName, err))
	}
}

// WriteFileNoCompression (just as WriteFile) writes a byte array into a new file in
// the zip archive. However, the content is not compressed.
func (w *Writer) WriteFileNoCompression(content []byte, fileName string) {
	f, err := w.writer.CreateHeader(&zip.FileHeader{Name: fileName, Method: zip.Store})
	if err != nil {
		panic(fmt.Errorf("Error opening file '%s' in target epub: %v", fileName, err))
	}
	n, err := f.Write(content)
	if err != nil || n != len(content) {
		panic(fmt.Errorf("Error writing file '%s' to target epub: %v", fileName, err))
	}
}

// CopyFile copies a file from a reader to a writer. (This is as a new file. The properties
// of the file are not taken over.)
func CopyFile(sourceReader *Reader, targetWriter *Writer, fileName string) {
	source, err := sourceReader.reader.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("Error reading from source file %s: %v", fileName, err))
	}
	target, err := targetWriter.writer.Create(fileName)
	if err != nil {
		panic(fmt.Errorf("Error creating plain file %s: %v", fileName, err))
	}
	_, err = io.Copy(target, source)
	if err != nil {
		panic(fmt.Errorf("Error copying content for file %s: %v", fileName, err))
	}
}

// CopyFileNoCompression (just as CopyFile) copies a file from a reader to a writer. The
// file content is not compressed.
func CopyFileNoCompression(sourceReader *Reader, targetWriter *Writer, fileName string) {
	source, err := sourceReader.reader.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("Error reading from source file %s: %v", fileName, err))
	}
	target, err := targetWriter.writer.CreateHeader(&zip.FileHeader{Name: fileName, Method: zip.Store})
	if err != nil {
		panic(fmt.Errorf("Error creating plain file %s: %v", fileName, err))
	}
	_, err = io.Copy(target, source)
	if err != nil {
		panic(fmt.Errorf("Error copying content for file %s: %v", fileName, err))
	}
}
