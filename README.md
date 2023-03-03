chatgpt服务的香港服务器server，主要是作为大陆与美国服务器的中间转发服务，使得大陆可以不用翻墙访问美国服务器
部署在 ininpop-spider 服务器上（43.135.74.185:bshgT64iUhsdfgk）

解决流式响应问题，配置nginx：
在location下添加

proxy_buffering off;
proxy_cache off;
proxy_set_header Connection '';
proxy_http_version 1.1;
chunked_transfer_encoding off;
