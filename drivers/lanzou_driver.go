package drivers

import (
	"fmt"
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/model"
	"github.com/Xhofe/alist/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

type Lanzou struct{}

func (driver Lanzou) Config() DriverConfig {
	return DriverConfig{
		Name:      "Lanzou",
		OnlyProxy: false,
	}
}

func (driver Lanzou) Items() []Item {
	return []Item{
		{
			Name:     "onedrive_type",
			Label:    "lanzou type",
			Type:     SELECT,
			Required: true,
			Values:   "cookie,url",
		},
		{
			Name:        "access_token",
			Label:       "cookie",
			Type:        STRING,
			Description: "about 15 days valid",
		},
		{
			Name:  "root_folder",
			Label: "root folder file_id",
			Type:  STRING,
		},
		{
			Name:  "site_url",
			Label: "share url",
			Type:  STRING,
		},
		{
			Name:  "password",
			Label: "share password",
			Type:  STRING,
		},
	}
}

func (driver Lanzou) Save(account *model.Account, old *model.Account) error {
	if account.OnedriveType == "cookie" {
		if account.RootFolder == "" {
			account.RootFolder = "-1"
		}
	}
	account.Status = "work"
	_ = model.SaveAccount(account)
	return nil
}

func (driver Lanzou) File(path string, account *model.Account) (*model.File, error) {
	path = utils.ParsePath(path)
	if path == "/" {
		return &model.File{
			Id:        account.RootFolder,
			Name:      account.Name,
			Size:      0,
			Type:      conf.FOLDER,
			Driver:    driver.Config().Name,
			UpdatedAt: account.UpdatedAt,
		}, nil
	}
	dir, name := filepath.Split(path)
	files, err := driver.Files(dir, account)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name == name {
			return &file, nil
		}
	}
	return nil, PathNotFound
}

func (driver Lanzou) Files(path string, account *model.Account) ([]model.File, error) {
	path = utils.ParsePath(path)
	var rawFiles []LanZouFile
	cache, err := conf.Cache.Get(conf.Ctx, fmt.Sprintf("%s%s", account.Name, path))
	if err == nil {
		rawFiles, _ = cache.([]LanZouFile)
	} else {
		file, err := driver.File(path, account)
		if err != nil {
			return nil, err
		}
		rawFiles, err = driver.GetFiles(file.Id, account)
		if err != nil {
			return nil, err
		}
		if len(rawFiles) > 0 {
			_ = conf.Cache.Set(conf.Ctx, fmt.Sprintf("%s%s", account.Name, path), rawFiles, nil)
		}
	}
	files := make([]model.File, 0)
	for _, file := range rawFiles {
		files = append(files, *driver.FormatFile(&file))
	}
	return files, nil
}

func (driver Lanzou) Link(path string, account *model.Account) (string, error) {
	file, err := driver.File(path, account)
	if err != nil {
		return "", err
	}
	log.Debugf("down file: %+v", file)
	downId := file.Id
	if account.OnedriveType == "cookie" {
		downId, err = driver.GetDownPageId(file.Id, account)
		if err != nil {
			return "", err
		}
	}
	link, err := driver.GetLink(downId)
	if err != nil {
		return "", err
	}
	return link, nil
}

func (driver Lanzou) Path(path string, account *model.Account) (*model.File, []model.File, error) {
	path = utils.ParsePath(path)
	log.Debugf("lanzou path: %s", path)
	file, err := driver.File(path, account)
	if err != nil {
		return nil, nil, err
	}
	if file.Type != conf.FOLDER {
		file.Url, _ = driver.Link(path, account)
		return file, nil, nil
	}
	files, err := driver.Files(path, account)
	if err != nil {
		return nil, nil, err
	}
	return nil, files, nil
}

func (driver Lanzou) Proxy(c *gin.Context, account *model.Account) {
	c.Request.Header.Del("Origin")
}

func (driver Lanzou) Preview(path string, account *model.Account) (interface{}, error) {
	return nil, NotSupport
}

var _ Driver = (*Lanzou)(nil)
