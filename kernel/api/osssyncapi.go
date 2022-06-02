package api

import (
	"errors"
	"fmt"
	"github.com/88250/gulu"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/siyuan-note/siyuan/kernel/filesys"
	"github.com/siyuan-note/siyuan/kernel/model"
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
	cloudDir := arg["syncDir"].(string)
	_, data, err := getCloudFileContent("sync", cloudDir, ".siyuan", "conf.json")
	if err != nil {
		ret.Data = map[string]interface{}{
			"v": -1,
		}
		return
	}
	conf := &filesys.DataConf{}
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

func getSiYuanWorkspace(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)
	syncFileInfo, data, err := getCloudFileContent("sync/"+model.Conf.Sync.CloudName, "index.json")
	m := make(map[string]interface{})
	if err = gulu.JSON.UnmarshalJSON(data, &m); nil != err {
		util.LogErrorf("unmarshal index failed: %s", err)
		err = errors.New(fmt.Sprintf("unmarshal index failed"))
	}
	syncSize := 0
	backupSize := 0

	for _, v := range m {
		syncSize += int(v.(map[string]interface{})["size"].(float64))
	}

	backupFileInfo, data, err := getCloudFileContent("backup", "index.json")
	if err = gulu.JSON.UnmarshalJSON(data, &m); nil != err {
		util.LogErrorf("unmarshal index failed: %s", err)
		err = errors.New(fmt.Sprintf("unmarshal index failed"))
	}
	for _, v := range m {
		backupSize += int(v.(map[string]interface{})["size"].(float64))
	}

	ret.Data = map[string]interface{}{
		"backup": map[string]interface{}{
			"size":    backupSize,
			"updated": storage.ParsePutTime(backupFileInfo.PutTime).Format("2006-01-02 15:04:05"),
		},
		// todo 复检大小
		"assetSize": 0,
		"size":      syncSize + backupSize,
		"sync": map[string]interface{}{
			"size":    syncSize,
			"updated": storage.ParsePutTime(syncFileInfo.PutTime).Format("2006-01-02 15:04:05"),
		},
	}
}

func listCloudSyncDirOSS(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)
	m := make(map[string]interface{})

	backupFileInfo, data, err := getCloudFileContent("backup", "index.json")
	if err = gulu.JSON.UnmarshalJSON(data, &m); nil != err {
		util.LogErrorf("unmarshal index failed: %s", err)
		err = errors.New(fmt.Sprintf("unmarshal index failed"))
	}
	backupSize := 0
	for _, v := range m {
		backupSize += int(v.(map[string]interface{})["size"].(float64))
	}

	ret.Data = map[string]interface{}{
		"dirs": []map[string]interface{}{
			{
				"name":    "main",
				"size":    backupSize,
				"updated": storage.ParsePutTime(backupFileInfo.PutTime).Format("2006-01-02 15:04:05"),
			},
		},
		"size": backupSize,
	}
}

func getSiYuanWorkspaceSync(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	m := make(map[string]interface{})
	_, data, err := getCloudFileContent("backup", "index.json")
	backupSize := 0
	if err = gulu.JSON.UnmarshalJSON(data, &m); nil != err {
		util.LogErrorf("unmarshal index failed: %s", err)
		err = errors.New(fmt.Sprintf("unmarshal index failed"))
	}
	for _, v := range m {
		backupSize += int(v.(map[string]interface{})["size"].(float64))
	}

	_, data, err = getCloudFileContent("sync", model.Conf.Sync.CloudName, ".siyuan", "conf.json")

	if err = gulu.JSON.UnmarshalJSON(data, &m); nil != err {
		util.LogErrorf("unmarshal index failed: %s", err)
		err = errors.New(fmt.Sprintf("unmarshal index failed"))
	}

	ret.Data = map[string]interface{}{
		"assetSize":  0,
		"backupSize": backupSize,
		"d":          m["device"],
		"v":          m["syncVer"],
	}

}

func getSiYuanFile(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)
	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}
	cloudDirPath := arg["path"].(string)

	uid := model.Conf.User.UserId
	mac := qbox.NewMac(QiniuAk, QiniuSk)
	key := path.Join(QiniuBase, uid, cloudDirPath)

	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	downloadURL := storage.MakePrivateURL(mac, QiniuDomain, key, deadline)
	ret.Data = map[string]interface{}{
		"url": downloadURL,
	}
}

func getCloudFileContent(elem ...string) (storage.FileInfo, []byte, error) {
	uid := model.Conf.User.UserId
	mac := qbox.NewMac(QiniuAk, QiniuSk)
	key := path.Join(QiniuBase, uid, filepath.Join(elem...))

	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	downloadURL := storage.MakePrivateURL(mac, QiniuDomain, key, deadline)

	resp, err := util.NewCloudFileRequest15s(model.Conf.System.NetworkProxy.String()).Get(downloadURL)
	if nil != err {
		util.LogErrorf("download request [%s] failed: %s", downloadURL, err)
	}
	if 404 == resp.StatusCode {
		util.LogInfof("first upload request [%s] ", downloadURL)
		return storage.FileInfo{}, nil, errors.New("404")
	}
	if 200 != resp.StatusCode {
		util.LogErrorf("download request [%s] status code [%d]", downloadURL, resp.StatusCode)
		err = errors.New(fmt.Sprintf("download file list failed [%d]", resp.StatusCode))
	}
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	fileInfo, sErr := bucketManager.Stat(QiniuBucket, key)
	if sErr != nil {
		fmt.Println(sErr)
	}
	bytes, err := resp.ToBytes()
	return fileInfo, bytes, err

}
