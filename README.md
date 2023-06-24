## Nginx log format to parse

```conf
log_format jsonformat escape=json
'{'
  '"time":$msec,'
  '"remote_addr": "$remote_addr",'
  '"request_method": "$request_method",'
  '"request_uri": "$request_uri",'
  '"uri": "$uri",'
  '"request_filename": "$request_filename",'
  '"body_bytes_sent": $body_bytes_sent,'
  '"http_referer": "$http_referer",'
  '"connection":"$connection",'
  '"request":"$request",'
  '"status":$status,'
  '"user_agent":"$http_user_agent",'
  '"request_time":$request_time,'
  '"upstream_addr":"$upstream_addr",'
  '"upstream_status":$upstream_status,'
  '"upstream_response_time": $upstream_response_time,'
  '"upstream_connect_time": $upstream_connect_time,'
  '"upstream_header_time": $upstream_header_time'
'}';
```