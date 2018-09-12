package main

// Hypervisor interface
type Hypervisor interface {
    start(string) error
}

// available hypervisors
var hypervisors = map[string]func() Hypervisor {
        "qemu-system-x86_64" : newQemu,
    } 
