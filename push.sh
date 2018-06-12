#!/bin/bash

if [[ $(git status -s) ]]
then 
	echo "There are pending changes, Please commit."
	exit 1;
fi

git add .
git commit -m "01:20 UPDATE - changes to layout"
git push -u origin master
