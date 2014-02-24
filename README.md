wmi
===

WMI interface for Go with a database/sql interface

# WARNING

Due to a [memory bug](https://github.com/mattn/go-ole/issues/13) only use the wmi package in its own executable and never within a package. This executable must disable the GC and should either exit after running a query or safely exit after running N queries (since memory use will go up). See [this github issue](https://github.com/StackExchange/wmi/issues/1) for more details.
