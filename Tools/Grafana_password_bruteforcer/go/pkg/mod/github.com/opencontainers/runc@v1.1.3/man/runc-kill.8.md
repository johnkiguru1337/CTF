% runc-kill "8"

# NAME
**runc-kill** - send a specified signal to container

# SYNOPSIS
**runc kill** [**--all**|**-a**] _container-id_ [_signal_]

# DESCRIPTION

By default, **runc kill** sends **SIGTERM** to the container's initial process
only.

A different signal can be specified either by its name (with or without the
**SIG** prefix), or its numeric value. Use **kill**(1) with **-l** option
to list available signals.

# OPTIONS
**--all**|**-a**
: Send the signal to all processes inside the container.

# EXAMPLES

The following will send a **KILL** signal to the init process of the
**ubuntu01** container:

	# runc kill ubuntu01 KILL

# SEE ALSO

**runc**(1).
