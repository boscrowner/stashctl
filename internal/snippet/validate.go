package snippet

import (
	"errors"
	"strings"
)

// Validation errors.
var (
	ErrEmptyTitle   = errors.New("snippet title must not be empty")
	ErrEmptyContent = errors.New("snippet content must not be empty")
	ErrTitleTooLong = errors.New("snippet title must not exceed 200 characters")
	ErrInvalidTag   = errors.New("tag must not contain spaces or commas")
)

// Validate checks that a Snippet has valid fields.
func Validate(s *Snippet) error {
	if strings.TrimSpace(s.Title) == "" {
		return ErrEmptyTitle
	}
	if len(s.Title) > 200 {
		return ErrTitleTooLong
	}
	if strings.TrimSpace(s.Content) == "" {
		return ErrEmptyContent
	}
	for _, tag := range s.Tags {
		if err := validateTag(tag); err != nil {
			return err
		}
	}
	return nil
}

func validateTag(tag string) error {
	if strings.ContainsAny(tag, " ,") {
		return ErrInvalidTag
	}
	return nil
}
