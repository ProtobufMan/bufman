package parse

//
//import (
//	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
//	"github.com/bufbuild/protocompile/protoutil"
//	"google.golang.org/protobuf/reflect/protoreflect"
//	"strings"
//)
//
//type DocumentGenerator struct {
//	messages map[string]*registryv1alpha1.Message
//}
//
//func (g *DocumentGenerator) getNestedName(name string, nestedPrefixes ...string) string {
//	parts := append(nestedPrefixes, name)
//	return strings.Join(parts, ".")
//}
//
//func (g *DocumentGenerator) GetMessages(messageDescriptors protoreflect.MessageDescriptors, nestedPrefixes ...string) []*registryv1alpha1.Message {
//	for i := 0; i < messageDescriptors.Len(); i++ {
//		messageDescriptor := messageDescriptors.Get(i)
//
//		message := g.GetMessage(messageDescriptor, nestedPrefixes...)
//
//
//	}
//
//	return messages
//}
//
//func (g *DocumentGenerator) GetMessage(messageDescriptor protoreflect.MessageDescriptor, nestedPrefixes ...string) *registryv1alpha1.Message {
//	if message, ok := g.messages[string(messageDescriptor.Name())]; ok {
//		return message
//	}
//
//	messageLocation := messageDescriptor.ParentFile().SourceLocations().ByDescriptor(messageDescriptor)
//	message := &registryv1alpha1.Message{
//		Name: string(messageDescriptor.Name()),
//		NestedName:           g.getNestedName(string(messageDescriptor.Name()), nestedPrefixes...),
//		FullName:             string(messageDescriptor.FullName()),
//		Description:          messageLocation.LeadingComments,
//		FilePath:             messageDescriptor.ParentFile().Path(),
//		IsMapEntry:           messageDescriptor.IsMapEntry(),
//		Fields:               nil,
//		Location:             &registryv1alpha1.Location{
//			StartLine:   int32(messageLocation.StartLine),
//			StartColumn: int32(messageLocation.StartColumn),
//			EndLine:     int32(messageLocation.EndLine),
//			EndColumn: int32(messageLocation.EndColumn),
//		},
//		MessageExtensions:    nil,
//		MessageOptions:       &registryv1alpha1.MessageOptions{
//			Deprecated: protoutil.ProtoFromMessageDescriptor(messageDescriptor).Options.GetDeprecated(),
//		},
//		ImplicitlyDeprecated: false,
//	}
//
//	// fields
//
//
//	// extensions
//
//	// nested message
//	messageDescriptor.Messages()
//
//
//	// nested enum
//	messageDescriptor.Enums()
//}
//
//func (g *DocumentGenerator) getFields(fieldDescriptors protoreflect.FieldDescriptors) {
//	for i := 0; i < fieldDescriptors.Len(); i++ {
//		fieldDescriptor := fieldDescriptors.Get(i)
//		fieldOptions := protoutil.ProtoFromFieldDescriptor(fieldDescriptor).Options
//		field := &registryv1alpha1.Field{
//			Name:            string(fieldDescriptor.Name()),
//			Description:     "",
//			Label:           "",
//			NestedType:      "",
//			FullType:        "",
//			Tag: uint32(fieldDescriptor.Number()),
//			MapEntry:        &registryv1alpha1.MapEntry{
//				KeyFullType:          "",
//				ValueNestedType:      "",
//				ValueFullType:        "",
//				ValueImportModuleRef: nil,
//			},
//			ImportModuleRef: nil,
//			Extendee:        "",
//			FieldOptions:    &registryv1alpha1.FieldOptions{
//				Deprecated: fieldOptions.GetDeprecated(),
//				Packed:     fieldOptions.Packed,
//				Ctype: int32(fieldOptions.GetCtype()),
//				Jstype:     int32(fieldOptions.GetJstype()),
//			},
//		}
//
//	}
//}
//
//func (g *DocumentGenerator) getServices(fileDescriptor protoreflect.FileDescriptor) {
//
//	for i := 0; i < fileDescriptor.Services().Len(); i++ {
//
//		serviceDescriptor := fileDescriptor.Services().Get(0)
//		serviceLocation := fileDescriptor.SourceLocations().ByDescriptor(serviceDescriptor)
//		service := &registryv1alpha1.Service{
//			Name:       string(serviceDescriptor.FullName().Name()),
//			NestedName: string(serviceDescriptor.FullName().Name()),
//			FullName:             string(serviceDescriptor.FullName()),
//			Description:          serviceLocation.LeadingComments,
//			FilePath:             serviceDescriptor.ParentFile().Path(),
//			Location:             &registryv1alpha1.Location{
//				StartLine:   int32(serviceLocation.StartLine),
//				StartColumn: int32(serviceLocation.StartColumn),
//				EndLine:     int32(serviceLocation.EndLine),
//				EndColumn: int32(serviceLocation.EndColumn),
//			},
//			Methods:              nil,
//			ServiceOptions:       &registryv1alpha1.ServiceOptions{
//				Deprecated: protoutil.ProtoFromServiceDescriptor(serviceDescriptor).Options.GetDeprecated(),
//			},
//			ImplicitlyDeprecated: ,
//		}
//
//
//		for i := 0; i < serviceDescriptor.Methods().Len(); i++ {
//			methodDescriptor := serviceDescriptor.Methods().Get(i)
//			methodLocation := fileDescriptor.SourceLocations().ByDescriptor(methodDescriptor)
//			inputMessageDescriptor := methodDescriptor.Input()
//			method := registryv1alpha1.Method{
//				Name: string(methodDescriptor.Name()),
//				Description:          methodLocation.LeadingComments,
//				Request:              &registryv1alpha1.MethodRequestResponse{
//					NestedType: "",
//					FullType:        string(methodDescriptor.FullName()),
//					Streaming:       methodDescriptor.IsStreamingClient(),
//					Message:         &registryv1alpha1.Message{
//						Name:                 string(inputMessageDescriptor.Name()),
//						NestedName:           "",
//						FullName: string(inputMessageDescriptor.FullName()),
//						Description:          ,
//						FilePath:             inputMessageDescriptor.ParentFile().Path(),
//						IsMapEntry:           inputMessageDescriptor.IsMapEntry(),
//						Fields:               nil,
//						Location:             nil,
//						MessageExtensions:    nil,
//						MessageOptions:       nil,
//						ImplicitlyDeprecated: false,
//					},
//					ImportModuleRef: inputMessageDescriptor,
//				},
//				Response:             &registryv1alpha1.MethodRequestResponse{
//					NestedType:      "",
//					FullType:        "",
//					Streaming:       false,
//					Message:         nil,
//					ImportModuleRef: nil,
//				},
//				MethodOptions:        nil,
//				ImplicitlyDeprecated: false,
//			}
//		}
//	}
//}
