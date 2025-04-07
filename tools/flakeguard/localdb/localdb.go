package localdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync" // Import sync for mutex
	"time"

	"github.com/rs/zerolog/log"
)

const defaultDBFileName = ".flaky_test_db.json"

// Entry represents one record in the local DB for a given test.
// Removed IsSkipped as SkippedAt.IsZero() is the source of truth.
type Entry struct {
	TestPackage string    `json:"test_package"`
	TestName    string    `json:"test_name"`
	JiraTicket  string    `json:"jira_ticket,omitempty"` // Allow empty string
	AssigneeID  string    `json:"jira_assignee_id,omitempty"`
	SkippedAt   time.Time `json:"skipped_at,omitempty"`
}

// DB is a simple in-memory map keyed by "pkg::testName" => Entry,
// and it contains the file path used for persistence.
// Added a mutex for safe concurrent access if needed, although current TUI might be single-threaded.
type DB struct {
	data map[string]Entry
	path string
	mu   sync.RWMutex // Read-Write mutex
}

// NewDBWithPath creates a new, empty DB using the provided path and returns a POINTER.
// If the provided path is empty, the default path is used.
func NewDBWithPath(path string) *DB { // Returns *DB
	if path == "" {
		path = DefaultDBPath()
	}
	log.Debug().Str("path", path).Msg("Initializing DB")
	// Return address of the new DB struct
	return &DB{
		data: make(map[string]Entry),
		path: path,
		// Mutex is zero-valued (unlocked) implicitly
	}
}

// LoadDB loads the JSON file from the default path into a DB, returning a POINTER.
func LoadDB() (*DB, error) { // Returns *DB
	return LoadDBWithPath("")
}

// LoadDBWithPath loads the JSON file from the specified path into a DB, returning a POINTER.
// Corrected locking strategy.
func LoadDBWithPath(path string) (*DB, error) { // Returns *DB
	db := NewDBWithPath(path) // Gets *DB pointer

	// Check file existence *before* locking
	fileInfo, err := os.Stat(db.path) // Use db.path directly
	if os.IsNotExist(err) {
		log.Info().Str("path", db.path).Msg("Local DB file not found, starting with empty DB.")
		return db, nil // File doesn't exist, return the initialized empty DB pointer
	}
	if err != nil {
		return db, fmt.Errorf("failed to stat local DB file '%s': %w", db.path, err)
	}
	if fileInfo.IsDir() {
		return db, fmt.Errorf("local DB path '%s' is a directory, not a file", db.path)
	}

	// File exists, lock for reading and decoding
	db.mu.Lock()         // Lock the pointed-to DB's mutex
	defer db.mu.Unlock() // Ensure unlock

	f, err := os.Open(db.path)
	if err != nil {
		return db, fmt.Errorf("failed to open local DB file '%s': %w", db.path, err)
	}
	defer f.Close()

	tempData := make(map[string]Entry)
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&tempData); err != nil {
		log.Error().Err(err).Str("path", db.path).Msg("Failed to decode local DB JSON. Returning empty DB state.")
		db.data = make(map[string]Entry) // Clear data on decode error
		return db, fmt.Errorf("failed to decode local DB JSON from '%s': %w", db.path, err)
	}

	// Success, assign data
	db.data = tempData // Modify the map in the pointed-to DB struct
	log.Info().Str("path", db.path).Int("entries", len(db.data)).Msg("Local DB loaded successfully.")
	return db, nil // Return the pointer to the populated DB
	// db.mu.Unlock() is called here via defer
}

// Save persists the DB contents to its file path in JSON format.
func (db *DB) Save() error {
	db.mu.RLock() // Lock for reading data
	defer db.mu.RUnlock()

	if db.path == "" {
		return fmt.Errorf("cannot save DB: path is not set")
	}

	log.Debug().Str("path", db.path).Int("entries", len(db.data)).Msg("Attempting to save local DB")

	// Create directory if it doesn't exist
	dir := filepath.Dir(db.path)
	if err := os.MkdirAll(dir, 0750); err != nil { // Use 0750 for permissions
		return fmt.Errorf("failed to create directory '%s' for local DB: %w", dir, err)
	}

	// Create or truncate the file
	f, err := os.Create(db.path)
	if err != nil {
		return fmt.Errorf("failed to create/truncate local DB file '%s': %w", db.path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // Prettify the JSON.
	if err := enc.Encode(db.data); err != nil {
		return fmt.Errorf("failed to encode local DB to JSON '%s': %w", db.path, err)
	}

	log.Info().Str("path", db.path).Msg("Local DB saved successfully.")
	return nil
}

// FilePath returns the file path where the DB is stored.
func (db *DB) FilePath() string {
	// No lock needed as path is immutable after creation
	return db.path
}

// GetEntry retrieves the full Entry record for a given test.
func (db *DB) GetEntry(testPackage, testName string) (Entry, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	key := makeKey(testPackage, testName)
	entry, exists := db.data[key]
	return entry, exists
}

// UpsertEntry updates or inserts a record with all relevant details.
// This is the preferred method for modifying DB state.
func (db *DB) UpsertEntry(testPackage, testName, jiraTicket string, skippedAt time.Time, assigneeID string) error {
	if testPackage == "" || testName == "" {
		return fmt.Errorf("cannot upsert entry with empty package or name")
	}

	db.mu.Lock() // Lock for writing
	defer db.mu.Unlock()

	key := makeKey(testPackage, testName)
	// Get existing entry or create a new one
	entry, _ := db.data[key] // We don't care about 'exists' here, just overwrite

	// Update ALL fields based on input parameters
	entry.TestPackage = testPackage // Ensure these are set even if entry was new
	entry.TestName = testName
	entry.JiraTicket = jiraTicket
	entry.AssigneeID = assigneeID
	entry.SkippedAt = skippedAt.UTC() // Always store in UTC

	// Store the updated/new entry back in the map
	db.data[key] = entry
	log.Debug().Str("key", key).Interface("entry", entry).Msg("Upserted DB entry")
	return nil
}

// RemoveEntry removes a record from the DB. Returns true if an entry was removed.
func (db *DB) RemoveEntry(testPackage, testName string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	key := makeKey(testPackage, testName)
	_, exists := db.data[key]
	if exists {
		delete(db.data, key)
		log.Debug().Str("key", key).Msg("Removed DB entry")
	}
	return exists
}

// GetAllEntries returns a slice of all entries currently in the DB.
func (db *DB) GetAllEntries() []Entry {
	db.mu.RLock()
	defer db.mu.RUnlock()
	entries := make([]Entry, 0, len(db.data))
	for _, entry := range db.data {
		entries = append(entries, entry)
	}
	return entries
}

// DefaultDBPath returns the default DB file path in the user's home directory.
// Renamed from getDefaultDBPath to be exported.
// If the user's home directory cannot be determined, it falls back to the current directory.
func DefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warn().Msg("Failed to get user home directory; using current directory for local DB.")
		home = "." // Use current directory as fallback
	}
	return filepath.Join(home, defaultDBFileName)
}

// makeKey is a helper to combine the package and test name into a single map key.
func makeKey(pkg, testName string) string {
	// Consider normalization? Lowercase? Trim spaces? For now, keep simple.
	return pkg + "::" + testName
}
