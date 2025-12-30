package sources

import (
    "encoding/json"
    "os/exec"
    "strings"
    "time"

    "logflow/cmd/agent/model"
)

type windowsEvent struct {
    TimeCreated      time.Time `json:"TimeCreated"`
    Id               uint32    `json:"Id"`
    LevelDisplayName string    `json:"LevelDisplayName"`
    ProviderName     string    `json:"ProviderName"`
    Message          string    `json:"Message"`
}

func normalizeLevel(level string) string {
    switch strings.ToLower(level) {
    case "information":
        return "info"
    case "warning":
        return "warn"
    case "error":
        return "error"
    case "critical":
        return "critical"
    default:
        return "unknown"
    }
}

func ReadApplicationLogs() ([]model.LogEvent, error) {
    cmd := exec.Command(
        "powershell",
        "-Command",
        `Get-WinEvent -LogName Application -MaxEvents 10 |
         Select TimeCreated, Id, LevelDisplayName, ProviderName, Message |
         ConvertTo-Json`,
    )

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var raw []windowsEvent
    if err := json.Unmarshal(output, &raw); err != nil {
        return nil, err
    }

    hostnameCmd := exec.Command("hostname")
    hostBytes, _ := hostnameCmd.Output()
    host := strings.TrimSpace(string(hostBytes))

    events := make([]model.LogEvent, 0)
    for _, e := range raw {
        events = append(events, model.LogEvent{
            Source:    "Application",
            Provider:  e.ProviderName,
            Level:     normalizeLevel(e.LevelDisplayName),
            Message:   e.Message,
            EventID:   e.Id,
            Timestamp: e.TimeCreated,
            Host:      host,
        })
    }

    return events, nil
}
