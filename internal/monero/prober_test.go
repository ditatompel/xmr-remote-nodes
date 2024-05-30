package monero

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"xmr-remote-nodes/internal/config"
	"xmr-remote-nodes/internal/database"
)

var testMySQL = true

func init() {
	// load test db config from OS environment variable
	//
	// Example:
	// TEST_DB_HOST=127.0.0.1 \
	// TEST_DB_PORT=3306 \
	// TEST_DB_USER=testuser \
	// TEST_DB_PASSWORD=testpass \
	// TEST_DB_NAME=testdb go test ./... -v
	//
	// To run benchmark only, add `-bench=. -run=^#` to the `go test` command
	config.DBCfg().Host = os.Getenv("TEST_DB_HOST")
	config.DBCfg().Port, _ = strconv.Atoi(os.Getenv("TEST_DB_PORT"))
	config.DBCfg().User = os.Getenv("TEST_DB_USER")
	config.DBCfg().Password = os.Getenv("TEST_DB_PASSWORD")
	config.DBCfg().Name = os.Getenv("TEST_DB_NAME")

	if err := database.ConnectDB(); err != nil {
		testMySQL = false
	}
}

// TODO: Add database test table and then clean it up

func TestProberRepo_CheckApi(t *testing.T) {
	if !testMySQL {
		fmt.Println("Skip test, not connected to database")
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
		fmt.Println("Skip bench, not connected to database")
		b.Skip("Skip bench, not connected to database")
	}
	repo := NewProber()
	for i := 0; i < b.N; i++ {
		repo.CheckApi("")
	}
}
