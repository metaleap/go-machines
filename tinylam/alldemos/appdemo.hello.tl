main args := ** concat ["Hello " name "!"]
  name    := args.findWhere (_ -> True) "World" // if !args "World" (else ** concat ["'" (args.@ 0) "'"])
