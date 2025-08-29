#!/usr/bin/env bash
set -euo pipefail

# Rewrites all commit authors/committers in the current repo to
#   Name:  social protocol labs
#   Email: opensource@socialprotocol.dev
#
# Run this script INSIDE the new verifier-kotlin repo after migration.
#
# Note: This operation rewrites history. Ensure you are on a throwaway clone
# or are comfortable force-pushing and breaking existing clones.

NAME="social protocol labs"
EMAIL="opensource@socialprotocol.dev"

if ! git rev-parse --git-dir >/dev/null 2>&1; then
  echo "Error: not inside a Git repository." >&2
  exit 1
fi

echo "Rewriting authors/committers to: ${NAME} <${EMAIL}> ..."

# Prefer git filter-repo if available (faster, safer)
if command -v git-filter-repo >/dev/null 2>&1; then
  git filter-repo --force --mailmap <(cat <<MAILMAP
${NAME} <${EMAIL}> <>
MAILMAP
  )
  echo "Done via git filter-repo."
  exit 0
fi

# Fallback to filter-branch (built-in, slower)
git filter-branch -f --env-filter "
export GIT_AUTHOR_NAME='${NAME}'
export GIT_AUTHOR_EMAIL='${EMAIL}'
export GIT_COMMITTER_NAME='${NAME}'
export GIT_COMMITTER_EMAIL='${EMAIL}'
" -- --all

echo "Done via git filter-branch. Remember to force-push: git push -f"

