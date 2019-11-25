main args env stdin := jsondom Result? Yay => json.toString | Nay => msg -> ERR msg stdin
  jsondom           := json.parsed stdin
