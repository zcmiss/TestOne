package config

import . "os"

// 是开发模式
var IsDeveloper = true

// 是移动端服务
var IsMobileService = false

// 是管理端服务
var IsBackendService = false

// 是采集任务服务
var IsCollectTaskService = false

// 初始化
func init() {
	if Getenv("IS_DEVELOPER") == "YES" {
		IsDeveloper = true
	}
	if Getenv("IS_MOBILE_SERVICE") == "YES" {
		IsMobileService = true
	} else {
		IsBackendService = true
	}
	if Getenv("IS_COLLECT_TASK_SERVICE") == "YES" {
		IsCollectTaskService = true
	}
}
