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
<form hx-get="/" hx-target="body" hx-swap="outerHTML" hx-push-url="true">
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">Hexadecimal</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <input type="hidden" name="w" value="hex" />
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input type="text" name="input" id="input" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="aabbccddeeff" name="input" />
      <button type="submit" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</form>
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
<form hx-get="/" hx-target="body" hx-swap="outerHTML" hx-push-url="true">
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">Base64</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <input type="hidden" name="w" value="base64" />
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input type="text" name="input" id="input" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="" name="input" />
      <button type="submit" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</form>
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
<form hx-get="/" hx-target="body" hx-swap="outerHTML" hx-push-url="true">
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">ASCII</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <input type="hidden" name="w" value="ascii" />
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input type="text" name="input" id="input" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="" name="input" />
      <button type="submit" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</form>
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
<form hx-get="/" hx-target="body" hx-swap="outerHTML" hx-push-url="true">
  <label for="email" class="block text-sm font-medium leading-6 text-gray-500">Bech32</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <input type="hidden" name="w" value="bech32" />
	<div class="flex flex-col w-full gap-1">
      <input type="text" class="block w-full bg-gray-800 rounded-md border-0 py-1.5 text-gray-50 shadow-sm ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" value="{{ .HRP }}" placeholder="cosmos" name="hrp" />
      <div class="relative flex flex-grow items-stretch focus-within:z-10">
        <input type="text" name="input" id="input" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ .BZ | encode }}" placeholder="" name="input" />
        <button type="submit" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
          Submit
        </button>
      </div>
    </div>
  </div>
</form>
	`)).Execute(w, struct {
		HRP string
		BZ  []byte
	}{
		HRP: hrp,
		BZ:  bz,
	})
	return template.HTML(w.String())
}
