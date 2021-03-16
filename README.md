# goTrace



## Intro

goTrace helps Go developer understand the dynamic behavior of their application by automatically:
- Instruments Go source
   * It traverses source AST tree
   * And injects trace collection API to the source
- Executes the target application and redirects its trace to ```stderr```
   * Go deadlock detector can be disabled on applications that suffer from deadlock
- Inserts traces into a MySQL database
   * Now you can query the database and study the behavior

After storing traces in the database, then we can interpret, summarize and retrieve traces in various formats to study Go dynamic behavior.

## Install

Steps:
- Patch the runtime
- Dependencies
- Build
- Features

### Patching Runtime
`goTrace-runtime.patch` has all the needed injections to the Go runtime in order to capture additional events like channel operations, waiting groups and mutexes.

Assuming your Go installation is in `/usr/local/go`, download Go 1.14.4 and unpack it into `/usr/local/go-new`.
```
 sudo -i
 mkdir -p /usr/local/go-new
 curl https://dl.google.com/go/go1.14.4.darwin-amd64.tar.gz | tar -xz -C /usr/local/go-new
 ```

Then, copy patch and apply it:
```
sudo patch -p1 -d /usr/local/go-new/go < goTrace-runtime.patch
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

### Build
First make sure you have set-up the Go environment variables correctly
```
export GOROOT=/usr/local/go
export GOPATH=<path-to>/goTrace
export PATH=$GOROOT/bin:$PATH
```
Then
```
cd goTrace/src
go build
```

#### Set variables
You may want to set the size of training data chunk size and the middle folder for storing CL-related files at the beginning of the *main.go*:
```
const WORD_CHUNK_LENGTH = 11
var CLOUTPATH = os.Getenv("GOPATH")+"/traces/clx"
```


### Features
GoTrace provides an extendable platform that anyone can extend and generate various reports, abstraction, training data, etc. from database rows.

#### Resource Report (RR)
In Go, *Goroutines* are light-weight independent processing units that preempt an actual thread (on real core) that are managed by Go runtime (i.e., *scheduler*).  As the philosophy of Go states:

> Do not communicate by sharing memory; intsead, share memory by communicating

Go benefits from the advantages of *shared* memory systems (like OpenMP) and *distributed* memory systems (like MPI). **Channels**, **Mutexes** and **WaitingGroups** are some of the main virtual resources that Go utilize to establish communication between goroutines. As they are key concepts in *concurrent Go* and potential root-cause of majority of concurrent bugs, goTrace summarize collected traces in form of **Resource Reports**(RR). RRs show the activities of individual resource were resources are:
- Channles
- Mutex (Normal and RW)
- WaitingGroups

```
./src -cmd=rr -app=<path-to-your-app>.go [-other options]
```
Generates reports like below (for more options on RR options, run `./src --help`)

```
Channel global ID: 5
Owner: N/A (e.g., created globaly)
Closed? No
```
| TS | Send | Recv |
| ---:| --- | --- |
| 92036 | G1: dl-triple-sol.go>main.main:29<br/> | - |
| 189693 | - | G19: dl-triple-sol.go>main.worker:18<br/> |
| 207166 | G19: dl-triple-sol.go>main.worker:18<br/> | - |
| 233748 | - | G23: dl-triple-sol.go>main.worker:18<br/> |
| 236827 | G23: dl-triple-sol.go>main.worker:18<br/> | - |
| 243652 | - | G22: dl-triple-sol.go>main.worker:18<br/> |
| 246988 | G22: dl-triple-sol.go>main.worker:18<br/> | - |
| 264384 | - | G21: dl-triple-sol.go>main.worker:18<br/> |
| 267027 | G21: dl-triple-sol.go>main.worker:18<br/> | - |
| 281524 | - | G20: dl-triple-sol.go>main.worker:18<br/> |
| 284013 | G20: dl-triple-sol.go>main.worker:18<br/> | - |
| 307901 | - | G1: dl-triple-sol.go>main.main:36<br/> |

#### Hierarchical Clustering (HAC)
Using FCA, we hierarchically cluster goroutines from different aspects(for more our information refer to our CLUSTER'19 paper):
- **CHNL**: Channel activities
- **GRTN**: Goroutines layout
- **MUTX**: Mutex-related events
- **WGRP**: WaitingGroup-related events
- **PROC**: Process-level activities
- **GCMM**: Go garbage collection and memory events
- **SYSC**: System calls
- **MISC**: Other events (e.g., user annotated)

```
./src -cmd=hac -app=<path-to-your-app>.go [CHNL GRTN MUTX WGRP]
```

`CHNL` only clusters goroutines based on attributes extracted from channel activities while `CHNL-GRTN-MUTX PROC` generates two clustering for goroutines:
- One with respect to combination of activities related to channels, goroutines and mutexes.
- And one with w.r.t. processes

Here is a demo:

![hac-demo](demo/hac-demo.png)

#### Word2Vec model
As we are gaining knowledge about applying machine learning methods into HPC correctness, **word2vec** idea in particular. goTrace generates training data from collected traces.
```
./src -cmd=word -app=<path-to-your-app>.go -out=<training-folder>
```

 For more information please refer to [DataBenchmark](DataBenchmark/)
