This package is a simple Go interface to the C readline(3) library.

Most of the signal handling capabilities of C readline are turned off. The user Go code takes full responsibility for performing all readline-related actions which need to be done when the process receives a signal. The benefit is that it allows the Go user code to implement a proper program shutdown procedure.

This package was forked from http://bitbucket.org/taruti/go-readline by Taru Karttunen, which in turn was forked from http://sigpipe.org/go/readline-go/ by Michael Elkins.
