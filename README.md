# Dolphin's codesearch UI

A server and associated JS app running (soon) on https://cs.dolphin-emu.org/

The application allows browsing through Dolphin's source code and browse
through xrefs: what uses this function, where is this function defined, etc. A
quick search feature allows regexp-matches over the whole codebase in a few
milliseconds.

The xrefs indexing is powered by [Google's Kythe project](http://kythe.io/).
The configuration for Dolphin's indexing pipeline is not part of this
repository — see [SADM](https://github.com/dolphin-emu/sadm) instead for
Buildbot and indexing scripts.

The regexp-search is powered by [Google's codesearch Go library](https://github.com/google/codesearch).
The library is extended to index from Kythe LevelDB tables instead of using
filesystem data.
