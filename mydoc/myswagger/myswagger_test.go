package myswagger

import (
	"testing"

	"github.com/xxjwxc/public/tools"
)

func TestDomain(t *testing.T) {
	SetHost("http://localhost:8080")
	SetBasePath("/v1")
	doc := NewDoc()
	var tag Tag
	tag.Name = "pet"
	tag.Description = "Everything about your Pets"
	tag.ExternalDocs = &ExternalDocs{
		Description: "Find out more",
		URL:         "https://github.com/xxjwxc/public",
	}
	doc.AddTag(tag)

	var def Definition
	def.Type = "object"
	def.Properties = make(map[string]Propertie)
	def.Properties["id"] = Propertie{
		Type:        "integer",
		Format:      "int64",
		Description: "des text",
	}
	def.Properties["status"] = Propertie{
		Type:        "string",
		Format:      "string",
		Description: "Order Status",
		Enum:        []string{"placed", "approved", "delivered"},
	}

	doc.AddDefinitions("Pet", def)

	var p Param
	p.Tags = []string{"pet"}
	p.Summary = "Add a new pet to the store"
	p.Description = "描述"
	// p.OperationID = "addPet"
	p.Parameters = []Element{Element{
		In:          "body", //  body, header, formData, query, path
		Name:        "body", //  body, header, formData, query, path
		Description: "Pet object that needs to be added to the store",
		Required:    true,
		Schema: Schema{
			Ref: "#/definitions/Pet",
		},
	}}
	doc.AddPatch("/pet", p, "post", "get")

	jsonsrc := doc.GetAPIString()

	tools.WriteFile("/Users/xxj/Downloads/out.json", []string{jsonsrc}, true)

}
