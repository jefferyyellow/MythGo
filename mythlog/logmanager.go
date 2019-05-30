package mythlog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type LogManager struct {
	fatalLog Log
	errorLog Log
	infoLog  Log
	warnLog  Log
	debugLog map[string]*Log
	timeNow  time.Time
}

// 致命错误日志
func (pLogManager *LogManager) FatalLog() *Log {
	return &pLogManager.fatalLog
}

// 错误日志
func (pLogManager *LogManager) ErrorLog() *Log {
	return &pLogManager.errorLog
}

// 信息日志
func (pLogManager *LogManager) InfoLog() *Log {
	return &pLogManager.infoLog
}

// 警告日志
func (pLogManager *LogManager) WarnLog() *Log {
	return &pLogManager.warnLog
}

// 调试日志
func (pLogManager *LogManager) DebugLog(name string) *Log {
	return pLogManager.debugLog[name]
}

// 初始化日志管理器
func (pLogManager *LogManager) Init() {
	pLogManager.debugLog = make(map[string]*Log)
}

// 设置时间
func (pLogManager *LogManager) SetTime(tTimeNow time.Time) {
	pLogManager.timeNow = tTimeNow
}

// 增加debug日志
func (pLogManager *LogManager) AddDebugLog(strDebugName string, pLog *Log) {
	if nil == pLog {
		return
	}
	pLogManager.debugLog[strDebugName] = pLog
}

// 根据日志内容和类型最终输出日志到显示端
func (pLogManager *LogManager) LogMessage(location bool, pLog *Log, logContent string, Type string) {
	if nil == pLog {
		return
	}
	var strlocation string
	if location {
		pc, filename, line, ok := runtime.Caller(2)
		if ok {
			strlocation = fmt.Sprintf("[ %s: %d ][%s]", filepath.Base(filename), line, runtime.FuncForPC(pc).Name())
		}
	}

	year, month, day := pLogManager.timeNow.Date()
	hour, min, sec := pLogManager.timeNow.Clock()
	log := fmt.Sprintf("[%04d-%02d-%02d %02d:%02d:%02d] %s : ", year, month, day, hour, min, sec, Type) + strlocation + logContent + "\n"
	pLog.DisplayLog(log)
}

// 记录信息日志(不带文件/函数/行号的定位)
func (pLogManager *LogManager) LogInfoMessage(format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLogManager.LogMessage(false, &pLogManager.infoLog, logContent, "Info")
}

// 记录信息日志(有文件/函数/行号的定位)
func (pLogManager *LogManager) LogInfoMessageLocation(format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLogManager.LogMessage(true, &pLogManager.infoLog, logContent, "Info")
}

// 记录警告日志(有文件/函数/行号的定位)
func (pLogManager *LogManager) LogWarnMessage(format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLogManager.LogMessage(true, &pLogManager.warnLog, logContent, "Warn")
}

// 记录错误日志(有文件/函数/行号的定位)
func (pLogManager *LogManager) LogErrorMessage(format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLogManager.LogMessage(true, &pLogManager.errorLog, logContent, "Error")
}

// 记录致命日志(有文件/函数/行号的定位)
func (pLogManager *LogManager) LogFatalMessage(format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLogManager.LogMessage(true, &pLogManager.fatalLog, logContent, "Fatal")
}

// 记录调试日志(不带文件/函数/行号的定位)
func (pLogManager *LogManager) LogDebugMessage(name string, format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLog := pLogManager.debugLog[name]
	if nil == pLog {
		return
	}
	pLogManager.LogMessage(false, pLog, logContent, "Debug")
}

// 记录调试日志(有文件/函数/行号的定位)
func (pLogManager *LogManager) LogDebugMessageLocation(name string, format string, a ...interface{}) {
	logContent := fmt.Sprintf(format, a...)
	pLog := pLogManager.debugLog[name]
	if nil == pLog {
		return
	}
	pLogManager.LogMessage(true, pLog, logContent, "Debug")
}

// 关闭日志
func (pLogManager *LogManager) Close() {
	pLogManager.fatalLog.Close()
	pLogManager.errorLog.Close()
	pLogManager.infoLog.Close()
	pLogManager.warnLog.Close()
	for _, pLog := range pLogManager.debugLog {
		pLog.Close()
	}
}
