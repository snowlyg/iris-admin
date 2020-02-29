package libs

import "time"

/**
 * 时间格式化
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func TimeFormat(time *time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
