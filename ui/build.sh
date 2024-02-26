#!/usr/bin/env bash

make() {
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
rm webdavjs.html webdavjs-ro.html

# Prepare
if [ ! -d webdav-js ]; then
  git clone https://github.com/dom111/webdav-js
  cd webdav-js
  pnpm i
  cd ..
fi

# Make original
cd webdav-js
pnpm run build
cd ..
make webdavjs.html

# Make read-only
cd webdav-js
git apply ../hide.patch
pnpm run build
cd ..
make webdavjs-ro.html
cd webdav-js
git apply --reverse ../hide.patch
cd ..

cp webdavjs* ../internal/server