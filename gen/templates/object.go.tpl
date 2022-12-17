{{- /*gotype: github.com/oringik/crypto-chateau/gen/ast.ObjectDefinition*/ -}}

type {{.Name | ToCamel}} struct {
{{- range .Fields}}
    {{ .Name | ToCamel }} {{ .Type | GoType }}
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
        arrBuf []byte
        {{- /* TODO: precalculate size based on static fields  */}}
        b = make([]byte, 0, {{ mul (len .Fields) 16 }})
    )

    size := conv.ConvertSizeToBytes(0)
    b = append(b, size...)

	{{- range .Fields}}
	{{- if not .Type.IsArray}}
	{{- template "marshal" dict "Type" .Type "Name" .Name "BufName" "b" "InputVar" (printf "o.%s" (.Name | ToCamel))}}
	{{- else}}
	arrBuf = make([]byte, 0, 128)
	{{- $inputVar := printf "el%s" (.Name | ToCamel) }}
	for _, {{$inputVar}} := range o.{{.Name | ToCamel}} {
		{{- template "marshal" dict "Type" .Type "Name" .Name "BufName" "arrBuf" "InputVar" $inputVar }}
	}
	{{- /*TODO: check if buf size exceess max payload bytes size: max(uint32) */}}
	b = append(b, conv.ConvertSizeToBytes(len(arrBuf))...)
	b = append(b, arrBuf...)
	{{- end}}
    {{- end}}

    size = conv.ConvertSizeToBytes(len(b)-len(size))
    for i := 0; i < len(size); i++ {
        b[i] = 0
    }

	return b
}

{{define "unmarshal"}}
{{- if eqType .Type.Type "uint64"}}
    buf, err = {{ .BufName }}.Slice(8)
	if err != nil {
		return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
	}
	{{ .OutputVar }} = conv.ConvertBytesToUint64(buf)
{{- else if eqType .Type.Type "uint32"}}
    buf, err = {{ .BufName }}.Slice(4)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToUint32(buf)
{{- else if eqType .Type.Type "uint16"}}
    buf, err = {{ .BufName }}.Slice(2)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToUint16(buf)
{{- else if eqType .Type.Type "uint8"}}
    buf, err = {{ .BufName }}.Slice(1)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToUint8(buf)
    {{- else if eqType .Type.Type "int64"}}
    buf, err = {{ .BufName }}.Slice(8)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt64(buf)
{{- else if eqType .Type.Type "int32"}}
    buf, err = {{ .BufName }}.Slice(4)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt32(buf)
{{- else if eqType .Type.Type "int16"}}
    buf, err = {{ .BufName }}.Slice(2)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt16(buf)
{{- else if eqType .Type.Type "int8"}}
    buf, err = {{ .BufName }}.Slice(1)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToInt8(buf)
{{- else if eqType .Type.Type "byte"}}
    buf, err = {{ .BufName }}.Slice(1)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToByte(buf)
{{- else if eqType .Type.Type "bool"}}
    buf, err = {{ .BufName }}.Slice(1)
    if err != nil {
        return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
    }
    {{ .OutputVar }} = conv.ConvertBytesToBool(buf)
{{- else if eqType .Type.Type "string"}}
    size, err = {{ .BufName }}.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read {{ .Name | ToCamel }} size")
	}
	buf, err = {{ .BufName }}.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
	}
    {{ .OutputVar }} = conv.ConvertBytesToString(buf)
{{- else if eqType .Type.Type "object" }}
    size, err = {{ .BufName }}.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read {{ .Name | ToCamel }} size")
	}
	buf, err = {{ .BufName }}.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read {{ .Name | ToCamel }}")
	}
    if err = {{ .OutputVar }}.Unmarshal(buf); err != nil {
        return errors.Wrap(err, "failed to unmarshal {{ .Name | ToCamel }}")
    }
{{- end }}
{{- end}}

func (o *{{.Name | ToCamel}}) Unmarshal(b *conv.BinaryIterator) error {
	var (
		err             error
		size, arrSize   int
		buf, arrBuf     *conv.BinaryIterator
	)

	{{- range .Fields}}
    {{- if not .Type.IsArray}}
    {{ $outputVar := printf "o.%s" (.Name | ToCamel) }}
    {{- template "unmarshal" dict "Type" .Type "Name" .Name "BufName" "b" "OutputVar" $outputVar}}
    {{- else}}
	size, err = b.NextSize()
	if err != nil {
		return errors.Wrap(err, "failed to read {{.Name | ToCamel}} size")
	}
	arrBuf, err = b.Slice(size)
	if err != nil {
		return errors.Wrap(err, "failed to read {{.Name | ToCamel}}")
	}
	for arrBuf.HasNext() {
        {{- $outputVar := printf "el%s" (.Name | ToCamel) }}
        var {{$outputVar}} {{.Type | GoType }}
        {{- template "unmarshal" dict "Type" .Type "Name" .Name "BufName" "arrBuf" "OutputVar" $outputVar}}
        {{- /*TODO: append won't work with static arrays like [16]byte */}}
		o.MagicObjectList = append(o.MagicObjectList, {{$outputVar}})
	}
	{{- end}}
    {{- end}}

	return nil
}
