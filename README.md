# svxlink-pipe

This command reads svxlink data from stdin and render it in a webpage.

If you want to read data from a different machine over ssh you can use this command:

```bash
( ssh -t <remote-host> 'cat /tmp/sql_state' & tail -n1 -f /var/log/svxserver.log ) | ./svxlink-pipe
```
