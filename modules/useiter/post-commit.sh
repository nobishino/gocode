#!/bin/zsh

echo "Running post-commit hook"

git log --oneline -1 > tmp
echo  "" >> tmp
gp share main.go go.mod go.sum >> tmp
echo  "" >> tmp
echo  "\`\`\`go" >> tmp
cat main.go >> tmp
echo  "\`\`\`" >> tmp

git commit --amend -m "$(cat tmp)"

cat tmp