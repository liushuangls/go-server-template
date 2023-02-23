package logger

import "testing"

func TestNewLogger(t *testing.T) {
	config := Config{
		FileDir:    "log",
		MaxSize:    1,
		MaxAge:     30,
		MaxBackups: 5,
		Compress:   false,
	}
	log := NewLogger(&config)
	for i := 0; i < 10; i++ {
		log.Infof("log :%s%d", "记录日志", i)
		log.Errorf("log :%s%d", "记录日志", i)
	}
}
