# JWPF

## Introduction
Jack Web Path Finder  is a tool that allows you to brute force directories and files of a WebServer

## Why
I like to make my own script and program, since it encrease my ability to solve daily problems and make all my stuff automated without depending from others code.

## Configure
### Install dependecies
To run the program, you first need golang. Take a look at the [Installation Guide](https://golang.org/doc/install)
### Compile
On Linux-based system, you can execute the `config.sh`  script.
- `sh config.sh`

After compiled the program, it is placed in the  `bin/` folder named `jwpf.o`.
Run the `install.sh` script  to place it under the `/usr/local/bin/` path. 

## Syntax
### Explaination
**jwpf** needs 3 positional arguments.

- `<url>` is the target we want to run the attack against on. 
- `<dictionary>` is the path to the file to use as a dictionary.
- `<threads>` number of threads to run at the same time.

### Example
`jwpf http://jackrendor.cf path/to/dict.txt 10`


## TODO
- Add a log system, since it outputs everything on stdin without saving it on file.
- Better handler of status code.

## Thanks to
- Me [telegram](https://t.me/jackrendor)
- And me   [linkedin](https://it.linkedin.com/in/jackrendor)
- And Golang Italia [Telegram Group](https://t.me/golangit)
