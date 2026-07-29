package database_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Boopitty/BotBattleOnline/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func findEnvFile() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	dir := filepath.Dir(filename)
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// Initializes database for testing
func setupTestDB(t *testing.T) database.DBTX {
	t.Helper()

	envPath := findEnvFile()
	if envPath == "" {
		t.Skip("Could not find .env file")
	}

	err := godotenv.Load(envPath)
	if err != nil {
		t.Skipf("Could not load godotenv from %s: %v", envPath, err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		t.Skip("DB_URL is not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return conn
}

func TestQueries_UserCRUD(t *testing.T) {
	id := uuid.New()
	testTime := time.Now()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db database.DBTX
		// Named input parameters for target function.
		arg         database.CreateUserParams
		wantMade    database.User
		wantMadeErr bool
		wantGot     database.User
		wantGotErr  bool
		wantDelErr  bool
	}{
		{
			name: "Default",
			db:   setupTestDB(t),
			arg: database.CreateUserParams{
				ID:             id,
				Username:       "example",
				HashedPassword: "1234",
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
			},
			wantMade: database.User{
				ID:             id,
				Username:       "example",
				HashedPassword: "1234",
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
			},
			wantMadeErr: false,
			wantGot: database.User{
				ID:             id,
				Username:       "example",
				HashedPassword: "1234",
				CreatedAt:      testTime,
				UpdatedAt:      testTime,
			},
			wantGotErr: false,
			wantDelErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := database.New(tt.db)

			// Test CreateUser()
			made, madeErr := q.CreateUser(t.Context(), tt.arg)
			if madeErr != nil {
				if !tt.wantMadeErr {
					t.Errorf("CreateUser() failed: %v", madeErr)
				}
				return
			}
			if tt.wantMadeErr {
				t.Fatal("CreateUser() succeeded unexpectedly")
			}

			if made.ID != tt.wantMade.ID {
				t.Fatalf("CreateUser() = %v, want %v", made, tt.wantMade)
			}

			// Test GetUser()
			got, gotErr := q.GetUser(t.Context(), tt.arg.Username)
			if gotErr != nil {
				if !tt.wantGotErr {
					t.Errorf("GetUser() failed: %v", gotErr)
				}
				return
			}
			if tt.wantGotErr {
				t.Fatal("GetUser() succeeded unexpectedly")
			}

			if got.ID != tt.wantGot.ID ||
				got.Username != tt.wantGot.Username ||
				got.HashedPassword != tt.wantGot.HashedPassword {
				t.Errorf("GetUser() = %v, want %v", got, tt.wantGot) // Error, not Fatal to continue with test to attempt to delete item
			}

			// Test DeleteUser()
			delErr := q.DeleteUser(t.Context(), tt.arg.ID)
			if delErr != nil {
				if !tt.wantDelErr {
					t.Errorf("DeleteUser() failed: %v", delErr)
				}
				return
			}
			if tt.wantDelErr {
				t.Fatalf("DeleteUser succeeded unexectedly")
			}
		})
	}
}
