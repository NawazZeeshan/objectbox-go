/*
 * Copyright 2018 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package templates

import (
	"text/template"
)

var BindingTemplate = template.Must(template.New("binding").Funcs(funcMap).Parse(
	`// Code generated by ObjectBox; DO NOT EDIT.

package {{.Package}}

import (
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	{{if .UsesFbUtils}}"github.com/objectbox/objectbox-go/objectbox/fbutils"{{end}}
)

{{range $entity := .Entities -}}
{{$entityNameCamel := $entity.Name | StringCamel -}}
type {{$entityNameCamel}}_EntityInfo struct {
	Id objectbox.TypeId
	Uid uint64
}

var {{$entity.Name}}Binding = {{$entityNameCamel}}_EntityInfo {
	Id: {{$entity.Id}}, 
	Uid: {{$entity.Uid}},
}

var {{$entity.Name}}_ = struct {
	{{range $property := $entity.Properties -}}
    {{$property.Name}} objectbox.TypeId
    {{end -}}
}{
	{{range $property := $entity.Properties -}}
    {{$property.Name}}: {{$property.Id}},
    {{end -}}
}

func ({{$entityNameCamel}}_EntityInfo) AddToModel(model *objectbox.Model) {
    model.Entity("{{$entity.Name}}", {{$entity.Id}}, {{$entity.Uid}})
    {{range $property := $entity.Properties -}}
    model.Property("{{$property.ObName}}", objectbox.PropertyType_{{$property.ObType}}, {{$property.Id}}, {{$property.Uid}})
    {{if len $property.ObFlags -}}
        model.PropertyFlags(
        {{- range $key, $flag := $property.ObFlags -}}
            {{if gt $key 0}} | {{end}}objectbox.PropertyFlags_{{$flag}}
        {{- end}})
        {{- /* TODO model.PropertyIndexId() && model.PropertyRelation() */}}
    {{end -}}
	{{if $property.Relation}}model.PropertyRelation("{{$property.Relation.Target}}", {{$property.Index.Id}}, {{$property.Index.Uid}})
	{{else if $property.Index}}model.PropertyIndex({{$property.Index.Id}}, {{$property.Index.Uid}})
    {{end -}}
    {{end -}}
    model.EntityLastPropertyId({{$entity.LastPropertyId.GetId}}, {{$entity.LastPropertyId.GetUid}})
}

func ({{$entityNameCamel}}_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*{{$entity.Name}}).{{$entity.IdProperty.Name}}, nil
}

func ({{$entityNameCamel}}_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*{{$entity.Name}}).{{$entity.IdProperty.Name}} = id
	return nil
}

func ({{$entityNameCamel}}_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) {
    {{if $entity.HasNonIdProperty}}obj := object.(*{{$entity.Name}}){{end -}}

    {{- range $property := $entity.Properties}}
        {{- if eq $property.FbType "UOffsetT"}}
            {{- if eq $property.GoType "string"}}
    var offset{{$property.Name}} = fbutils.CreateStringOffset(fbb, obj.{{$property.Name}})
            {{- else if eq $property.GoType "[]byte"}}
    var offset{{$property.Name}} = fbutils.CreateByteVectorOffset(fbb, obj.{{$property.Name}})
            {{- else -}}
            TODO offset creation for the {{$property.Name}}, type ${{$property.GoType}} is not implemented
            {{- end -}}
        {{end}}
    {{- end}}

    // build the FlatBuffers object
    fbb.StartObject({{$entity.LastPropertyId.GetId}})
    {{range $property := $entity.Properties -}}
    fbb.Prepend{{$property.FbType}}Slot({{$property.FbSlot}},
        {{- if eq $property.FbType "UOffsetT"}} offset{{$property.Name}}, 0)
        {{- else if eq $property.Name $entity.IdProperty.Name}} id, 0)
        {{- else if eq $property.GoType "bool"}} obj.{{$property.Name}}, false)
        {{- else if eq $property.GoType "int"}} int32(obj.{{$property.Name}}), 0)
        {{- else if eq $property.GoType "uint"}} uint32(obj.{{$property.Name}}), 0)
        {{- else}} obj.{{$property.Name}}, 0)
        {{- end}}
    {{end -}}
}

func ({{$entityNameCamel}}_EntityInfo) ToObject(bytes []byte) interface{} {
	table := &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	return &{{$entity.Name}}{
	{{- range $property := $entity.Properties}}
		{{$property.Name}}: {{if eq $property.GoType "bool"}} table.GetBoolSlot({{$property.FbvTableOffset}}, false)
        {{- else if eq $property.GoType "int"}} int(table.GetUint32Slot({{$property.FbvTableOffset}}, 0))
        {{- else if eq $property.GoType "uint"}} uint(table.GetUint32Slot({{$property.FbvTableOffset}}, 0))
		{{- else if eq $property.GoType "rune"}} rune(table.GetInt32Slot({{$property.FbvTableOffset}}, 0))
		{{- else if eq $property.GoType "string"}} fbutils.GetStringSlot(table, {{$property.FbvTableOffset}})
        {{- else if eq $property.GoType "[]byte"}} fbutils.GetByteVectorSlot(table, {{$property.FbvTableOffset}})
		{{- else}} table.Get{{$property.GoType | StringTitle}}Slot({{$property.FbvTableOffset}}, 0)
        {{- end}},
	{{- end}}
	}
}

func ({{$entityNameCamel}}_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*{{$entity.Name}}, 0, capacity)
}

func ({{$entityNameCamel}}_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	return append(slice.([]*{{$entity.Name}}), object.(*{{$entity.Name}}))
}

type {{$entity.Name}}Box struct {
	*objectbox.Box
}

func BoxFor{{$entity.Name}}(ob *objectbox.ObjectBox) *{{$entity.Name}}Box {
	return &{{$entity.Name}}Box{
		Box: ob.Box({{$entity.Id}}),
	}
}

func (box *{{$entity.Name}}Box) Put(object *{{$entity.Name}}) ({{$entity.IdProperty.GoType}}, error) {
	return box.Box.Put(object)
}

func (box *{{$entity.Name}}Box) PutAsync(object *{{$entity.Name}}) ({{$entity.IdProperty.GoType}}, error) {
	return box.Box.PutAsync(object)
}

func (box *{{$entity.Name}}Box) PutAll(objects []*{{$entity.Name}}) ([]{{$entity.IdProperty.GoType}}, error) {
	return box.Box.PutAll(objects)
}

func (box *{{$entity.Name}}Box) Get(id {{$entity.IdProperty.GoType}}) (*{{$entity.Name}}, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*{{$entity.Name}}), nil
}

func (box *{{$entity.Name}}Box) GetAll() ([]*{{$entity.Name}}, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*{{$entity.Name}}), nil
}

func (box *{{$entity.Name}}Box) Remove(object *{{$entity.Name}}) (err error) {
	return box.Box.Remove(object.{{$entity.IdProperty.Name}})
}

{{end -}}`))
