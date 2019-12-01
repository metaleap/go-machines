main args env             := [onInput, \n, -0, concat "Enter text to reverse, empty aborts:" [\n]]

\n                        := 10

onInput prevstate inputLn := [ nextstate , reversed (\n.+> inputLn) ]
  nextstate               := inputLn ? [] => same | .. => _ _ -> prevstate
