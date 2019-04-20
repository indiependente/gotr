# gotr
Golang implementation of the traceroute command (ICMP based)
 
 ### How to install
 ```
 git clone https://github.com/indiependente/gotr
 cd gotr
 make deps
 make build
 ```
 
 ### Usage
 ```
 $ ./gotr -h
GoTR - Golang implementation of the traceroute command (ICMP based) (root access required)

  Usage:
    GoTR [host]

  Positional Variables:
    host (Required) - hostname to traceroute
  Flags:
       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -t --ttl  Number of allowed hops
```
  ### Example
  ```
$ sudo ./gotr google.com -t 5
Password:
2019-04-20 20:19:22 Launching traceroute against google.com (172.217.168.174) 👁‍🗨
	#HOP	REMOTE IP		MSGLENGTH	NAMES
	1:	192.168.1.1		[36 bytes]	[liveboxplus]

	2:	192.169.255.254		[72 bytes]	[ip-192-169-255-254.ip.secureserver.net.]

	3:	10.145.83.2		[32 bytes]

	4:	172.217.168.174		[8 bytes]	[mad07s10-in-f14.1e100.net.]

Destination reached 🎉
Time elapsed : 50ms
```
