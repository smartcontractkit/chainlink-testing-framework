package cmd

import "github.com/spf13/cobra"

// UserMapping holds user information.
type UserMapping struct {
	JiraUserID string `json:"jira_user_id"`
	UserName   string `json:"user_name"`
	PillarName string `json:"pillar_name"`
}

// UserTestMapping holds a single test pattern for a user.
type UserTestMapping struct {
	JiraUserID string `json:"jira_user_id"`
	Pattern    string `json:"pattern"`
}

// Shared command flags
var (
	userMappingPath     string // path to user mapping file
	userTestMappingPath string // path to user test mapping file
)

// InitCommonFlags initializes the common flags for both commands.
func InitCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&userMappingPath, "user-mapping", "user_mapping.json", "Path to the user mapping JSON file")
	cmd.Flags().StringVar(&userTestMappingPath, "user-test-mapping", "user_test_mapping.json", "Path to the user test mapping JSON file")
}
