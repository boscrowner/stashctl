package snippet_test

import (
	"strings"
	"testing"

	"github.com/user/stashctl/internal/snippet"
)

func makeAttachment(t *testing.T, snippetID, name, path, mime string) snippet.Attachment {
	t.Helper()
	a, err := snippet.NewAttachment(snippetID, name, path, mime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return a
}

func TestNewAttachmentValid(t *testing.T) {
	a := makeAttachment(t, "snip-1", "notes.txt", "/tmp/notes.txt", "text/plain")
	if a.ID == "" {
		t.Error("expected non-empty ID")
	}
	if a.SnippetID != "snip-1" {
		t.Errorf("expected snippet_id snip-1, got %s", a.SnippetID)
	}
	if a.MIMEType != "text/plain" {
		t.Errorf("expected mime text/plain, got %s", a.MIMEType)
	}
}

func TestNewAttachmentDefaultMIME(t *testing.T) {
	a := makeAttachment(t, "snip-1", "data.bin", "/tmp/data.bin", "")
	if a.MIMEType != "application/octet-stream" {
		t.Errorf("expected default mime, got %s", a.MIMEType)
	}
}

func TestNewAttachmentEmptySnippetID(t *testing.T) {
	_, err := snippet.NewAttachment("", "file.txt", "/tmp/file.txt", "")
	if err == nil {
		t.Error("expected error for empty snippet_id")
	}
}

func TestNewAttachmentEmptyName(t *testing.T) {
	_, err := snippet.NewAttachment("snip-1", "", "/tmp/file.txt", "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNewAttachmentNameTooLong(t *testing.T) {
	long := strings.Repeat("x", 129)
	_, err := snippet.NewAttachment("snip-1", long, "/tmp/file.txt", "")
	if err == nil {
		t.Error("expected error for name exceeding 128 chars")
	}
}

func TestNewAttachmentEmptyPath(t *testing.T) {
	_, err := snippet.NewAttachment("snip-1", "file.txt", "", "")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestAttachmentsFor(t *testing.T) {
	a1 := makeAttachment(t, "snip-1", "a.txt", "/tmp/a.txt", "")
	a2 := makeAttachment(t, "snip-2", "b.txt", "/tmp/b.txt", "")
	a3 := makeAttachment(t, "snip-1", "c.txt", "/tmp/c.txt", "")

	all := []snippet.Attachment{a1, a2, a3}
	result := snippet.AttachmentsFor(all, "snip-1")
	if len(result) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(result))
	}
}

func TestRemoveAttachment(t *testing.T) {
	a1 := makeAttachment(t, "snip-1", "a.txt", "/tmp/a.txt", "")
	a2 := makeAttachment(t, "snip-1", "b.txt", "/tmp/b.txt", "")

	updated, ok := snippet.RemoveAttachment([]snippet.Attachment{a1, a2}, a1.ID)
	if !ok {
		t.Error("expected removal to succeed")
	}
	if len(updated) != 1 {
		t.Errorf("expected 1 attachment after removal, got %d", len(updated))
	}
}

func TestRemoveAttachmentNotFound(t *testing.T) {
	a1 := makeAttachment(t, "snip-1", "a.txt", "/tmp/a.txt", "")
	updated, ok := snippet.RemoveAttachment([]snippet.Attachment{a1}, "nonexistent")
	if ok {
		t.Error("expected removal to fail for unknown id")
	}
	if len(updated) != 1 {
		t.Error("slice should be unchanged")
	}
}
