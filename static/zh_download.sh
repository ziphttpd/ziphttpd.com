#!/bin/bash
set -e

if [ "$1" ]; then ZH_HOME=$1; else ZH_HOME=$(cd $(dirname $0); pwd); fi
ZH_SRC=${ZH_HOME}/src

echo "ZH_HOME: ${ZH_HOME}"
echo "ZH_SRC: ${ZH_SRC}"

if [ ! -d "${ZH_HOME}" ]; then mkdir ${ZH_HOME}; fi
if [ ! -d "${ZH_SRC}" ]; then mkdir ${ZH_SRC}; fi

function SETUP_PROJECT() {
	PRJ=$1
	echo "PROJECT BUILD START: ${PRJ}"
	PRJ_DIR=${ZH_SRC}/${PRJ}
	if [ ! -d "${PRJ_DIR}" ]; then git clone https://github.com/ziphttpd/${PRJ}; fi
	cd ${PRJ_DIR}
	git pull
	if [ -f "${PRJ_DIR}/setup.sh" ]; then bash ${PRJ_DIR}/setup.sh "${ZH_HOME}"; fi
	echo "PROJECT BUILD   END: ${PRJ}"
}

cd ${ZH_SRC}

SETUP_PROJECT "ziphttpd.unix"
SETUP_PROJECT "ziphttpd"
SETUP_PROJECT "zhsig"
SETUP_PROJECT "selector"

#cd ${ZH_HOME}
#./zhget -host ziphttpd.com -group ziphttpd-V1
#./ziphttpd.sh start
#./selector.sh start

exit 0
