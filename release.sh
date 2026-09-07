#!/usr/bin/env bash

# 发版：按 <上游版本>-<上游 commit>-adaptation 的格式打标签并推送，触发 GitHub Actions 构建产物
# 用法：./release.sh [上游分支，默认 v2]

set -e

UPSTREAM_REF="${1:-v2}"
RELEASE_BRANCH="v2-adaptation"
SUFFIX="adaptation"
REMOTE="origin"

current_branch=$(git rev-parse --abbrev-ref HEAD)
if [ "$current_branch" != "$RELEASE_BRANCH" ]; then
    echo "Current branch is ${current_branch}, expected ${RELEASE_BRANCH}"
    exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
    echo "Working tree is dirty, commit or stash first"
    exit 1
fi

if ! git rev-parse --verify --quiet "$UPSTREAM_REF" > /dev/null; then
    echo "Upstream ref ${UPSTREAM_REF} not found"
    exit 1
fi

base_version=$(git describe --tags --abbrev=0 --match 'v2*' --exclude "*-${SUFFIX}" "$UPSTREAM_REF")
upstream_commit=$(git rev-parse --short=7 "$UPSTREAM_REF")
tag="${base_version}-${upstream_commit}-${SUFFIX}"

if git rev-parse --verify --quiet "refs/tags/${tag}" > /dev/null; then
    echo "Tag ${tag} already exists"
    exit 1
fi

echo "Tag:      ${tag}"
echo "Commit:   $(git rev-parse --short HEAD) ($(git log -1 --format=%s))"
echo "Upstream: ${UPSTREAM_REF} -> ${upstream_commit}"
echo ""
read -r -p "Press Enter to tag and push, Ctrl-C to abort "

git tag -a "$tag" -m "release ${tag}"
git push "$REMOTE" "$tag"

echo ""
echo "Pushed ${tag}, watch the release workflow:"
echo "  https://github.com/liasica/swag/actions"
