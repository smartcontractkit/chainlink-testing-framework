package localdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultDBFileName = ".flaky_test_db.json"

// Entry represents one record in the local DB for a given test.
type Entry struct {
	TestPackage string `json:"test_package"`
	TestName    string `json:"test_name"`
	JiraTicket  string `json:"jira_ticket"`
}

// DB is a simple in-memory map keyed by "pkg::testName" => Entry
type DB struct {
	data map[string]Entry
}

// NewDB returns a new, empty DB.
func NewDB() DB {
	return DB{
		data: make(map[string]Entry),
	}
}

// LoadDB loads the JSON file at ~/.flaky_test_db.json into a DB.
// If the file does not exist, an empty DB is returned.
func LoadDB() (DB, error) {
	path := getDBPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist => return empty DB
		return NewDB(), nil
	}

	f, err := os.Open(path)
	if err != nil {
		return DB{}, fmt.Errorf("failed to open local DB file: %w", err)
	}
	defer f.Close()

	db := NewDB()
	if err := json.NewDecoder(f).Decode(&db.data); err != nil {
		return DB{}, fmt.Errorf("failed to decode local DB JSON: %w", err)
	}
	return db, nil
}

// SaveDB writes the DB contents to ~/.flaky_test_db.json in JSON format.
func SaveDB(db DB) error {
	path := getDBPath()
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create local DB file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // optional prettify
	if err := enc.Encode(db.data); err != nil {
		return fmt.Errorf("failed to encode local DB to JSON: %w", err)
	}
	return nil
}

// FilePath returns the path where the local DB is stored (e.g. ~/.flaky_test_db.json).
func FilePath() string {
	return getDBPath()
}

// Get retrieves the Jira ticket ID for (testPackage, testName), if it exists.
func (db *DB) Get(testPackage, testName string) (string, bool) {
	key := makeKey(testPackage, testName)
	e, ok := db.data[key]
	if !ok {
		return "", false
	}
	return e.JiraTicket, true
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

// getDBPath returns the path to the local DB file in the user's home directory.
func getDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current dir if we cannot get home
		home = "."
	}
	return filepath.Join(home, defaultDBFileName)
}

// makeKey is a helper to combine package+testName into a single map key.
func makeKey(pkg, testName string) string {
	return pkg + "::" + testName
}
