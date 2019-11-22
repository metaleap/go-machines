main args env             := [onInput, \n, concat "Enter text to reverse:" [\n]]

\n                        := 10

onInput oldState inputLn  := [oldState, [111, 33, \n]]
