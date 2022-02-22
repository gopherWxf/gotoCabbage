package log_test

import (
	"gotoCabbage/utils/log"
	"testing"
)

func TestSetup(t *testing.T) {
	log.Setup(&log.Settings{
		Path:       "logs",
		Name:       "filename",
		Ext:        "log",
		TimeFormat: "2006-01-02",
		MaxBackups: 20, //最多MaxBackups个文件
		MaxSize:    10, //文件最大MaxSize(k)
		Cnt:        1,  //从filename_cnt开始，后期做config配置
	})
	for i := 0; i < 10000; i++ {
		log.Error("ERROR")
	}
	log.Fatal("stop")
}
