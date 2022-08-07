package plugin

import (
	"fmt"
	// "log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vincenzopalazzo/cln4go/plugin"
)

type PluginState struct {
	Server   *http.Server
	Password string
}

type Login struct {
	Password string `json:"password" binding:"required"`
}

type TLSFiles struct {
	Ca        string `json:"ca"`
	ClientKey string `json:"client_key"`
	Client    string `json:"client"`
}

func readTlsFile() (TLSFiles, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return TLSFiles{}, err
	}
	clientKey, err := os.ReadFile(home + "/.lightning/testnet/client-key.pem")
	if err != nil {
		return TLSFiles{}, err
	}
	ca, err := os.ReadFile(home + "/.lightning/testnet/ca.pem")
	if err != nil {
		return TLSFiles{}, err
	}
	client, err := os.ReadFile(home + "/.lightning/testnet/client.pem")
	if err != nil {
		return TLSFiles{}, err
	}
	files := TLSFiles{Ca: string(ca), ClientKey: string(clientKey), Client: string(client)}
	return files, nil
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
		data, err := readTlsFile()
		if err != nil {
			c.JSON(404, gin.H{"error from tls files": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": &data})
	})

	plugin.State.Server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	var serverError error

	go func() {
		// service connections
		if err := plugin.State.Server.ListenAndServe(); err != nil {
			serverError = fmt.Errorf("gin framework error: %s", err)
		}
		serverError = nil
	}()
	if serverError != nil {
		return nil, serverError
	}
	return map[string]any{"message": "Server up and running,/ listen and serve on 0.0.0.0:8080 "}, nil
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

func (instance *OnShutdown[PluginState]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) {
	os.Exit(0)
}
