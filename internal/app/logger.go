package app

import "github.com/DrusGalkin/forum-client/pkg/logger"

func RunLogger(choose bool) {
	if err := logger.InitLogger(choose); err != nil {
		panic(err)
	}
	defer logger.Logger.Sync()

	logger.Logger.Info("Логгер запущен на микросервисе forum-client")
}
