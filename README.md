# randrctl2

[![Build Status](https://travis-ci.org/edio/randrctl2.svg?branch=master)](https://travis-ci.org/edio/randrctl2)

An attempt to rewrite [randrctl](https://github.com/edio/randrctl) in more suitable language using low(er)-level X api

## Goals

- better distribution<br/>
  statically compiled binary
- no dependency to _xrandr_<br/>
  new version uses X bindings for Go and talks to X RandR extension directly
- improved support for RandR protocol<br/>
  no need to reverse xrandr behavior and try parsing its output
- performance<br/>
  Single run should take way less than 100ms (with _randrctl_ it is almost a second)
- better cli<br/>
  _spf13/cobra_ offers bash and zsh completion generation
