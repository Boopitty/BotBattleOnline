package auth_test

import (
	"testing"

	"github.com/Boopitty/BotBattleOnline/internal/auth"
)

func TestHashPass(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		want     bool // check if password is different
		wantErr  bool
	}{
		{
			name:     "Default",
			password: "abc123",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "Empty String",
			password: "",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "Special Characters",
			password: "@#${}",
			want:     true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.HashPass(tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPass() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HashPass() succeeded unexpectedly")
			}

			if (got == tt.password) && !tt.want {
				t.Errorf("HashPass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPassHash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		want     bool
		hash     string
		wantErr  bool
	}{
		{
			name:     "Default",
			password: "abc123",
			hash:     "abc123",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "Incorrect Password",
			password: "abc123",
			hash:     "123abc",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "All Empty",
			password: "",
			hash:     "",
			want:     true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := auth.HashPass(tt.hash)
			if err != nil {
				t.Fatal("failed to hash password for test")
			}

			got, gotErr := auth.CheckPassHash(tt.password, hashedPassword)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckPassHash() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckPassHash() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("CheckPassHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
