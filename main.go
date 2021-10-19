package main

import (
    "flag"
    "fmt"
    "os"
)

const (
    INVENTORY = "/tmp/sshwho.json"
    KEYPATH   = "/tmp/keys"
)

type Config struct {
    InventoryInPath     string
    KeyPath             string
    Item                string

    LogFile             string
    Server              string

    AsJson              bool
}

func ProcessInventory(c *Config) {
    inventory := NewInventory(c.KeyPath, c.InventoryInPath)
    if c.AsJson {
        inventory.AsJson()
    } else {
        inventory.AsText()
    }
}

func ProcessAuthLog(c *Config) {
    if c.LogFile != "" {
        inventory := NewInventory(c.KeyPath, "")
        if c.KeyPath == KEYPATH {
            inventory = NewInventory(c.KeyPath, c.InventoryInPath)
        }

        var logs []Log

        if c.Server != "" {
           client := NewSshClient(c.Server)
           logs, _ = NewLogs(c.LogFile, "Accepted", client)
           client.Close()
        } else {
            logs, _ = NewLogs(c.LogFile, "Accepted", nil)
        }

        for _, l := range logs {
            message := l.Format("")
            for  _, fingerprint := range inventory.Fingerprints {
                if l.Fingerprint == fingerprint.SHA256  || l.Fingerprint == fingerprint.MD5 {
                    message = l.Format(fingerprint.Username)
                    fmt.Println(message)
                }
            }
        }
    }
}

func main() {
    // Default configuration
    var c = Config{ AsJson: false }
    
    default_keypath := os.Getenv("SSHWHO_KEYPATH")
    if default_keypath == "" {
        default_keypath = KEYPATH
    }

    default_inv := os.Getenv("SSHWHO_INVENTORY")
    if default_inv == "" {
        default_inv = INVENTORY
    }
    
    // Gather first argument and remove it
    // It enables to modify os.Args before calling flag.Parse()
    // in order to have git-style command line argument
    if len(os.Args) <= 1 {
        fmt.Printf("usage: sshwho <action> [ <option> ]+\n\n")
        fmt.Printf("Use -h for more information\n")
        os.Exit(1)
    }
    index := 1
    c.Item = os.Args[index]
    os.Args = append(os.Args[:index], os.Args[index+1:]...)


    // inventory
    flag.StringVar(&c.KeyPath, "k", default_keypath, "Key path as directory or file")
    flag.StringVar(&c.InventoryInPath, "i", default_inv, "JSON file where to read inventory data from")

    // log file
    flag.StringVar(&c.Server, "s", "", "SSH connection string to remote server")
    flag.StringVar(&c.LogFile, "f", "", "Auth log file path to analyze")

    // output
    flag.BoolVar(&c.AsJson, "j", false, "JSON output")


    flag.Usage = func() {
        fmt.Printf("usage: sshwho <action> [ <option> ]+\n\n")

        fmt.Printf("Actions available\n")
        fmt.Printf("-----------------\n- inv\n- analyze\n\n")

        fmt.Printf("Options associated to actions\n\n")

        fmt.Printf("inv:\n")
        fmt.Printf("  -k string\n")
        fmt.Printf("    	Key path as directory or file (default %s)\n", default_keypath)
        fmt.Printf("  -i string\n")
        fmt.Printf("    	JSON file where to read inventory data from (default %s)\n", default_inv) 
        fmt.Printf("  -j	JSON output\n\n")

        fmt.Printf("analyze:\n")
        fmt.Printf("  -k string\n")
        fmt.Printf("    	Key path as directory or file (default %s)\n", default_keypath)
        fmt.Printf("  -i string\n")
        fmt.Printf("    	JSON file where to read inventory data from (default %s)\n", default_inv) 
        fmt.Printf("  -f    string\n")
        fmt.Printf("        Auth log file path to analayze\n")
        fmt.Printf("  -s    string\n")
        fmt.Printf("        SSH connection string to remote server such as [user@]host[:port]\n")

    }

    flag.Parse()
    if c.Item == "inv" {
        ProcessInventory(&c)
    } else if c.Item == "a" || c.Item == "analyze" {
        ProcessAuthLog(&c)
    } else {
        flag.Usage()
    }
}
