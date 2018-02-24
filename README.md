# gotr
Golang implementation of the traceroute command (ICMP based)
 
 ### Usage
 `sudo ./gotr hostname [TTL]`
  - It needs to be run as superuser because it has to send ICMP messages
  - `hostname` can be both an IP address or an hostname
  - `TTL` sets the maximum number of hops between the source and the target host
  
  ### Example
  ```
$ sudo ./gotr google.com
2018/02/24 19:12:05 Launching traceroute against google.com (216.58.206.142)
						          #HOP	REMOTE IP			MSGLENGTH
2018/02/24 19:12:05 	1:		199.91.137.242		[32 bytes]
2018/02/24 19:12:05 	2:		195.66.224.125		[32 bytes]
2018/02/24 19:12:05 	3:		108.170.246.129		[36 bytes]
2018/02/24 19:12:05 	4:		216.239.56.193		[36 bytes]
2018/02/24 19:12:05 	5:		216.58.206.142		[8 bytes]	[lhr25s15-in-f14.1e100.net.]
2018/02/24 19:12:05 Time elapsed : 35ms
```
