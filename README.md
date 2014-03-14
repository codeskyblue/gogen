gails
=====

beego + xorm skeleton, for quickly REST-API develop.

## how to use
	wget https://raw.github.com/shxsun/gails/master/gails
	chmod +x gails

	mv gails ~/bin/ # or /usr/local/bin/
	export PATH=$HOME/bin:$PATH

cd into gopath dir, to create the first app

	cd $GOPATH/src
	gails myapi user

start you app

	cd myapi; bee run
