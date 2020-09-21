package views

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"github.com/legenove/nano-server-sdk/servers"
	"github.com/legenove/utils"
)

type docInfo struct {
	PackageName string
	Title       string
	Path        string
}

func GetAllSwaggerFileByPath(_fpath string) []string {
	var res = []string{}
	if !utils.PathExists(_fpath) {
		return res
	}
	list, err := ioutil.ReadDir(_fpath)
	if err != nil {
		return res
	}
	for _, n := range list {
		if n.IsDir() {
			continue
		}
		ns := strings.Split(n.Name(), ".")
		if len(ns) > 0 && utils.IsInStringSlice([]string{"yml", "yaml", "json"}, ns[len(ns)-1]) {
			p := filepath.Join(_fpath, n.Name())
			res = append(res, p)
		}
	}
	return res
}

// /doc [get]
func TemplateDocApi(c *gin.Context) {
	indexPath := path.Join(servers.Server.DocDir, "swagger")
	if !utils.PathExists(indexPath) {
		panic(servers.ErrSchemaOptionNotFound)
	}
	docFileNames := GetAllSwaggerFileByPath(indexPath)
	infos := make([]docInfo, 0, len(docFileNames))
	for _, fname := range docFileNames {
		fname := filepath.Base(fname)
		pName := strings.Split(fname, ".")[0]
		infos = append(infos, docInfo{
			Path:        "/doc/swagger/" + pName + ".json",
			Title:       "[ " + pName + " ]接口文档",
			PackageName: pName,
		})
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"docInfos":    infos,
		"docLen":      len(infos) - 1,
		"serverTitle": servers.Server.GetServerTitle(),
		"serverGroup": servers.Server.GetServerGroup(),
		"serverName":  servers.Server.GetServerName(),
	})
}

// /doc/file [get]
func TemplateDocFileApi(c *gin.Context) {
	indexPath := path.Join(servers.Server.DocDir, "swagger")
	if !utils.PathExists(indexPath) {
		panic(servers.ErrSchemaOptionNotFound)
	}
	docFileNames := GetAllSwaggerFileByPath(indexPath)
	infos := make([]docInfo, 0, len(docFileNames))
	for _, fname := range docFileNames {
		fname := filepath.Base(fname)
		pName := strings.Split(fname, ".")[0]
		infos = append(infos, docInfo{
			Path:        "/doc/swagger/" + pName + ".json",
			Title:       "[ " + pName + " ]接口文档",
			PackageName: pName,
		})
	}

	c.HTML(http.StatusOK, "filedoc.tmpl", gin.H{
		"docInfos":    infos,
		"docLen":      len(infos) - 1,
		"serverTitle": servers.Server.GetServerTitle(),
		"serverGroup": servers.Server.GetServerGroup(),
		"serverName":  servers.Server.GetServerName(),
	})
}

// /doc/swagger/:name [get]
func TemplateDocSwaggerApi(c *gin.Context) {
	packageName := c.Param("name")

	if packageName == "" {
		panic(servers.ErrSchemaOptionNotFound)
	}
	_type := "json"
	if strings.HasSuffix(packageName, ".json") {
		packageName = packageName[:len(packageName)-5]
	} else if strings.HasSuffix(packageName, ".yaml") {
		_type = "yaml"
		packageName = packageName[:len(packageName)-5]
	} else if strings.HasSuffix(packageName, ".yml") {
		_type = "yaml"
		packageName = packageName[:len(packageName)-4]
	}
	tpath := path.Join(servers.Server.DocDir, "swagger", packageName+".yaml")

	data, _ := ioutil.ReadFile(tpath)
	if _type == "json" {
		dat, _ := yaml.YAMLToJSON(data)
		c.Status(200)
		c.Writer.Write(dat)
	} else {
		c.Status(200)
		c.Writer.Write(data)
	}
}

// /doc/proto/:name [get]
func TemplateDocProtoApi(c *gin.Context) {
	packageName := c.Param("name")

	if packageName == "" {
		panic(servers.ErrSchemaOptionNotFound)
	}
	if strings.HasSuffix(packageName, ".proto") {
		packageName = packageName[:len(packageName)-6]
	}
	tpath := path.Join(servers.Server.DocDir, "proto", packageName, packageName+".proto")
	if !utils.FileExists(tpath) {
		panic(servers.ErrSchemaOptionNotFound)
	}
	data, _ := ioutil.ReadFile(tpath)
	c.Status(200)
	c.Writer.Write(data)
}
