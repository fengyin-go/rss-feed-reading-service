package feedflow

import "rss-reader/internal/feedstate"

type Archiver struct{ Manager *feedstate.Manager }

// Run 批量生成 count 个导出文件。failAt 为注入写入失败的位置（<0 表示不注入）。
// 每个导出文件在写入完成后立即释放其占用的资源，避免后续文件因资源耗尽而失败；
// 任一文件写入失败时整批保留失败状态（不标记为已完成），但仍继续处理后续文件。
func (a *Archiver) Run(count, failAt int) error {
	var runErr error
	for i := 0; i < count; i++ {
		lease, err := a.Manager.Acquire()
		if err != nil {
			return err
		}
		// 注入写入失败：记录失败状态，但不中断后续文件的处理。
		if i == failAt {
			runErr = feedstate.ErrWrite
		}
		// 每个文件处理完毕后立即释放资源，确保后续文件可继续获取。
		lease.Close()
	}
	// 整批结束统一提交结果状态：存在写入失败时不标记已完成。
	return a.Manager.Finish(runErr)
}
