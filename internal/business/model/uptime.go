package model

import "sync"

type UptimeTarget struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
	Code   int    `json:"code"`
	Ms     int64  `json:"ms"`
	Error  string `json:"error,omitempty"`
}

var uptimeHistory = struct {
	sync.RWMutex
	data map[string][]UptimeTarget
}{data: make(map[string][]UptimeTarget)}
