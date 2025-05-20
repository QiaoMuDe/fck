@echo off
chcp 65001 > nul
setlocal

:: 设置输出文件名
set OUTPUT_FILE=输出文件名.exe
:: 设置项目名
set PROJECT_NAME=项目名
:: 入口文件
set ENTRY_FILE=main.go
:: 是否使用vendor模式构建
set USE_VENDOR_IN_BUILD=false


echo 检测 Go 编译器是否安装...
go version > nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未检测到 Go 编译器，请确保已安装并配置好环境变量。
    exit /b 1
)

:: 检查当前目录是否存在 go.mod 文件
if not exist "go.mod" (
    echo 错误: 当前目录下未找到 go.mod 文件，请确保在项目根目录下运行脚本。
    exit /b 1
)

:: 检查入口文件是否存在
if not exist "%ENTRY_FILE%" (
    echo 错误: 未找到入口文件 %ENTRY_FILE%, 或修改脚本中的入口文件名。
    exit /b 1
)

REM :: 通过 go mod tidy 确保依赖项的一致性
go mod tidy > mod_tidy.log 2>&1
if %errorlevel% neq 0 (
    echo 错误: go mod tidy 执行失败，请检查错误信息。
    type mod_tidy.log
    del mod_tidy.log 
)
del mod_tidy.log

REM :: 通过 go vet 检查代码
go vet %ENTRY_FILE% > vet.log 2>&1
if %errorlevel% neq 0 (
    echo 错误: go vet 检查失败，请查看以下问题：
    type vet.log
    del vet.log
    exit /b 1
)
del vet.log

REM :: 通过 gofmt 格式化代码
gofmt -w . > fmt.log 2>&1
if %errorlevel% neq 0 (
    echo 错误: gofmt 格式化失败，请查看以下问题：
    type fmt.log
    del fmt.log 
)
del fmt.log

:: 获取 Git 版本信息
echo 正在获取 Git 版本信息...

:: 使用 PowerShell 获取 Git 版本信息
for /f "tokens=*" %%i in ('PowerShell -Command "git describe --tags --always --dirty 2>$null"') do set GIT_VERSION=%%i
if "%GIT_VERSION%"=="unknown" (
    echo 警告: 无法获取 Git 版本信息，可能是当前目录不是 Git 仓库。
    exit /b 1
)

for /f "tokens=*" %%i in ('PowerShell -Command "git rev-parse --short HEAD 2>$null"') do set GIT_COMMIT=%%i
for /f "tokens=*" %%i in ('PowerShell -Command "git log -1 --format=%%cd --date=iso 2>$null"') do set GIT_COMMIT_TIME=%%i
for /f "tokens=*" %%i in ('PowerShell -Command "[System.TimeZoneInfo]::ConvertTimeBySystemTimeZoneId([DateTime]::Parse('%GIT_COMMIT_TIME%'), 'UTC').ToString('yyyy-MM-ddTHH:mm:ssZ')"') do set FORMAT_TIME=%%i
for /f "tokens=*" %%i in ('PowerShell -Command "if ((git status --porcelain 2>$null).Length -gt 0) { 'dirty' } else { 'clean' }"') do set GIT_STATUS=%%i

:: 获取构建时间
echo 正在获取构建时间...
for /f "tokens=*" %%i in ('PowerShell -Command "[System.TimeZoneInfo]::ConvertTimeBySystemTimeZoneId((Get-Date), 'UTC').ToString('yyyy-MM-ddTHH:mm:ssZ')"') do set BUILD_TIME=%%i

echo 正在构建程序...
:: 编译注入的参数
set LD_FLAGS=-X "gitee.com/MM-Q/verman.appName=%PROJECT_NAME%" -X "gitee.com/MM-Q/verman.gitVersion=%GIT_VERSION%" -X "gitee.com/MM-Q/verman.gitCommit=%GIT_COMMIT%" -X "gitee.com/MM-Q/verman.gitCommitTime=%FORMAT_TIME%" -X "gitee.com/MM-Q/verman.buildTime=%BUILD_TIME%" -X "gitee.com/MM-Q/verman.gitTreeState=%GIT_STATUS%" -s -w

:: 编译程序
set BUILD_CMD=go build -ldflags "%LD_FLAGS%"
if "%USE_VENDOR_IN_BUILD%"=="true" (
    set BUILD_CMD=%BUILD_CMD% -mod=vendor
)
%BUILD_CMD% -o %OUTPUT_FILE% %ENTRY_FILE% > build.log 2>&1
if %errorlevel% neq 0 (
    echo 错误: 程序构建失败，请检查错误信息。
    type build.log
    del build.log
    exit /b 1
)

echo ----------------------------------------------------------
echo %PROJECT_NAME% 构建成功，输出文件为 %OUTPUT_FILE%
echo 版本信息:
echo Git 仓库状态: %GIT_STATUS%
echo Git 版本: %GIT_VERSION%
echo Git 提交: %GIT_COMMIT%
echo Git 提交时间: %FORMAT_TIME%
echo 构建时间: %BUILD_TIME%
echo ----------------------------------------------------------

:: 清理临时文件
del build.log
endlocal