main args env := [on, 0, -0, ""]
  on _ jsrc   := [same , jsrc.parsed Result? Yay => json.stringified | Nay => concat "err: "]
