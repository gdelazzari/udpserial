# udpserial

An UDP to serial port bridge.

Allows to tunnel serial ports through UDP packets, i.e. the packets received by UDP will be output as serial data and (viceversa) incoming serial data will be sent as UDP packets.

[![Animated demo](https://raw.githubusercontent.com/gdelazzari/udpserial/master/assets/demo.mp4)](https://raw.githubusercontent.com/gdelazzari/udpserial/master/assets/demo.mp4)

This project was developed and tested on Linux only.

It has not been used in production, although it proved to be quite stable during extensive testing in a prototyping environment.

One known issue is the absence of a limit for the incoming buffer, which can lead to out of memory conditions in case of DoS attacks (or exaggerated incoming traffic) that can make the daemon crash. Due to the usage of the [Go](https://go.dev/) programming language, the software should be memory safe, albeit I haven't checked the code in a long while (it has been 5 years since I've worked on this).

The [Web UI](https://github.com/gdelazzari/udpserial/tree/master/panel), which was written with the [Vue.js](https://vuejs.org/) framework, is using old dependencies and could take some upgrades of the packages.

Finally, the Web UI is not password protected, so it is not meant to be publicly exposed.

This software is no longer in use, but improvements are welcome.

## Features

- Lightweight (~6MB RAM and low CPU usage)
- Web UI for configuration and monitoring
- Multiple options for serial to UDP packetization:
  - automatic (with timeout from last character)
  - manually specified string (e.g. `\n\r`) for known protocols
- Can handle an unlimited number of serial ports in parallel
- Port configuration includes:
  - baudrate
  - data bits
  - stop bits
- UDP stream configuration includes:
  - output address and port for incoming serial data
  - listening address and port for outgoing serial data
- Includes a real-time plot of each port activity (in bytes/s)
- Logging to file, console and Web UI

## Building and running

Build the web panel interface, which requires a working [NPM](https://www.npmjs.com/) installation:
```console
$ cd panel
$ npm install
$ npm run build
$ cd ..
```

Build the main daemon, which requires a working [Go](https://go.dev/) installation:
```console
$ go build
```

Copy `definitions.json.example` to `definitions.json` and adjust the file by listing the serial ports you want to expose and the baudrates you want to support.

Then the daemon can be launched by running the statically compiled executable `udpserial`. The web interface will be listening on `0.0.0.0:8080`, and can be used to inspect the current traffic in real-time, configure UDP tunnels and check the daemon logging output.

The daemon can be configured as a systemd service to ensure it is always running in the background.

## Usage

Configure serial tunnels under the *CONFIGURATION* tab, the various options should be self-explanatory. After you added, edited or unlinked ports, clicking on *APPLY CONFIGURATION* will make the changes effective and persistent, by saving them in a `config.json` file under the current directory.

The *STATISTICS* tab allows to inspect in real-time the traffic on the various channels, measured in bytes/s.

The *DIAGNOSTIC* tab allows to inspect the system log.

See the animated demo above for the software in action.

## License

Copyright (C) 2022  Giacomo De Lazzari

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

