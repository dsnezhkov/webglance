## Webglance 

### Release

### Source

`git clone https://github.com/dsnezhkov/webglance`

`go build`

*Dependencies:*
`go get -u github.com/raff/godet

or 

`go run main.com`

### Oveview

Based on:
[Chrome Dev Tools Protocol Viewer](https://chromedevtools.github.io/devtools-protocol/)

Advantages: 

- Full browser JS support, while significantly increased speed in comparison to Selnium/Webdriver. 
- No reliance on PhantomJS / NodeJS or GTK/QT libraries. Just Chrome browser.

Tons  of useful remote clients for Chrome DevTools [here](https://github.com/ChromeDevTools/awesome-chrome-devtools)

Webglance uses  [godet](https://github.com/raff/godet) client as I found it to be easily understood. 

Webglance takes a list of urls from file (format per line : `http(s)://host[:port]` ) and attempts to screenshot the content on that url.

Users can spacify delay for screenshots (page load times) and output directory for PDFs/PNGs.

### Operaton 

Note: So far tested on MacOS. 

1. Launch headless chrome, with remote debugger. On MacOS you can use aliases to make path easier to call

Example of setup (MacOS):

`alias google-chrome='/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome`

`google-chrome --headless --remote-debugging-port=9222`

```
DevTools listening on ws://127.0.0.1:9222/devtools/browser/6010b06d-a804-4d1e-a749-020c8609f7a9
```

2. Webglance client attachment

```bash
$ ./webglance --help
Usage of ./webglance:
  -hostdbg string
    	Chrome debugger <host:port>.  (default "localhost:9222")
  -outdir string
    	Results output directory. (default "output")
  -pagedelay int
    	Delay Page load before screen.  (default 5)
  -urlfile string
    	Path to Url file.  (default "urls.txt")
```

Contents of `urls.txt`

```
https://google.com:443
#https://ibm.com
https://msn.com
https://amazon.com
```

```bash
$ ./webglance -urlfile=urls.txt  -pagedelay=5 -hostdbg=localhost:9222 -outdir=output
2018/02/15 01:19:22 Taking defaults (urls.txt) on url file
2018/02/15 01:19:22 Taking defaults (localhost:9222) on debugger location
2018/02/15 01:19:22 Taking defaults (output) on output location
2018/02/15 01:19:22 Attempting to connect to Chrome Debugger on localhost:9222
2018/02/15 01:19:22 Connected.
2018/02/15 01:19:22 Directory output exists. Saving PDF/PNG there
2018/02/15 01:19:22 [+] Screening (https://google.com:443)
2018/02/15 01:19:22 read message: read tcp 127.0.0.1:64138->127.0.0.1:9222: use of closed network connection
2018/02/15 01:19:22 permanent network error
2018/02/15 01:19:28 Commented (#https://ibm.com) Skipping.
2018/02/15 01:19:28 [+] Screening (https://msn.com)
2018/02/15 01:19:28 read message: read tcp 127.0.0.1:64140->127.0.0.1:9222: use of closed network connection
2018/02/15 01:19:28 permanent network error
2018/02/15 01:19:35 [+] Screening (https://amazon.com)
2018/02/15 01:19:35 read message: read tcp 127.0.0.1:64142->127.0.0.1:9222: use of closed network connection
2018/02/15 01:19:35 permanent network error
```
Contents of output:

```
-rw-r--r--@ 1 501  20  1360680 Feb 15 01:19 amazon.com.pdf
-rw-r--r--@ 1 501  20   696487 Feb 15 01:19 amazon.com.png
-rw-r--r--  1 501  20    78458 Feb 15 01:19 google.com_443.pdf
-rw-r--r--  1 501  20    48844 Feb 15 01:19 google.com_443.png
-rw-r--r--  1 501  20   281509 Feb 15 01:19 msn.com.pdf
-rw-r--r--  1 501  20   984270 Feb 15 01:19 msn.com.png
```

