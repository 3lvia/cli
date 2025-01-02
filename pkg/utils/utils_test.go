package utils

import (
	"embed"
	"path"
	"testing"
)

//go:embed _test/*.tmpl
var templates embed.FS

func TestRemoveZeroValues(t *testing.T) {
	t.Parallel()

	slice := []string{"", "a", "", "b", "c", "", "d", ""}
	expected := []string{"a", "b", "c", "d"}

	result := RemoveZeroValues(slice)

	if len(result) != len(expected) {
		t.Fatalf("Expected %v, but got %v", expected, result)
	}

	for i, value := range expected {
		if result[i] != value {
			t.Fatalf("Expected %v, but got %v", expected, result)
		}
	}
}

func TestWriteFileWithTemplate1(t *testing.T) {
	t.Parallel()

	const (
		fileName     = "file.txt"
		templateFile = "_test/" + fileName + ".tmpl"
	)

	directory := t.TempDir()

	filePath, err := WriteFileWithTemplate(
		directory,
		fileName,
		templateFile,
		templates,
		struct {
			Name string
		}{
			Name: "World",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if filePath != path.Join(directory, fileName) {
		t.Fatalf("Expected %s, but got %s", path.Join(directory, fileName), filePath)
	}
}

func TestWriteFileWithTemplate2(t *testing.T) {
	t.Parallel()

	const (
		fileName     = "Dockerfile"
		templateFile = "_test/" + fileName + ".tmpl"
	)

	directory := t.TempDir()

	filePath, err := WriteFileWithTemplate(
		directory,
		fileName,
		templateFile,
		templates,
		struct {
			ImageName string
			ImageTag  string
		}{
			ImageName: "alpine",
			ImageTag:  "latest",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if filePath != path.Join(directory, fileName) {
		t.Fatalf("Expected %s, but got %s", path.Join(directory, fileName), filePath)
	}
}

func TestPromptYesNo(t *testing.T) {
	t.Parallel()

	answer, err := PromptYesNo("Do you like testing?", true)
	if err != nil {
		t.Fatal(err)
	}

	if !answer {
		t.Fatalf("Expected %v, but got %v", true, answer)
	}
}
