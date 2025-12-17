#!/bin/bash

backend='software_design_backend'

# 关闭后端服务
sudo pkill -f ${backend}

if [ $? -eq 0 ]; then
	echo 'successfully close backend'
else
	echo 'failed to close backend'
	exit 1
fi

# 关闭前端服务
sudo nginx -s stop

if [ $? -eq 0 ]; then
	echo 'successfully close fronted'
else
	echo 'failed to close fronted'
	exit 1
fi

echo 'successfully stop project'

exit 0
