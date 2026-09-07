#!/usr/bin/env bash

# 把上游模块路径替换为本 fork 的模块路径
# 同时处理两种形态：
#   1. import 路径 github.com/swaggo/swag
#   2. 生成的定义名 github_com_swaggo_swag（点号被 swag 转成下划线）

set -e

OLD_MODULE="github.com/swaggo/swag"
NEW_MODULE="github.com/liasica/swag"
OLD_DEFINITION="github_com_swaggo_swag"
NEW_DEFINITION="github_com_liasica_swag"

echo "Replacing ${OLD_MODULE} -> ${NEW_MODULE}"
echo "Replacing ${OLD_DEFINITION} -> ${NEW_DEFINITION}"
echo ""

files=$(grep -rlE "${OLD_MODULE}|${OLD_DEFINITION}" . \
    --exclude-dir=.git \
    --exclude-dir=vendor \
    --exclude-dir=node_modules \
    --exclude=rename.sh \
    --binary-files=without-match \
    2>/dev/null || true)

if [ -z "$files" ]; then
    echo "Nothing to replace"
    exit 0
fi

echo "Files to update:"
echo "$files"
echo ""

echo "$files" | xargs sed -i '' \
    -e "s|${OLD_MODULE}|${NEW_MODULE}|g" \
    -e "s|${OLD_DEFINITION}|${NEW_DEFINITION}|g"

# 替换后 import 分组内的字母序会被打乱，重新格式化
echo ""
echo "Running gofmt..."
unformatted=$(gofmt -l . 2>/dev/null || true)
if [ -n "$unformatted" ]; then
    echo "$unformatted" | xargs gofmt -w
    echo "$unformatted"
fi

echo ""
echo "Done. Affected files: $(echo "$files" | wc -l | tr -d ' ')"
echo ""
echo "Next steps:"
echo "  go mod tidy"
echo "  for d in example/*/; do (cd \"\$d\" && go mod tidy); done"
