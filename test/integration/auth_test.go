package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ticketer/test/integration/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
	// 1. Setup DB and App
	db := testutils.SetupTestDB(t)
	app := testutils.SetupTestApp(t, db)

	// Need to register auth.AuthMiddleware for /auth/me manually or wrap the mux
	// Actually, app.Mux already has /auth/me registered with the middleware inside auth.Handler.RegisterRoutes!
	
	// 2. Simulate Google OAuth Callback using internal service
	// We just call the HandleCallback directly to skip the actual external Google redirect
	// For HandleCallback to work, it needs to exchange a "code". 
	// But HandleCallback internally calls google oauth2 service. Wait! 
	// If it does, we can't test it directly unless we mock it.
	// Let's test the endpoint /auth/me by manually creating a user and a token.
	
	// Instead, let's create a user and generate a token manually for testing endpoints.
	ctx := context.Background()
	
	// Insert user directly via db or service
	// Since HandleCallback is tightly coupled to Google, let's just make a user in DB
	// and use JWT signing directly. Wait, HandleCallback is in Service, so maybe
	// we just insert into db.
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, name, role, oauth_provider)
		VALUES ('11111111-1111-1111-1111-111111111111', 'test@example.com', 'Test User', 'admin', 'google')
	`)
	require.NoError(t, err)

	// The app uses the JWT secret from env or default
	// To test /auth/me we need a token. 
	// Wait, we can test /auth/logout!
	
	// Test Unauthorized Logout
	req, err := http.NewRequest(http.MethodPost, "/auth/logout", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	app.Mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Test Authorized Logout
	cookie := testutils.GenerateTestToken("11111111-1111-1111-1111-111111111111", "admin")
	reqAuth, err := http.NewRequest(http.MethodPost, "/auth/logout", nil)
	require.NoError(t, err)
	reqAuth.AddCookie(cookie)

	rrAuth := httptest.NewRecorder()
	app.Mux.ServeHTTP(rrAuth, reqAuth)
	assert.Equal(t, http.StatusOK, rrAuth.Code)

	var resp map[string]string
	err = json.Unmarshal(rrAuth.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "logged out", resp["message"])
}
