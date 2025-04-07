// mapping/mapping.go
package mapping

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"
)

// UserMapping defines the structure for mapping Jira User IDs to Pillar Names.
type UserMapping struct {
	JiraUserID string `json:"jira_user_id"`
	PillarName string `json:"pillar_name"`
	// Add other relevant user details if needed (e.g., Slack ID, Name)
}

// UserTestMapping defines the structure for mapping test patterns to Jira User IDs.
type UserTestMapping struct {
	JiraUserID string `json:"jira_user_id"`
	Pattern    string `json:"pattern"`
}

// UserTestMappingWithRegex holds the original mapping and its compiled regex.
type UserTestMappingWithRegex struct {
	UserTestMapping
	CompiledRegex *regexp.Regexp
}

// LoadUserMappings reads the user mapping file and returns a map for easy lookup.
func LoadUserMappings(path string) (map[string]UserMapping, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// It might be acceptable for this file to be optional, return empty map
		if os.IsNotExist(err) {
			log.Warn().Str("path", path).Msg("User mapping file not found, proceeding without it.")
			return make(map[string]UserMapping), nil
		}
		return nil, fmt.Errorf("failed to read user mapping file '%s': %w", path, err)
	}

	var mappings []UserMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user mapping file '%s': %w", path, err)
	}

	// Convert array to map for efficient lookup
	userMap := make(map[string]UserMapping)
	for _, user := range mappings {
		userMap[user.JiraUserID] = user
	}
	log.Info().Str("path", path).Int("count", len(userMap)).Msg("Loaded user mappings")
	return userMap, nil
}

// LoadUserTestMappings reads the user test mapping file, compiles regex patterns,
// and returns a slice of mappings with their compiled regexes.
func LoadUserTestMappings(path string) ([]UserTestMappingWithRegex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// It might be acceptable for this file to be optional, return empty slice
		if os.IsNotExist(err) {
			log.Warn().Str("path", path).Msg("User test mapping file not found, proceeding without auto-assignment.")
			return []UserTestMappingWithRegex{}, nil
		}
		return nil, fmt.Errorf("failed to read user test mapping file '%s': %w", path, err)
	}

	var mappings []UserTestMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user test mapping file '%s': %w", path, err)
	}

	compiledMappings := make([]UserTestMappingWithRegex, 0, len(mappings))
	for _, m := range mappings {
		re, err := regexp.Compile(m.Pattern)
		if err != nil {
			// Log error but continue processing other patterns
			log.Error().Err(err).Str("pattern", m.Pattern).Msg("Failed to compile regex pattern, skipping this mapping")
			continue
		}
		compiledMappings = append(compiledMappings, UserTestMappingWithRegex{
			UserTestMapping: m,
			CompiledRegex:   re,
		})
	}

	log.Info().Str("path", path).Int("count", len(compiledMappings)).Msg("Loaded and compiled user test mappings")
	return compiledMappings, nil
}

// FindAssigneeIDForTest iterates through the compiled test mappings and returns the
// Jira User ID from the *first* matching pattern. The order in the mapping file matters.
// It typically matches against the testPackage.
func FindAssigneeIDForTest(testPath string, compiledMappings []UserTestMappingWithRegex) (string, error) {
	if testPath == "" {
		return "", fmt.Errorf("testPath cannot be empty")
	}
	for _, compiledMapping := range compiledMappings {
		if compiledMapping.CompiledRegex.MatchString(testPath) {
			log.Debug().Str("testPath", testPath).Str("pattern", compiledMapping.Pattern).Str("assignee", compiledMapping.JiraUserID).Msg("Found matching assignee pattern")
			return compiledMapping.JiraUserID, nil // Return the first match
		}
	}
	log.Debug().Str("testPath", testPath).Msg("No matching assignee pattern found")
	return "", nil // No match found
}
