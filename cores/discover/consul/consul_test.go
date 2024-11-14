package consul

import (
	"testing"
)

var (
	client *Client
)

func TestMain(m *testing.M) {
	client = GetConsulClient()
	m.Run()
}

func TestClient_GetAddressWithTag(t *testing.T) {
	type args struct {
		service string
		tag     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test-with-no-tag", args{"openai", ""}, "10.211.55.5:29915"},
		{"test-with-tag", args{"openai", "v1.0.1"}, "10.211.55.5:29915"},
		{"test-with-ip", args{"10.11.43.113", ""}, "10.11.43.113"},
		{"test-with-ip-ports", args{"10.11.43.113:28080", ""}, "10.11.43.113:28080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetAddressWithTag(tt.args.service, tt.args.tag); got != tt.want {
				t.Errorf("Client.GetAddressWithTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
