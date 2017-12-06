// Pipe - A small and beautiful blogging platform written in golang.
// Copyright (C) 2017, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package console

import (
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"os"
	"io"
)

func uploadAction(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	session := util.GetSession(c)
	if nil == session {
		result.Code = -1
		result.Msg = "please login before upload"

		return
	}

	form, err := c.MultipartForm()
	if nil != err {
		msg := "parse upload file header failed"
		logger.Errorf(msg + ": " + err.Error())

		result.Code = -1
		result.Msg = msg

		return
	}

	file := form.File["file"][0]
	f, err := file.Open()
	if nil != err {
		msg := "open upload file failed"
		logger.Errorf(msg + ": " + err.Error())

		result.Code = -1
		result.Msg = msg

		return
	}
	defer f.Close()

	tempDir := os.TempDir()
	logger.Trace("temp dir path is [" + tempDir + "]")
	zipFilePath := filepath.Join(tempDir, session.UName+"-md.zip")
	zipFile,err := os.Create(zipFilePath)
	if nil != err {
		logger.Errorf("create temp file [" + zipFilePath + "] failed: " + err.Error())

		result.Code = -1
		result.Msg = "create temp file failed"

		return
	}
	_, err = io.Copy(zipFile, f)
	if nil !=err{
		logger.Errorf("write temp file [" + zipFilePath + "] failed: " + err.Error())

		result.Code = -1
		result.Msg = "write temp file failed"

		return
	}
	defer zipFile.Close()

	unzipPath := filepath.Join(tempDir, session.UName+"-md")
	if err = os.RemoveAll(unzipPath);nil!=err {
		logger.Errorf("remove temp dir [" + unzipPath + "] failed: " + err.Error())

		result.Code = -1
		result.Msg = "remove temp dir failed"

		return
	}
	if err = os.Mkdir(unzipPath, 0755);nil!=err{
		logger.Errorf("make temp dir [" + unzipPath + "] failed: " + err.Error())

		result.Code = -1
		result.Msg = "make temp dir failed"

		return
	}
	if err = util.Zip.Unzip(zipFilePath, unzipPath);nil!=err{
		logger.Errorf("unzip [" + zipFilePath + "] to [" + unzipPath + "] failed: " + err.Error())

		result.Code = -1
		result.Msg = "unzip failed"

		return
	}



	data := map[string]interface{}{}
	data["succMap"] = succMap
	data["errFiles"] = errFiles
	result.Data = data
}
