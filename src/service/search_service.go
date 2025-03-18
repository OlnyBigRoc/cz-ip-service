package service

import (
	"cz-ip-service/src/constant"
	"cz-ip-service/src/czdb"
	"cz-ip-service/src/utils"
	"cz-ip-service/src/vo"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SearchService struct {
	Config        *vo.Config
	LatestVersion string
	IPV4Search    *czdb.DbSearcher
	IPV6Search    *czdb.DbSearcher
}

func NewSearchService() *SearchService {
	service := &SearchService{}
	service.Config = &vo.Config{}
	err := env.Parse(service.Config)
	if err != nil {
		panic(err)
	}
	log, err := service.UpdateDBFile()
	if err != nil {
		panic(err)
	}
	utils.Log.Infof("czdb init complete, log: %v", log)
	service.LatestVersion, err = service.GetCZLatestVersion()
	if err != nil {
		panic(err)
	}
	utils.Log.Infof("czdb init complete, latestVersion: %v", service.LatestVersion)
	return service
}

func (s *SearchService) CheckIpVersion(ip string) constant.DbType {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// 无效的 IP 地址
		return constant.Invalid
	}
	if parsedIP.To4() != nil {
		return constant.IPV4
	}
	return constant.IPV6
}

func (s *SearchService) _DownloadDBFile() (datePath string, err error) {
	datePath = time.Now().Format("20060102")
	czdbFile := filepath.Join(s.Config.DbPath, datePath, "czdb.zip")
	dbPath := filepath.Join(s.Config.DbPath, datePath)
	// 检查文件夹是否存在
	exist, err := utils.IsDirExist(filepath.Join(s.Config.DbPath, datePath))
	if err != nil {
		return
	}
	// 如果文件夹不存在，则创建文件夹
	if !exist {
		err = utils.CreateDir(dbPath)
		if err != nil {
			return
		}
	} else {
		return
	}
	// 下载文件
	err = utils.DownloadFile(fmt.Sprintf("%s%s", constant.DBFileURL, s.Config.FileKey), czdbFile)
	if err != nil {
		return
	}
	// 解压文件
	err = utils.DecompressZip(czdbFile, dbPath)
	if err != nil {
		return
	}
	return
}
func (s *SearchService) _DeleteHistoryFile(dateDir string) error {
	// 检查文件夹是否存在
	exist, err := utils.IsDirExist(s.Config.DbPath)
	if err != nil {
		return err
	}
	// 删除除newPath外的所有文件
	if exist {
		var files []os.DirEntry
		files, err = utils.ListDir(s.Config.DbPath)
		if err != nil {
			return err
		}
		for _, file := range files {
			var info fs.FileInfo
			info, err = file.Info()
			if err != nil {
				return err
			}
			if info.Name() == dateDir {
				continue
			}
			err = utils.DeleteDir(filepath.Join(s.Config.DbPath, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *SearchService) _LoadDBFile(newDir string) error {
	// 检查文件是否存在
	existV4, err := utils.IsFileExist(filepath.Join(s.Config.DbPath, newDir, s.Config.V4File))
	if err != nil {
		return err
	}
	existV6, err := utils.IsFileExist(filepath.Join(s.Config.DbPath, newDir, s.Config.V6File))
	if err != nil {
		return err
	}
	if existV4 {
		s.IPV4Search, err = czdb.NewDbSearcher(filepath.Join(s.Config.DbPath, newDir, s.Config.V4File), constant.MEMORY, s.Config.SecretKey)
		if err != nil {
			panic(err)
		}

	}
	if existV6 {
		s.IPV6Search, err = czdb.NewDbSearcher(filepath.Join(s.Config.DbPath, newDir, s.Config.V6File), constant.MEMORY, s.Config.SecretKey)
		if err != nil {
			panic(err)
		}

	}
	return nil
}

func (s *SearchService) _UnLoadDBFile() error {
	if s.IPV6Search != nil {
		s.IPV6Search.Close()
		s.IPV6Search = nil
	}
	if s.IPV4Search != nil {
		s.IPV4Search.Close()
		s.IPV4Search = nil
	}
	return nil
}

func (s *SearchService) Search(ctx *gin.Context, ip string) (*vo.IPInfo, error) {
	start := time.Now()
	if strings.Trim(ip, " ") == "" {
		ip = ctx.ClientIP()
	}
	// 检查IP类型
	ipType := s.CheckIpVersion(ip)
	if ipType == constant.Invalid {
		return nil, nil
	}
	var err error
	var ipInfoStr string
	if ipType == constant.IPV4 {
		ipInfoStr, err = s.IPV4Search.Search(ip)
		if err != nil {
			return nil, err
		}
	} else if ipType == constant.IPV6 {
		ipInfoStr, err = s.IPV6Search.Search(ip)
		if err != nil {
			return nil, err
		}
	}
	ipInfoVO := vo.NewInfo(ipInfoStr, ip)
	ipInfoVO.Time = time.Since(start).Microseconds()
	return ipInfoVO, nil
}

func (s *SearchService) UpdateDBFile() (string, error) {
	// 检查是否需要更新数据库
	cv, err := s.CheckUpdate()
	if err != nil {
		return fmt.Sprintf("检查更新数据库错误：%v", err.Error()), err
	}
	if !cv.IsUpdate { // 不需要更新
		return cv.Msg, nil
	}

	startTime := time.Now()
	// 下载最新的数据库
	datePath, err := s._DownloadDBFile()
	if err != nil {
		return fmt.Sprintf("下载数据库文件错误：%v", err.Error()), err
	}
	// 卸载数据库
	err = s._UnLoadDBFile()
	if err != nil {
		return fmt.Sprintf("卸载数据库错误：%v", err.Error()), err
	}
	// 加载数据库
	err = s._LoadDBFile(datePath)
	if err != nil {
		return fmt.Sprintf("加载数据库错误：%v", err.Error()), err
	}
	// 删除历史文件
	err = s._DeleteHistoryFile(datePath)
	if err != nil {
		return fmt.Sprintf("删除历史文件错误：%v", err.Error()), err
	}
	// 计算耗时
	costTime := time.Since(startTime)
	// 更新版本号
	s.LatestVersion = cv.NewVersion
	// 输出耗时
	return fmt.Sprintf("更新数据库成功，耗时：%v", costTime), nil
}

// CheckUpdate 检查更新 string 不为空
func (s *SearchService) CheckUpdate() (*vo.CheckVersion, error) { //
	// 获取最新版本号
	cv := &vo.CheckVersion{
		IsUpdate:   false,
		Msg:        "",
		OldVersion: s.LatestVersion,
		NewVersion: "",
	}
	latestVersion, err := s.GetCZLatestVersion()
	if err != nil {
		return cv, err
	}
	// 比较版本号
	if latestVersion != s.LatestVersion {
		// 版本号不同，返回更新信息
		cv.IsUpdate = true
		cv.NewVersion = latestVersion
		cv.Msg = fmt.Sprintf("最新版本号：%v，当前版本号：%v，需要更新", latestVersion, s.LatestVersion)
		return cv, nil
	}
	cv.Msg = fmt.Sprintf("最新版本号：%v，当前版本号：%v，不需要更新", latestVersion, s.LatestVersion)
	return cv, nil
}
func (s *SearchService) GetCZLatestVersion() (string, error) {
	// https://cz88.net/api/communityIpVersions/getLatestVersion
	bodyMap := map[string]interface{}{
		"APPCODE": s.Config.DeveloperKey,
	}
	body, err := utils.PostJsonForBody("https://cz88.net/api/communityIpVersions/getLatestVersion", bodyMap)
	if err != nil {
		return "", err
	}
	version := vo.Version{}
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", err
	}
	if err = version.CheckError(); err != nil {
		return "", err
	}
	return version.Data, nil
}
