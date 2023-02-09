package main

// get processes listening on ports
//
// linux:   netstat -ltnp, filter LISTEN, $7: - or pid/cmd
// windows: netstat -ano, filter LISTENING, $5: pid
