main args   := args ? => fac | [] => ERR "missing: int arg to factorial"
  fac arg _ := (num.parseUnsigned arg) Result? Yay => num.factorial _ | Nay => ERR _
