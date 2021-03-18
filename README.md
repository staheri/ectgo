# ECT GO
This is a prototype tool for automatic collection of concurrency tracing.
Details of the idea and design is under review in 6th Workshop on Formal Integrated Development Environment ([F-IDE 2021](https://cister-labs.pt/f-ide2021/))

## Build
Follow Install for downloading dependencies and building a new runtime. Then execute below program
```
cd [path-to-ectgo]
go build
```
Current build command is not tested and might crash.I am actively working on engineering a nice build for **ectgo**. Please refer to [patch](https://github.com/staheri/ectgo/ectgo_runtime_v1_15_6.patch) for the details of enhancement to the original tracer package.
Please refer to [goatlib](https://github.com/staheri/goatlib)] for the implemented API of automatic trace collection and storage.

Email me [staheri@cs.utah.edu](MAILTO:staheri@cs.utah.edu) if you want to try this out. I would be more than happy to assist.


## Overview
![Overview](overview.png)


## Install
Currently work for Go V1.15.6 (macOS) but extensible to other versions and systems (replace below download link with the ones compatible with your system)
Steps:
- Patch the runtime
- Dependencies
- Build
- Features

### Patching Runtime
`goTrace_runtime_v1_15_6.patch` has all the needed injections to the Go runtime in order to capture additional events like channel operations, waiting groups and mutexes.

Assuming your Go installation is in `/usr/local/go`, download Go 1.15.6 and unpack it into `/usr/local/go-new`.
```
 sudo -i
 mkdir -p /usr/local/go-new
 curl https://dl.google.com/go/go1.15.6.darwin-amd64.tar.gz | tar -xz -C /usr/local/go-new
 ```

Then, copy patch and apply it:
```
sudo patch -p1 -d /usr/local/go-new/go < goTrace_runtime_v1_15_6.patch
```

Now you can build the new runtime
```
 sudo -i
 cd /usr/local/go-new/go/src
 export GOROOT_BOOTSTRAP=/usr/local/go #or choose yours
 ./make.bash
 ```

Finally, `export PATH` or `use ln -s` command to make this Go version actual in your system:
```
 export PATH=/usr/local/go-new/go/bin:$PATH
 ```
or (assuming your PATH set to use /usr/local/go)
```
	sudo mv /usr/local/go /usr/local/go-orig
	sudo ln -nsf /usr/local/go-new/go /usr/local/go
```
NOTE: return your previous installation by `sudo ln -nsf /usr/local/go-orig /usr/local/go`


### Dependencies
GoTrace uses different libraries and drivers. *[TODO] Use Go Modules/Vendors to automatically detect dependencies and versions*

#### Libraries

- Fine tables: [github.com/jedib0t/go-pretty/table](https://github.com/jedib0t/go-pretty)
  `go get github.com/jedib0t/go-pretty`
- Go MySQL driver: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
  `go get github.com/go-sql-driver/mysql`
- AST traversal: [golang.org/x/tools/go/ast/astutil](https://golang.org/x/tools/go/ast/astutil)
  `golang.org/x/tools/go/ast/astutil`
- There might be more libraries needed indirectly.

#### Database

- MySQL: [Install on Mac](https://dev.mysql.com/doc/mysql-osx-excerpt/5.7/en/osx-installation-pkg.html)

### Build and Usage
[working on it]
