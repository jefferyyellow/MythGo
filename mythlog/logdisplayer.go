package mythlog

import (
	"bufio"
	"fmt"
	"os"
)

// LogDisplay是日志接口
type LogDisplayer interface {
	DisplayLog(strLog string)
	Close()
}

// 标准输出显示
type StdLogDisplayer struct {
}

func (pDisplayer *StdLogDisplayer) DisplayLog(strLog string) {
	fmt.Printf("%s", strLog)
}

func (pDisplayer *StdLogDisplayer) Close() {
}

// 回转文件显示
type RollFileDisplayer struct {
	baseFileName string        // 基本文将名字
	curFileSize  int32         // 当前文件大小
	maxFileSize  int32         // 最大文件大小
	maxBackNum   int16         // 大小的回转文件数量
	writer       *bufio.Writer // 缓存写入接口
	filePointer  *os.File      // 文件句柄
}

// 初始化基本名字，最大的文件大小和回转文件数量
func (pDisplayer *RollFileDisplayer) Init(baseFileName string, maxFileSize int32, maxBackNum int16) {

	filePointer, err := os.OpenFile(baseFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if nil != err {
		fmt.Printf("Init log file: %s, %v", baseFileName, err)
		return
	}

	pDisplayer.writer = bufio.NewWriter(filePointer)
	pDisplayer.filePointer = filePointer
	curFielSize, _ := filePointer.Seek(0, os.SEEK_END)
	pDisplayer.curFileSize = int32(curFielSize)
	pDisplayer.baseFileName = baseFileName
	pDisplayer.maxFileSize = maxFileSize
	pDisplayer.maxBackNum = maxBackNum
}

func (pDisplayer *RollFileDisplayer) DisplayLog(strlog string) {
	pDisplayer.writer.WriteString(strlog)
	pDisplayer.curFileSize += int32(len(strlog))
	// 当文件达到最大值，回转
	if pDisplayer.curFileSize > pDisplayer.maxFileSize {
		pDisplayer.rollOver()
	}
}

// 文件回转
func (pDisplayer *RollFileDisplayer) rollOver() {
	// 如果不需要回转，直接返回
	if pDisplayer.maxBackNum <= 0 {
		return
	}

	pDisplayer.writer.Flush()
	if nil != pDisplayer.filePointer {
		pDisplayer.filePointer.Close()
	}
	// 将最后一个删除
	newFileName := fmt.Sprintf("%s.%d", pDisplayer.baseFileName, pDisplayer.maxBackNum)
	os.Remove(newFileName)
	pDisplayer.curFileSize = 0

	// 从后到前开始重命名
	for i := pDisplayer.maxBackNum - 1; i >= 1; i-- {
		oldFileName := fmt.Sprintf("%s.%d", pDisplayer.baseFileName, i)
		os.Rename(oldFileName, newFileName)
		newFileName = oldFileName
	}
	os.Rename(pDisplayer.baseFileName, newFileName)

	filePointer, err := os.OpenFile(pDisplayer.baseFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// 如果打开玩家没有问题，重新赋值
	if nil == err {
		pDisplayer.writer = bufio.NewWriter(filePointer)
		pDisplayer.filePointer = filePointer
	}
}

// 结束的时候一定要调用，否则出现日志可能丢失的问题
func (pDisplayer *RollFileDisplayer) Close() {
	pDisplayer.writer.Flush()
	pDisplayer.curFileSize = 0
	if nil != pDisplayer.filePointer {
		pDisplayer.filePointer.Close()
	}
}
