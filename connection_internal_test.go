package sqlitecloud

import (
	"reflect"
	"strings"
	"testing"
)

func TestConnectionCommandsAuthPreference(t *testing.T) {
	tests := []struct {
		name     string
		config   SQCloudConfig
		wantCmd  string
		wantArgs []interface{}
	}{
		{"api key wins", SQCloudConfig{ApiKey: "k", Token: "t", Username: "u", Password: "p"}, "AUTH APIKEY ?;", []interface{}{"k"}},
		{"token next", SQCloudConfig{Token: "t", Username: "u", Password: "p"}, "AUTH TOKEN ?;", []interface{}{"t"}},
		{"user/pass fallback", SQCloudConfig{Username: "u", Password: "p"}, "AUTH USER ? PASSWORD ?;", []interface{}{"u", "p"}},
		{"hashed password", SQCloudConfig{Username: "u", Password: "hash", PasswordHashed: true}, "AUTH USER ? HASH ?;", []interface{}{"u", "hash"}},
		{"no credentials", SQCloudConfig{}, "", []interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuf, gotArgs := connectionCommands(tt.config)
			if tt.wantCmd == "" && gotBuf != "" {
				t.Fatalf("expected no auth command, got %q", gotBuf)
			}
			if tt.wantCmd != "" && !strings.Contains(gotBuf, tt.wantCmd) {
				t.Fatalf("expected %q in buffer, got %q", tt.wantCmd, gotBuf)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Fatalf("args mismatch: want %v got %v", tt.wantArgs, gotArgs)
			}
		})
	}
}
