package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestImportCmd(t *testing.T) {
	env := setupTest(t)

	type entry struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Language string   `json:"language"`
		Tags     []string `json:"tags"`
	}
	data := []entry{
		{Title: "Greet", Content: "fmt.Println(\"hi\")", Language: "go", Tags: []string{"go"}},
		{Title: "Loop", Content: "for i := range x {}", Language: "go", Tags: []string{"go", "loop"}},
	}
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(f).Encode(data)
	f.Close()

	var out bytes.Buffer
	env.cmd.SetOut(&out)
	env.cmd.SetArgs([]string{"import", f.Name()})
	if err := env.cmd.Execute(); err != nil {
		t.Fatalf("import cmd failed: %v", err)
	}
	if !strings.Contains(out.String(), "Imported: 2") {
		t.Errorf("expected import count in output, got: %s", out.String())
	}
}

func TestImportCmdBadFile(t *testing.T) {
	env := setupTest(t)
	env.cmd.SetArgs([]string{"import", "/no/such/file.json"})
	if err := env.cmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}
