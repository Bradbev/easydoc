# /internal

[Up](../readme.md)


## /internal/markdown

Helper functions that wrap around the `blackfriday` package.

## /internal/search

Simple text indexer. The current implementation is very basic

- Files are cached in memory on first search
- Simple line-by-line regex matching

This is inefficient in both memory and processing space, but it is
unlikely Easydoc will be given enough data that this matters :)

## /internal/walker

Helper functions to walk directories and find `*.md` files.
