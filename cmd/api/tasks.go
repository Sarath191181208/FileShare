package api

import "strings"

func (app *Application) DeleteFromDBTask() {
	// Fetch the metadata with the oldest upload date
	metadata, err := app.Models.MetaData.FetchTop()
	if err != nil {
		app.Logger.Println(err)
		return
	}

	// Delete the metadata from the database
	err = app.Models.MetaData.Delete(metadata.ID)
	if err != nil {
		app.Logger.Println(err)
		return
	}
	app.Logger.Printf("Deleted metadata with id %d", metadata.ID)

	// Delete the file from the file store
	key := metadata.FileUrl[strings.LastIndex(metadata.FileUrl, "/")+1:]
	app.Logger.Printf("Deleting file with key %s", key)
	err = app.FileStore.Delete(key)
	if err != nil {
		app.Logger.Printf("S3 Delete Err: %v", err)
		return
	}
}
