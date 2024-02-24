#!/usr/bin/env bash

rm webdavjs.html
cat > webdavjs.html <<EOL
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>GoWebDAV</title>
  </head>
  <body></body>
  <script>

EOL

curl -s https://cdn.jsdelivr.net/gh/dom111/webdav-js/src/webdav-min.js >> webdavjs.html

echo -e "\n</script><style>\n" >> webdavjs.html

curl -s https://cdn.jsdelivr.net/gh/dom111/webdav-js/assets/css/style-min.css >> webdavjs.html

echo -e "\n</style></html>" >> webdavjs.html