package wmi

import "testing"

func TestQuery(t *testing.T) {
	var dst []Win32_Process
	err := Query("SELECT * FROM Win32_Process", &dst)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFieldMismatch(t *testing.T) {
	type s struct {
		Name        string
		HandleCount uint32
		Blah        uint32
	}
	var dst[]s
	err := Query("SELECT Name, HandleCount FROM Win32_Process", &dst)
	if err == nil || err.Error() != `wmi: cannot load field "Blah" into a "uint32": no such struct field` {
		t.Error("Expected err field mismatch")
	}
}
