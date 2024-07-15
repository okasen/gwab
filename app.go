package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

// App struct
type App struct {
	ctx       context.Context
	novel     Novel
	directory string
}

type Novel struct {
	Title string `json:"Title"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.novel = Novel{Title: "TBA"}
}

func (a *App) SetTitle(value string) Novel {
	a.novel.Title = value
	return a.novel
}

func (a *App) Reset() {
	a.novel.Title = ""
}

func (a *App) Title() string {
	fmt.Println(a.novel.Title)
	return fmt.Sprintf(a.novel.Title)
}

func (a *App) SelectDirectory() {
	options := runtime.OpenDialogOptions{
		CanCreateDirectories: true,
	}
	directory, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		fmt.Println("Something went wrong while choosing a file directory.")
	}
	a.directory = directory
}

func handleError(msg string, e error) {
	if e != nil {
		fmt.Println(msg, ":", e)
	}
}

func (a *App) Save(text string) {
	// TODO: use io and deferred file saving to make file management
	options := runtime.SaveDialogOptions{
		DefaultFilename:      fmt.Sprintf("%s.html", a.novel.Title),
		CanCreateDirectories: true}
	if a.directory != "" {
		options.DefaultDirectory = a.directory
	}
	filepath, err := runtime.SaveFileDialog(a.ctx, options)
	handleError("something went wrong while saving", err)
	file, err := os.Create(filepath)

	defer func(file *os.File) {
		err = file.Close()
		handleError("something went wrong while closing file", err)
	}(file)

	written, err := file.WriteString(text)
	handleError("something went wrong while saving", err)
	if written == 0 && text != "" {
		err = errors.New("written file is unexpectedly empty")
		handleError("something went wrong while saving", err)
	}
	fmt.Println(a.novel.Title, "has been saved! File path:", filepath)
}
