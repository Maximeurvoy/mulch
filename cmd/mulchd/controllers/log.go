package controllers

import "github.com/OnitiFR/mulch/common"
import "github.com/OnitiFR/mulch/cmd/mulchd/server"

// LogController sends logs to client
func LogController(req *server.Request) {
	req.Stream.Infof("Hi! You will receive all logs for all targets.")
	req.SetTarget(common.MessageAllTargets)
	// nothing to do, just wait forever…
	select {}
}
