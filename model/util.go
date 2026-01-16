package model

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Generic type constraint
type HasID interface {
	GetID() uint
}

// Function to extract IDs from a slice of structs
func ExtractIDs[T HasID](data []T) []uint {
	var ids []uint
	for _, item := range data {
		ids = append(ids, item.GetID())
	}

	return ids
}

func NormalizeSpaces(str string) string {
	// 1. Trim leading/trailing spaces:
	str = strings.TrimSpace(str)

	// 2. Replace multiple spaces with single spaces:
	spaceRegex := regexp.MustCompile(`\s+`)
	str = spaceRegex.ReplaceAllString(str, " ")

	return str
}

func TrimAllSpace(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// User
func HashPassword(password string) ([]byte, error) {
	// Choose a suitable cost factor for bcrypt
	cost := 12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return bytes, err
}

func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // true if the passwords match
}

type ValidationError struct {
	ItemName string
	Reasons  []string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s is invalid: %v", e.ItemName, e.Reasons)
}

func ValidateStringLength(itemName, value string, min, max int) *ValidationError {
	minLength := min
	maxLength := max
	reasons := []string{} // Start with an empty list of reasons
	if len(value) < minLength {
		reasons = append(reasons, fmt.Sprintf("%s is too short (minimum %d characters)", itemName, minLength))
	}

	if len(value) > maxLength {
		reasons = append(reasons, fmt.Sprintf("%s is too long (maximum %d characters)", itemName, maxLength))
	}

	if len(reasons) > 0 {
		return &ValidationError{ItemName: itemName, Reasons: reasons}
	}

	return nil
}

func ValidatePassword(password string) *ValidationError {
	minLength := 8        // 8-12 characters (the absolute minimum for basic security)
	maxLength := 64       // 64-128 characters (allows for passphrases, but some systems have limits)
	reasons := []string{} // Start with an empty list of reasons

	if len(password) < minLength {
		reasons = append(reasons, fmt.Sprintf("Password is too short (minimum %d characters)", minLength))
	}

	if len(password) > maxLength {
		reasons = append(reasons, fmt.Sprintf("Password is too long (maximum %d characters)", maxLength))
	}

	if len(reasons) > 0 {
		return &ValidationError{ItemName: "Password", Reasons: reasons}
	}

	return nil // Password passed validation
}

// Image

type ImageExtension int

const (
	JPEG ImageExtension = iota
	PNG
	GIF
	WebP
	TIFF
	BMP
)

var formatExts = map[string]ImageExtension{
	"jpg":  JPEG,
	"jpeg": JPEG,
	"png":  PNG,
	"gif":  GIF,
	"webp": WebP,
	"tif":  TIFF,
	"tiff": TIFF,
}

var formatNames = map[ImageExtension]string{
	JPEG: "jpg",
	PNG:  "png",
	GIF:  "gif",
	WebP: "webp",
	TIFF: "tiff",
}

func (f ImageExtension) String() string {
	return formatNames[f]
}

// FormatFilename takes a filename and returns a new filename with the correct extension
func FormatFilename(filename string) (string, error) {
	filename = sanitizeFilename(filepath.Base(filename), '_')
	ext := filepath.Ext(filename)
	base := strings.Trim(filename, ext)

	mExt, ok := formatExts[strings.ToLower(strings.TrimPrefix(ext, "."))]
	if !ok {
		return "", errors.New("unsupported image format")
	}

	return base + "." + mExt.String(), nil
}

var validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_.-"

func sanitizeFilename(filename string, replaceChar rune) string {
	var sanitized strings.Builder
	for _, r := range filename {
		// unicode.IsGraphic: Ensures we only include printable characters.
		// strings.ContainsRune: Checks if the rune (character) is within your validChars set.
		if unicode.IsGraphic(r) && strings.ContainsRune(validChars, r) {
			// Valid character
			sanitized.WriteRune(r)
		} else {
			// Invalid character
			sanitized.WriteRune(replaceChar)
		}
	}

	return sanitized.String()
}
