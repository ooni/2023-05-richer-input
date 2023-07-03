# Architecture

We depend upon the semi-official [probe-engine](https://github.com/ooni/probe-engine)
that exports [probe-cli](https://github.com/ooni/probe-cli) internals on a best
effort basis to the OONI community.

## The cmd directory

The [cmd/ooniprobe](cmd/ooniprobe/) directory contains a minimal "ooniprobe"
client implementing the "runx" command defined in [DESIGN.md](DESIGN.md).

## The pkg directory

The [analysis](pkg/analysis/) package contains common code for data analysis.

The [dsl](pkg/dsl/) package contains an internal and external DSL used
to implement richer input for some nettests. This package is alternative to
(and likely better than) [mininettest](pkg/mininettest/).

The [experiment](pkg/experiment/) package reimplements the IM nettests
to use the "mini nettests" defined in [DESIGN.md](DESIGN.md).

The [mininettest](pkg/mininettest/) package implements the "mini nettests".

The [modelx](pkg/modelx/) package contains data structures and interfaces
specific of this PoC (and called modelx because they conceptually extend the
model package of the probe-engine).

The [ooniprobe/interpreter](pkg/ooniprobe/interpreter/) package implements the
"interpreter" defined in [DESIGN.md](DESIGN.md).

The [x](pkg/x/) package contains experimental packages.

## The testdata directory

This directory contains the location and the scripts on which the
"interpreter" depends. You can use these files to test the PoC.
