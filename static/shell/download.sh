#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
LANG=en_US.UTF-8
GOFLY_VERSION="go-fly-pro0.5.5-linux64"
SETUP_PATH="/gofly"

Red_Error(){
	echo '=================================================';
	printf '\033[1;31;40m%b\033[0m\n' "$@";
	GetSysInfo
	exit 1;
}
Yellow_echo(){
	printf '\033[1;33m%b\033[0m\n' "$@";
}
GetSysInfo(){
	if [ -s "/etc/redhat-release" ];then
		SYS_VERSION=$(cat /etc/redhat-release)
	elif [ -s "/etc/issue" ]; then
		SYS_VERSION=$(cat /etc/issue)
	fi
	SYS_INFO=$(uname -a)
	SYS_BIT=$(getconf LONG_BIT)
	MEM_TOTAL=$(free -m|grep Mem|awk '{print $2}')
	CPU_INFO=$(getconf _NPROCESSORS_ONLN)

	echo  ${SYS_VERSION}
	echo  Bit:${SYS_BIT} Mem:${MEM_TOTAL}M Core:${CPU_INFO}
	echo  ${SYS_INFO}
	echo  "请截图以上报错信息向官网客服求助"
}


if [ $(whoami) != "root" ];then
	echo "请使用root权限执行GOFLY下载命令！"
	exit 1;
fi
is64bit=$(getconf LONG_BIT)
if [ "${is64bit}" != '64' ];then
	Red_Error "抱歉, 当前GOFLY版本不支持32位系统, 请使用64位系统!";
fi


echo "
+----------------------------------------------------------------------
| GOFLY-LIVE-CHAT FOR Linux
+----------------------------------------------------------------------
| Copyright © 2020-2099 BT-SOFT() All rights reserved.
+----------------------------------------------------------------------
| The GOFLY URL will be http://SERVER_IP:8081 when installed.
+----------------------------------------------------------------------
"
while [ "$go" != 'y' ] && [ "$go" != 'n' ]
do
	read -p "你想要安装 GOFLY 到 ${SETUP_PATH} 目录?(y/n): " go;
done

if [ "$go" == 'n' ];then
	exit;
fi

if [ ! -d ${SETUP_PATH} ];then
  mkdir ${SETUP_PATH}
fi
cd ${SETUP_PATH}

echo '---------------------------------------------';
echo "开始下载GOFLY压缩包...";

wget -O ${GOFLY_VERSION}".zip" "https://gofly-1304282073.cos.ap-nanjing.myqcloud.com/${GOFLY_VERSION}.zip"
unzip -o ${GOFLY_VERSION}".zip" -d ${GOFLY_VERSION}
chmod 0777 -R ${GOFLY_VERSION}
cd "${SETUP_PATH}/${GOFLY_VERSION}"

echo '---------------------------------------------';
echo "开始配置MySQL数据库，请手动创建好[数据库名称，用户名、密码]...";
echo "请在[宝塔面板=>数据库]创建MySQL数据库名称，用户名、密码...";
#检测数据库连接信息
read -p "MySQL服务地址(默认127.0.0.1):" MYSQL_SERVER
read -p "MySQL端口号(默认3306):" MYSQL_PORT
read -p "MySQL数据库名(自定义):" MYSQL_DB
read -p "MySQL用户名(自定义):" MYSQL_USER
read -p "MySQL密码(自定义):" MYSQL_PASS

SQL_RESULT=`mysql -h ${MYSQL_SERVER} -P ${MYSQL_PORT} -u ${MYSQL_USER} -p${MYSQL_PASS} -e "show databases"`
echo "连接数据库测试结果:${SQL_RESULT}"

CONFIG_TEXT='
{\n
	"Server":"=SE=",\n
	"Port":"=PO=",\n
	"Database":"=DA=",\n
	"Username":"=US=",\n
	"Password":"=PA="\n
}
'

echo -e ${CONFIG_TEXT}|sed "s/=SE=/${MYSQL_SERVER}/g"|sed "s/=PO=/${MYSQL_PORT}/g"|sed "s/=DA=/${MYSQL_DB}/g"|sed "s/=US=/${MYSQL_USER}/g"|sed "s/=PA=/${MYSQL_PASS}/g" > config/mysql.json
echo '---------------------------------------------';
echo "${SETUP_PATH}/${GOFLY_VERSION}/config/mysql.json数据库配置文件内容...";
echo -e `cat config/mysql.json`

getIpAddress=$(curl -sS --connect-timeout 10 -m 60 https://www.bt.cn/Api/getIpAddress)
LOCAL_IP=$(ip addr | grep -E -o '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}' | grep -E -v "^127\.|^255\.|^0\." | head -n 1)


echo "=================================================================="
Yellow_echo "恭喜您，下载成功！"
Yellow_echo "下载保存地址：执行 cd ${SETUP_PATH}/${GOFLY_VERSION}/ ,切换进入该目录"
echo "=================================================================="
Yellow_echo "1. 创建数据库 ${MYSQL_DB}，如果创建数据库失败，请进[宝塔面板]手动创建"
Yellow_echo "2. 编辑 config/mysql.json 可配置MySQL数据库连接信息"
Yellow_echo "3. 执行 ./go-fly-pro install ，确保上面配置信息正确，自动创建导入数据库相关表"
Yellow_echo "4. 执行 ./go-fly-pro server [-d]，监听端口,开启服务，-d为守护进程模式"
Yellow_echo "5. 外网地址: http://${getIpAddress}:8081/login"
Yellow_echo "6. 内网地址: http://${LOCAL_IP}:8081/login"
Yellow_echo "7. 后台主账号用户名:kefu2 密码:123  坐席账号:kefu3 密码:123"
Yellow_echo "8. 若无法访问，请检查[宝塔防火墙/服务器安全组]是否有放行[8081]端口"
Yellow_echo "9. 如有端口冲突，执行 ./go-fly-pro server -p 新端口 -d"
Yellow_echo "10. 关闭服务可以使用此命令 kill -9 \$(pidof 'go-fly-pro')"
