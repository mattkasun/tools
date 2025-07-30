# tools
set of utility functions
* PrettyByteSize - returns a human readable string of int bytes  
    
    ${999 => 999 \color{green}\space B}$  
    ${2048 => 2.0 \color{green}\space KiB}$  
    ${1058575 => 1.01 \color{green}\space MiB}$

## packages
### logging
slog helpers
* create logger with default, disard, text or json handler
* options to set output, loglevel, timeformat, include source, truncate source