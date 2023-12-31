package gencommon

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

// SanitizeSourceFile ensures a valid source context.
func SanitizeSourceFile(srcFile string) string {
	pwd := os.Getenv("PWD")
	if srcFile == "" {
		log.Fatal("this command should be run in a go:generate context or with the correct source file flag set.")
	}

	if !path.IsAbs(srcFile) {
		srcFile = path.Join(pwd, srcFile)
	}
	return srcFile
}

// SanitizeOutFile ensures there is a valid destination file.
// This should be run AFTER SanitizeSourceFile, using the result of it as the srcFile argument.
func SanitizeOutFile(flagVal, srcFile, genName string) string {
	srcFileDir, srcFileName := path.Split(srcFile)
	if flagVal == "" {
		return path.Join(srcFileDir, strings.TrimSuffix(srcFileName, ".go")+"."+genName+".go")
	}
	if !path.IsAbs(flagVal) {
		return path.Join(srcFileDir, flagVal)
	}

	return flagVal
}

// Write writes a nicely-formatted template.
func Write(tmpl *template.Template, templateData any, destination string) error {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, templateData)
	if err != nil {
		return err
	}
	tempResult := buf.Bytes()
	result, err := format.Source(tempResult)
	if err == nil {
		// process imports; format and add/remove if needed.
		result, err = imports.Process("", result, nil)
	}
	if err != nil {
		log.Printf("[WARN] - formatting of source file failed with error: %+v", err)
		// gracefully recover- the actual file is broken but it is helpful to see
		// what is wrong with the output itself.
		result = tempResult
	}

	f, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteAt(result, 0)
	return err
}
