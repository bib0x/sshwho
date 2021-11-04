package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"

    "golang.org/x/crypto/ssh"
)

const (
  Month         int = iota + 1
  Day
  Hour
  Server
  _
  Type
  User
  Ip
  Port
  Version
  Algo
  Hash
)

type Log struct {
   Date        string
   Target      string
   Username    string
   From        string
   Port        int
   Algo        string
   Fingerprint string
   Type        string
}

type Logs []Log

func NewLog(line string) (Log, error) {
    regPattern := `(?P<month>.*)\s+(?P<day>.*) (?P<hour>.*) (?P<hostname>.*) (?P<_>.*): Accepted (?P<type>.*) for (?P<user>.*) from (?P<ip>.*) port (?P<port>.*) (?P<version>.*): (?P<algo>.*) (?P<data>.*)`
    re := regexp.MustCompile(regPattern)
    matches := re.FindStringSubmatch(line)

    if len(matches) < Type {
        return Log{}, fmt.Errorf("Malformed log line, no type field.")
    }

    var match Log
    if matches[Type] == "publickey" || matches[Type] == "password" {
        port, err := strconv.Atoi(matches[Port])
        if err != nil { return Log{}, err }
        match = Log{
            Date: fmt.Sprintf("%s %s %s", matches[Month], matches[Day], matches[Hour]),
            Target: matches[Server],
            Username: matches[User],
            From: matches[Ip],
            Port: port,
            Type: matches[Type],
        }
        if matches[Type] == "publickey" {
            match.Algo = matches[Algo]
            match.Fingerprint = matches[Hash]
        }
    } else {
        err := fmt.Errorf("Log type not managed.")
        return Log{}, err
    }
    return match, nil
}

func NewLogs(logfile string, pattern string, client *ssh.Client) ([]Log, error) {
    if client != nil {
        return NewLogsFromSSH(logfile, pattern, client)
    }
    return NewLogsFromLocal(logfile, pattern)
}

func NewLogsFromSSH(logfile string, pattern string, client *ssh.Client) ([]Log, error) {
    session, err := client.NewSession()
    if err != nil { panic(err) }
    defer session.Close()

    var buf bytes.Buffer
    session.Stdout = &buf

    cmd := fmt.Sprintf("grep '%s' %s", pattern, logfile)
    err = session.Run(cmd)
    if err != nil { panic(err) }

    logs := make([]Log, 0, 10)
    lines := strings.Split(buf.String(), "\n")
    for _, line := range lines {
        if strings.Index(line, pattern) >= 0 {
            log, _ := NewLog(line)
            if (Log{}) != log {
                logs = append(logs, log)
            }
        }
    }
    return logs, nil
}

func NewLogsFromLocal(logfile string, pattern string) ([]Log, error) {
    f, err := os.Open(logfile)
    if err != nil { return nil, err }
    defer  f.Close()

    logs := make([]Log, 0, 10)
    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Index(line, pattern) >= 0 {
            log, _ := NewLog(line)
            if (Log{}) != log {
                logs = append(logs, log)
            }
        }
    }
    return logs, nil
}

func (m *Log) Format(username string) string {
    u := m.Username
    if len(username) > 0 {
        u = username
    }

    using := m.Type
    if m.Type == "publickey" {
        using = m.Fingerprint
    }

    return fmt.Sprintf("[+] %s %s connected to %s from %s using %s", m.Date, u, m.Target, m.From, using) 
}
