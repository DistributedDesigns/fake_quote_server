Fake Quote Server
=====

Mimics the behavior of the day trading quote server, which is only accessible on UVic VPN. Most of the implementation references this blog [post][go-tcp-server-blog]. Probably useful to local dev and testing.

AFAIK the quote server has no published API or other docs ([this][quote-server-client] is all I'm going on) so this implementation probably doesn't handle error conditions correctly. I'll update it as I figure out more of the behavior.

#### Installing from Docker
_TBD_

#### Local Usage
Start the server and stablish a socket connection and send a formatted req `stock_symbol,user_id`:
```bash
go run server.go

# in another window
echo "XYZ,cool_user" | nc localhost 4444

# server sends
729.99,XYZ,cool_user,1484459366,WFlaY29vbF91c2Vy77+9
```
Return format is `quote,stock_symbol,user_id,timestamp,cryptokey`.

[go-tcp-server-blog]: https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go
[quote-server-client]: http://www.ece.uvic.ca/~seng462/ProjectWebSite/ClientThread.py
