package plugin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/vincenzopalazzo/cln4go/client"
	"github.com/vincenzopalazzo/cln4go/plugin"
)

type PluginState struct {
	Server   *http.Server
	Password string
	Client   client.Client
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
