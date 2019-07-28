workflow "build and test" {
  on = "push"
  resolves = ["test"]
}

action "build" {
  uses = "actions-contrib/go@master"
  args = ["build", "-v", "./..."]
}

action "test"{
  uses = "actions-contrib/go@master"
  args = ["test", "-v", "./..."]
  needs = "build"
}
