package mythlog

const MAX_DISPLAY_SIZE int = 4

type Log struct {
	name          string
	displayerList [MAX_DISPLAY_SIZE]LogDisplayer
	displayerSize int
}

// 增加Displayer
func (pLog *Log) AddDisplayer(pDisplayer LogDisplayer) {
	if pLog.displayerSize >= MAX_DISPLAY_SIZE {
		return
	}
	if nil == pDisplayer {
		return
	}

	pLog.displayerList[pLog.displayerSize] = pDisplayer
	pLog.displayerSize++
}

// 得到指定索引的Displayer
func (pLog *Log) GetDisplayer(nIndex int) LogDisplayer {
	if nIndex >= MAX_DISPLAY_SIZE || nIndex < 0 {
		return nil
	}
	return pLog.displayerList[nIndex]
}

// 显示日志
func (pLog *Log) DisplayLog(strLog string) {
	for i := 0; i < pLog.displayerSize; i++ {
		pLog.displayerList[i].DisplayLog(strLog)
	}
}

// 日志关闭（这个很重要，退出时调用，出问题时调用，因为有些Displayer有缓存，需要刷入磁盘）
func (pLog *Log) Close() {
	for i := 0; i < pLog.displayerSize; i++ {
		pLog.displayerList[i].Close()
	}
}
