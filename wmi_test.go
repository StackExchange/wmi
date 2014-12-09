package wmi

import (
	"testing"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

//23456789 123456789 123456789 123456789 123456789 123456789 123456789 123456789
func TestInvokeProcess(t *testing.T) {
	var err error
	var args []interface{}
	/*err := InvokeMethod("fakeclass", "fakemethod", nil)
			if err == nil {
				t.Fatal("Expected error when invoking fakemethod on fakeclass")
			}

			err := InvokeMethod("StdRegProv", "GetStringValue", nil)
			if err != nil {
			t.Fatal(err)
		}

		//Powershell testing
		#gwmi StdRegProv -List | select -ExpandProperty Methods
		$HKCU = [uint32]'0x80000001'
		$HKLM = [uint32]'0x80000002'
		$classREG = [wmiclass]"StdRegProv"
		#$classREG.EnumValues($HKLM,"System\CurrentControlSet\Control\Session Manager\Environment")
		$classREG.GetStringValue($HKLM,"System\CurrentControlSet\Control\Session Manager\Environment","Path")


		#gwmi Win32_Process -List | select -ExpandProperty Methods
		$classPS = [wmiclass]"Win32_Process"
		$result = $classPS.Create("cmd.exe",$null,$null)
		$result | gm *

		$classSVC = gwmi -class Win32_Service -Filter "Name='WSearch'"
		$classSVC.StartService() | GM *

	  ////using new CIM methods in WMF 4/5
		//$dsclcm = Get-CimClass -Namespace ROOT\Microsoft\Windows\DesiredStateConfiguration MSFT_DSCLocalConfigurationManager
		//icim -CimClass $dsclcm -MethodName GetConfigurationStatus -ComputerName localhost | select -ExpandProperty ConfigurationStatus

	  //using wmi methods
		$dsc = gwmi *DSC* -list -Namespace root\Microsoft\Windows\DesiredStateConfiguration
		$dsclcm = @($dsc | ?{$_.Name -eq 'MSFT_DSCLocalConfigurationManager'})[0]
		$dsclcm | gm * -force

		$result = $dsclcm.GetMetaConfiguration()
		$result.MetaConfiguration

		$result = $dsclcm.GetConfigurationStatus(0)
		$result.ConfigurationStatus

	*/

	//NOTE: These require running 'go test' as an administrator, otherwise you will get access denied errors.

	var ReturnValue int = -100 //see http://msdn.microsoft.com/en-us/library/aa393673(v=vs.85).aspx
	/*
		args = []interface{}{}
		err = InvokeMethod("Win32_Service.Name='Audiosrv'", "StopService", args, &ReturnValue)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("StopService Result: %v", ReturnValue)

		time.Sleep(2 * time.Second) //Increase delay to confirm service is stopped

		ReturnValue = -100 //see http://msdn.microsoft.com/en-us/library/aa393660(v=vs.85).aspx
		args = []interface{}{}
		err = InvokeMethod("Win32_Service.Name='Audiosrv'", "StartService", args, &ReturnValue)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("StartService Result: %v", ReturnValue)

		ReturnValue = -100 //see http://msdn.microsoft.com/en-us/library/aa389388(v=vs.85).aspx
		var ProcessID int
		//todo: replace nil with Win32_ProcessStartup class http://msdn.microsoft.com/en-us/library/aa394375(v=vs.85).aspx
		//also find a way to pass parameters to program. ex: cmd.exe /C echo hello
		args = []interface{}{"cmd.exe", "c:\\", nil, &ProcessID}
		err = InvokeMethod("Win32_Process", "Create", args, &ReturnValue)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Create Process Result: %v ProcessID: %v", ReturnValue, ProcessID)

		time.Sleep(2 * time.Second)

		ReturnValue = -100     //see http://msdn.microsoft.com/en-us/library/aa393907(v=vs.85).aspx
		args = []interface{}{} //this works
		//args = []interface{}{uint(1)} //this throws exception. exit code returned by terminated process
		err = InvokeMethod(fmt.Sprintf("Win32_Process.Handle=%v", ProcessID), "Terminate", args, &ReturnValue)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Terminate ProcessID %v Result: %v", ProcessID, ReturnValue)

	*/

	//DSC Testing: requires WMF 5.0 preview. Test on ny-rdp01 or using above ps1 tests
	//See MSFT_DSCLocalConfigurationManager mof file at "C:\Windows\System32\wbem\DscCore.mof" on ny-rdp01 for method signatures
	//also http://www.codeproject.com/Articles/27914/WMI-MOF-Basics and wbemtest.exe see http://technet.microsoft.com/en-us/library/ee692770.aspx
	ReturnValue = -100 //see http://msdn.microsoft.com/en-us/library/dn469252(v=vs.85).aspx
	//[implemented,static : ToSubclass DisableOverride] uint32 GetMetaConfiguration([Out,EmbeddedInstance("MSFT_DSCMetaConfiguration") : ToSubclass] MSFT_DSCMetaConfiguration MetaConfiguration);
	var MetaConfigurationRAW ole.VARIANT        //This may need to be an instance of MSFT_DSCMetaConfiguration instead of string pointer
	args = []interface{}{&MetaConfigurationRAW} //return is EmbeddedInstance("MSFT_DSCMetaConfiguration") see http://msdn.microsoft.com/en-us/library/aa390377(v=vs.85).aspx
	err = InvokeMethod("MSFT_DSCLocalConfigurationManager", "GetMetaConfiguration", args, &ReturnValue, nil, `root\Microsoft\Windows\DesiredStateConfiguration`)
	if err != nil {
		t.Fatal(err)
	}
	MetaConfiguration := MetaConfigurationRAW.ToIDispatch() //NOTE: this was throwing "Object is not connected to server." until I commented out ole.CoUninitialize() in InvokeMethod. Need to call ole.CoInitializeEx outside InvokeMethod unless the variant type is converted to a native struct
	FrequencyMins, err := oleutil.GetProperty(MetaConfiguration, "ConfigurationModeFrequencyMins")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("MSFT_DSCLocalConfigurationManager GetMetaConfiguration Result: %v FrequencyMins: %v", ReturnValue, FrequencyMins.Value())

	ReturnValue = -100 //see http://msdn.microsoft.com/en-us/library/aa393907(v=vs.85).aspx
	var resultRAW ole.VARIANT
	//[implemented,static : ToSubclass DisableOverride] uint32 GetConfigurationStatus([In] uint32 Flags,[Out,EmbeddedInstance("MSFT_DSCConfigurationStatus") : ToSubclass] MSFT_DSCConfigurationStatus ConfigurationStatus[]);
	args = []interface{}{nil, &resultRAW} //nil is input flags, but uint(0) doesnt work. 0 is flag for previous, 1 = ?? 2 = ??
	err = InvokeMethod("MSFT_DSCLocalConfigurationManager", "GetConfigurationStatus", args, &ReturnValue, nil, `root\Microsoft\Windows\DesiredStateConfiguration`)
	if err != nil {
		t.Fatal(err)
	}
	//RawResult VT = VT_ARRAY | VT_VARIANT
	resultRAWArray := resultRAW.ToArray() //result is an Array... we need the first result in the array converted to an IDispatch
	resultRAWArrayType, _ := resultRAWArray.GetType()
	resultRAWArrayDimensions, _ := resultRAWArray.GetDimensions()
	resultRAWArraySize, _ := resultRAWArray.GetSize()
	resultRAWArrayTotalElements, _ := resultRAWArray.TotalElements(0)
	t.Logf("GetConfigurationStatus\nresultRAW: %#v \nArrayType: %#v \nArrayDim: %#v \nArraySize: %#v \nArrayTotal: %#v", resultRAW, resultRAWArrayType, resultRAWArrayDimensions, resultRAWArraySize, resultRAWArrayTotalElements)
	var result []ole.VARIANT
	err = resultRAWArray.ToVariantArray(&result)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v cool", "we")
	t.Logf("GetConfigurationStatus Result len %v cap %v", len(result), cap(result))
	itemRAW := result[0]
	t.Logf("GetConfigurationStatus itemRAW %#v", itemRAW)
	item := itemRAW.ToIDispatch()
	t.Logf("GetConfigurationStatus item %#v", item)
	NumberOfResources, err := oleutil.GetProperty(item, "NumberOfResources") //item is wmi class MSFT_DSCConfigurationStatus
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("MSFT_DSCLocalConfigurationManager GetConfigurationStatus ReturnValue: %v NumberOfResources: %#v", ReturnValue, NumberOfResources.Value())

	t.FailNow() //debuging: force showing of t.LogF messages while testing
}

type Win32_Process_Create struct {
	ProcessId   uint32
	ReturnValue uint32
}

type MSFT_DSCConfigurationStatus struct {
}

/* Disabled original tests and benchmarks to speed things up

func TestQuery(t *testing.T) {
	var dst []Win32_Process
	q := CreateQuery(&dst, "")
	err := Query(q, &dst)
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

func TestStrings(t *testing.T) {
	printed := false
	f := func() {
		var dst []Win32_Process
		zeros := 0
		q := CreateQuery(&dst, "")
		for i := 0; i < 5; i++ {
			err := Query(q, &dst)
			if err != nil {
				t.Fatal(err, q)
			}
			for _, d := range dst {
				v := reflect.ValueOf(d)
				for j := 0; j < v.NumField(); j++ {
					f := v.Field(j)
					if f.Kind() != reflect.String {
						continue
					}
					s := f.Interface().(string)
					if len(s) > 0 && s[0] == '\u0000' {
						zeros++
						if !printed {
							printed = true
							j, _ := json.MarshalIndent(&d, "", "  ")
							t.Log("Example with \\u0000:\n", string(j))
						}
					}
				}
			}
			fmt.Println("iter", i, "zeros:", zeros)
		}
		if zeros > 0 {
			t.Error("> 0 zeros")
		}
	}

	fmt.Println("Disabling GC")
	debug.SetGCPercent(-1)
	f()
	fmt.Println("Enabling GC")
	debug.SetGCPercent(100)
	f()
}

func TestNamespace(t *testing.T) {
	var dst []Win32_Process
	q := CreateQuery(&dst, "")
	err := QueryNamespace(q, &dst, `root\CIMV2`)
	if err != nil {
		t.Fatal(err)
	}
	dst = nil
	err = QueryNamespace(q, &dst, `broken\nothing`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateQuery(t *testing.T) {
	type TestStruct struct {
		Name  string
		Count int
	}
	var dst []TestStruct
	output := "SELECT Name, Count FROM TestStruct WHERE Count > 2"
	tests := []interface{}{
		&dst,
		dst,
		TestStruct{},
		&TestStruct{},
	}
	for i, test := range tests {
		if o := CreateQuery(test, "WHERE Count > 2"); o != output {
			t.Error("bad output on", i, o)
		}
	}
	if CreateQuery(3, "") != "" {
		t.Error("expected empty string")
	}
}

func _TestMany(t *testing.T) {
	limit := 5000
	fmt.Println("running until:", limit)
	fmt.Println("No panics mean it succeeded. Other errors are OK.")
	runtime.GOMAXPROCS(2)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < limit; i++ {
			if i%25 == 0 {
				fmt.Println(i)
			}
			var dst []Win32_PerfRawData_PerfDisk_LogicalDisk
			q := CreateQuery(&dst, "")
			err := Query(q, &dst)
			if err != nil {
				fmt.Println("ERROR disk", err)
			}
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i > -limit; i-- {
			if i%25 == 0 {
				fmt.Println(i)
			}
			var dst []Win32_OperatingSystem
			q := CreateQuery(&dst, "")
			err := Query(q, &dst)
			if err != nil {
				fmt.Println("ERROR OS", err)
			}
		}
		wg.Done()
	}()
	wg.Wait()
}


type Win32_Process struct {
	CSCreationClassName        string
	CSName                     string
	Caption                    string
	CommandLine                *string
	CreationClassName          string
	CreationDate               *time.Time
	Description                string
	ExecutablePath             *string
	ExecutionState             *uint16
	Handle                     string
	HandleCount                uint32
	InstallDate                *time.Time
	KernelModeTime             uint64
	MaximumWorkingSetSize      *uint32
	MinimumWorkingSetSize      *uint32
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
	Status                     *string
	TerminationDate            *time.Time
	ThreadCount                uint32
	UserModeTime               uint64
	VirtualSize                uint64
	WindowsVersion             string
	WorkingSetSize             uint64
	WriteOperationCount        uint64
	WriteTransferCount         uint64
}

type Win32_PerfRawData_PerfDisk_LogicalDisk struct {
	AvgDiskBytesPerRead          uint64
	AvgDiskBytesPerRead_Base     uint32
	AvgDiskBytesPerTransfer      uint64
	AvgDiskBytesPerTransfer_Base uint32
	AvgDiskBytesPerWrite         uint64
	AvgDiskBytesPerWrite_Base    uint32
	AvgDiskQueueLength           uint64
	AvgDiskReadQueueLength       uint64
	AvgDiskSecPerRead            uint32
	AvgDiskSecPerRead_Base       uint32
	AvgDiskSecPerTransfer        uint32
	AvgDiskSecPerTransfer_Base   uint32
	AvgDiskSecPerWrite           uint32
	AvgDiskSecPerWrite_Base      uint32
	AvgDiskWriteQueueLength      uint64
	Caption                      string
	CurrentDiskQueueLength       uint32
	Description                  string
	DiskBytesPerSec              uint64
	DiskReadBytesPerSec          uint64
	DiskReadsPerSec              uint32
	DiskTransfersPerSec          uint32
	DiskWriteBytesPerSec         uint64
	DiskWritesPerSec             uint32
	FreeMegabytes                uint32
	Frequency_Object             uint64
	Frequency_PerfTime           uint64
	Frequency_Sys100NS           uint64
	Name                         string
	PercentDiskReadTime          uint64
	PercentDiskReadTime_Base     uint64
	PercentDiskTime              uint64
	PercentDiskTime_Base         uint64
	PercentDiskWriteTime         uint64
	PercentDiskWriteTime_Base    uint64
	PercentFreeSpace             uint32
	PercentFreeSpace_Base        uint32
	PercentIdleTime              uint64
	PercentIdleTime_Base         uint64
	SplitIOPerSec                uint32
	Timestamp_Object             uint64
	Timestamp_PerfTime           uint64
	Timestamp_Sys100NS           uint64
}

type Win32_OperatingSystem struct {
	BootDevice                                string
	BuildNumber                               string
	BuildType                                 string
	Caption                                   string
	CodeSet                                   string
	CountryCode                               string
	CreationClassName                         string
	CSCreationClassName                       string
	CSDVersion                                string
	CSName                                    string
	CurrentTimeZone                           int16
	DataExecutionPrevention_Available         bool
	DataExecutionPrevention_32BitApplications bool
	DataExecutionPrevention_Drivers           bool
	DataExecutionPrevention_SupportPolicy     uint8
	Debug                                     bool
	Description                               string
	Distributed                               bool
	EncryptionLevel                           uint32
	ForegroundApplicationBoost                uint8
	FreePhysicalMemory                        uint64
	FreeSpaceInPagingFiles                    uint64
	FreeVirtualMemory                         uint64
	InstallDate                               time.Time
	LargeSystemCache                          uint32
	LastBootUpTime                            time.Time
	LocalDateTime                             time.Time
	Locale                                    string
	Manufacturer                              string
	MaxNumberOfProcesses                      uint32
	MaxProcessMemorySize                      uint64
	MUILanguages                              []string
	Name                                      string
	NumberOfLicensedUsers                     uint32
	NumberOfProcesses                         uint32
	NumberOfUsers                             uint32
	OperatingSystemSKU                        uint32
	Organization                              string
	OSArchitecture                            string
	OSLanguage                                uint32
	OSProductSuite                            uint32
	OSType                                    uint16
	OtherTypeDescription                      string
	PAEEnabled                                bool
	PlusProductID                             string
	PlusVersionNumber                         string
	PortableOperatingSystem                   bool
	Primary                                   bool
	ProductType                               uint32
	RegisteredUser                            string
	SerialNumber                              string
	ServicePackMajorVersion                   uint16
	ServicePackMinorVersion                   uint16
	SizeStoredInPagingFiles                   uint64
	Status                                    string
	SuiteMask                                 uint32
	SystemDevice                              string
	SystemDirectory                           string
	SystemDrive                               string
	TotalSwapSpaceSize                        uint64
	TotalVirtualMemorySize                    uint64
	TotalVisibleMemorySize                    uint64
	Version                                   string
	WindowsDirectory                          string
}

*/
