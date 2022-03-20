# [ SSHWHO ]

## Overview

This program aims to parse `auth.log` file to identify connected users
according to a public keys inventory.

It can be useful to audit `root` user connection on a server based 
on a set of known keys.

## Build

### Manual

```bash
$ go build && go install
```

### NixOS

```bash
$ cat default.nix
{ lib, buildGoModule, fetchFromGitHub }:

buildGoModule rec {
  pname = "sshwho";
  version = "0.1.0";

  src = fetchFromGitHub {
    owner = "bib0x";
    repo = "sshwho";
    rev = "v${version}";
    sha256 = "10qzzbrnx7yzvrxvshx0nmd8skx9w4wbkbak00f82s1nyxii9jrz";
  };

  vendorSha256 = "03jnps2wg4n4bw2sy8r9mm83apnlw4fs1cn9nsmbbrr41s8i4b54";

  meta = with lib; {
    description = "SSH auth.log analyzer based on SSH public keys inventory";
    homepage = "https://github.com/bib0x/sshwho";
    license = licenses.mit;
    maintainers = with maintainers; [ bib0x ];
    platforms = platforms.linux;
  };

}
```

## Configuration

### Environment variable

```
export SSHWHO_INVENTORY=/tmp/sshwho.json
export SSHWHO_KEYPATH=$HOME/dev/lab/keys
```

## Example

### Inventory

```
$ export SSHWHO_INVENTORY=/tmp/sshwho.json
$ export SSHWHO_KEYPATH=/tmp/keys

# Create inventory based on default configuration
$ sshwho inv

# Create inventory based on a file
$ sshwho inv -k ~/.ssh/authorized_keys

# Create inventory based on a directory with several key files
$ sshwho inv -k /tmp/keys

# Create a JSON inventory based on a directory
$ sshwho inv -k /tmp/keys -j

# Persist an inventory for cache
$ sshwho inv -k /tmp/keys -j /tmp/sshwho_cache.json
```

### Analyze logs

```
# Analyze remote /home/user/auth.log file over SSH
$ sshwho analyze -f /home/user/auth.log -k /tmp/keys -s user@10.0.0.2
[+] Oct 17 21:24:10 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 17 21:41:15 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:05:59 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:06:10 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:06:38 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:06:44 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:09:50 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:21:45 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:25:28 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:25:38 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 10:34:42 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 10:53:19 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Analyze local file using a cache file
$ sshwho analyze -f `pwd`/auth.log -i /tmp/sshwho_cache.json
[+] Oct 17 21:24:10 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 17 21:41:15 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:05:59 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:06:10 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:06:38 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:06:44 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:09:50 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:21:45 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:25:28 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 09:25:38 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 10:34:42 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
[+] Oct 18 10:53:19 user@virtlab connected to vmtest from 10.0.0.5 using SHA256:wfO7Fxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Todo

- Fix bug and CLI
- Improve CLI help usage()
- Compile and release package
- Add Tests
