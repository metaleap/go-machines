main stdin  := jsondom Result? Nay => ERR | Yay => json.toString
  jsondom   := json.parsed stdin
