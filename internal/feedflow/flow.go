package feedflow

import "rss-reader/internal/feedstate"

type Creator struct{ Builder *feedstate.Builder }

// Create 构建一条订阅连接描述符。
//
// 失败语义：构建失败时返回错误且不返回半成品描述符（value 为 nil），
// 缓存中也不会残留任何半成品（见 Builder.Build），因此后续对相同订阅连接
// 或其他订阅连接的查询都不会因本次失败而继续报错或读到脏数据。
func (c *Creator) Create(id string, fail bool) (value *feedstate.Descriptor, err error) {
	defer func() {
		if r := recover(); r != nil {
			// 失败已被 Builder.Build 处理：缓存里没有半成品。
			// 这里把 panic 转成错误，并确保不向调用方返回半成品。
			value = nil
			err = feedstate.ErrBuild
		}
	}()
	return c.Builder.Build(id, fail), nil
}
