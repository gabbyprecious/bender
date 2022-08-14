package plugin

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vincenzopalazzo/cln4go/plugin"
	// "strconv"
)

type PluginState struct {
	Server   *http.Server
	Password string
}

type Login struct {
	Password string `json:"password" binding:"required"`
}

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

type StartServer[T PluginState] struct{}

func (instance *StartServer[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) (map[string]any, error) {
	if plugin.State.Server != nil {
		return nil, fmt.Errorf("Server is already running")
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/tls", func(c *gin.Context) {
		var login Login
		if err := c.BindJSON(&login); err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		if plugin.State.Password != login.Password {
			c.JSON(403, gin.H{"status": "Access Denied"})
			return
		}
		targetPath, err := instance.createDownloadDir(plugin)
		if err != nil {
			c.JSON(403, err)
			return
		}
		fileName := "certificates.zip"
		zipTargetPath := filepath.Join(targetPath, fileName)
		err = zipSource(targetPath, zipTargetPath)
		if err != nil {
			c.JSON(403, err)
			return
		}
		c.FileAttachment(zipTargetPath, fileName)
		os.RemoveAll(targetPath)
	})
	portValue, found := plugin.GetOpt("bender_port")

	/// TODO: disable the plugin if we do not insert the port
	if !found || portValue == "-1" {
		portValue = "9080"
	}

	plugin.State.Server = &http.Server{
		Addr:    fmt.Sprintf(":%v", portValue),
		Handler: router,
	}

	go func() {
		// service connections
		if err := plugin.State.Server.ListenAndServe(); err != nil {
			//plugin.Log("info", fmt.Sprintf("error received from server: %s", err))
			plugin.State.Server = nil
		}
	}()

	return map[string]any{"message": "Server up and running listen and serve on 0.0.0.0:" + fmt.Sprintf("%v", portValue)}, nil
}

type SetPassword[T PluginState] struct{}

func (instance *SetPassword[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) (map[string]any, error) {
	plugin.State.Password = fmt.Sprintf(request["password"].(string))
	return map[string]interface{}{"message": "Password set", "password": plugin.State.Password}, nil
}

type Hello[T PluginState] struct{}

func (instance *Hello[PluginState]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) (map[string]any, error) {
	return map[string]any{"message": "hello from go 1.18"}, nil
}

type OnShutdown[T PluginState] struct{}

func (instance *OnShutdown[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := plugin.State.Server.Shutdown(ctx); err != nil {
		panic(err)
	}
	os.Exit(0)
}
