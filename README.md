# 🐧 Penguinguide

Penguinguide is a friendly helper for people who are new to Linux. It installs and removes packages, shows system information, checks network and WiFi, and includes a guided quickstart that explains what is happening as you go.

The main idea is simple. Do useful things on your system and learn the real commands at the same time.

---

## Features

* Detects your Linux distribution and package family
* Install, remove, search, and inspect packages while showing native commands
* System summary with hostname, distribution, kernel, memory, and load
* Network overview including default gateway, DNS servers, and interface addresses
* WiFi details such as signal strength, band, channel hints, and security
* Latency and bandwidth checks with a simple speed test
* WiFi doctor that combines wireless checks and a quick speed test
* Quickstart mode that walks through common tasks interactively
* `quickstart --script` that prints plain shell commands for teaching or notes

---

## Install

Clone and build:

    git clone https://github.com/yourname/penguinguide
    cd penguinguide
    go build
    ./penguinguide

You can also run it directly with Go:

    go run .

---

## Usage examples

Show a system summary:

    penguinguide sys

Install a package:

    penguinguide install htop

Explain and preview before running a change:

    penguinguide install htop --dry-run --explain

WiFi information and guidance:

    penguinguide sys wifi

WiFi doctor with health check and quick speed test:

    penguinguide wifi-doctor

Quickstart guided tour:

    penguinguide quickstart

Generate native commands for documentation or training:

    penguinguide quickstart --script

---

## Project goals

Penguinguide aims to:

* Make Linux less confusing for beginners
* Teach commands instead of hiding them
* Work across common distributions without complex setup
* Stay small, readable, and easy to extend
* Encourage people to explore their system with confidence

If this tool helps someone say “oh, now that makes sense”, then it is doing its job.

---

## Contributing

Ideas, issues, and pull requests are welcome. If you would like to add new commands or improve the explanations, see `CONTRIBUTING.md` for a short guide on how the project is structured and how to extend it.

---

## License

Penguinguide is released under the MIT License.

