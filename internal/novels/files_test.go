package novels

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockWriter struct{ mock.Mock }

func (m *MockWriter) Write(filepath string, text string, novel *Novel, e ErrorHandler) int {
	args := m.Called(filepath, text, novel)
	return args.Int(0)
}

type MockSaveDialog struct{ mock.Mock }

func (m *MockSaveDialog) Handle(novel *Novel, ctx context.Context, e ErrorHandler) string {
	args := m.Called(novel, ctx)
	return args.String(0)
}

type MockErrorHandler struct{ mock.Mock }

func (m *MockErrorHandler) Handle(msg string, e error, novel *Novel, level string) {
	m.Called(msg, e, novel, level)
	return
}

func TestSavingCallsFileCreationWithFullHTMLContents(t *testing.T) {
	novel := Novel{Title: "TestFile"}
	ctx := context.Context(context.Background())
	text := "<p>Test Paragraph!</p>"
	expectedFilePath := "file/path/TestFile.html"
	mw := MockWriter{}
	msd := MockSaveDialog{}
	meh := MockErrorHandler{}
	meh.On("Handle", mock.Anything, nil, &novel, "").Return()
	mw.On("Write", expectedFilePath, text, &novel).Return(len(text))
	msd.On("Handle", &novel, ctx).Return(expectedFilePath)
	Save(&novel, ctx, text, &mw, &msd, &meh)
	msd.AssertCalled(t, "Handle", &novel, ctx)
	mw.AssertCalled(t, "Write", expectedFilePath, text, &novel)
}

func TestSavingEmptyFileDoesNotCallErrorHandler(t *testing.T) {
	novel := Novel{Title: "TestFile"}
	ctx := context.Context(context.Background())
	text := "" // empty text file
	expectedFilePath := "file/path/TestFile.html"
	mw := MockWriter{}
	msd := MockSaveDialog{}
	meh := MockErrorHandler{}
	mw.On("Write", expectedFilePath, text, &novel).Return(0) // 0 length because no text
	msd.On("Handle", &novel, ctx).Return(expectedFilePath)
	meh.On("Handle", mock.Anything, nil, &novel, "").Return(nil)
	Save(&novel, ctx, text, &mw, &msd, &meh)
	meh.AssertNotCalled(t, "Handle", mock.Anything, nil, &novel, "")
}

func TestErrorIsThrowWhenNothingIsWrittenButTextPresent(t *testing.T) {
	novel := Novel{Title: "TestFile"}
	ctx := context.Context(context.Background())
	text := "<p>There is text here!</p>" // non-empty text file
	expectedFilePath := "file/path/TestFile.html"
	mw := MockWriter{}
	msd := MockSaveDialog{}
	meh := MockErrorHandler{}
	mw.On("Write", expectedFilePath, text, &novel).Return(0) // 0 length should throw error
	msd.On("Handle", &novel, ctx).Return(expectedFilePath)
	meh.On("Handle", "something went wrong while saving", errors.New("written file is unexpectedly empty"), &novel, "error")
	Save(&novel, ctx, text, &mw, &msd, &meh)
	meh.AssertCalled(t, "Handle", "something went wrong while saving", errors.New("written file is unexpectedly empty"), &novel, "error")
}

func TestNovelDebugLogHoldsAllHandledErrorsAndLogs(t *testing.T) {
	novel := Novel{Title: "TestForDebugFile"}
	assert.Equal(t, novel.Debug, *new([]string), "The Debug log should currently be empty")

	realErrorHandler := NovelErrorHandler{}

	text := "<p>There is text here!</p>"
	ctx := context.Context(context.Background())
	expectedFilePath := "file/path/TestFile.html"
	mw := MockWriter{}
	msd := MockSaveDialog{}
	mw.On("Write", expectedFilePath, text, &novel).Return(0) // 0 length should throw error
	msd.On("Handle", &novel, ctx).Return(expectedFilePath)
	Save(&novel, ctx, text, &mw, &msd, realErrorHandler)

	expectedLogs := []string{"debug- Saving the novel TestForDebugFile !", "error- something went wrong while saving: written file is unexpectedly empty", "debug- TestForDebugFile has been saved! File path: file/path/TestFile.html"}
	assert.Equal(t, novel.Debug, expectedLogs, "The Debug log should have events and errors logged")
}
