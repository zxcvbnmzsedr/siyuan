package api

import (
	"errors"
	"fmt"
	"github.com/88250/gulu"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/siyuan-note/siyuan/kernel/filesys"
	"github.com/siyuan-note/siyuan/kernel/util"
	"net/http"
	"path"
	"path/filepath"
	"time"
)

const (
	QiniuAk     = ""
	QiniuSk     = ""
	QiniuBucket = "siyuan-backup"
	QiniuDomain = ""
	QiniuBase   = "siyuan"
)

func getOssUploadToken(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}
	uid := c.Query("uid")
	filename := arg["name"].(string)
	dirPath := arg["dirPath"].(string)
	mac := qbox.NewMac(QiniuAk, QiniuSk)
	key := filepath.Join(QiniuBase, uid, dirPath, filename)
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", QiniuBucket, key),
	}
	token := putPolicy.UploadToken(mac)
	ret.Data = map[string]interface{}{
		"token": token,
	}
}

func getCloudSyncVer(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}
	uid := c.Query("uid")
	cloudDir := arg["syncDir"].(string)
	key := filepath.Join(QiniuBase, uid, "sync", cloudDir, ".siyuan", "conf.json")

	mac := qbox.NewMac(QiniuAk, QiniuSk)
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	downloadURL := storage.MakePrivateURL(mac, QiniuDomain, key, deadline)

	resp, err := util.NewCloudFileRequest15s("").Get(downloadURL)
	if nil != err {
		util.LogErrorf("download request [%s] failed: %s", downloadURL, err)
		return
	}
	if 404 == resp.StatusCode {
		util.LogInfof("first upload request [%s] ", downloadURL)
		ret.Data = map[string]interface{}{
			"v": -1,
		}
		return

	}
	if 200 != resp.StatusCode {
		util.LogErrorf("download request [%s] status code [%d]", downloadURL, resp.StatusCode)
		err = errors.New(fmt.Sprintf("download file list failed [%d]", resp.StatusCode))
		return
	}
	conf := &filesys.DataConf{}
	data, err := resp.ToBytes()
	if err = gulu.JSON.UnmarshalJSON(data, &conf); nil != err {
		util.LogErrorf("unmarshal index failed: %s", err)
		err = errors.New(fmt.Sprintf("unmarshal index failed"))
		return
	}
	ret.Data = map[string]interface{}{
		"v": conf.SyncVer,
	}
}

func getCloudFileListOSS(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}
	uid := c.Query("uid")
	cloudDirPath := arg["dirPath"].(string)

	mac := qbox.NewMac(QiniuAk, QiniuSk)
	key := path.Join(QiniuBase, uid, cloudDirPath, "index.json")
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	downloadURL := storage.MakePrivateURL(mac, QiniuDomain, key, deadline)
	ret.Data = map[string]interface{}{
		"url": downloadURL,
	}
}
