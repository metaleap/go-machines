main args   := args ? .. => fac | [] => ERR "missing: int arg to factorial"
  fac arg _ := arg.parseUnsigned Result? Yay => num.factorial | Nay => ERR
