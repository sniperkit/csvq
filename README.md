# csvq

SQL-like query language for csv

[![Build Status](https://travis-ci.org/mithrandie/csvq.svg?branch=master)](https://travis-ci.org/mithrandie/csvq)
[![codecov](https://codecov.io/gh/mithrandie/csvq/branch/master/graph/badge.svg)](https://codecov.io/gh/mithrandie/csvq)

csvq is a command line tool to operate CSV files. 
You can read, update, delete CSV records with SQL-like query.

You can also execute multiple operations sequentially in managed transactions by passing a procedure or using the interactive shell.
In the multiple operations, you can use variables, cursors, temporary tables, and other features. 

## Install

### Install executable binary

1. Download an archive file from [release page](https://github.com/mithrandie/csvq/releases).
2. Extract the downloaded archive and add a binary file in it to your path.

### Build from source

#### Requirements

Go 1.9 or later (ref. [Getting Started - The Go Programming Language](https://golang.org/doc/install))

#### Build with either of the following two ways

##### Use go get

1. ```$ go get -u github.com/mithrandie/csvq```

##### Build with strict dependencies

1. Install Glide (ref. [Glide: Vendor Package Management for Golang](https://github.com/Masterminds/glide))
2. ```$ go get -d github.com/mithrandie/csvq```
3. Change directory to $GOPATH/github.com/mithrandie/csvq
4. ```$ glide install```
5. ```$ go install```


## Usage

```shell
# Simple query
csvq "select id, name from `user.csv`"
csvq "select id, name from user"

# Specify data delimiter as tab character
csvq -d "\t" "select count(*) from `user.csv`"

# Load no-header-csv
csvq --no-header "select c1, c2 from user"

# Load from redirection or pipe
csvq "select * from stdin" < user.csv
cat user.csv | csvq "select *"

# Output in JSON format
csvq -f json "select integer(id) as id, name from user"

# Output to a file
csvq -o new_user.csv "select id, name from user"

# Load statements from file
$ cat statements.sql
VAR @id := 0;
SELECT @id := @id + 1 AS id,
       name
  FROM user;

$ csvq -s statements.sql

# Execute statements in the interactive shell
$ csvq
csvq > UPDATE users SET name = 'Mildred' WHERE id = 2;
1 record updated on "/home/mithrandie/docs/csv/users.csv".
csvq > COMMIT;
Commit: file "/home/mithrandie/docs/csv/users.csv" is updated.
csvq > EXIT;

# Show help
csvq -h
```

More details >> [https://mithrandie.github.io/csvq](https://mithrandie.github.io/csvq)