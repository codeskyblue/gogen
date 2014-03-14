#!/bin/bash -
#
TMPDIR=$(mktemp -d)
test -z "$TMPDIR" && exit 1
trap '/bin/rm -fr ${TMPDIR:?}' EXIT

PROGRAM=$(basename $0)
if test $# -ne 2
then
	cat <<EOF
Usage:
	$PROGRAM <dir> table
Example:
	$PROGRAM ./ book

Author: skyblue,  caution: don't overite your files
EOF
exit 1
fi

mkdir -p "$1"; cd "$1"
git clone https://github.com/shxsun/gails-default $TMPDIR

# $PROGRAM [shxsun/gails] book name:string
table=$2
shift 2
Table=$(echo "$table" | perl -pe 's/.*/\u$&/')
fields="$@" # not done yet

mkdir -p controllers conf models
cpr(){
	cp -arv "$TMPDIR/$1" "$2"
}
cpr controllers/base_controller.go controllers
cpr controllers/base_controller.go controllers
cpr controllers/user.go controllers/${table}.go
cpr models/init.go models/init.go
cpr models/user.go models/${table}.go
cpr conf/app.conf conf
cpr main.go ./
cpr main_test.go ./
cpr .gitignore ./

P=${PWD#${GOPATH}/src/}
find -type f  | xargs -i sed -i -e "s/Book/${Table}/g" -e "s/book/${table}/g" -e "s;github.com/shxsun/gails;$P;g" {}
