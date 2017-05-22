# !/bin/sh
LATEST_CODEGEN=`git ls-remote https://github.com/capitalone/oas-nodegen.git | grep refs/heads/master | cut -f 1`
sed \
-e "s:oas-nodegen.git#[a-zA-Z0-9]\{40\}:oas-nodegen.git#$LATEST_CODEGEN:g" \
package.json > package.json.tmp
cp package.json.tmp package.json
rm package.json.tmp
npm install