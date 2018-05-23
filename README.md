logptn: Text Log Pattern Analyzer
====================

## Setup

```
$ go get github.com/m-mizutani/logptn
```

## Usage

```
$ logptn <log file1> [log file2 ...]
```

Example:

```
$ logptn /var/log/secure
     4 [4ffb267b] Feb  1 *:*:* pylon sshd[*]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=*  user=root
     4 [845f4659] Feb  1 *:*:* pylon sshd[*]: Failed password for root from * port * ssh2
     6 [847ccf35] Feb  1 *:*:* pylon sshd[*]: Connection closed by * [preauth]
     3 [de051cd9] Feb  1 08:*:* pylon sshd[*]: Invalid user * from *
     2 [8e9e2a13] Feb  1 08:*:* pylon sshd[*]: input_userauth_request: invalid user * [preauth]
     2 [22190c74] Feb  1 08:*:* pylon sshd[*]: pam_unix(sshd:auth): check pass; user unknown
     2 [83fba2bf] Feb  1 08:*:* pylon sshd[*]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=192.168.0.3
     2 [f1ba83ea] Feb  1 08:*:* pylon sshd[*]: Failed password for invalid user * from 192.168.0.3 port * ssh2
     4 [e4a6f815] Feb  1 08:*:01 pylon CRON[*]: pam_unix(cron:session): session opened for user root by (uid=0)
     4 [5256845b] Feb  1 08:*:01 pylon CRON[*]: pam_unix(cron:session): session closed for user root
```

simple JSON output:

```json
$ ./logptn /var/log/secure -d sjson | jq . 
{
  "formats": [
    {
      "segments": [
        "Feb  1 ",
        null,
        ":",
        null,
        ":",
        null,
        " pylon sshd[",
        null,
        "]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=",
        null,
        "  user=root"
      ],
      "count": 4
    },
(snip)
```

Output as HTML heatmap:

```bash
$ ./logptn /var/log/secure -d heatmap -o mymap.html
$ open mymap.html # Mac
```

<img src="./doc/sample.png" width=1024 />


## Options

```
Usage:
  logptn [OPTIONS]

Application Options:
  -o, --output=                          Output file, '-' means stdout (default: -)
  -d, --dumper=[text|json|sjson|heatmap]
  -t, --threshold=
  -s, --delimiters=
  -c, --content=[log|format]
      --enable-regex

Help Options:
  -h, --help                             Show this help message
```
