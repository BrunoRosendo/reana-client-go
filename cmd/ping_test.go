package cmd

import (
	"net/http"
	"reanahub/reana-client-go/utils"
	"testing"

	"github.com/spf13/viper"
)

var pingServerPath = "/api/you"

func TestPing(t *testing.T) {
	serverResponse := `{
				"email": "john.doe@example.org",
				"reana_server_version": "0.9.0a5"
			}`
	expected := []string{
		"REANA server version: 0.9.0a5",
		"Authenticated as: <john.doe@example.org>",
	}
	testCmdRun(t, "ping", pingServerPath, serverResponse, http.StatusOK, expected)
}

func TestUnreachableServer(t *testing.T) {
	viper.Set("server-url", "https://unreachable.invalid")
	t.Cleanup(func() {
		viper.Reset()
	})

	rootCmd := NewRootCmd()
	output, err := utils.ExecuteCommand(rootCmd, "ping", "-t", "1234")

	if err == nil {
		t.Errorf("Expected an error, instead got '%s'", output)
	}

	expectedErr := "'https://unreachable.invalid' not found, please verify the provided server URL or check your internet connection"
	if utils.HandleApiError(err).Error() != expectedErr {
		t.Errorf("Expected server not found error, instead got '%s'", err.Error())
	}
}
