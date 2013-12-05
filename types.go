package wmi

import "time"

type Win32_PerfRawData_Tcpip_NetworkInterface struct {
	BytesReceivedPerSec             uint32
	BytesSentPerSec                 uint32
	BytesTotalPerSec                uint64
	Caption                         string
	CurrentBandwidth                uint32
	Description                     string
	Frequency_Object                uint64
	Frequency_PerfTime              uint64
	Frequency_Sys100NS              uint64
	Name                            string
	OutputQueueLength               uint32
	PacketsOutboundDiscarded        uint32
	PacketsOutboundErrors           uint32
	PacketsPerSec                   uint32
	PacketsReceivedDiscarded        uint32
	PacketsReceivedErrors           uint32
	PacketsReceivedNonUnicastPerSec uint32
	PacketsReceivedPerSec           uint32
	PacketsReceivedUnicastPerSec    uint32
	PacketsReceivedUnknown          uint32
	PacketsSentNonUnicastPerSec     uint32
	PacketsSentPerSec               uint32
	PacketsSentUnicastPerSec        uint32
	Timestamp_Object                uint64
	Timestamp_PerfTime              uint64
	Timestamp_Sys100NS              uint64
}

type Win32_PerfRawData_PerfOS_Processor struct {
	C1TransitionsPerSec   uint64
	C2TransitionsPerSec   uint64
	C3TransitionsPerSec   uint64
	Caption               string
	DPCRate               uint32
	DPCsQueuedPerSec      uint32
	Description           string
	Frequency_Object      uint64
	Frequency_PerfTime    uint64
	Frequency_Sys100NS    uint64
	InterruptsPerSec      uint32
	Name                  string
	PercentC1Time         uint64
	PercentC2Time         uint64
	PercentC3Time         uint64
	PercentDPCTime        uint64
	PercentIdleTime       uint64
	PercentInterruptTime  uint64
	PercentPrivilegedTime uint64
	PercentProcessorTime  uint64
	PercentUserTime       uint64
	Timestamp_Object      uint64
	Timestamp_PerfTime    uint64
	Timestamp_Sys100NS    uint64
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
