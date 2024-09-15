## How are we handling concurrent file uploading to S3?

As we use the `uploader.Upload` from aws it gives us intelligent buffering large files into smaller chunks 
and sending them in parallel across multiple goroutines. 
You can configure the buffer size and concurrency through the Uploader's parameters.
Additional functional options can be provided to configure the individual upload.
These options are copies of the Uploader instance Upload is called from.
Modifying the options will not impact the original Uploader instance.
Use the WithUploaderRequestOptions helper function to pass in request 
options that will be applied to all API operations made with this uploader.
It is safe to call this method concurrently across goroutines.

## Background Job for File Deletion
We created a helper function in our application to handle the background task 
```go
// file /cmd/api/api.go
func (app *Application) Background(fn func()) {
	// Launch a background goroutine.
	go func() {
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
        app.Logger.Printf("Recovered from a panic: %v", err)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}
```

It is used like this:
```go
// file: /cmd/main.go
app.Background(func() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    for v := range ticker.C {
        app.Logger.Printf("Running delete from db task at %v", v)
        app.DeleteFromDBTask()
    }
})
```
