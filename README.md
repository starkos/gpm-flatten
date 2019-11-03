# gpm-flatten

Like most music services, [Google Play Music][gm] will occasionally lose the streaming rights to certain songs. When that happens, those songs simply disappear from your library as if they never existed. Yikes. Also like most music services, GPM will on occasion swap out songs in your collection for live or censored versions, i.e. ones with lower licensing fees.

This project is part of my current gaffer tape solution to detect these unwanted changes to my music collection. Use it like this:

1. Export your Google Play Music collection using [Google Takeout][gt]. Download and unzip the archive somewhere local.

2. Flatten it using **gpm-flatten** (this project):

```
$ gpm-flatten -source <path to Takeout export> -dest <output path>
```

3. Compare the results to your previous export. I check mine into a Git repository and then look at the deltas there, but you can use [any file comparison tool][fc].

4. Repeat as desired to detect changes.

## Getting It

Binaries for Windows and macOS are available for [the latest release](https://github.com/starkos/gpm-flatten/releases).

If you would rather build from source, **gpm-flatten** is written in [Go](https://golang.org/doc/install):

```
$ go get github.com/starkos/gpm-flatten
```

## Stay in touch

* Twitter - [@starkos](https://twitter.com/starkos)

## License

[MIT](https://opensource.org/licenses/MIT)

[fc]: https://en.wikipedia.org/wiki/Comparison_of_file_comparison_tools
[gm]: https://play.google.com/music/listen
[gt]: https://takeout.google.com/settings/takeout
