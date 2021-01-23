@echo off
pushd .

if "%1"=="" (
	set ZH_HOME=%~dp0
) else (
	set ZH_HOME=%1
)
set ZH_SRC=%ZH_HOME%src

echo "ZH_HOME: %ZH_HOME%"
echo "ZH_SRC: %ZH_SRC%"

if not exist %ZH_HOME% mkdir %ZH_HOME%
if not exist %ZH_SRC% mkdir %ZH_SRC%

cd %ZH_SRC%

call :SETUP_PROJECT ziphttpd.windows
call :SETUP_PROJECT ziphttpd
call :SETUP_PROJECT zhsig
call :SETUP_PROJECT selector
call :SETUP_PROJECT updater
call :SETUP_PROJECT Launcher

popd

exit /B 0

:SETUP_PROJECT
	set PRJ=%1
	echo "PROJECT BUILD START: %PRJ%"
	set PRJ_DIR=%ZH_SRC%\%PRJ%

	cd %ZH_SRC%
	if not exist %PRJ_DIR% git clone https://github.com/ziphttpd/%PRJ%

	cd %PRJ_DIR%
	git pull
	if exist %PRJ_DIR%\setup.cmd %PRJ_DIR%\setup.cmd %ZH_HOME%
	echo "PROJECT BUILD   END: %PRJ%"

exit /b
