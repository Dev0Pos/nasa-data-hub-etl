package database

import (
	"testing"
)

func TestValidateInitMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		want    InitMode
		wantErr bool
	}{
		{
			name:    "create mode",
			mode:    "create",
			want:    InitModeCreate,
			wantErr: false,
		},
		{
			name:    "revive mode",
			mode:    "revive",
			want:    InitModeRevive,
			wantErr: false,
		},
		{
			name:    "auto mode",
			mode:    "auto",
			want:    InitModeAuto,
			wantErr: false,
		},
		{
			name:    "case insensitive",
			mode:    "CREATE",
			want:    InitModeCreate,
			wantErr: false,
		},
		{
			name:    "invalid mode",
			mode:    "invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty mode",
			mode:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateInitMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInitMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateInitMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
