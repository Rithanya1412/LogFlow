package model

import "time"

type LogEvent struct {
    Source    string    // Application, System
    Provider  string    // App or service name
    Level     string    // info, warn, error, critical
    Message   string
    EventID   uint32
    Timestamp time.Time
    Host      string
}
