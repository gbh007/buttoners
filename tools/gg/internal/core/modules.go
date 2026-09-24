package core

import "slices"

const (
	ModuleLogs   = "logs"
	ModuleTraces = "traces"

	ModuleRedis = "redis"

	ModuleHTTPClient = "http_client"
	ModuleHTTPServer = "http_server"

	ModuleGRPCClient = "grpc_client"
	ModuleGRPCServer = "grpc_server"

	ModuleQueueReader = "queue_reader"
	ModuleQueueWriter = "queue_writer"
)

type ModuleType = string

func (c Core) HasModule(name ModuleType) bool {
	return slices.Contains(c.EnabledModules, name)
}

func (c Core) HasAnyModule(names ...ModuleType) bool {
	for _, name := range names {
		if slices.Contains(c.EnabledModules, name) {
			return true
		}
	}

	return false
}
