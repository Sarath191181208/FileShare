## Problem

Build a file-sharing platform that allows users to upload, manage, and share
files. The system should be able to handle multiple users, store metadata in
PostgreSQL, manage file uploads to S3 (or locally if cloud integration is not
available), implement caching for file metadata. **The project must be built in
Go** and should demonstrate proficiency in handling concurrency and performance
optimizations.

The problem is clearly written in the following
[doc](https://drive.google.com/file/d/1zeOOxV8rMPXlVkRl236omBBQW_f1EW9g/view)

## Setup gotchas

- postgres `citext` object.

## SQL Injection attacks 
The app is safe against injection attacks are we are passing parametrized quries to the db using the function `QueryRowContext`. Therefore go internally handles the sql injection attacks by sanitizing the input.

## How does the app handle concurrent requests ?

We are using the `net/http` package which is provided as default by golang, as
this is the case the concurrent requests are handled internally by the http
package given to us by go. The internal code of the parallelism is as follows:

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

This is clearly explained in the following
[post](https://stackoverflow.com/questions/40610398/golang-concurrent-http-request-handling)

## What is `.tmux` ?

It is just a simple bash script to init the project to initialize tmux with the
docker config and other dev tools easily rather than creating them one by one

## What's **version** in user model

We’ve included a **version** number column, which we will increment each time a
user record is updated. This will allow us to use optimistic locking to prevent
race conditions when updating user records, in the same way that we did with
movies earlier in the book.

## How are passwords compared and what preventive measures have we taken?

- Passwords are compared using `bcrypt.CompareHashAndPassword()` which re-hashes
  using the same salt and cost param.
- The compare function is safe against timing attacks because of
  `subtle.ConstantTimeCompare()` is an internal func of
  `bcrypt.CompareHashAndPassword`.

## Signup request (`/signup`)

The curl request will look as follows:

```bash
curl -X POST localhost:4000/register \
  -H "Content-Type: application/json" \
  -d '{"email": "sarath@gmail.com", "password": "12345678"}'
```

The example response is as follows:

```
{
  "user": {
    "id": 1,
    "created_at": "2024-09-14T10:28:51Z",
    "email": "sarath@gmail.com"
  }
}
```

## Background worker template

```go
// Launch a background goroutine template
go func() {
    // Run a deferred function which uses recover() to catch any panic, and log an
    // error message instead of terminating the application.
    defer func() {
            if err: = recover();
            err != nil {
                app.logger.PrintError(fmt.Errorf("%s", err), nil)
            }
        }()
    // run the process
    err = app.Process()
    if err != nil {
        app.logger.PrintError(err, nil)
    }
}()
```

## Where can I find the aws-go documentation?
[https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/)
