package main

import (
	"context"
	"fmt"
	"gwab/internal/novels"
	_ "gwab/internal/novels"
)

// App struct
type App struct {
	ctx       context.Context
	directory string
	novel     novels.Novel
	ops       novels.Operations
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.novel = novels.Novel{Title: "TBA"}
	a.ops = novels.Ops()
}

func (a *App) SetTitle(value string) novels.Novel {
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

func (a *App) Save(text string) {
	fw := novels.FileWriter{}
	nsd := novels.NovelSaveDialog{}
	neh := novels.NovelErrorHandler{}
	a.ops.Save(&a.novel, a.ctx, text, &fw, &nsd, &neh)
}

func (a *App) DebugLog(level string, export bool) {
	neh := novels.NovelErrorHandler{}
	a.ops.DebugLog(&a.novel, level, export, &neh)
}
