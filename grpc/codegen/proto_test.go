package codegen

import (
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/expr"
	"goa.design/goa/v3/grpc/codegen/testdata"
)

func TestProtoFiles(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"protofiles-unary-rpcs", testdata.UnaryRPCsDSL, testdata.UnaryRPCsProtoCode},
		{"protofiles-unary-rpc-no-payload", testdata.UnaryRPCNoPayloadDSL, testdata.UnaryRPCNoPayloadProtoCode},
		{"protofiles-unary-rpc-no-result", testdata.UnaryRPCNoResultDSL, testdata.UnaryRPCNoResultProtoCode},
		{"protofiles-server-streaming-rpc", testdata.ServerStreamingRPCDSL, testdata.ServerStreamingRPCProtoCode},
		{"protofiles-client-streaming-rpc", testdata.ClientStreamingRPCDSL, testdata.ClientStreamingRPCProtoCode},
		{"protofiles-bidirectional-streaming-rpc", testdata.BidirectionalStreamingRPCDSL, testdata.BidirectionalStreamingRPCProtoCode},
		{"protofiles-same-service-and-message-name", testdata.MessageWithServiceNameDSL, testdata.MessageWithServiceNameProtoCode},
		{"protofiles-method-with-reserved-proto-name", testdata.MethodWithReservedNameDSL, testdata.MethodWithReservedNameProtoCode},
		{"protofiles-multiple-methods-same-return-type", testdata.MultipleMethodsSameResultCollectionDSL, testdata.MultipleMethodsSameResultCollectionProtoCode},
		{"protofiles-method-with-acronym", testdata.MethodWithAcronymDSL, testdata.MethodWithAcronymProtoCode},
		{"protofiles-custom-package-name", testdata.ServiceWithPackageDSL, testdata.ServiceWithPackageCode},
		{"protofiles-struct-meta-type", testdata.StructMetaTypeDSL, testdata.StructMetaTypePackageCode},
		{"protofiles-default-fields", testdata.DefaultFieldsDSL, testdata.DefaultFieldsPackageCode},
		{"protofiles-custom-message-name", testdata.CustomMessageNameDSL, testdata.CustomMessageNamePackageCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunGRPCDSL(t, c.DSL)
			fs := ProtoFiles("", expr.Root)
			if len(fs) != 1 {
				t.Fatalf("got %d files, expected one", len(fs))
			}
			sections := fs[0].SectionTemplates
			if len(sections) < 3 {
				t.Fatalf("got %d sections, expected at least three", len(sections))
			}
			code := sectionCode(t, sections[1:]...)
			if code != c.Code {
				t.Errorf("%s: got\n%s\ngot vs. expected:\n%s", c.Name, code, codegen.Diff(t, code, c.Code))
			}
			fpath := codegen.CreateTempFile(t, code)
			if err := protoc(fpath, nil); err != nil {
				t.Fatalf("error occurred when compiling proto file %q: %s", fpath, err)
			}
		})
	}
}

func TestMessageDefSection(t *testing.T) {
	cases := []struct {
		Name string
		DSL  func()
		Code string
	}{
		{"user-type-with-primitives", testdata.MessageUserTypeWithPrimitivesDSL, testdata.MessageUserTypeWithPrimitivesMessageCode},
		{"user-type-with-alias", testdata.MessageUserTypeWithAliasMessageDSL, testdata.MessageUserTypeWithAliasMessageCode},
		{"user-type-with-nested-user-types", testdata.MessageUserTypeWithNestedUserTypesDSL, testdata.MessageUserTypeWithNestedUserTypesCode},
		{"result-type-collection", testdata.MessageResultTypeCollectionDSL, testdata.MessageResultTypeCollectionCode},
		{"user-type-with-collection", testdata.MessageUserTypeWithCollectionDSL, testdata.MessageUserTypeWithCollectionCode},
		{"array", testdata.MessageArrayDSL, testdata.MessageArrayCode},
		{"map", testdata.MessageMapDSL, testdata.MessageMapCode},
		{"primitive", testdata.MessagePrimitiveDSL, testdata.MessagePrimitiveCode},
		{"with-metadata", testdata.MessageWithMetadataDSL, testdata.MessageWithMetadataCode},
		{"with-security-attributes", testdata.MessageWithSecurityAttrsDSL, testdata.MessageWithSecurityAttrsCode},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			RunGRPCDSL(t, c.DSL)
			fs := ProtoFiles("", expr.Root)
			if len(fs) != 1 {
				t.Fatalf("got %d files, expected one", len(fs))
			}
			sections := fs[0].SectionTemplates
			if len(sections) < 3 {
				t.Fatalf("got %d sections, expected at least three", len(sections))
			}
			code := sectionCode(t, sections[:2]...)
			msgCode := sectionCode(t, sections[3:]...)
			if msgCode != c.Code {
				t.Errorf("%s: got\n%s\ngot vs. expected:\n%s", c.Name, msgCode, codegen.Diff(t, msgCode, c.Code))
			}
			fpath := codegen.CreateTempFile(t, code+msgCode)
			if err := protoc(fpath, nil); err != nil {
				t.Fatalf("error occurred when compiling proto file %q: %s", fpath, err)
			}
		})
	}
}
