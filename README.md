# sr - **s**how **r**ecent

`sr` is a command line tool to show the most recently modified files and directories (taking into account subdirectories).

`sr` traverses the given directory recursively, sorts files and subdirectories depending on their last
modified dates of the files inside and shows you the results.

```bash
$ sr -h
Usage of sr:
  -L int
        the max depth of the directory tree; -1 if no depth limit (default -1)
  -d    show dates in digital format
  -j    show results in json format
  -n    turn colors off
  -r    reverse the order of items
```

The result has the format as the following:
```
<parent_file_or_dir_name> <last_modified_date_of_parent_or_child> <the_number_of_children> <last_modified_child_name>
```


## Examples of usage

```bash
$ sr
README.md  06 May 23 15:29 +0200  0   README.md
.git/      06 May 23 15:28 +0200  26  .git/index
LICENSE    06 May 23 15:26 +0200  0   LICENSE
sr.go      22 Apr 23 08:22 +0200  0   sr.go
go.mod     16 Apr 23 20:44 +0200  0   go.mod
go.sum     16 Apr 23 20:30 +0200  0   go.sum

$ sr -d
README.md  2023-05-06T15:29:05+0200  0   README.md
.git/      2023-05-06T15:28:34+0200  26  .git/index
LICENSE    2023-05-06T15:26:43+0200  0   LICENSE
sr.go      2023-04-22T08:22:58+0200  0   sr.go
go.mod     2023-04-16T20:44:25+0200  0   go.mod
go.sum     2023-04-16T20:30:27+0200  0   go.sum

$ sr -j -d
[{"name":"README.md","isDir":false,"child":"README.md","lastModified":"2023-05-06T15:29:05+0200","numChildren":0},{"name":".git","isDir":true,"child":".git/index","lastModified":"2023-05-06T15:28:34+0200","numChildren":26},{"name":"LICENSE","isDir":false,"child":"LICENSE","lastModified":"2023-05-06T15:26:43+0200","numChildren":0},{"name":"sr.go","isDir":false,"child":"sr.go","lastModified":"2023-04-22T08:22:58+0200","numChildren":0},{"name":"go.mod","isDir":false,"child":"go.mod","lastModified":"2023-04-16T20:44:25+0200","numChildren":0},{"name":"go.sum","isDir":false,"child":"go.sum","lastModified":"2023-04-16T20:30:27+0200","numChildren":0}]
```


## Installation

```bash
$ go build && go install
```