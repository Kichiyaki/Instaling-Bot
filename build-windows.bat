@echo off
rm -rf bot-instaling.exe
go build -ldflags "-H windowsgui" -o bot-instaling.exe