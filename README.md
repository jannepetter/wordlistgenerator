
# WordListGenerator - WLG

Creates wordlists

## usage/commands

See help about the commands
```javascript
go run wlg -h
```

For adding words to your file
```javascript
go run wlg base -h
```

For building your words in your file
```javascript
go run wlg add -h
```

Premade strategies for creating wordlist
```javascript
go run wlg quick -h
```

Mangle your wordlist file
```javascript
go run wlg mangle -h
```

Build executable for faster runtime. Navigate to folder wordlistgenerator and run
```javascript
go build
```
It will create wlg.exe which will be executable with windows
```javascript
wlg.exe -h
```

With linux
```javascript
./wlg.exe -h
```
