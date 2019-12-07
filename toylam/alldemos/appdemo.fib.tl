main args env := args ? .. => fac | [] => ERR "missing: int arg to fibonacci" same
  fac arg _   := arg.parseUnsigned Result? Yay => num.fibonacci | Nay => msg -> ERR msg same
