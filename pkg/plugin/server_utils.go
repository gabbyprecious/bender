package plugin

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

func zipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (instance *StartServer[T]) createDownloadDir(plugin *plugin.Plugin[PluginState]) (string, error) {

	clnPath, found := plugin.GetConf("lightning-dir")
	if !found {
		return "", fmt.Errorf("lightning-dir not found in config")
	}
	sourcePath := strings.Join([]string{clnPath.(string), "bender"}, "/")
	err := os.MkdirAll(sourcePath, 0755)
	if err != nil {
		return "", err
	}

	targetPath := filepath.Join(clnPath.(string))

	files := []string{"ca.pem", "client-key.pem", "client.pem"}

	for _, file := range files {
		targetFile := filepath.Join(targetPath, file)
		sourceFile := filepath.Join(sourcePath, file)
		_, err := copy(targetFile, sourceFile)
		if err != nil {
			return "", err
		}
	}
	return sourcePath, nil
}
