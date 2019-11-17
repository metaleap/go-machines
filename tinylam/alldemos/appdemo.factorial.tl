main args   := args ? => fac | [] => ERR "missing: int arg to factorial"
  fac arg _ := num.factorial (num.parseUnsigned arg)
