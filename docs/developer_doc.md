## Problem 
Build a file-sharing platform that allows users to upload, manage, and share
files. The system should be able to handle multiple users, store metadata in PostgreSQL, manage file
uploads to S3 (or locally if cloud integration is not available), implement caching for file metadata.
**The project must be built in Go** and should demonstrate proficiency in handling concurrency and
performance optimizations.


The problem is clearly written in the following [doc](https://drive.google.com/file/d/1zeOOxV8rMPXlVkRl236omBBQW_f1EW9g/view)
## How does the app handle concurrent requests ? 

We are using the `net/http` package which is provided as default by golang,
as this is the case the concurrent requests are handled internally by the 
http package given to us by go. The internal code of the parallelism is as follows:

```go 
l, _ := net.Listen("tcp", addr)
for {
    rw, _ := l.Accept()
    conn := &conn{
        server: srv,
        rwc:    rwc,
    }
    go s.serve(conn)
```

This is clearly explained in the following [post](https://stackoverflow.com/questions/40610398/golang-concurrent-http-request-handling)

## What is `.tmux` ? 
It is just a simple bash script to init the project to initialize tmux with the docker config and other dev tools
easily rather than creating them one by one
