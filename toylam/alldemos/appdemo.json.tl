main args env := [on, 0, -0, ""]
  on _ jsrc   := [same , (json.parsed jsrc) Result? Yay => json.stringified | Nay => same]
