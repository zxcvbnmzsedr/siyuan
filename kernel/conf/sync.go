// SiYuan - Build Your Eternal Digital Garden
// Copyright (c) 2020-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package conf

import (
	"path/filepath"

	"github.com/siyuan-note/siyuan/kernel/util"
)

type Sync struct {
	CloudName  string    `json:"cloudName"`  // 云端同步目录名称
	Enabled    bool      `json:"enabled"`    // 是否开启同步
	Uploaded   int64     `json:"uploaded"`   // 最近上传时间
	Downloaded int64     `json:"downloaded"` // 最近下载时间
	Synced     int64     `json:"synced"`     // 最近同步时间
	Stat       string    `json:"stat"`       // 最近同步统计信息
	QiniuSync  QiniuSync `json:"qiniuSync"`  // 三方同步秘钥
}

type QiniuSync struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}

func NewSync() *Sync {
	return &Sync{
		CloudName: "main",
		Enabled:   true,
	}
}

func (s *Sync) GetSaveDir() string {
	return filepath.Join(util.WorkspaceDir, "sync")
}
