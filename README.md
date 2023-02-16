## knock

`knock program` runs a program and restarts it
when its binary file has been changed.

It is mainly used to restart dev servers written in compiled languages
(hi, Go!).

Contrary to many similar solutions, the trigger for a restart is a change
of a target binary, not the source code change. It is a good assumption,
as the source code is very often compiled-on-save by another tool or ide.

Enjoy!
