# sr - **s**how **r**ecent

`sr <dir>` is a command line tool to show you the most recently modified files and directories (taking into account subdirectories).

`sr` traverses the given directory recursively, sorts files and subdirectories depending on their last
modified dates of the files inside and shows you the result.

```bash
$ sr -h
Usage of sr:
  -L int
        the max depth of the directory tree; -1 if no depth limit (default -1)
  -a    show all files and directories (including hidden ones)
  -d    show dates in digital format
  -f    show dates in fuzzy format (e.g. '12 minutes ago')
  -j    show results in json format
  -n    turn colors off
  -r    reverse the order of items
  -u    show time in UTC
```

The result has the format as the following:
```
<parent_file_or_dir_name> <last_modified_date_of_parent_or_child> <the_number_of_children> <last_modified_child_name>
```


## Examples of usage

```bash
$ sr .
README.md  21 Sep 23 21:53 +0200  0  README.md
sr.go      21 Sep 23 21:52 +0200  0  sr.go
LICENSE    06 May 23 15:26 +0200  0  LICENSE
go.mod     16 Apr 23 20:44 +0200  0  go.mod
go.sum     16 Apr 23 20:30 +0200  0  go.sum

$ sr . -a
README.md  21 Sep 23 21:53 +0200  0   README.md
sr.go      21 Sep 23 21:52 +0200  0   sr.go
.git/      21 Sep 23 21:43 +0200  66  .git/index
LICENSE    06 May 23 15:26 +0200  0   LICENSE
go.mod     16 Apr 23 20:44 +0200  0   go.mod
go.sum     16 Apr 23 20:30 +0200  0   go.sum

$ sr -f
README.md  30 seconds ago         0  README.md
sr.go      4 minutes ago          0  sr.go
sr         7 minutes ago          0  sr
go.mod     12 minutes ago         0  go.mod
go.sum     12 minutes ago         0  go.sum
LICENSE    06 May 23 15:26 +0200  0  LICENSE

$ sr -d
README.md  2023-09-21T21:54:27+0200  0  README.md
sr.go      2023-09-21T21:52:25+0200  0  sr.go
LICENSE    2023-05-06T15:26:43+0200  0  LICENSE
go.mod     2023-04-16T20:44:25+0200  0  go.mod
go.sum     2023-04-16T20:30:27+0200  0  go.sum

$ sr -j -d
[{"name":"README.md","isDir":false,"child":"README.md","lastModified":"2023-09-21T21:54:41+0200","numChildren":0},{"name":"sr.go","isDir":false,"child":"sr.go","lastModified":"2023-09-21T21:52:25+0200","numChildren":0},{"name":"LICENSE","isDir":false,"child":"LICENSE","lastModified":"2023-05-06T15:26:43+0200","numChildren":0},{"name":"go.mod","isDir":false,"child":"go.mod","lastModified":"2023-04-16T20:44:25+0200","numChildren":0},{"name":"go.sum","isDir":false,"child":"go.sum","lastModified":"2023-04-16T20:30:27+0200","numChildren":0}]
```


## Installation

```bash
$ go build && go install
```