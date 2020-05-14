projectDir=$1
branch=`git rev-parse --abbrev-ref HEAD`

if [ "$branch" = "master" ]; then
  if [[ -z "$(git diff HEAD^ HEAD $projectDir)" ]]; then
    echo "NO"
  else
    echo "YES"
  fi
elif git diff --name-only origin/master...$branch  | grep "^${projectDir}" ; then
  echo "YES"
else
  echo "NO"
fi