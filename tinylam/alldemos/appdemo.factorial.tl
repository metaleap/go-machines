main args   := if ?args (num.factorial (num.parseUnsigned (args.-<))) (ERR errmsg)
  errmsg    := "missing: int arg to factorial"
