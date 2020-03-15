@echo off
rm -rf bot-instaling.exe
go build -ldflags "-H windowsgui" -o build/bot-instaling.exe