
# Falco tracing framework

## GSOC 2019 project
This repository is part of a GSOC 2019 [project](https://summerofcode.withgoogle.com/projects/#5280508706029568). [Here](https://github.com/mLavacca/falco-tracing-framework/issues/1) it is possible to find a summarization of my GSOC experience.

## Abstract
The effectiveness of Falco relies on the assumption that it is able to detect all the events in the system, analyze them and produce an output accordingly, therefore a major requirement that this kind of system is required to satisfy, is the capacity to detect all the events, without missing a single one of them. In order to ensure this kind of reliability, it is necessary that Falco is efficient and not affected by performance constraints that can inficiate its effectiveness under both low and high load conditions.

## Introduction
The analysis of the system has been divided into two different parts:
-   stack traces analysis: this kind of analysis traces all the possible branches in the execution flow of Falco and creates a stack trace for each of them, producing a profiling graph;
-   rules checking analysis: all the rules that are loaded into Falco are traced separately and they are profiled thanks to the measurement of the main metrics (counter, latency).

## The analysis system
The analysis system is composed by some different software entities:
-   [Falco-plugin](falco-plugin): some  `C++`  modules that provide functions for tracing operations are built with Falco. A  `C`  wrapper is also provided to allow  `libscap`  to be traced.
-   Falco-instrumentation: function calls to insert in Falco (whose build is conditioned by macros) allow the insertion of tracing points in Falco.
-   [Falco-tracer](falco-tracer): the tracer (written in  `go`) that gets a configuration file (written in  `yaml`), and performs the following operations:
    -   record: launches sysdig with  `-w`  flag (write result on disk) and lets it run for n seconds;
    -   offline-report: launches Falco n times with  `-e`  flag (get events from the  `.scap`  file produced in the previous phase); at the end of each iteration, it loads the metrics written on file by the Falco plugin, keeping up to date a data structure of the average metrics. Last, it creates many different output files needed by the metrics graphic representation tools.
- [Rules-plotter](rules-plotter): a tool for drawing plots of the Falco rules metrics.


## The output format of the system

In order to have a meaningful graphical representation of the system, Falco-tracer produces three different output files that will be used by different tools to produce plots and graphs:

-   custom  `.json`  file representing the system’s stacktrace tree and the rule metrics;
-   `.dot`  file ([dot](https://www.graphviz.org/)  flowchart) representing in the dot format all the possible Falco’s stacktraces;
-   `.folded`  file ([flamegraph](https://github.com/brendangregg/FlameGraph)) representing a flamegraph of Falco.

## The architecture
![architecture](https://github.com/mLavacca/falco-tracing-framework/blob/media/tracing_architecture.png)

## Build
1. Download [Falco](https://github.com/mLavacca/falco/tree/tracing) and [sysdig](https://github.com/mLavacca/sysdig/tree/tracing) tracing branches in this main directory
2. Build [Falco](https://falco.org/docs/installation) with the cmake flag `-DFALCO_TRACE_FLAG="stacktrace"` if you want to profile stack traces or `-DFALCO_TRACE_FLAG="stacktrace"rules"` if you want to profile rules metrics.
3. Build falco-tracer:
```
cd falco-trace
./compile.sh
```
