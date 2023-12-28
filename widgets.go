package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"html/template"
	"net/url"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

type Hex struct{}

func (h Hex) ID() string {
	return "hex"
}

func (h Hex) Parse(s string, _ url.Values) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	s = strings.ReplaceAll(s, " ", "")
	return hex.DecodeString(s)
}

func (h Hex) HTML(b []byte, _ url.Values) template.HTML {
	w := new(bytes.Buffer)
	template.Must(template.New("hex").Funcs(template.FuncMap{
		"encode": func(b []byte) string {
			return hex.EncodeToString(b)
		},
	}).Parse(`
<div>
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">Hexadecimal</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input id="input-hex" type="text" onfocusin="updateInput('hex')" name="input-hex" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="aabbccddeeff" name="input-hex" />
      <button type="submit" name="w" value="hex" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</div>
    `)).Execute(w, b)
	return template.HTML(w.String())
}

type Base64 struct{}

func (b Base64) ID() string {
	return "base64"
}

func (b Base64) Parse(s string, _ url.Values) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (b Base64) HTML(bz []byte, _ url.Values) template.HTML {
	w := new(bytes.Buffer)
	template.Must(template.New("base64").Funcs(template.FuncMap{
		"encode": func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		},
	}).Parse(`
<div>
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">Base64</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input id="input-base64" type="text" onfocusin="updateInput('base64')" name="input-base64" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="" name="input-base64" />
      <button type="submit" name="w" value="base64" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</div>
    `)).Execute(w, bz)
	return template.HTML(w.String())
}

type ASCII struct{}

func (b ASCII) ID() string {
	return "ascii"
}

func (b ASCII) Parse(s string, _ url.Values) ([]byte, error) {
	return []byte(s), nil
}

func (b ASCII) HTML(bz []byte, _ url.Values) template.HTML {
	w := new(bytes.Buffer)
	template.Must(template.New("base64").Funcs(template.FuncMap{
		"encode": func(b []byte) string {
			return string(b)
		},
	}).Parse(`
<div>
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">ASCII</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input id="input-ascii" type="text" onfocusin="updateInput('ascii')" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="" name="input-ascii" />
      <button type="submit" name="w" value="ascii" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</div>
    `)).Execute(w, bz)
	return template.HTML(w.String())
}

type Bech32 struct{}

func (b Bech32) ID() string {
	return "bech32"
}

func (b Bech32) Parse(s string, params url.Values) ([]byte, error) {
	hrp, decoded, err := bech32.DecodeAndConvert(s)
	if err != nil {
		return nil, err
	}
	if params.Get("hrp") == "" {
		params.Set("hrp", hrp)
	}

	return decoded, nil
}

func (b Bech32) HTML(bz []byte, params url.Values) template.HTML {
	hrp := params.Get("hrp")
	if hrp == "" {
		hrp = "bytez"
	}

	w := new(bytes.Buffer)
	template.Must(template.New("bech32").Funcs(template.FuncMap{
		"encode": func(b []byte) string {
			s, err := bech32.ConvertAndEncode(hrp, b)
			if err != nil {
				return err.Error()
			}

			return s
		},
	}).Parse(`
<div>
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">Bech32</label>
  <div class="mt-2 flex rounded-md shadow-sm">
	<div class="flex flex-col w-full gap-1">
      <input id="input-hrp" type="text" onfocusin="updateInput('bech32');" class="block w-full bg-gray-800 rounded-md border-0 py-1.5 text-gray-50 shadow-sm ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" value="{{ .HRP }}" placeholder="cosmos" name="hrp" />
      <div class="relative flex flex-grow items-stretch focus-within:z-10">
        <input id="input-bech32" type="text" onfocusin="updateInput('bech32')" name="input-bech32" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ .BZ | encode }}" placeholder="" name="input-bech32" />
        <button name="w" value="bech32" type="submit" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
          Submit
        </button>
      </div>
    </div>
  </div>
</div>
	`)).Execute(w, struct {
		HRP string
		BZ  []byte
	}{
		HRP: hrp,
		BZ:  bz,
	})
	return template.HTML(w.String())
}
