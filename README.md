# logwatch

A `tail -f` like logfile watcher written in Go. It can watch an arbitrary
number of files in parallel and will output all lines of all files that
match the pattern defined for the respective file.

Usage:

    go install
    logwatch <file1> <pattern1> [<file2> <pattern2> ...]
