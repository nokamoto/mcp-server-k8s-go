//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var g0 = sh.RunCmd("go")

var Default = Build

func Build() error {
	mg.SerialDeps(Fmt, Import, Tidy, Test, Install)
	return nil
}

func Fmt() error {
	return g0("run", "mvdan.cc/gofumpt", "-w", ".")
}

func Import() error {
	return g0("run", "golang.org/x/tools/cmd/goimports", "-w", ".")
}

func Tidy() error {
	return g0("mod", "tidy")
}

func Test() error {
	return g0("test", "./...")
}

func Install() error {
	return g0("install", "./cmd/mcp-server-k8s-go")
}
