package novels

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"strings"
	"time"
)

type saveFunc func(novel *Novel, context context.Context, text string, writer Writer, s SaveDialog, e ErrorHandler)

type DebugLogFunc func(novel *Novel, level string, export bool, e ErrorHandler)

type Operations struct {
	Save     saveFunc
	DebugLog DebugLogFunc
}

func Ops() Operations {
	return Operations{
		Save:     Save,
		DebugLog: DebugLog,
	}
}

type Novel struct {
	Title string   `json:"Title"`
	Debug []string `json:"Debug"`
}

type Writer interface {
	Write(filepath string, text string, novel *Novel, e ErrorHandler) int
}

type SaveDialog interface {
	Handle(novel *Novel, ctx context.Context, e ErrorHandler) string
}

type ErrorHandler interface {
	Handle(msg string, e error, novel *Novel, level string)
}

type NovelSaveDialog struct{}

type FileWriter struct{}

type NovelErrorHandler struct{}

func (f *FileWriter) Write(filepath string, text string, novel *Novel, e ErrorHandler) int {

	file, err := os.Create(filepath)

	defer func(file *os.File) {
		err = file.Close()
		e.Handle("something went wrong while closing file", err, novel, "warn")
	}(file)

	written, err := file.WriteString(text)
	e.Handle("something went wrong while saving", err, novel, "error")
	return written
}

func (dialog NovelSaveDialog) Handle(novel *Novel, ctx context.Context, e ErrorHandler) string {
	options := runtime.SaveDialogOptions{
		DefaultFilename:      fmt.Sprintf("%s.html", novel.Title),
		CanCreateDirectories: true}
	filepath, err := runtime.SaveFileDialog(ctx, options)
	e.Handle("something went wrong while saving", err, novel, "error")
	return filepath
}

func (h NovelErrorHandler) Handle(msg string, e error, novel *Novel, level string) {
	if e != nil {
		log := fmt.Sprintf("%s- %s: %v", level, msg, e)
		fmt.Println(log)
		novel.Debug = append(novel.Debug, log)
	}
}

func Save(novel *Novel, ctx context.Context, text string, writer Writer, s SaveDialog, e ErrorHandler) {
	novel.Debug = append(novel.Debug, fmt.Sprintf("debug- Saving the novel %s !", novel.Title))
	filepath := s.Handle(novel, ctx, e)
	written := writer.Write(filepath, text, novel, e)
	if written == 0 && text != "" {
		e.Handle("something went wrong while saving", errors.New("written file is unexpectedly empty"), novel, "error")
	}
	log := fmt.Sprintf("debug- %s has been saved! File path: %s", novel.Title, filepath)
	fmt.Println(log)
	novel.Debug = append(novel.Debug, log)
}

func DebugLog(novel *Novel, level string, export bool, e ErrorHandler) {
	novel.Debug = append(novel.Debug, "debug- Generating Debug Log")
	if export {
		dir, err := os.UserHomeDir()
		e.Handle("something went wrong while getting user home directory", err, novel, "error")
		path := fmt.Sprintf("%v/%s-%s-%v-%s.txt", dir, novel.Title, "debug", time.Now().Unix(), level)
		file, err := os.Create(path)
		defer func(file *os.File) {
			err = file.Close()
			e.Handle("something went wrong while closing file", err, novel, "warn")
		}(file)
		e.Handle("Error generating file for debug log", err, novel, "warn")
		for _, entry := range novel.Debug {
			if level == "" {

				fmt.Println(entry)
				_, err = file.WriteString(fmt.Sprintf(entry + "\n"))
				continue
			}
			if prefix := fmt.Sprintf("%s-", level); strings.HasPrefix(entry, prefix) {
				_, err = file.WriteString(fmt.Sprintf(entry + "\n"))
				fmt.Println(entry)
			}
		}
	}
}
