package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"
)

func FingerprintHash(pubkey ssh.PublicKey, algo string) string {
    if algo == "md5" {
        return ssh.FingerprintLegacyMD5(pubkey)
    }
    return ssh.FingerprintSHA256(pubkey)
}

// user@host:port
// host:port
// host
func SplitSshConnStr(connStr string) (string, string, int) {
    user := "root"
    port := 22
    host := connStr

    // user@host
    if strings.Contains(host, "@") {
        parts := strings.Split(host, "@")
        user = parts[0]
        prefix := fmt.Sprintf("%s@", user) 
        host = strings.TrimPrefix(host, prefix)
    }

    // host:port
    if strings.Contains(host, ":") {
        parts := strings.Split(host, ":")
        p, err := strconv.Atoi(parts[1])
        if err != nil { panic(err) }
        port = p
        suffix := fmt.Sprintf(":%v", port)
        host = strings.TrimSuffix(host, suffix)
    }

    return user, host, port
}

func NewSshClient(connStr string) *ssh.Client {
    user, host, port := SplitSshConnStr(connStr)

    socket := os.Getenv("SSH_AUTH_SOCK")
    conn, err := net.Dial("unix", socket)
    if err != nil { panic(err) }

    agentClient := agent.NewClient(conn)
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.PublicKeysCallback(agentClient.Signers),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    target := fmt.Sprintf("%s:%v", host, port)
    client, err := ssh.Dial("tcp", target, config)
    if err != nil { panic(err) }

    return client
}
