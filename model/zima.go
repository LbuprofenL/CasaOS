/*
 * @Author: LinkLeong link@icewhale.org
 * @Date: 2022-05-13 18:15:46
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-08-01 18:32:57
 * @FilePath: /CasaOS/model/zima.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package model

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/maruel/natural"
)

type Path struct {
	Name       string                 `json:"name"`   // File name or document name
	Path       string                 `json:"path"`   // Full path to file or folder
	IsDir      bool                   `json:"is_dir"` // Is it a folder
	Date       time.Time              `json:"date"`
	Size       int64                  `json:"size"` // File Size
	Type       string                 `json:"type,omitempty"`
	Label      string                 `json:"label,omitempty"`
	Write      bool                   `json:"write"`
	Extensions map[string]interface{} `json:"extensions"`
}

type DeviceInfo struct {
	LanIpv4     []string `json:"lan_ipv4"`
	Port        int      `json:"port"`
	DeviceName  string   `json:"device_name"`
	DeviceModel string   `json:"device_model"`
	DeviceSN    string   `json:"device_sn"`
	Initialized bool     `json:"initialized"`
	OS_Version  string   `json:"os_version"`
	Hash        string   `json:"hash"`
}

// PathList 定义 Path 切片类型，用于实现排序功能
type PathList []Path

// getFileType 获取文件类型，用于 type 排序
func getFileType(p *Path) string {
	if p.IsDir {
		return "000_folder" // 使用前缀确保文件夹排在前面
	}

	// 从文件名提取扩展名
	ext := strings.ToLower(filepath.Ext(p.Name))
	if ext == "" {
		return "999_file" // 无扩展名文件排在有扩展名文件后面
	}
	return ext[1:] // 去掉点号返回扩展名
}

// Sort 为 PathList 实现排序功能
func (p PathList) Sort(orderBy, orderDirection string) {
	if orderBy == "" {
		return
	}

	sort.Slice(p, func(i, j int) bool {
		// 文件夹优先显示逻辑（除非按 type 排序）
		if orderBy != "type" {
			if p[i].IsDir && !p[j].IsDir {
				return true
			}
			if !p[i].IsDir && p[j].IsDir {
				return false
			}
		}

		var result bool
		switch orderBy {
		case "name":
			// 使用自然排序算法
			result = natural.Less(p[i].Name, p[j].Name)
		case "size":
			if p[i].Size == p[j].Size {
				// 大小相同时按名称排序
				result = natural.Less(p[i].Name, p[j].Name)
			} else {
				result = p[i].Size < p[j].Size
			}
		case "modified":
			if p[i].Date.Equal(p[j].Date) {
				// 时间相同时按名称排序
				result = natural.Less(p[i].Name, p[j].Name)
			} else {
				result = p[i].Date.Before(p[j].Date)
			}
		case "type":
			typeI := getFileType(&p[i])
			typeJ := getFileType(&p[j])
			if typeI == typeJ {
				// 类型相同时按名称排序
				result = natural.Less(p[i].Name, p[j].Name)
			} else {
				result = typeI < typeJ
			}
		default:
			// 未知排序字段，按名称排序
			result = natural.Less(p[i].Name, p[j].Name)
		}

		// 处理降序
		if orderDirection == "desc" {
			result = !result
		}

		return result
	})
}
