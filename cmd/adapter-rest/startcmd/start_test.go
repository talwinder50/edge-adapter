/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package startcmd

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type mockServer struct{}

func (s *mockServer) ListenAndServe(host string, handler http.Handler) error {
	return nil
}

func TestListenAndServe(t *testing.T) {
	var w HTTPServer
	err := w.ListenAndServe("wronghost", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "address wronghost: missing port in address")
}

func TestStartCmdContents(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	require.Equal(t, "start", startCmd.Use)
	require.Equal(t, "Start adapter-rest", startCmd.Short)
	require.Equal(t, "Start adapter-rest inside the edge-adapter", startCmd.Long)

	checkFlagPropertiesCorrect(t, startCmd, hostURLFlagName, hostURLFlagShorthand, hostURLFlagUsage)
}

func TestStartCmdWithBlankArg(t *testing.T) {
	t.Run("test blank host url arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		args := []string{"--" + hostURLFlagName, ""}
		startCmd.SetArgs(args)

		err := startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "host-url value is empty", err.Error())
	})
}

func TestStartCmdWithMissingArg(t *testing.T) {
	t.Run("test missing host url arg", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		err := startCmd.Execute()

		require.Error(t, err)
		require.Equal(t,
			"Neither host-url (command line flag) nor ADAPTER_REST_HOST_URL (environment variable) have been set.",
			err.Error())
	})
}

func TestStartCmdWithBlankEnvVar(t *testing.T) {
	t.Run("test blank host env var", func(t *testing.T) {
		startCmd := GetStartCmd(&mockServer{})

		err := os.Setenv(hostURLEnvKey, "")
		require.NoError(t, err)

		err = startCmd.Execute()
		require.Error(t, err)
		require.Equal(t, "ADAPTER_REST_HOST_URL value is empty", err.Error())
	})
}

func TestStartCmdValidArgs(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	args := []string{"--" + hostURLFlagName, "localhost:8080"}
	startCmd.SetArgs(args)

	err := startCmd.Execute()

	require.Nil(t, err)
}

func TestHealthCheck(t *testing.T) {
	b := &httptest.ResponseRecorder{}
	healthCheckHandler(b, nil)

	require.Equal(t, http.StatusOK, b.Code)
}

func TestHydraLoginHandler(t *testing.T) {
	r := &httptest.ResponseRecorder{}
	hydraLoginHandler(r, nil)

	require.Equal(t, http.StatusOK, r.Code)
}

func TestOidcCallbackHandler(t *testing.T) {
	r := &httptest.ResponseRecorder{}
	oidcCallbackHandler(r, nil)

	require.Equal(t, http.StatusOK, r.Code)
}

func TestHydraConsentHandler(t *testing.T) {
	r := &httptest.ResponseRecorder{}
	hydraConsentHandler(r, nil)

	require.Equal(t, http.StatusOK, r.Code)
}

func TestGetPresentationRequestHandler(t *testing.T) {
	r := &httptest.ResponseRecorder{}
	getPresentationRequestHandler(r, nil)

	require.Equal(t, http.StatusOK, r.Code)
}

func TestPresentationResponseHandler(t *testing.T) {
	r := &httptest.ResponseRecorder{}
	presentationResponseHandler(r, nil)

	require.Equal(t, http.StatusOK, r.Code)
}

func TestUserInfoHandler(t *testing.T) {
	r := &httptest.ResponseRecorder{}
	userInfoHandler(r, nil)

	require.Equal(t, http.StatusOK, r.Code)
}

func TestStartCmdValidArgsEnvVar(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	setEnvVars(t)

	defer unsetEnvVars(t)

	err := startCmd.Execute()
	require.NoError(t, err)
}

func TestTLSSystemCertPoolInvalidArgsEnvVar(t *testing.T) {
	startCmd := GetStartCmd(&mockServer{})

	setEnvVars(t)

	defer unsetEnvVars(t)
	require.NoError(t, os.Setenv(tlsSystemCertPoolEnvKey, "wrongvalue"))

	err := startCmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid syntax")
}

func TestTestResponse(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		testResponse(&stubWriter{})
	})
}

func TestUIHandler(t *testing.T) {
	t.Run("handle base path", func(t *testing.T) {
		handled := false
		uiHandler(uiEndpoint, func(_ http.ResponseWriter, _ *http.Request, path string) {
			handled = true
			require.Equal(t, uiEndpoint+"/index.html", path)
		})(nil, &http.Request{URL: &url.URL{Path: uiEndpoint}})
		require.True(t, handled)
	})
	t.Run("handle subpaths", func(t *testing.T) {
		const expected = uiEndpoint + "/css/abc123.css"
		handled := false
		uiHandler(uiEndpoint, func(_ http.ResponseWriter, _ *http.Request, path string) {
			handled = true
			require.Equal(t, expected, path)
		})(nil, &http.Request{URL: &url.URL{Path: expected}})
		require.True(t, handled)
	})
}

func setEnvVars(t *testing.T) {
	err := os.Setenv(hostURLEnvKey, "localhost:8080")
	require.NoError(t, err)
}

func unsetEnvVars(t *testing.T) {
	err := os.Unsetenv(hostURLEnvKey)
	require.NoError(t, err)
}

func checkFlagPropertiesCorrect(t *testing.T, cmd *cobra.Command, flagName, flagShorthand, flagUsage string) {
	flag := cmd.Flag(flagName)

	require.NotNil(t, flag)
	require.Equal(t, flagName, flag.Name)
	require.Equal(t, flagShorthand, flag.Shorthand)
	require.Equal(t, flagUsage, flag.Usage)
	require.Equal(t, "", flag.Value.String())

	flagAnnotations := flag.Annotations
	require.Nil(t, flagAnnotations)
}

type stubWriter struct {
}

func (s *stubWriter) Write(p []byte) (n int, err error) {
	return -1, errors.New("test")
}