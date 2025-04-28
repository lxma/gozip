package gozip

import (
    "bufio"
    "bytes"
    "github.com/stretchr/testify/assert"
    "os"
    "slices"
    "testing"
)

const cTestSourceFile = "testsource.zip"
const cTestTargetFile = "testtarget.zip"

func TestReader(t *testing.T) {
    reader := OpenReader(cTestSourceFile)
    defer reader.Close()

    contents := reader.GetContents()
    slices.Sort(contents)
    assert.True(t, slices.Equal(contents, []string{"fileone.txt", "filetwo.txt"}), "file contents should be detected correctly")

    assert.Equal(t, "Hello\n", string(reader.ReadFile("fileone.txt")), "file should be returned correctly")
    assert.Equal(t, "World\n", string(reader.ReadFile("filetwo.txt")), "file should be returned correctly")
}

func getFileContents(t *testing.T, filename string) []byte {
    file, err := os.Open(filename)
    if err != nil {
        assert.FailNow(t, err.Error())
    }
    defer file.Close()

    stats, statsErr := file.Stat()
    if statsErr != nil {
        assert.FailNow(t, err.Error())
    }

    contents := make([]byte, stats.Size())
    bufferedReader := bufio.NewReader(file)
    _, err = bufferedReader.Read(contents)
    if statsErr != nil {
        assert.FailNow(t, err.Error())
    }
    return contents
}

func TestWriter(t *testing.T) {
    writer := OpenWriter(cTestTargetFile)
    writer.WriteFile([]byte("Hello World 1"), "hello1.txt")
    writer.WriteFile([]byte("Hello World 2"), "hello2.txt")
    writer.WriteFileNoCompression([]byte("Hello World 3"), "hello3.txt")
    writer.Close()

    reader := OpenReader(cTestTargetFile)
    contents := reader.GetContents()
    slices.Sort(contents)
    assert.True(t, slices.Equal(contents, []string{"hello1.txt", "hello2.txt", "hello3.txt"}), "all files should have been written")
    assert.Equal(t, "Hello World 1", string(reader.ReadFile("hello1.txt")), "file should be returned correctly")
    assert.Equal(t, "Hello World 2", string(reader.ReadFile("hello2.txt")), "file should be returned correctly")
    assert.Equal(t, "Hello World 3", string(reader.ReadFile("hello3.txt")), "file should be returned correctly")
    reader.Close()

    zipFileAsBytes := getFileContents(t, cTestTargetFile)
    indexOfUncompressedContents := bytes.Index(zipFileAsBytes, []byte("Hello World 3"))
    assert.True(t, indexOfUncompressedContents > 0, "Third file should be written uncompressed")
}

func TestMain(m *testing.M) {
    m.Run()
    os.Remove(cTestTargetFile)
}
