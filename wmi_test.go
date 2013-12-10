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
	var dst []s
	err := Query("SELECT Name, HandleCount FROM Win32_Process", &dst)
	if err == nil || err.Error() != `wmi: cannot load field "Blah" into a "uint32": no such struct field` {
		t.Error("Expected err field mismatch")
	}
}

type Win32_Process struct {
	CSCreationClassName        string
	CSName                     string
	Caption                    string
	CommandLine                string
	CreationClassName          string
	CreationDate               time.Time
	Description                string
	ExecutablePath             string
	ExecutionState             uint16
	Handle                     string
	HandleCount                uint32
	InstallDate                time.Time
	KernelModeTime             uint64
	MaximumWorkingSetSize      uint32
	MinimumWorkingSetSize      uint32
	Name                       string
	OSCreationClassName        string
	OSName                     string
	OtherOperationCount        uint64
	OtherTransferCount         uint64
	PageFaults                 uint32
	PageFileUsage              uint32
	ParentProcessId            uint32
	PeakPageFileUsage          uint32
	PeakVirtualSize            uint64
	PeakWorkingSetSize         uint32
	Priority                   uint32
	PrivatePageCount           uint64
	ProcessId                  uint32
	QuotaNonPagedPoolUsage     uint32
	QuotaPagedPoolUsage        uint32
	QuotaPeakNonPagedPoolUsage uint32
	QuotaPeakPagedPoolUsage    uint32
	ReadOperationCount         uint64
	ReadTransferCount          uint64
	SessionId                  uint32
	Status                     string
	TerminationDate            time.Time
	ThreadCount                uint32
	UserModeTime               uint64
	VirtualSize                uint64
	WindowsVersion             string
	WorkingSetSize             uint64
	WriteOperationCount        uint64
	WriteTransferCount         uint64
}
