package wmi

import "testing"

func TestQuery(t *testing.T) {
	var dst []Win32_Process
	err := Query("SELECT * FROM Win32_Process", &dst)
	if err != nil {
		t.Fatal(err)
	}
}
