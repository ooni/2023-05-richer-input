# 2023-05: Richer Input PoC

We define "richer input" as the possibility of delivering to OONI-Probe
nettests complex input. As of 2023-05-30, OONI nettests either do not
take any input or take a string as input.

This repository contains a proof-of-concept redesign of OONI probe,
where we implement richer input and explore how it can simplify the
implementation.

We include a stripped down version of `ooniprobe` that only supports
the `runx` command to support experimental runs. This command only
allows to run nettests given (a) a known probe location and (b) a script
telling the OONI Probe what to do exactly. This is the most basic
functionality that OONI Probe could provide.

See below for examples of how you can invoke the `runx` command.

See [DESIGN.md](DESIGN.md) and [ARCHITECTURE.md](ARCHITECTURE.md)
for more information.

## Building

Obtain go1.20.4 using these commands:

```console
go get golang.org/dl/go1.20.4@latest
~/go/bin/go1.20.4 download
```

Then build using:

```console
~/go/bin/go1.20.4 build -v ./cmd/ooniprobe
```

## Running

Try:

```console
./ooniprobe --logfile LOG.txt \
	--location testdata/location.json \
	--script testdata/plan.jsonc
```

To run a reasonably complete OONI measurements.

The [testdata/planc.json](testdata/plan.jsonc) contains several
comments explaining its structure. See also the [DESIGN.md](DESIGN.md)
document for more information on the rationale.
