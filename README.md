logptn: Text Log Pattern Analyzer
====================

## Setup

```
$ go get github.com/m-mizutani/logptn
```

## Usage

```
$ logptn <log file>
```

Example:

```
$ logptn /var/log/auth.log
2018/05/03 15:09:54 Jul  9 *:*:* pylon sshd[*]: Invalid user * from *
2018/05/03 15:09:54 Jul  9 *:*:* pylon sshd[*]: input_userauth_request: invalid user * [preauth]
2018/05/03 15:09:54 Jul  9 *:*:* pylon sshd[*]: pam_unix(sshd:auth): check pass; user unknown
2018/05/03 15:09:54 Jul  9 *:*:* pylon sshd[*]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=*
2018/05/03 15:09:54 Jul  9 *:*:* pylon sshd[*]: Failed password for invalid user * from * port * ssh2
2018/05/03 15:09:54 Jul  9 *:*:* pylon CRON[*]: pam_unix(cron:session): session closed for user *
```
