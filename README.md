<!--- PROJECT LOGO --->
![project_logo](https://github.com/amcajal/8_bit_hubble_golang/blob/master/project/doc/media/8_bit_hubble_golang_logo.png)

**8-bit Hubble (Go)** is a complete reimplementation—and ongoing **extension**—of the original [**8-bit Hubble**](https://github.com/amcajal/8_bit_hubble)

For background on the tribute to Hubble, classic games, and the “Bob Ross” layering idea, see the [original project README](https://github.com/amcajal/8_bit_hubble) and [wiki](https://github.com/amcajal/8_bit_hubble/wiki).

---

## Overview

**8-bit Hubble** (Golang version) generates random **pixel-art galaxy** images: star sprites in several sizes, optional nebula and background treatments, occasional large “special” sprites such as planets, etc. Each run is driven by a **pseudo-random number generator**; you can fix a **seed** to reproduce the same image later.

The program is delivered as:

| Entry point | Role |
|-------------|------|
| `cmd/gui` | Desktop UI ([Gio](https://gioui.org/)): output path, filename, seed, optional TOML config, animated WebP, and paletted PNG |
| `cmd/cli` | Same generator from the command line (`-o`, `-n`, `-s`, `-c`, `-a`, `-p`) |

Output is a **500×500** image by default: static **PNG**, optionally **256-colour paletted** PNG, or **animated WebP** when blinking stars are enabled.

Example output:

![example_2](https://github.com/amcajal/8_bit_hubble_golang/blob/master/examples/8bhg_galaxy_example_2.png)![example_3](https://github.com/amcajal/8_bit_hubble_golang/blob/master/examples/8bhg_galaxy_example_3.png)![example_4](https://github.com/amcajal/8_bit_hubble_golang/blob/master/examples/8bhg_galaxy_example_4.png)![example_5](https://github.com/amcajal/8_bit_hubble_golang/blob/master/examples/8bhg_galaxy_example_5.png)![webp_example](https://github.com/amcajal/8_bit_hubble_golang/blob/master/examples/8bhg_animated.webp)

---

## Relationship to the original C project

| | Original C ([8_bit_hubble](https://github.com/amcajal/8_bit_hubble)) | This Go project |
|---|--------|-----------------|
| **Goal** | Safety-critical style C codebase, GTK GUI, libpng | Idiomatic Go, standard `image` encoding, modern GUI stack |
| **Scope** | Static PNG galaxies, three star size classes | Same core loop **plus** planets, nebula, animation, tuning file, and more |
| **Sprites** | Spreadsheet → Python → generated C arrays | PNG assets embedded in the binary (`go:embed`) |

This repository is **not** a line-for-line port. It preserves the **procedural galaxy** concept and many algorithmic ideas, but goes down its own way.

---

## Features kept from the original

- **Procedural galaxies** with random placement, colours, and layering—still inspired by a painter’s-algorithm style build-up (see the original [Bob Ross algorithm notes](https://github.com/amcajal/8_bit_hubble/wiki/Page-6:-Bob-Ross-Algorithm)).
- **500×500** output canvas.
- **Seeded RNG** so runs are reproducible when you pass the same seed.
- **Sprite size tiers** (small / medium / large) painted in layers with random counts within configured limits.
- **CLI and GUI** workflows: choose where to save, optional seed, generate.
- **GPL v3** licensing and authorship of the original **8-bit Hubble** concept.

---

## What is new in the Go version

### User-facing features

- **Animated galaxies** (`-a` / **Animated** in the GUI): multi-frame **WebP** with stars that blink via opacity cycling (planets stay static).
- **Paletted export** (`-p` / **Paletted**): Floyd–Steinberg dithering to the Plan 9 **256-colour** palette.
- **TOML configuration** (`-c` / **Config file**): override layer counts, sprite counts per layer, paint probabilities, rotation/scale chances, and more—see [`config_example.toml`](config_example.toml).
- **Richer GUI**: folder/file pickers, collapsible advanced options, in-app help/about, and a result overlay showing the seed used.

### Generation and visuals

- **Special sprites** (e.g. planets) as a fourth category, with larger assets than the original star set.
- **Procedural nebula** (coloured Perlin-style noise, blended layers) instead of only a flat “cosmic” fill.
- **Cosmic Latte** star-field background (noise-based), optional per run.
- **Per-sprite hue rotation** for colour variation (`internal/palette`).
- **Random rotation and 2× nearest-neighbour scaling** on sprites (configurable chances in TOML).
- **PNG-only pipeline**: no libpng; encoding via Go’s `image/png` and [nativewebp](https://github.com/HugoSmits86/nativewebp) for animation.

### Developer experience

- **Pure Go** build (`go build` / `go run`); no C toolchain, GTK, or Python sprite conversion step.
- **Sprites as ordinary PNG files** under `internal/sprites/assets/`, embedded at compile time—create art in any editor, add the file, register it in `static.go`.
- **Tests** in packages such as `internal/palette`.
- Layout centred on `internal/` (`galaxy`, `sprites`, `param`, `rngm`, `palette`).

Historical note: early experiments in this repo used base64-in-source or remote PNG fetching for sprites; the current design uses **`go:embed`** on PNG assets, which keeps authoring simple and avoids libpng.

---

## Quickstart

### Requirements

- [Go](https://go.dev/) 1.26+ (see [`go.mod`](go.mod))

### GUI

```bash
git clone https://github.com/amcajal/8_bit_hubble_golang.git
cd 8_bit_hubble_golang
go run ./cmd/gui
```

Fill in the output directory and image name, optionally set a seed and advanced flags, then **Generate Galaxy!**

### CLI

```bash
go run ./cmd/cli --help

# Default output: ./8bh_galaxy.png (seed = current time)
go run ./cmd/cli

# Reproducible run with custom path
go run ./cmd/cli -s 42 -o ./output -n my_galaxy

# Animated WebP with config overrides
go run ./cmd/cli -a -s 99 -n blinking_stars -c config_example.toml
```

Build standalone binaries:

```bash
go build -o 8bHubbleGolang ./cmd/cli
go build -o 8bHubbleGolang-gui ./cmd/gui
```

---

## License

**8-bit Hubble** was created by Alberto Martin Cajal. This Go reimplementation and extension is released under **GNU GPL v3.0**—see [`LICENSE`](LICENSE) and the [GPL text](https://www.gnu.org/licenses/gpl-3.0.en.html).

Third-party Go modules (Gio, TOML parser, bild, nativewebp, etc.) are used under their respective licenses; see module sources and `go.mod`.

---

## Contact

Original author contact details: [8_bit_hubble README — Contact](https://github.com/amcajal/8_bit_hubble#contact).
