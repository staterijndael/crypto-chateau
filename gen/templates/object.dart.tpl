{{define "marshal"}}
{{- if eqType .Type.Type "uint64"}}
	{{ .BufName }}.addAll(ConvertUint64ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "uint32"}}
	{{ .BufName }}.addAll(ConvertUint32ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "uint16"}}
	{{ .BufName }}.addAll(ConvertUint16ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "uint8"}}
	{{ .BufName }}.addAll(ConvertUint8ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "int64"}}
	{{ .BufName }}.addAll(ConvertInt64ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "int32"}}
	{{ .BufName }}.addAll(ConvertInt32ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "int16"}}
	{{ .BufName }}.addAll(ConvertInt16ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "int8"}}
	{{ .BufName }}.addAll(ConvertInt8ToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "int"}}
	{{ .BufName }}.addAll(ConvertIntToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "byte"}}
	{{ .BufName }}.addAll(ConvertByteToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "bool"}}
	{{ .BufName }}.addAll(ConvertBoolToBytes({{ .InputVar }}));
{{- else if eqType .Type.Type "string"}}
	{{ .BufName }}.addAll(ConvertSizeToBytes({{ .InputVar }}.codeUnits.length));
	{{ .BufName }}.addAll(ConvertStringToBytes({{ .InputVar }}));
{{- else }}
		{{ .BufName }}.addAll({{ .InputVar }}.Marshal());
{{- end}}
{{- end}}

class {{.Name | ToCamel}} implements Message {
  {{- range .Fields}}
  {{ .Type | DartType }} {{ .Name | ToCamel }};
  {{- end}}

  {{- if not (eq (.Fields | GetSliceLength) 0) }}
  {{.Name | ToCamel}}({
    {{- range .Fields}}
    required this.{{ .Name | ToCamel }},
    {{- end}}
  });
  {{else}}
  {{.Name | ToCamel}}();
  {{end}}

  {{.Name | ToCamel}} Copy() => {{.Name | ToCamel}}({{printf "%s" (FillDefaultObjectParams (.Name))}});

  Uint8List Marshal() {
      List<int> b = [];

      List<int> size = ConvertSizeToBytes(0);
      b.addAll(size);

      {{- range .Fields}}
      {{- if not .Type.IsArray}}
      {{- template "marshal" dict "Type" .Type "Name" .Name "BufName" "b" "InputVar" (printf "%s!" (.Name | ToCamel))}}
      {{- else}}
      {{- $arrBufName := printf "arrBuf%s" (.Name | ToCamel) }}
      List<int> {{$arrBufName}} = [];
      {{- $inputVar := printf "el%s" (.Name | ToCamel) }}
      for (var {{$inputVar}} in {{.Name | ToCamel}}!) {
          {{- template "marshal" dict "Type" .Type "Name" .Name "BufName" $arrBufName "InputVar" $inputVar }}
      }
      {{- /*TODO: check if buf size exceess max payload bytes size: max(uint32) */}}
      b.addAll(ConvertSizeToBytes({{$arrBufName}}.length));
      b.addAll({{$arrBufName}});
      {{- end}}
      {{- end}}
      size = ConvertSizeToBytes(b.length - size.length);
      for (int i = 0; i < size.length; i++) {
      	b[i] = size[i];
      }

      return Uint8List.fromList(b);
  }

  {{define "unmarshal"}}
  {{if eqType .Type.Type "uint64"}}
      binaryCtx.buf = {{ .BufName }}.slice(8);
  	{{ .OutputVar }} = ConvertBytesToUint64(binaryCtx.buf);
  {{ else if eqType .Type.Type "uint32"}}
      binaryCtx.buf = {{ .BufName }}.slice(4);
      {{ .OutputVar }} = ConvertBytesToUint32(binaryCtx.buf);
  {{ else if eqType .Type.Type "uint16"}}
      binaryCtx.buf = {{ .BufName }}.slice(2);
      {{ .OutputVar }} = ConvertBytesToUint16(binaryCtx.buf);
  {{ else if eqType .Type.Type "uint8"}}
      binaryCtx.buf = {{ .BufName }}.slice(1);
      {{ .OutputVar }} = ConvertBytesToUint8(binaryCtx.buf);
  {{- else if eqType .Type.Type "int64"}}
      binaryCtx.buf = {{ .BufName }}.slice(8);
      {{ .OutputVar }} = ConvertBytesToInt64(binaryCtx.buf);
  {{ else if eqType .Type.Type "int32"}}
      binaryCtx.buf = {{ .BufName }}.slice(4);
      {{ .OutputVar }} = ConvertBytesToInt32(binaryCtx.buf);
  {{ else if eqType .Type.Type "int16"}}
      binaryCtx.buf = {{ .BufName }}.slice(2);
      {{ .OutputVar }} = ConvertBytesToInt16(binaryCtx.buf);
  {{ else if eqType .Type.Type "int8"}}
      binaryCtx.buf = {{ .BufName }}.slice(1);
      {{ .OutputVar }} = ConvertBytesToInt8(binaryCtx.buf);
  {{- else if eqType .Type.Type "int64"}}
      binaryCtx.buf = {{ .BufName }}.slice(8);
      {{ .OutputVar }} = ConvertBytesToInt(binaryCtx.buf);
  {{ else if eqType .Type.Type "byte"}}
      binaryCtx.buf = {{ .BufName }}.slice(1);
      {{ .OutputVar }} = ConvertBytesToByte(binaryCtx.buf);
  {{ else if eqType .Type.Type "bool"}}
      binaryCtx.buf = {{ .BufName }}.slice(1);
      {{ .OutputVar }} = ConvertBytesToBool(binaryCtx.buf);
  {{ else if eqType .Type.Type "string"}}
      binaryCtx.size = {{ .BufName }}.nextSize();
  	binaryCtx.buf = {{ .BufName }}.slice(binaryCtx.size);
      {{ .OutputVar }} = ConvertBytesToString(binaryCtx.buf);
  {{ else if eqType .Type.Type "object" }}
      binaryCtx.size = {{ .BufName }}.nextSize();
  	binaryCtx.buf = {{ .BufName }}.slice(binaryCtx.size);
      {{ .OutputVar }}!.Unmarshal(binaryCtx.buf);
  {{ end }}
  {{ end}}

  void Unmarshal(BinaryIterator b) {
  	BinaryCtx binaryCtx = BinaryCtx();

  	{{- range .Fields}}
      {{- if not .Type.IsArray}}
      {{ $outputVar := (.Name | ToCamel) }}
      {{- template "unmarshal" dict "Type" .Type "Name" .Name "BufName" "b" "OutputVar" $outputVar}}
      {{- else}}
  	binaryCtx.size = b.nextSize();

  	binaryCtx.arrBuf = b.slice(binaryCtx.size);
  	binaryCtx.pos = 0;

  	{{.Name | ToCamel}}.extend(binaryCtx.size, {{.Type | FillDefaultValue}});
  	while (binaryCtx.arrBuf.hasNext()) {
  	  	   {{$outputVar := printf "el%s" (.Name | ToCamel) }}
  	  	   {{if eq .Type.ObjectName ""}}
  	  	       {{ DartType .Type true }} {{$outputVar}};
  	  	   {{else -}}
  	  	       {{ DartType .Type true }} {{printf "%s = %s(%s)" $outputVar (DartType .Type true) (FillDefaultObjectParams (DartType .Type true))}};
  	  	   {{end}}
          {{- template "unmarshal" dict "Type" .Type "Name" .Name "BufName" "binaryCtx.arrBuf" "OutputVar" $outputVar }}
           {{.Name | ToCamel}}![binaryCtx.pos] = {{$outputVar}};
           binaryCtx.pos++;
  	}
  	{{- end}}
      {{- end}}
  }

}