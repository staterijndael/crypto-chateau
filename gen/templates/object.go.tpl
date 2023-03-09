{{- /*gotype: github.com/oringik/crypto-chateau/gen/ast.ObjectDefinition*/ -}}

type {{.Name | ToCamel}} struct {
{{- range .Fields}}
    {{if or (not .Type.IsArray) (eq (.Type.ArrSize) 0)}}
    {{ .Name | ToCamel }} {{ .Type | GoType }}
    {{else}}
    {{ .Name | ToCamel }} *{{ .Type | GoType }}
    {{end}}
{{- end}}
}

var _ message.Message = (*{{.Name | ToCamel}})(nil)

{{define "marshal"}}
{{- if eqType .Type.Type "uint64"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertUint64ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "uint32"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertUint32ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "uint16"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertUint16ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "uint8"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertUint8ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "int64"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertInt64ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "int32"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertInt32ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "int16"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertInt16ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "int8"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertInt8ToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "int"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertIntToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "byte"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertByteToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "bool"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertBoolToBytes({{ .InputVar }})...)
{{- else if eqType .Type.Type "string"}}
	{{ .BufName }} = append({{ .BufName }},conv.ConvertSizeToBytes(len([]byte({{ .InputVar }})))...)
	{{ .BufName }} = append({{ .BufName }},conv.ConvertStringToBytes({{ .InputVar }})...)
{{- else }}
	{{ .BufName }} = append({{ .BufName }},{{ .InputVar }}.Marshal()...)
{{- end}}
{{- end}}

func (o *{{.Name | ToCamel}}) Marshal() []byte {
    var (
        {{- /* TODO: precalculate size based on static fields  */}}
        b = make([]byte, 0, {{ mul (len .Fields) 16 }})
    )

    size := conv.ConvertSizeToBytes(0)
    b = append(b, size...)

	{{- range .Fields}}
	{{- if not .Type.IsArray}}
	{{- template "marshal" dict "Type" .Type "Name" .Name "BufName" "b" "InputVar" (printf "o.%s" (.Name | ToCamel))}}
	{{- else}}
	{{- $arrBufName := printf "arrBuf%s" (.Name | ToCamel) }}
	{{$arrBufName}} := make([]byte, 0, 128)
	{{- $inputVar := printf "el%s" (.Name | ToCamel) }}
	for _, {{$inputVar}} := range o.{{.Name | ToCamel}} {
		{{- template "marshal" dict "Type" .Type "Name" .Name "BufName" $arrBufName "InputVar" $inputVar }}
	}
	{{- /*TODO: check if buf size exceess max payload bytes size: max(uint32) */}}
	b = append(b, conv.ConvertSizeToBytes(len({{$arrBufName}}))...)
	b = append(b, {{$arrBufName}}...)
	{{- end}}
    {{- end}}

    size = conv.ConvertSizeToBytes(len(b)-len(size))
    for i := 0; i < len(size); i++ {
        b[i] = size[i]
    }

	return b
}

{{define "unmarshal"}}
{{if eqType .Type.Type "uint64"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(8)
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
	}
	{{ .OutputVar }} = conv.ConvertBytesToUint64(binaryCtx.buf)
{{ else if eqType .Type.Type "uint32"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(4)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToUint32(binaryCtx.buf)
{{ else if eqType .Type.Type "uint16"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(2)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToUint16(binaryCtx.buf)
{{ else if eqType .Type.Type "uint8"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(1)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToUint8(binaryCtx.buf)
{{- else if eqType .Type.Type "int64"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(8)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt64(binaryCtx.buf)
{{ else if eqType .Type.Type "int32"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(4)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt32(binaryCtx.buf)
{{ else if eqType .Type.Type "int16"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(2)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt16(binaryCtx.buf)
{{ else if eqType .Type.Type "int8"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(1)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt8(binaryCtx.buf)
{{ else if eqType .Type.Type "int"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(8)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt(binaryCtx.buf)
{{ else if eqType .Type.Type "byte"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(1)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToByte(binaryCtx.buf)
{{ else if eqType .Type.Type "bool"}}
    binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(1)
    if binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToBool(binaryCtx.buf)
{{ else if eqType .Type.Type "string"}}
    binaryCtx.size, binaryCtx.err = {{ .BufName }}.NextSize()
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }} size")
	}
	binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(binaryCtx.size)
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
	}
    {{ .OutputVar }} = conv.ConvertBytesToString(binaryCtx.buf)
{{ else if eqType .Type.Type "object" }}
    binaryCtx.size, binaryCtx.err = {{ .BufName }}.NextSize()
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }} size")
	}
	binaryCtx.buf, binaryCtx.err = {{ .BufName }}.Slice(binaryCtx.size)
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{ .Name | ToCamel }}")
	}
    if binaryCtx.err = {{ .OutputVar }}.Unmarshal(binaryCtx.buf); binaryCtx.err != nil {
        return errors.Wrap(binaryCtx.err, "failed to unmarshal {{ .Name | ToCamel }}")
    }
{{ end }}
{{ end}}

func (o *{{.Name | ToCamel}}) Unmarshal(b *conv.BinaryIterator) error {
	binaryCtx := struct {
		err           error
		size, arrSize, pos int
		buf, arrBuf   *conv.BinaryIterator
	}{}
	
	binaryCtx.err = nil

	{{- range .Fields}}
    {{- if not .Type.IsArray}}
    {{ $outputVar := printf "o.%s" (.Name | ToCamel) }}
    {{- template "unmarshal" dict "Type" .Type "Name" .Name "BufName" "b" "OutputVar" $outputVar}}
    {{- else}}
	binaryCtx.size, binaryCtx.err = b.NextSize()
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{.Name | ToCamel}} size")
	}
	binaryCtx.arrBuf, binaryCtx.err = b.Slice(binaryCtx.size)
	if binaryCtx.err != nil {
		return errors.Wrap(binaryCtx.err, "failed to read {{.Name | ToCamel}}")
	}
	{{- if not (eq .Type.ArrSize 0) }}
	binaryCtx.pos = 0
	{{- end}}
	for binaryCtx.arrBuf.HasNext() {
        {{- $outputVar := printf "el%s" (.Name | ToCamel) -}}
        var {{$outputVar}} {{ GoType .Type true }}
        {{- template "unmarshal" dict "Type" .Type "Name" .Name "BufName" "binaryCtx.arrBuf" "OutputVar" $outputVar }}
        {{if eq .Type.ArrSize 0 -}}
        o.{{.Name | ToCamel}} = append(o.{{.Name | ToCamel}}, {{$outputVar}})
		{{ else -}}
		o.{{.Name | ToCamel}}[binaryCtx.pos] = {{$outputVar}}
		binaryCtx.pos++
		{{- end}}
	}
	{{- end}}
    {{- end}}

    return nil
}

func (o *{{.Name | ToCamel}}) Copy() message.Message{
    return {{printf "&%s{}" (.Name | ToCamel)}}
}
