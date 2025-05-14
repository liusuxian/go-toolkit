/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-07-13 20:17:18
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2025-05-13 12:25:11
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package gtkoss

// DeleteObjects 删除多个对象
func (s *AliyunOSS) DeleteObjects(objectKeys ...string) (err error) {
	_, err = s.bucket.DeleteObjects(objectKeys)
	return
}
