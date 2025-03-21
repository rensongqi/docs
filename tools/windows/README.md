- [1 批处理脚本](#1-批处理脚本)
  - [1.1 if 判断](#11-if-判断)
  - [1.2 for 循环](#12-for-循环)
  - [1.3 函数调用](#13-函数调用)
  - [1.4 常用语句](#14-常用语句)

# 1 批处理脚本

| 常用命令 | 参数                                            | 作用                                                         |
| -------- | ----------------------------------------------- | ------------------------------------------------------------ |
| xcopy    | `XCOPY source [destination]`                    |                                                              |
|          | /S（大小写均可）                                | 复制目录和子目录，除了空的。                                 |
|          | /E                                              | 复制目录和子目录，包括空的。 与 /S /E 相同。可以用来修改 /T。 |
|          | /Q                                              | 复制时不显示文件名。                                         |
|          | /T                                              | 创建目录结构，但不复制文件。不包括空目录或子目录。/T /E 包括空目录和子目录。 |
|          | /Y                                              | /Y 禁止提示以确认改写一个现存目标文件。                      |
| mklink   | `MKLINK [[/D] | [/H] | [/J]] LINK Target`       |                                                              |
|          | /D                                              | 创建目录符号链接。默认为文件符号链接。                       |
|          | /H                                              | 创建硬链接而非符号链接。                                     |
|          | /J                                              | 创建目录联接。                                               |
| set      | `/A expression` || `/P variable=[promptString]` |                                                              |
|          | /A                                              | 命令行开关指定等号右边的字符串为被评估的数字表达式           |
|          | /P                                              | 命令行开关允许将变量数值设成用户输入的一行输入               |
|          |                                                 |                                                              |

**xcopy**

```bat
@REM windows打开xcopy需要在环境变量中添加如下条目
@REM C:\WINDOWS\SYSTEM32\XCOPY.EXE

@REM 无论是拷贝文件或是拷贝目录，都需要指定拷贝目标文件名或目录名
@REM 拷贝文件
echo f | xcopy /y E:\mussy\test1\ceshi.txt E:\mussy\test2\ceshi222.txt

@REM 拷贝目录
echo d | xcopy /y E:\mussy\test1\ E:\mussy\test2\test1
```

**mklink**

```bat
@REM 为目录创建软链接
mklink /j link E:\mussy\test1\link

@REM 为文件创建链接
mklink ceshi.txt E:\mussy\test1\ceshi.txt

@REM 为文件创建硬链接
mklink /H ceshi.txt E:\mussy\test1\ceshi.txt
```

**set**

```bat
@REM 设置环境变量
set path=dst=E:\\test
set folders=test1 test2 test3 test4

@REM 把用户的输入作为值传给变量
@echo off
set /p input=请输入字母A或B:
if "%input%"=="A" goto A
if "%input%"=="B" goto B
pause
exit

:A
echo 您输入的字母是A
pause
exit

:B
echo 您输入的字母是B
pause
exit
```



## 1.1 if 判断

```bat
@REM 一行实现
if exist release-dev\\Unreal\\SimOne\\SimOne.uproject del /F  /Q release-dev\\Unreal\\SimOne\\SimOne.uproject

@REM 多行实现
for %%a in (%folders%) do (
	if not exist %dst%\%%a (
		md %dst%\%%a
	)
)
```



## 1.2 for 循环

参考文章：[bat命令之for命令详解 ](https://www.cnblogs.com/kevin-yuan/p/3641847.html)

```bat
set folders=test1 test2 test3 test4
set dst=E:\\test
@REM 循环创建目录 一行实现
for %%a in (%folders%) do (md %dst%\\%%a)

@REM 循环创建目录 多行实现
for %%a in (%folders%) do (
	md %dst%\\%%a
)

@REM 搜索当前目录下有哪些文件
for %%i in (*.*) do echo "%%i"

@REM 打印文件中的内容 默认以" "为分隔符，且只会打印每行分隔符之前的第一列
set variable=E:\\test\test.txt
@REM 文件内容如下
::This is a new world
::may the coda
::be with you
::So do i
for /f %%a in (%variable%) do (
    echo %%a
)

@REM 指定分隔符
for /f "delims=\n" %%a in (%variable%) do (
    echo %%a
)
```



## 1.3 函数调用

批处理脚本中函数是以:实现的，简单函数如下

```bat
@echo off
:funcA
    echo This is A

:funcB
    echo This is B

@REM 由于bat脚本是行运行，所以俩函数都会被输出，输出结果如下
This is A
This is B
```

如果想控制只输出某一个函数，则应该如下

```bat
@echo off
call :funcB
goto:eof

:funcA
    echo This is A
goto:eof

:funcB
    echo This is B
goto:eof

@REM 这样在整个bat脚本之初就先call调用B函数，且执行完call之后直接goto:eof, 那么就实现了只调用一个函数的功能，输出如下
This is B
@REM 如果call调用完B函数，后边没有加goto:eof，那么则A函数也会被调用一次，因为bat是按行处理的，这种情况下输出如下
This is B
This is A
```



## 1.4 常用语句

```bat
@echo off ---> 关闭回显，默认是on

%~dp0 ---> 更改当前目录为批处理本身的目录 (https://blog.csdn.net/qq_22642239/article/details/88549969)

goto:eof ---> CMD返回并将等待下一条命令
```



