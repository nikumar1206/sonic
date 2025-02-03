package main

// basically keep a stack
// stack item can be a fully valid json (nested json exists)
// stack item is popped everytime we encounter a closing bracket or a closing ], depending on stack opener
// stack item is added when we encounter an opening bracket or opening [
