package linkingservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AdminClient struct {
	AdminURL string
}

func NewAdminClient(adminURL string) *AdminClient {
	if strings.TrimSpace(adminURL) == "" {
		return nil
	}

	return &AdminClient{AdminURL: adminURL}
}

func NewAdminClientFromOutput(out *Output) *AdminClient {
	if out == nil {
		return nil
	}

	return NewAdminClient(out.LocalAdminURL)
}

func (c *AdminClient) SetOwnerOrg(ctx context.Context, owner, orgID string) error {
	if c == nil || strings.TrimSpace(c.AdminURL) == "" {
		return fmt.Errorf("linking service admin URL is not configured")
	}

	payload, err := json.Marshal(map[string]string{
		"workflowOwner": owner,
		"orgID":         orgID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal linking service admin request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		strings.TrimSuffix(c.AdminURL, "/")+"/admin/link",
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to create linking service admin request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call linking service admin endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected linking service admin status: %d", resp.StatusCode)
	}

	return nil
}
