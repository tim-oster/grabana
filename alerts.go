package grabana

import (
	"context"
	"encoding/json"
	"errors"
	alert "github.com/K-Phoen/grabana/ngalert"
	"net/http"
	"net/url"

	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/sdk"
)

// Api doc: https://editor.swagger.io/?url=https://raw.githubusercontent.com/grafana/grafana/main/pkg/services/ngalert/api/tooling/api.json

const customDashboardRefKey = "customDashboardRef"

// ErrAlertNotFound is returned when the requested alert can not be found.
var ErrAlertNotFound = errors.New("alert not found")

type alertRef struct {
	Uid       string
	Title     string
	FolderUid string
	RuleGroup string
}

// ConfigureAlertManager updates the alert manager configuration.
func (client *Client) ConfigureAlertManager(ctx context.Context, manager *alertmanager.Manager) error {
	buf, err := manager.MarshalIndentJSON()
	if err != nil {
		return err
	}

	resp, err := client.sendJSON(ctx, http.MethodPost, "/api/alertmanager/grafana/config/api/v1/alerts", buf)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusAccepted {
		return client.httpError(resp)
	}

	return nil
}

// ListAlertsForDashboard fetches a list of alerts linked to the given dashboard.
func (client *Client) ListAlertsForDashboard(ctx context.Context, dashboardUID string) ([]alertRef, error) {
	resp, err := client.get(ctx, "/api/v1/provisioning/alert-rules")
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, client.httpError(resp)
	}

	var alerts []sdk.NgAlert
	if err := decodeJSON(resp.Body, &alerts); err != nil {
		return nil, err
	}

	var refs []alertRef
	for _, a := range alerts {
		uid := a.Annotations["__dashboardUid__"]
		if uid == "" {
			uid = a.Annotations[customDashboardRefKey]
		}
		if uid != dashboardUID {
			continue
		}
		refs = append(refs, alertRef{
			Uid:       a.Uid,
			Title:     a.Title,
			FolderUid: a.FolderUID,
			RuleGroup: a.RuleGroup,
		})
	}
	return refs, nil
}

func (client *Client) UpsertAlert(ctx context.Context, alertDefinition alert.Alert, datasourcesMap map[string]string) error {
	err := alertDefinition.HookDatasource(datasourcesMap)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(alertDefinition.Builder)
	if err != nil {
		return err
	}

	var (
		method = http.MethodPost
		path   = "/api/v1/provisioning/alert-rules"
	)
	if alertDefinition.Builder.Uid != "" {
		method = http.MethodPut
		path += "/" + url.PathEscape(alertDefinition.Builder.Uid)
	}

	resp, err := client.sendJSON(ctx, method, path, buf)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return client.httpError(resp)
	}

	return nil
}

func (client *Client) DeleteAlert(ctx context.Context, uid string) error {
	resp, err := client.delete(ctx, "/api/v1/provisioning/alert-rules/"+url.PathEscape(uid))
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNoContent {
		return client.httpError(resp)
	}

	return nil
}
