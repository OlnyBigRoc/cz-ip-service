package service

import (
	"cz-ip-service/src/vo"
	"testing"
)

func getSearchService() *SearchService {
	ss := NewSearchService()
	ss.Config = &vo.Config{}
	ss.Config.SecretKey = ""
	ss.Config.FileKey = ""
	ss.Config.DbPath = "./cz_db"
	ss.Config.V4File = "cz88_public_v4.czdb"
	ss.Config.V6File = "cz88_public_v6.czdb"
	return ss
}
func TestDownloadDBFile(t *testing.T) {
	ss := getSearchService()
	newDir, err := ss._DownloadDBFile()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(newDir)
}

func TestDeleteDBFile(t *testing.T) {
	ss := getSearchService()
	err := ss._DeleteHistoryFile("20250306")
	if err != nil {
		t.Error(err)
		return
	}
}
