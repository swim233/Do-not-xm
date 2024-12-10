package counter

import (
	"learn/units/logger"
	"os"
)

type Counter struct{}

func XmCounter(chatID int64) {
	file, err := os.Open("xmcounter.csv")
	if err != nil {
		logger.Error("无法加载csv文件:%s", err)
		return
	}
	defer file.Close()
	if _, err := os.Stat("xmcounter.csv"); os.IsNotExist(err) {
		// 如果 .env 文件不存在，创建并写入默认值
		logger.Info(".env 文件不存在，正在创建...")

		// 创建并打开 .env 文件
		file, err := os.Create("xmcounter.csv")
		if err != nil {
			logger.Error("创建 csv 文件失败: %s", err)
		}
		defer file.Close()
	}
}
