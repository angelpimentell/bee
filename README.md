# Bee - Port scan tool 🐝
Bee is a tool to scan ports by TCP and UDP protocols
## Building from Source Code
First, we compile the source code with the lightest size
```bash
go build -ldflags "-s -w" bee.go
```
Then, we compress the binary, this step is optional
```bash
upx bee
```
## Using Bee
```bash
./bee -i 10.0.0.1 -p 1-65535
```
### More examples
```bash
./bee -i 10.0.0.1 -p 80,443
./bee -i 10.0.0.1 -p 21,22,80-100
./bee -i 10.0.0.1 -p 80-100 -pU -pT
```
