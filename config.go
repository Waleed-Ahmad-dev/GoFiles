package main

import "time"

// Config: The root folder we are browsing
// "." means the current folder where the program is running
const RootFolder = "."
const TrashFolder = ".trash"
const TrashRetention = 30 * 24 * time.Hour // 30 Days