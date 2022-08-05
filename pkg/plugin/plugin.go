package plugin

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

type LOGIN struct {
	PASSWORD string `json:"password" binding:"required"`
}

type TLSFiles struct {
	CA        string `json:"ca"`
	CLIENTKEY string `json:"client_key"`
	CLIENT    string `json:"client"`
}

func readTlsFile() (TLSFiles, error){
	home, err := os.UserHomeDir()
	if err != nil {
	  return TLSFiles{}, err
	}
	clientKey, err := os.ReadFile(home + "/.lightning/testnet/client-key.pem")
	ca, err := os.ReadFile(home + ".lightning/testnet/ca.pem")
	client, err := os.ReadFile(home + ".lightning/testnet/client.pem")
	if err != nil{
	  return TLSFiles{}, err
	}
	files := TLSFiles{CA: string(ca), CLIENTKEY: string(clientKey), CLIENT: string(client)}
	return files, nil
}

func startServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/tls", func(c *gin.Context) {
		var login LOGIN
		c.BindJSON(&login)
		if os.Getenv("password") != login.PASSWORD {
			c.JSON(200, gin.H{"status": "Access Denied"})
		}
		data, _ := readTlsFile()
		c.JSON(200, gin.H{"data": &data})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type SetPassword[T PluginState] struct{}

func (instance *SetPassword[T]) Call(plugin *plugin.Plugin[T], request map[string]string) (map[string]string, error) {
	m := make(map[string]string)
	m["password"] = request["password"]
	os.Setenv("password", m["password"])
	return m, nil
}

type PluginState struct{}

type Hello[T PluginState] struct{}

func (instance *Hello[T]) Call(plugin *plugin.Plugin[T], request map[string]any) (map[string]any, error) {
	return map[string]any{"message": "hello from cln4go.template"}, nil
}

type OnShutdown[T PluginState] struct{}

func (instance *OnShutdown[T]) Call(plugin *plugin.Plugin[T], request map[string]any) {
	os.Exit(0)
}
