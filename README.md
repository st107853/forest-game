# Build and install
## Build 
1. Install necessary compilers (MinGW): `sudo apt install mingw-w64`
2. Build with necessary flags: `CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"`
3. Verify that build succeeded `file forest-game.exe`
4. Copy forest-game.exe to your Windows (R) machine (????)
5. Enjoy
