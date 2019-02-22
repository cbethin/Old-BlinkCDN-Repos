package blink

type Process struct {
  Number        int
  Callback      func()
  HasCompleted  bool
}
