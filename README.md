gosysstat
=========

Uses
====

A CLI which reports CPU stats to your console. Uses https://github.com/ossareh/libgosysstat

Mental Model
============

     +------+
     | main |
     +------+
        |
        |-- fileHandle = os.Open("/proc/<supported target>")
        |-- core.StatProcessor(
        |     supportedTarget.NewProcessor(fileHandle),
        |     chan
        |   )
        |
     +-----------------+
     | resultFormatter |
     +-----------------+
             |
             \-- fmt.Println(result)

