main args := if ?args (factorial (num.parseUnsigned (args.at 0))) (ERR errmsg)
  errmsg  := "missing: int arg to factorial"

factorial n := if (n.EQ 0) 1 (else n.MUL (factorial (n.SUB 1)) )
