package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "golang.org/x/crypto/ssh"
)

type KeyFingerprint struct {
    Username string
    Filename string
    MD5      string
    SHA256   string
}

type Inventory struct {
    Fingerprints []KeyFingerprint
    Keypath      string
    Path         string
}

func NewKeyFingerprint(keypath string, keydata []byte) (KeyFingerprint, error) {
    pk, comment, _, _, err := ssh.ParseAuthorizedKey(keydata)
    
    if err != nil {
        return KeyFingerprint{}, err
    }
    
    if len(comment) == 0 {
        comment = "<unknown>"
    }

    return KeyFingerprint{
        Filename: filepath.Base(keypath),
        Username: comment,
        MD5: FingerprintHash(pk, "md5"),
        SHA256: FingerprintHash(pk, "sha256"),
    }, nil
}

func NewInventory(keypath, invpath string) Inventory {
    // if inventory exists (use cache)
    if _, err := os.Stat(invpath); err == nil {
        inventory, err := NewInventoryFromFile(invpath)
        if err == nil {
            return inventory
        }
    }
    // else create an inventory from scratch
    inventory, err := NewInventoryFromScratch(keypath, invpath)
    if err != nil {
        inventory = Inventory{}
    }
    return inventory
}

func NewInventoryFromScratch(keypath string, invpath string) (Inventory, error) {
    kpath, _ := filepath.Abs(keypath)
    ipath, _ := filepath.Abs(invpath)

    inventory := Inventory{
        Keypath: kpath,
        Path: ipath,
    }

    err := filepath.Walk(kpath, func(path string, info os.FileInfo, err error) error {
        if err != nil { return err }  
        if info.Mode().IsRegular() {
            // Open keyfile
            f, err := os.Open(path)
            if err != nil { return err }
            defer f.Close()
            // Read data
            scanner := bufio.NewScanner(f)
            for scanner.Scan() {
                if fp, err := NewKeyFingerprint(path, scanner.Bytes()); err == nil {
                    inventory.Fingerprints = append(inventory.Fingerprints, fp) 
                }
            }
        }
        return nil 
    })

    return inventory, err
}

func NewInventoryFromFile(invpath string) (Inventory, error) {
    path, _ := filepath.Abs(invpath)
    inventory := Inventory{ Path: path }

    raw, err := ioutil.ReadFile(path) 
    if err !=nil {
        return inventory, err
    }

    json.Unmarshal(raw, &inventory.Fingerprints)
    return inventory, nil
}

func (i *Inventory) ToJsonFile() error {
    if len(i.Fingerprints) == 0 {
         return fmt.Errorf("Empty inventory (no keys found in %s)", i.Keypath)
    }
    data, _ := json.Marshal(i.Fingerprints)
    return ioutil.WriteFile(i.Path, data, 0644)
}

func (i *Inventory) AsText() {
    if len(i.Fingerprints) > 0 {
        for _, fp := range i.Fingerprints {
            fmt.Printf("%s: %s %s %s\n", fp.Filename, fp.Username, fp.MD5, fp.SHA256)
        }
    } else {
        fmt.Printf("Empty inventory (no keys found in %s)\n", i.Keypath)
    }
}

func (i *Inventory) AsJson() {
    data, _ := json.Marshal(i.Fingerprints)
    if len(data) > 0 {
        fmt.Printf("%s\n", string(data))
    }
}

func (i *Inventory) Info() {
    status := "Found"
    if _, err := os.Stat(i.Path); err != nil {
        status = "Not Found"
    }

    fmt.Printf("Inventory: %s\n", i.Path)
    fmt.Printf("Status:    %s\n", status)
}
