@echo off
rm -rf dist/bot-instaling.exe
go build -ldflags "-H windowsgui" -o dist/bot-instaling.exe