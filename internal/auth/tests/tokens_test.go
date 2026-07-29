package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Boopitty/BotBattleOnline/internal/auth"
	"github.com/google/uuid"
)

func TestCreateGetValidateJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID       uuid.UUID
		tokenSecret  string
		expiresIn    time.Duration
		wantMadeErr  bool
		headers      http.Header
		wantGotErr   bool
		wantValidErr bool
	}{
		{
			name:         "default",
			userID:       uuid.New(),
			tokenSecret:  "SuperSecret",
			expiresIn:    time.Duration(60) * time.Second,
			wantMadeErr:  false,
			headers:      make(http.Header),
			wantGotErr:   false,
			wantValidErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// MakeJWT()
			made, madeErr := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if madeErr != nil {
				if !tt.wantMadeErr {
					t.Errorf("MakeJWT() failed: %v", madeErr)
				}
				return
			}
			if tt.wantMadeErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}

			if made == "" {
				t.Fatalf("MakeJWT() = %v", made)
			}

			// GetBearerToken()
			tt.headers.Add("Authorization", made)
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantGotErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantGotErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}

			// ValidateJWT()
			valid, validErr := auth.ValidateJWT(got, tt.tokenSecret)
			if validErr != nil {
				if !tt.wantValidErr {
					t.Errorf("ValidateJWT() failed: %v", validErr)
				}
				return
			}
			if tt.wantValidErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}

			if valid != tt.userID {
				t.Errorf("ValidateJWT() = %v, want %v", valid, tt.userID)
			}
		})
	}
}
