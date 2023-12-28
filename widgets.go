package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/protocolbuffers/protoscope"
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
  <label class="block text-sm font-medium leading-6 text-gray-500">Hexadecimal</label>
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
  <label class="block text-sm font-medium leading-6 text-gray-500">Base64</label>
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
  <label class="block text-sm font-medium leading-6 text-gray-500">ASCII</label>
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
  <label class="block text-sm font-medium leading-6 text-gray-500">Bech32</label>
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

type Binary struct{}

func (h Binary) ID() string {
	return "binary"
}

func (h Binary) Parse(s string, _ url.Values) ([]byte, error) {
	s = strings.ReplaceAll(s, " ", "")

	if len(s)%8 != 0 {
		s = strings.Repeat("0", 8-len(s)%8) + s
	}

	var res []byte
	for i := 0; i < len(s); i += 8 {
		b, err := strconv.ParseUint(s[i:i+8], 2, 8)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, b)
		res = append(res, buf[7])
	}

	return res, nil
}

func (h Binary) HTML(b []byte, _ url.Values) template.HTML {
	w := new(bytes.Buffer)
	template.Must(template.New("binary").Funcs(template.FuncMap{
		"encode": func(b []byte) string {
			var binaryString strings.Builder
			for _, byt := range b {
				binaryString.WriteString(fmt.Sprintf("%08b", byt))
			}
			return binaryString.String()
		},
	}).Parse(`
<div>
  <label class="block text-sm font-medium leading-6 text-gray-500">Binary</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input id="input-binary" type="text" onfocusin="updateInput('binary')" name="input-binary" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="110011" name="input-binary" />
      <button type="submit" name="w" value="binary" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</div>
    `)).Execute(w, b)
	return template.HTML(w.String())
}

type Decimal struct{}

func (h Decimal) ID() string {
	return "decimal"
}

func (h Decimal) Parse(s string, _ url.Values) ([]byte, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, err
	}
	res := make([]byte, 8)
	binary.BigEndian.PutUint64(res, n)
	for len(res) > 0 && res[0] == 0 {
		res = res[1:]
	}
	return res, nil
}

func (h Decimal) HTML(b []byte, _ url.Values) template.HTML {
	w := new(bytes.Buffer)
	template.Must(template.New("decimal").Funcs(template.FuncMap{
		"encode": func(b []byte) string {
			if len(b) > 8 {
				return "# too large"
			}
			if len(b) < 8 {
				b = append(make([]byte, 8-len(b)), b...)
			}
			return strconv.FormatUint(binary.BigEndian.Uint64(b), 10)
		},
	}).Parse(`
<div>
  <label class="block text-sm font-medium leading-6 text-gray-500">Decimal</label>
  <div class="mt-2 flex rounded-md shadow-sm">
    <div class="relative flex flex-grow items-stretch focus-within:z-10">
      <input id="input-decimal" type="text" onfocusin="updateInput('decimal')" name="input-decimal" class="block w-full bg-gray-800 rounded-none rounded-l-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700 placeholder:text-gray-600 focus:ring-2 focus:ring-inset focus:ring-indigo-800 sm:text-sm sm:leading-6" value="{{ . | encode }}" placeholder="110011" name="input-decimal" />
      <button type="submit" name="w" value="decimal" class="relative -ml-px inline-flex items-center gap-x-1.5 rounded-r-md px-3 py-2 text-sm font-semibold text-gray-50 ring-1 ring-inset ring-gray-700 hover:bg-gray-900">
        Submit
      </button>
    </div>
  </div>
</div>
    `)).Execute(w, b)
	return template.HTML(w.String())
}

type Protobuf struct{}

func (h Protobuf) HTML(b []byte, _ url.Values) template.HTML {
	w := new(bytes.Buffer)

	rendered := protoscope.Write(b, protoscope.WriterOptions{})

	template.Must(template.New("protobuf").Parse(`
<div>
  <label class="block text-sm font-medium leading-6 text-gray-500">Protobuf</label>
  <div class="mt-2 flex">
	<pre class="w-full overflow-x-scroll bg-gray-800 rounded-md border-0 py-1.5 px-3 text-gray-50 ring-1 ring-inset ring-gray-700">{{ . }}</pre>
  </div>
</div>
    `)).Execute(w, rendered)
	return template.HTML(w.String())
}
