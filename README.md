# Tom -- Terminal for Open-Meteo

An interactive terminal user interface (TUI) weather app.
Based on data from [Open-Meteo](https://open-meteo.com/).

Very early work in progress!

![screenshot forecast](https://github.com/mlange-42/tom/assets/44003176/f20d793a-2d68-412a-b3f6-fb6a166f488c)
![screenshot plots](https://github.com/mlange-42/tom/assets/44003176/043acfa8-60c7-4d77-8deb-f78cb8adef18)

## Installation

Pre-compiled binaries for Linux, Windows and MacOS are available in the
[Releases](https://github.com/mlange-42/tom/releases).

> Alternatively, install the latest development version of Tom using [Go](https://go.dev):
> ```shell
> go install github.com/mlange-42/tom@main
> ```

## Usage

Run `tom` with the name of a place:

```
tom Buxtehude
```

Set the default location, so you can later simply run `tom` without arguments:

```
tom Buxtehude --default
```

Get help:

```
tom -h
```
