#!/bin/bash

OLDVERSION=$(grep 'github.com/gomorpheus/morpheus-fling/releases/download/' README.md | awk -F[/:] '{print $9}')
vi README.md
NEWVERSION=$(grep 'github.com/gomorpheus/morpheus-fling/releases/download/' README.md | awk -F[/:] '{print $9}')
if [ "$OLDVERSION" == "$NEWVERSION" ]
then
    echo 'Old version $OLDVERSION same as $NEWVERSION'
    read ANS
fi
git tag -a $NEWVERSION -m v$NEWVERSION
echo "Ready to push $NEWVERSION (cntrl-c to quit)?"
read ANS
git commit -m "Version $NEWVERSION"
git push origin master
git tag --force  -a $NEWVERSION -m v$NEWVERSION
git push origin --tags