# reloader

A simple program that restarts an executable with the given arguments whenever
the binary changes.

## Background

While starting web development with go, I discovered fairly quickly that it was
annoying to interrupt and restart my web executable over and over again, so
instead I wrote this quick helper program to manage that process.

## Usage

```sh
reloader [executable path] [args]
```

So for example, if you had an executable called `webserver` that typically takes
a port and host argument:

```sh
reloader webserver -port 8080 -host 0.0.0.0
```

This will launch the `webserver` executable with the given args, showing any
output from the command. If you do a new build and the webserver executable
is overwritten (or if it's filesystem modified team changes for any other
reason) the reloader program will send an interrupt signal to the webserver.
Once the webserver process exits the reloader program will then restart it.

To exit simply hit Ctrl-c. Reloader will send the interrupt signal to your
program and then exit once it has exited.

## Contributing

Contributions are greatly appreciated! While not always necessary, it can
certainly save you some effort by contacting me first, through GitHub Issues,
to ensure that your improvement isn't already being worked on by someone else.
That'll also let me make sure I'm ready to review and include your contribution
once it's ready. It's certainly ok to make a contribution without discussing it
first as well, but it can be beneficial to you to do so.

When you're ready to write some code, please fork my repository, make changes
on your fork, and then submit a pull request.

We ask that all contributors to this project adhere to the code of conduct
contained in the root of this project's source code.

## License

The MIT License (MIT)

Copyright (c) 2015 Alan Johnson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
