# Easydoc

Easydoc is a very simple Markdown server. Configure Easydoc to serve a directory
and it will find all `*.md` files under that directory structure. These files are
shown in a left hand Table of Contents view.
Markdown is rendered in the right hand panel when files are clicked from the TOC.
The corpus of files can be searched with a simple regex.

## Purpose

The main purpose of Easydoc is to support high level documentation in a source code
repository. Simply organise your code as you normally would, and then place `*.md`
files near to the code that they provide documentation for. Easydoc's own source follows
this pattern.

When documentation is placed near to code, it is easy to find when reading source code -
just look in the current directory for MD files. Look upward to find more documentation.
MD files near the leaf of the directory tree should be the most specific and detailed
level for documentation. As you move up toward the root of the directory tree, documentation
must become higher and higher level.

One of the biggest criticisms of documentation is that it doesn't reflect what the code
is meaning. This is a real problem. The best way to solve this problem is to write
documentation from a very high level and describe **what** the code is doing. Avoid
describing **how** it is done - this will allow the implementation to change without
destroying all the documentation. Documenting the **how** is generally best done in
the code itself with comments. In these cases, it's probably best to have a short
description in the MD file referring to the code.
This is of course (like all things) a balance. Some code will be best described in an MD
file because it is easier to write as proper documentation instead of code comments.
For example, protocols, or a complex allocator might fit better as real documentation.

## URL structure

Easydoc has a simple URL scheme

- **root**/path/to/file.md will render `file.md` alone
- **root**/#/path/to/file.md will render a Table Of Contents to the left, with `file.md` in the right
  pane. Note the `#` at the start of the URL.
- **root**/#/?search=<tosearch> will show a TOC on the left and search results on the right.

The general principle of Easydoc's URLs is that they can be bookmarked and shared with others.

This structure lets Markdown files link with other files in the tree using absolute paths -
[example](/internal/InternalReadme.md). Take note that there is no # in this URL.

## Flags

Easydoc accepts some commandline flags. (cmd/easydoc/easydoc.go)

| flag    | usage                                                                                                                                                                     |
| ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| root    | The root path to begin serving from. `easydoc.json` will be used as configuration in this path                                                                            |
| rootUrl | The root of the full path that Easydoc is served from. Default is `http://localhost:8080`. If you are behind a proxy, you might need something like `https://mysite/docs` |
| port    | The interface and port to bind on. Defaults to `localhost:8080`                                                                                                           |

## easydoc.json

Further configuration is stored in `easydoc.json`, which has the following known keys.

| key             | usage                                                                                                                                                           |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ignore          | An array of strings. Each string follows the same format as .gitignore [rules](https://git-scm.com/docs/gitignore)                                              |
| externalUrlBase | String. When present, each page will show a header with a link formed by prepending this value and the page location. This can be used to redirect to eg github |

## .gitignore

`.gitignore` will be loaded from **root** to restrict where files are indexed from

## Help

- [Markdown cheatsheet](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)
