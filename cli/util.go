package cli

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/pkg/browser"
)

func OpenURL(target string) error {
	_, err := url.ParseRequestURI(target)
	if err != nil {
		return err
	}
	return browser.OpenURL(target)
}

var (
	ScanDefaultString string
)

func Scan(message string) (string, error) {
	tmp := "/tmp"
	if runtime.GOOS == "windows" {
		tmp = os.Getenv("TEMP")
	}
	l, err := readline.NewEx(&readline.Config{
		Prompt:            message,
		HistoryFile:       filepath.Join(tmp, "crowi.txt"),
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		return "", err
	}
	defer l.Close()

	var line string
	for {
		if ScanDefaultString == "" {
			line, err = l.Readline()
		} else {
			line, err = l.ReadlineWithDefault(ScanDefaultString)
		}
		if err == readline.ErrInterrupt {
			if len(line) <= len(ScanDefaultString) {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return line, nil
	}
	return "", errors.New("canceled")
}

func FileContent(fname string) string {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func TempFile(filename string) (*os.File, error) {
	return os.Create(filepath.Join(os.TempDir(), filename))
}

func Underline(message, target string) error {
	if message == "" {
		return nil
	}
	if target == "" {
		return errors.New("target is empty")
	}
	link := color.New(color.Underline).SprintFunc()
	fmt.Printf("%s %s\n", message, link(target))
	return nil
}
