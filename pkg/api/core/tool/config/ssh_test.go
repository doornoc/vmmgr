package config

import "testing"

func TestCollectConfig(t *testing.T) {
	sshHosts, err := CollectConfig(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sshHosts)
}
