package routes

import (
	"context"

	"log-service/data"
	"log-service/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	// Read the body from the request
	input := req.GetLogEntry()

	// Write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// Insert the log entry into the database
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// Return the response
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}
