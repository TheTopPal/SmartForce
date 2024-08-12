package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func checkArgs() (string, string, string, error) {
	if len(os.Args) != 4 {
		return "", "", "", fmt.Errorf("Usage: go run main.go <directory> <old_text> <new_text>")
	}
	return os.Args[1], os.Args[2], os.Args[3], nil
}

func dirExists(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	return nil
}

func createLogFile() (*os.File, error) {
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	logFile, err := os.Create(fmt.Sprintf("log-%s.txt", timestamp))
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func replaceTextInFiles(dir, oldText, newText string, logFile *os.File) error {
	var replaced bool
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			text := string(data)
			offsets := []int{}
			for i := 0; i < len(text); i++ {
				if strings.HasPrefix(text[i:], oldText) {
					offsets = append(offsets, i)
				}
			}
			if len(offsets) > 0 {
				text = strings.ReplaceAll(text, oldText, newText)
				err = os.WriteFile(path, []byte(text), info.Mode())
				if err != nil {
					return err
				}
				log.Printf("Replaced text in file %s\n", path)
				replaced = true
				logFile.WriteString(fmt.Sprintf("File: %s\n", path))
				logFile.WriteString("Replacements:\n")
				for _, offset := range offsets {
					logFile.WriteString(fmt.Sprintf("  %d: \"%s\" -> \"%s\"\n", offset, oldText, newText))
				}
				logFile.WriteString("\n")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if !replaced {
		logFile.WriteString("No replacements made\n")
	}
	return nil
}

func main() {
	dir, oldText, newText, err := checkArgs()
	if err != nil {
		log.Fatal(err)
	}

	err = dirExists(dir)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := createLogFile()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	err = replaceTextInFiles(dir, oldText, newText, logFile)
	if err != nil {
		log.Fatal(err)
	}
}
