gosysstat
=========

Uses
====

Simple CLI for reading stats from the /proc filesystem on Linux. Uses https://github.com/ossareh/libgosysstat

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
