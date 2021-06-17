@echo off&setlocal enabledelayedexpansion
title build coremail_address_list_export
chcp 65001 >nul
@set version=v1.1
@set output=coremail_address_list_export
@set  build_dir=build
@if not exist %build_dir% (
        @mkdir %build_dir%
        @echo create folder %build_dir%
    ) 
@echo output:%output%
@echo build version:%version%
@set outfilename=%output%_windows_386_%version%.exe
@echo build windows/386 …… %build_dir%/%outfilename%
@set GOOS=windows&&set GOARCH=386&& go build  -trimpath -ldflags "-w -s" -buildmode=pie -o %build_dir%/%outfilename%  main.go

@set outfilename=%output%_windows_amd64_%version%.exe
@echo build windows/amd64 …… %build_dir%/%outfilename%
@set GOOS=windows&&set GOARCH=amd64&& go build -trimpath  -ldflags "-w -s" -buildmode=pie -o %build_dir%/%outfilename%  main.go

@set outfilename=%output%_linux_386_%version%
@echo build linux/386 …… %build_dir%/%outfilename%
@set GOOS=linux&&set GOARCH=386&& go build -trimpath  -ldflags "-w -s"  -o %build_dir%/%outfilename%  main.go

@set outfilename=%output%_linux_amd64_%version%
@echo build linux/amd64 …… %build_dir%/%outfilename%
@set GOOS=linux&&set GOARCH=amd64&& go build -trimpath  -ldflags "-w -s" -buildmode=pie -o %build_dir%/%outfilename%  main.go

@set outfilename=%output%_darwin_amd64_%version%
@echo build darwin/amd64 …… %build_dir%/%outfilename%
@set GOOS=darwin&&set GOARCH=amd64&& go build  -trimpath -ldflags "-w -s" -buildmode=pie -o %build_dir%/%outfilename%   main.go
@echo build finished!
@pause