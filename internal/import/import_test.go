package importer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	importer "github.com/user/stashctl/internal/import"
)

type jsonSnippet struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
}

func writeTempJSON(t *testing.T, data interface{}) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestFromJSONValid(t *testing.T) {
	data := []jsonSnippet{
		{Title: "Hello", Content: "fmt.Println()", Language: "go", Tags: []string{"go"}},
		{Title: "World", Content: "print()", Language: "python"},
	}
	path := writeTempJSON(t, data)
	snippets, res, err := importer.FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snippets) != 2 {
		t.Errorf("expected 2 snippets, got %d", len(snippets))
	}
	if res.Imported != 2 || res.Skipped != 0 {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestFromJSONSkipsInvalid(t *testing.T) {
	data := []jsonSnippet{
		{Title: "", Content: "something"},
		{Title: "Valid", Content: "code"},
	}
	path := writeTempJSON(t, data)
	snippets, res, err := importer.FromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snippets) != 1 {
		t.Errorf("expected 1 snippet, got %d", len(snippets))
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestFromFileUnsupportedFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	os.WriteFile(path, []byte("data"), 0644)
	_, _, err := importer.FromFile(path)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestFromFileMissing(t *testing.T) {
	_, _, err := importer.FromFile("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
