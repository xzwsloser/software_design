#!/bin/bash

backend='software_design_backend'

echo 'start service'

# 启动后端项目
pushd $(pwd)/backend

if [ -f "$(pwd)/bin/${backend}" ]; then
	echo 'backend program exists'
else
	echo 'backend program not exists'
	make depend
	make build
	if [ $? -eq 0 ]; then
		echo 'successfully build backend'
	else
		echo 'failed to build backend'
		exit 1
	fi
fi

make run > /dev/null 2>&1 &

if [ $? -eq  0 ]; then
	echo 'successfully start backend'
else
	echo 'failed to start backend'
	exit 1
fi

popd 

# 启动前端项目
pushd $(pwd)/fronted/run

if [ -d "$(pwd)/dist" ]; then
	echo 'successfully find dist dir'
else
	echo 'failed to find dist dir'
	exit 1
fi

if [ -f "$(pwd)/nginx.conf" ]; then
	echo 'successfully find nginx.conf'
else
	echo 'failed to find nginx.conf!'
	exit 1
fi

sudo chmod -R 755 $(pwd)/dist
sudo nginx -c $(pwd)/nginx.conf

if [ $? -eq 0 ]; then
	echo 'successfully start fronted'
else
	echo 'failed to start fronted'
	exit 1
fi

echo 'succesfully start project'
popd 

