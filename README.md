# wat

> Fast File Watcher

Run specific tasks when specific files are added, changed or deleted

## Usage

Install by doing:

```go
go get -u github.com/raphamorim/wat
```

Then simply run:

```bash
wat <path> <command>
```

Example using `date` command:

```bash
wat src/ date
```

Returns for every file update inside `src` folder:

```
Waiting changes...
Sat Oct 29 15:34:23 BRST 2016
```

## FAQs

#### EMFILE: Too many opened files.?

This is because of your system's max opened file limit. For OSX the default is very low (256). Temporarily increase your limit with `ulimit -n 10480`, the number being the new max limit.

In some versions of OSX the above solution doesn't work. In that case try `launchctl limit maxfiles 10480 10480` and restart your terminal. See here.
