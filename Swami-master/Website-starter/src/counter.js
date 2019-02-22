function latencyCounter(data, counter) {
  var latency = data.time *1000.0
  if  (latency > counter) {
    counter ++
  }

  return counter
}
