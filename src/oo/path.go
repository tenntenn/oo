package oo

import (
	"os"
	"path/filepath"
)

type Path string

func (p Path) String() string {
	if string(p) == "." {
		return ""
	}

	return string(p)
}

func (p Path) Dir() Path {
	return Path(filepath.Dir(string(p)))
}

func (p Path) Ext() Path {
	return Path(filepath.Ext(string(p)))
}

func (p Path) Abs() Path {
	abs, err := filepath.Abs(string(p))
	if err != nil {
		return p
	}

	return Path(abs)
}

func (p Path) Base() Path {
	return Path(filepath.Base(string(p)))
}

func (p Path) Rel() Path {
	cd := Path(filepath.Dir("."))
	rel, err := filepath.Rel(cd.Abs().String(), p.Abs().String())
	if err != nil {
		return p
	}

	return Path(rel)
}

func (p Path) cd() Path {
	pre := Path(filepath.Dir(".")).Abs()

	os.Chdir(p.String())

	return pre
}
