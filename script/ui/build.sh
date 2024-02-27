#!/usr/bin/env bash

set -e

make_html() {
  cat > $1 <<EOL
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
  cat webdav-js/dist/webdav.js  >> $1
  echo -e "\n</script><style>\n" >> $1
  cat webdav-js/dist/webdav.css >> $1
  echo -e "\n</style></html>" >> $1
}

# Clean
cd "$(dirname "$0")"
rm -f webdavjs.html webdavjs-ro.html

# Prepare
if [ ! -d webdav-js ]; then
  git clone --depth 1 https://github.com/Jipok/webdav-js # TODO: check commit is 63f2817b15f7b1309da886a67341ed626a838b16
  cd webdav-js
  pnpm i
  cd ..
fi

# Make original
cd webdav-js
set +e
pnpm run build >/dev/null 2>&1
set -e
cd ..
make_html webdavjs.html

mv webdavjs* ../../internal/server