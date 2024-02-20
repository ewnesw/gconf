# GCONF

WIP simple conf manager which track, copy and push given directories on changes

## Description :
Gconf basically takes input paths and track them, on creation it adds new folder to the tracking list and add each files to a cloning folder and push it to you repo
This is basically just a project to test and train on Golang 

Next steps to improve this project are (not in priority order don't worry):
- Daemonize the project
- clone files
- improve code
- test it ( probably the most urgent one tbh )

## Usage :
### Build project 
Build the project using : ```go build -o <output>``` in the project directory

### Run project:
gconf add < paths >

Example : 
```
gconf add /tmp
gconf add /tmp /var/save
```

WARNING: for the initialisation you need to specify each subfolders you want to track, for example if you want to track /tmp AND /tmp/example you need to specify both like 
```
gconf add /tmp /tmp/example
```

## Credits :
Thanks to : 
- github.com/go-git/go-git/v5
- github.com/fsnotify/fsnotify

## License :
Do anything you want
