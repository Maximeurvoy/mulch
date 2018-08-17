package main

import (
	"fmt"

	"github.com/Xfennec/mulch"
)

// Log provides error/warning/etc helpers for a Hub
type Log struct {
	target string
	hub    *Hub
}

// NewLog creates a new log for the provided target and hub
// note: mulch.MessageNoTarget is an acceptable target
func NewLog(target string, hub *Hub) *Log {
	return &Log{
		target: target,
		hub:    hub,
	}
}

// Log is a low-level function for sending a Message
func (log *Log) Log(message *mulch.Message) {
	// TODO: use our own *log.Logger (see log.go in Nosee project)
	fmt.Printf("%s(%s): %s\n", message.Type, message.Target, message.Message)
	message.Target = log.target
	log.hub.Broadcast(message)
}

// Error sends a MessageError Message
func (log *Log) Error(message string) {
	log.Log(mulch.NewMessage(mulch.MessageError, log.target, message))
}

// Errorf sends a formated string MessageError Message
func (log *Log) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Error(msg)
}

// Warning sends a MessageWarning Message
func (log *Log) Warning(message string) {
	log.Log(mulch.NewMessage(mulch.MessageWarning, log.target, message))
}

// Warningf sends a formated string MessageWarning Message
func (log *Log) Warningf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Warning(msg)
}

// Info sends an MessageInfo Message
func (log *Log) Info(message string) {
	log.Log(mulch.NewMessage(mulch.MessageInfo, log.target, message))
}

// Infof sends a formated string MessageInfo Message
func (log *Log) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Info(msg)
}

// Trace sends an MessageTrace Message
func (log *Log) Trace(message string) {
	log.Log(mulch.NewMessage(mulch.MessageTrace, log.target, message))
}

// Tracef sends a formated string MessageTrace Message
func (log *Log) Tracef(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Trace(msg)
}

// Success sends an MessageSuccess Message
func (log *Log) Success(message string) {
	log.Log(mulch.NewMessage(mulch.MessageSuccess, log.target, message))
}

// Successf sends a formated string MessageSuccess Message
func (log *Log) Successf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Success(msg)
}

// Failure sends an MessageFailure Message
func (log *Log) Failure(message string) {
	log.Log(mulch.NewMessage(mulch.MessageFailure, log.target, message))
}

// Failuref sends a formated string MessageFailure Message
func (log *Log) Failuref(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Failure(msg)
}

// SetTarget change the current "sending" target
func (log *Log) SetTarget(target string) {
	// You can't send to "*", only listen. But NoTarget does the same
	// since since everybody receives it.
	if target == mulch.MessageAllTargets {
		target = mulch.MessageNoTarget
	}
	log.target = target
}
