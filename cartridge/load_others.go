//go:build !unix

package cartridge

import "io/ioutil"

func Load(fn string) (*Data, error) {
	mem, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	res := &Data{
		m: mem,
	}

	if err = res.parse(); err != nil {
		res.Close()
		return nil, err
	}

	return res, nil
}

func (d *Data) unload() {
	d.m = nil
}
