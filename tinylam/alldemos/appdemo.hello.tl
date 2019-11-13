main args := ** concat ["Hello " name "!"]
  name    := args.findWhere (ever true) "World" // if !args "World" (else ** concat ["'" (args.@ 0) "'"])
