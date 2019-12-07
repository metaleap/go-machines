main args env := args ? .. => fac | [] => ERR "missing: int arg to factorial" same
  fac arg _   := arg.parseUnsigned Result? Yay => num.factorial | Nay => msg -> ERR msg same
