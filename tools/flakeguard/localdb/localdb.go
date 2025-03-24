package localdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

const defaultDBFileName = ".flaky_test_db.json"

// Entry represents one record in the local DB for a given test.
type Entry struct {
	TestPackage string    `json:"test_package"`
	TestName    string    `json:"test_name"`
	JiraTicket  string    `json:"jira_ticket"`
	IsSkipped   bool      `json:"is_skipped,omitempty"`
	SkippedAt   time.Time `json:"skipped_at,omitempty"`
}

// DB is a simple in-memory map keyed by "pkg::testName" => Entry,
// and it contains the file path used for persistence.
type DB struct {
	data map[string]Entry
	path string
}

// NewDB creates a new, empty DB using the default file path.
func NewDB() DB {
	return DB{
		data: make(map[string]Entry),
		path: getDefaultDBPath(),
	}
}

// NewDBWithPath creates a new, empty DB using the provided path.
// If the provided path is empty, the default path is used.
func NewDBWithPath(path string) DB {
	if path == "" {
		path = getDefaultDBPath()
	}
	return DB{
		data: make(map[string]Entry),
		path: path,
	}
}

// LoadDB loads the JSON file from the default path into a DB.
// If the file does not exist, an empty DB is returned.
func LoadDB() (DB, error) {
	return LoadDBWithPath("")
}

// LoadDBWithPath loads the JSON file from the specified path into a DB.
// If path is empty, the default path is used.
func LoadDBWithPath(path string) (DB, error) {
	if path == "" {
		path = getDefaultDBPath()
	}
	db := DB{
		data: make(map[string]Entry),
		path: path,
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist => return empty DB.
		return db, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return db, fmt.Errorf("failed to open local DB file: %w", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&db.data); err != nil {
		return db, fmt.Errorf("failed to decode local DB JSON: %w", err)
	}
	return db, nil
}

// Save persists the DB contents to its file path in JSON format.
func (db *DB) Save() error {
	f, err := os.Create(db.path)
	if err != nil {
		return fmt.Errorf("failed to create local DB file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // Optional: prettify the JSON.
	if err := enc.Encode(db.data); err != nil {
		return fmt.Errorf("failed to encode local DB to JSON: %w", err)
	}
	return nil
}

// FilePath returns the file path where the DB is stored.
func (db *DB) FilePath() string {
	return db.path
}

// Get retrieves the Jira ticket ID for (testPackage, testName), if it exists.
func (db *DB) Get(testPackage, testName string) (string, bool) {
	key := makeKey(testPackage, testName)
	entry, exists := db.data[key]
	return entry.JiraTicket, exists
}

// Set updates or inserts a record for (testPackage, testName) => jiraTicket.
func (db *DB) Set(testPackage, testName, jiraTicket string) {
	key := makeKey(testPackage, testName)
	db.data[key] = Entry{
		TestPackage: testPackage,
		TestName:    testName,
		JiraTicket:  jiraTicket,
	}
}

// UpdateTicketStatus updates or inserts the ticket status for (testPackage, testName).
// It accepts a boolean for whether the ticket is skipped and a timestamp in UTC.
func (db *DB) UpdateTicketStatus(testPackage, testName string, isSkipped bool, skippedAt time.Time) {
	key := makeKey(testPackage, testName)
	entry, exists := db.data[key]
	if !exists {
		entry = Entry{
			TestPackage: testPackage,
			TestName:    testName,
		}
	}
	entry.IsSkipped = isSkipped
	if isSkipped {
		entry.SkippedAt = skippedAt.UTC() // Ensure the timestamp is saved as UTC.
	} else {
		// If not skipped, clear the timestamp.
		entry.SkippedAt = time.Time{}
	}
	db.data[key] = entry
}

// GetAllEntries returns all entries in the DB.
func (db *DB) GetAllEntries() []Entry {
	entries := make([]Entry, 0, len(db.data))
	for _, entry := range db.data {
		entries = append(entries, entry)
	}
	return entries
}

// getDefaultDBPath returns the default DB file path in the user's home directory.
// If the user's home directory cannot be determined, it falls back to the current directory.
func getDefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warn().Msg("Failed to get user home directory; using current directory for local DB.")
		home = "."
	}
	return filepath.Join(home, defaultDBFileName)
}

// makeKey is a helper to combine the package and test name into a single map key.
func makeKey(pkg, testName string) string {
	return pkg + "::" + testName
}
