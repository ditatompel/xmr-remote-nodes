package monero

import (
	"testing"
)

// TODO: Add database test table and then clean it up

func TestProberRepo_CheckApi(t *testing.T) {
	if !testMySQL {
		t.Skip("Skip test, not connected to database")
	}
	tests := []struct {
		name    string
		apiKey  string
		want    Prober
		wantErr bool
	}{
		{
			name:    "Empty key",
			apiKey:  "",
			want:    Prober{},
			wantErr: true,
		},
		{
			name:    "Invalid key",
			apiKey:  "invalid",
			want:    Prober{},
			wantErr: true,
		},
	}

	repo := NewProber()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CheckApi(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProberRepo.CheckApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func BenchmarkProberRepo_CheckApi(b *testing.B) {
	if !testMySQL {
		b.Skip("Skip bench, not connected to database")
	}
	repo := NewProber()
	for i := 0; i < b.N; i++ {
		repo.CheckApi("")
	}
}
