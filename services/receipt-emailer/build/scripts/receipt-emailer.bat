@rem
@rem Copyright 2015 the original author or authors.
@rem
@rem Licensed under the Apache License, Version 2.0 (the "License");
@rem you may not use this file except in compliance with the License.
@rem You may obtain a copy of the License at
@rem
@rem      https://www.apache.org/licenses/LICENSE-2.0
@rem
@rem Unless required by applicable law or agreed to in writing, software
@rem distributed under the License is distributed on an "AS IS" BASIS,
@rem WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@rem See the License for the specific language governing permissions and
@rem limitations under the License.
@rem

@if "%DEBUG%" == "" @echo off
@rem ##########################################################################
@rem
@rem  receipt-emailer startup script for Windows
@rem
@rem ##########################################################################

@rem Set local scope for the variables with windows NT shell
if "%OS%"=="Windows_NT" setlocal

set DIRNAME=%~dp0
if "%DIRNAME%" == "" set DIRNAME=.
set APP_BASE_NAME=%~n0
set APP_HOME=%DIRNAME%..

@rem Resolve any "." and ".." in APP_HOME to make it shorter.
for %%i in ("%APP_HOME%") do set APP_HOME=%%~fi

@rem Add default JVM options here. You can also use JAVA_OPTS and RECEIPT_EMAILER_OPTS to pass JVM options to this script.
set DEFAULT_JVM_OPTS=

@rem Find java.exe
if defined JAVA_HOME goto findJavaFromJavaHome

set JAVA_EXE=java.exe
%JAVA_EXE% -version >NUL 2>&1
if "%ERRORLEVEL%" == "0" goto init

echo.
echo ERROR: JAVA_HOME is not set and no 'java' command could be found in your PATH.
echo.
echo Please set the JAVA_HOME variable in your environment to match the
echo location of your Java installation.

goto fail

:findJavaFromJavaHome
set JAVA_HOME=%JAVA_HOME:"=%
set JAVA_EXE=%JAVA_HOME%/bin/java.exe

if exist "%JAVA_EXE%" goto init

echo.
echo ERROR: JAVA_HOME is set to an invalid directory: %JAVA_HOME%
echo.
echo Please set the JAVA_HOME variable in your environment to match the
echo location of your Java installation.

goto fail

:init
@rem Get command-line arguments, handling Windows variants

if not "%OS%" == "Windows_NT" goto win9xME_args

:win9xME_args
@rem Slurp the command line arguments.
set CMD_LINE_ARGS=
set _SKIP=2

:win9xME_args_slurp
if "x%~1" == "x" goto execute

set CMD_LINE_ARGS=%*

:execute
@rem Setup the command line

set CLASSPATH=%APP_HOME%\lib\receipt-emailer-1.0.0.jar;%APP_HOME%\lib\idempotency.jar;%APP_HOME%\lib\kmongo-4.0.1.jar;%APP_HOME%\lib\kmongo-core-4.0.1.jar;%APP_HOME%\lib\kmongo-jackson-mapping-4.0.1.jar;%APP_HOME%\lib\kmongo-property-4.0.1.jar;%APP_HOME%\lib\kmongo-shared-4.0.1.jar;%APP_HOME%\lib\kmongo-id-jackson-4.0.1.jar;%APP_HOME%\lib\jackson-module-loader-0.1.0.jar;%APP_HOME%\lib\kmongo-data-4.0.1.jar;%APP_HOME%\lib\kreflect-1.0.0.jar;%APP_HOME%\lib\kmongo-id-4.0.1.jar;%APP_HOME%\lib\jackson-module-kotlin-2.11.0.jar;%APP_HOME%\lib\kotlin-reflect-1.3.72.jar;%APP_HOME%\lib\kotlin-stdlib-1.3.71.jar;%APP_HOME%\lib\aws-java-sdk-lambda-1.11.63.jar;%APP_HOME%\lib\aws-lambda-java-events-3.1.0.jar;%APP_HOME%\lib\aws-lambda-java-core-1.2.1.jar;%APP_HOME%\lib\sendgrid-java-4.5.0.jar;%APP_HOME%\lib\gson-2.8.6.jar;%APP_HOME%\lib\kotlin-stdlib-common-1.3.71.jar;%APP_HOME%\lib\annotations-13.0.jar;%APP_HOME%\lib\aws-java-sdk-core-1.11.63.jar;%APP_HOME%\lib\jmespath-java-1.11.63.jar;%APP_HOME%\lib\joda-time-2.8.1.jar;%APP_HOME%\lib\java-http-client-4.3.7.jar;%APP_HOME%\lib\jackson-dataformat-cbor-2.6.6.jar;%APP_HOME%\lib\jackson-databind-2.11.0.jar;%APP_HOME%\lib\jackson-core-2.11.0.jar;%APP_HOME%\lib\jackson-annotations-2.11.0.jar;%APP_HOME%\lib\bcprov-jdk15on-1.65.jar;%APP_HOME%\lib\httpclient-4.5.13.jar;%APP_HOME%\lib\commons-logging-1.2.jar;%APP_HOME%\lib\ion-java-1.0.1.jar;%APP_HOME%\lib\httpcore-4.4.14.jar;%APP_HOME%\lib\mongodb-driver-sync-4.0.3.jar;%APP_HOME%\lib\bson4jackson-2.9.2.jar;%APP_HOME%\lib\commons-codec-1.11.jar;%APP_HOME%\lib\mongodb-driver-core-4.0.3.jar;%APP_HOME%\lib\bson-4.0.3.jar

@rem Execute receipt-emailer
"%JAVA_EXE%" %DEFAULT_JVM_OPTS% %JAVA_OPTS% %RECEIPT_EMAILER_OPTS%  -classpath "%CLASSPATH%" com.cshep4.kripto.receiptemailer.Handler %CMD_LINE_ARGS%

:end
@rem End local scope for the variables with windows NT shell
if "%ERRORLEVEL%"=="0" goto mainEnd

:fail
rem Set variable RECEIPT_EMAILER_EXIT_CONSOLE if you need the _script_ return code instead of
rem the _cmd.exe /c_ return code!
if  not "" == "%RECEIPT_EMAILER_EXIT_CONSOLE%" exit 1
exit /b 1

:mainEnd
if "%OS%"=="Windows_NT" endlocal

:omega
