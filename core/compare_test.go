package core_test

import (
	. "github.com/HeavyWombat/dyff/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Compare", func() {
	Describe("Difference between YAMLs", func() {
		Context("Given two simple YAML structures", func() {
			It("should return that a string value was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: fOObAr
      version: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/name"))
				Expect(difference.From).To(BeIdenticalTo("foobar"))
				Expect(difference.To).To(BeIdenticalTo("fOObAr"))
			})

			It("should return that an integer was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: 1
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: 2
      version: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/name"))
				Expect(difference.From).To(BeIdenticalTo(1))
				Expect(difference.To).To(BeIdenticalTo(2))
			})

			It("should return that a float was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      level: 3.14159265359
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      level: 2.7182818284
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/level"))
				Expect(difference.From).To(BeNumerically("~", 3.14, 3.15))
				Expect(difference.To).To(BeNumerically("~", 2.71, 2.72))
			})

			It("should return that a boolean was modified", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      enabled: false
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      enabled: true
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/enabled"))
				Expect(difference.From).To(BeIdenticalTo(false))
				Expect(difference.To).To(BeIdenticalTo(true))
			})

			It("should return that one value was added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(ADDITION))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/version"))
				Expect(difference.From).To(BeNil())
				Expect(difference.To).To(BeIdenticalTo("v1"))
			})

			It("should return that one value was removed", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				difference := result[0]
				Expect(difference.Kind).To(BeIdenticalTo(REMOVAL))
				Expect(difference.Path.String()).To(BeIdenticalTo("/some/yaml/structure/version"))
				Expect(difference.From).To(BeIdenticalTo("v1"))
				Expect(difference.To).To(BeNil())
			})

			It("should return two diffs if one value was removed and another one added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      version: v1
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      name: foobar
      release: v1
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(2))

				Expect(result[0].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/version"))
				Expect(result[0].From).To(BeIdenticalTo("v1"))
				Expect(result[0].To).To(BeNil())

				Expect(result[1].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[1].Path.String()).To(BeIdenticalTo("/some/yaml/structure/release"))
				Expect(result[1].From).To(BeNil())
				Expect(result[1].To).To(BeIdenticalTo("v1"))
			})
		})

		Context("Given two YAML structures with simple lists", func() {
			It("should return that a string list entry was added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
      - three
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeNil())
				Expect(result[0].To).To(BeIdenticalTo("three"))
			})

			It("should return that an integer list entry was added", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
      - 3
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeNil())
				Expect(result[0].To).To(BeIdenticalTo(3))
			})

			It("should return that a string list entry was removed", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
      - three
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - one
      - two
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeIdenticalTo("three"))
				Expect(result[0].To).To(BeNil())
			})

			It("should return that an integer list entry was removed", func() {
				from := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
      - 3
`)

				to := getYamlFromString(`---
some:
  yaml:
    structure:
      list:
      - 1
      - 2
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(1))

				Expect(result[0].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/some/yaml/structure/list"))
				Expect(result[0].From).To(BeIdenticalTo(3))
				Expect(result[0].To).To(BeNil())
			})
		})

		Context("Given two YAML structures with complext content", func() {
			It("should return all differences in there", func() {
				from := getYamlFromString(`---
instance_groups:
- name: web
  instances: 1
  resource_pool: concourse_resource_pool
  networks:
  - name: concourse
    static_ips: 192.168.1.1
  jobs:
  - release: concourse
    name: atc
    properties:
      postgresql_database: &atc-db atc
      external_url: http://192.168.1.100:8080
      development_mode: true
  - release: concourse
    name: tsa
    properties: {}

  - name: db
    instances: 1
    resource_pool: concourse_resource_pool
    networks: [{name: concourse}, {name: testnet}]
    persistent_disk: 10240
    jobs:
    - release: concourse
      name: postgresql
      properties:
        databases:
        - name: *atc-db
          role: atc
          password: supersecret
`)

				to := getYamlFromString(`---
instance_groups:
- name: web
  instances: 1
  resource_pool: concourse_resource_pool
  networks:
  - name: concourse
    static_ips: 192.168.0.1
  jobs:
  - release: concourse
    name: atc
    properties:
      postgresql_database: &atc-db atc
      external_url: http://192.168.0.100:8080
      development_mode: false
  - release: concourse
    name: tsa
    properties: {}
  - release: custom
    name: logger

  - name: db
    instances: 2
    resource_pool: concourse_resource_pool
    networks: [{name: concourse}]
    persistent_disk: 10240
    jobs:
    - release: concourse
      name: postgresql
      properties:
        databases:
        - name: *atc-db
          role: atc
          password: "zwX#(;P=%hTfFzM["
`)

				result := CompareDocuments(from, to)
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(BeEquivalentTo(7))

				Expect(result[0].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[0].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/networks/name=concourse/static_ips"))
				Expect(result[0].From).To(BeIdenticalTo("192.168.1.1"))
				Expect(result[0].To).To(BeIdenticalTo("192.168.0.1"))

				Expect(result[1].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[1].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=atc/properties/external_url"))
				Expect(result[1].From).To(BeIdenticalTo("http://192.168.1.100:8080"))
				Expect(result[1].To).To(BeIdenticalTo("http://192.168.0.100:8080"))

				Expect(result[2].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[2].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=atc/properties/development_mode"))
				Expect(result[2].From).To(BeIdenticalTo(true))
				Expect(result[2].To).To(BeIdenticalTo(false))

				Expect(result[3].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[3].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=db/instances"))
				Expect(result[3].From).To(BeIdenticalTo(1))
				Expect(result[3].To).To(BeIdenticalTo(2))

				Expect(result[4].Kind).To(BeIdenticalTo(REMOVAL))
				Expect(result[4].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=db/networks/name=testnet"))
				Expect(result[4].From).To(BeEquivalentTo(yaml.MapSlice{yaml.MapItem{Key: "name", Value: "testnet"}}))
				Expect(result[4].To).To(BeNil())

				Expect(result[5].Kind).To(BeIdenticalTo(MODIFICATION))
				Expect(result[5].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=db/jobs/name=postgresql/properties/databases/name=atc/password"))
				Expect(result[5].From).To(BeIdenticalTo("supersecret"))
				Expect(result[5].To).To(BeIdenticalTo("zwX#(;P=%hTfFzM["))

				Expect(result[6].Kind).To(BeIdenticalTo(ADDITION))
				Expect(result[6].Path.String()).To(BeIdenticalTo("/instance_groups/name=web/jobs/name=logger"))
				Expect(result[6].From).To(BeNil())
				Expect(result[6].To).To(BeEquivalentTo(yaml.MapSlice{yaml.MapItem{Key: "release", Value: "custom"}, yaml.MapItem{Key: "name", Value: "logger"}}))
			})
		})
	})
})
