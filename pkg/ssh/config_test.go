package ssh

import (
	"fmt"
	"testing"
)

func TestParseConfigTarget(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{"localhost", "localhost"},
		{"user@192.168.1.1", "user@192.168.1.1"},
		{"user@192.168.1.1@2222", "user@192.168.1.1 -p 2222"},
		{"user@192.168.1.1@2222@/home/user/.ssh/id_rsa", "user@192.168.1.1 -p 2222 -i /home/user/.ssh/id_rsa"},
	}

	for _, test := range tests {
		got := ParseConfig(test.host).Target()
		if got != test.want {
			t.Errorf("%s = %s, want %s", fmt.Sprintf("ParseConfig(%s).Target()", test.host), got, test.host)
		}
	}
}
