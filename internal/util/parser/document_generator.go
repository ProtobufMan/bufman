package parser

import (
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule/bufmoduleprotocompile"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/bufbuild/protocompile/linker"
	"github.com/bufbuild/protocompile/protoutil"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type DocumentGenerator interface {
	GenerateDocument(packageName string) *registryv1alpha1.PackageDocumentation
}

type documentGeneratorImpl struct {
	commitName            string // 当前的commitName，用于判断是否是外部依赖
	packageLinkerMap      map[string]linker.Files
	linkers               linker.Files
	parserAccessorHandler bufmoduleprotocompile.ParserAccessorHandler
	messageSet            map[string]*registryv1alpha1.Message
}

func NewDocumentGenerator(commitName string, links linker.Files, parserAccessorHandler bufmoduleprotocompile.ParserAccessorHandler) DocumentGenerator {
	return &documentGeneratorImpl{
		commitName:            commitName,
		linkers:               links,
		parserAccessorHandler: parserAccessorHandler,
		packageLinkerMap:      map[string]linker.Files{},
		messageSet:            map[string]*registryv1alpha1.Message{},
	}
}

// isDependency 判断是否是外部依赖
func (g *documentGeneratorImpl) isDependency(fileDescriptor protoreflect.FileDescriptor) bool {
	return g.parserAccessorHandler.Commit(fileDescriptor.Path()) != g.commitName
}

func (g *documentGeneratorImpl) toProtoLocation(loc protoreflect.SourceLocation) *registryv1alpha1.Location {
	return &registryv1alpha1.Location{
		StartLine:   int32(loc.StartLine),
		StartColumn: int32(loc.StartColumn),
		EndLine:     int32(loc.EndLine),
		EndColumn:   int32(loc.EndColumn),
	}
}

func (g *documentGeneratorImpl) getPackageLinkers(packageName string) linker.Files {
	if packageLinkers, ok := g.packageLinkerMap[packageName]; ok {
		return packageLinkers
	}

	packageLinkers := make(linker.Files, 0, len(g.linkers))
	for i := 0; i < len(g.linkers); i++ {
		if string(g.linkers[i].Package()) == packageName {
			packageLinkers = append(packageLinkers, g.linkers[i])
		}
	}
	g.packageLinkerMap[packageName] = packageLinkers
	return packageLinkers
}

func (g *documentGeneratorImpl) GenerateDocument(packageName string) *registryv1alpha1.PackageDocumentation {

	packageDocument := &registryv1alpha1.PackageDocumentation{
		Name:     packageName,
		Services: g.GetPackageServices(packageName),
		Enums:    g.GetPackageEnums(packageName),
		Messages: g.GetPackageMessages(packageName),
		// TODO FileExtensions: nil,
	}

	return packageDocument
}

func (g *documentGeneratorImpl) getNestedName(name string, nestedPrefixes ...string) string {
	parts := append(nestedPrefixes, name)
	return strings.Join(parts, ".")
}

func (g *documentGeneratorImpl) GetPackageMessages(packageName string) []*registryv1alpha1.Message {
	packageLinkers := g.getPackageLinkers(packageName)
	var messages []*registryv1alpha1.Message

	for i := 0; i < len(packageLinkers); i++ {
		packageLink := packageLinkers[i]
		messages = append(messages, g.GetMessages(packageLink.Messages())...)
	}

	return messages
}

func (g *documentGeneratorImpl) GetMessages(messageDescriptors protoreflect.MessageDescriptors, nestedPrefixes ...string) []*registryv1alpha1.Message {
	messages := make([]*registryv1alpha1.Message, 0, messageDescriptors.Len())

	for i := 0; i < messageDescriptors.Len(); i++ {
		messageDescriptor := messageDescriptors.Get(i)
		message := g.GetMessage(messageDescriptor, nestedPrefixes...)
		messages = append(messages, message)
	}

	return messages
}

func (g *documentGeneratorImpl) GetMessage(messageDescriptor protoreflect.MessageDescriptor, nestedPrefixes ...string) *registryv1alpha1.Message {
	if message, ok := g.messageSet[string(messageDescriptor.FullName())]; ok {
		return message
	}

	// get location info
	messageLocation := messageDescriptor.ParentFile().SourceLocations().ByDescriptor(messageDescriptor)
	// get options
	messageOptions := protoutil.ProtoFromMessageDescriptor(messageDescriptor).Options

	message := &registryv1alpha1.Message{
		Name:        string(messageDescriptor.Name()),
		NestedName:  g.getNestedName(string(messageDescriptor.Name()), nestedPrefixes...),
		FullName:    string(messageDescriptor.FullName()),
		Description: messageLocation.LeadingComments,
		FilePath:    messageDescriptor.ParentFile().Path(),
		IsMapEntry:  messageDescriptor.IsMapEntry(),
		Location:    g.toProtoLocation(messageLocation),
		MessageOptions: &registryv1alpha1.MessageOptions{
			Deprecated: messageOptions.GetDeprecated(),
		},
		// TODO ImplicitlyDeprecated:
	}

	// TODO fill message fields

	// TODO fill message extensions

	// TODO handle message nested message

	// TODO handle message nested enum

	// 记录message
	g.messageSet[string(messageDescriptor.FullName())] = message

	return message
}

func (g *documentGeneratorImpl) GetPackageEnums(packageName string) []*registryv1alpha1.Enum {
	packageLinkers := g.getPackageLinkers(packageName)
	var enums []*registryv1alpha1.Enum

	for i := 0; i < len(packageLinkers); i++ {
		packageLink := packageLinkers[i]
		enums = append(enums, g.GetEnums(packageLink.Enums())...)
	}

	return enums
}

func (g *documentGeneratorImpl) GetEnums(enumDescriptors protoreflect.EnumDescriptors, nestedPrefixes ...string) []*registryv1alpha1.Enum {
	enums := make([]*registryv1alpha1.Enum, 0, enumDescriptors.Len())

	for i := 0; i < enumDescriptors.Len(); i++ {
		enumDescriptor := enumDescriptors.Get(i)
		enum := g.GetEnum(enumDescriptor, nestedPrefixes...)
		enums = append(enums, enum)
	}

	return enums
}

func (g *documentGeneratorImpl) GetEnum(enumDescriptor protoreflect.EnumDescriptor, nestedPrefixes ...string) *registryv1alpha1.Enum {
	// get location info
	enumLocation := enumDescriptor.ParentFile().SourceLocations().ByDescriptor(enumDescriptor)
	// get options
	enumOptions := protoutil.ProtoFromEnumDescriptor(enumDescriptor).GetOptions()

	enum := &registryv1alpha1.Enum{
		Name:        string(enumDescriptor.Name()),
		NestedName:  g.getNestedName(string(enumDescriptor.Name()), nestedPrefixes...),
		FullName:    string(enumDescriptor.FullName()),
		Description: enumLocation.LeadingComments,
		FilePath:    enumDescriptor.ParentFile().Path(),
		Location:    g.toProtoLocation(enumLocation),
		EnumOptions: &registryv1alpha1.EnumOptions{
			Deprecated: enumOptions.GetDeprecated(),
			AllowAlias: enumOptions.GetAllowAlias(),
		},
		// TODO ImplicitlyDeprecated: false,
	}

	// enum values
	enumValues := make([]*registryv1alpha1.EnumValue, 0, enumDescriptor.Values().Len())
	for i := 0; i < enumDescriptor.Values().Len(); i++ {
		enumValueDescriptor := enumDescriptor.Values().Get(i)
		enumValueLocation := enumValueDescriptor.ParentFile().SourceLocations().ByDescriptor(enumValueDescriptor)
		enumValueOptions := protoutil.ProtoFromEnumValueDescriptor(enumValueDescriptor).GetOptions()

		enumValue := &registryv1alpha1.EnumValue{
			Name:        string(enumValueDescriptor.Name()),
			Number:      int32(enumValueDescriptor.Number()),
			Description: enumValueLocation.LeadingComments,
			EnumValueOptions: &registryv1alpha1.EnumValueOptions{
				Deprecated: enumValueOptions.GetDeprecated(),
			},
		}

		enumValues = append(enumValues, enumValue)
	}
	enum.Values = enumValues

	return enum
}

func (g *documentGeneratorImpl) GetPackageServices(packageName string) []*registryv1alpha1.Service {
	packageLinkers := g.getPackageLinkers(packageName)
	var services []*registryv1alpha1.Service

	for i := 0; i < len(packageLinkers); i++ {
		packageLink := packageLinkers[i]
		services = append(services, g.GetServices(packageLink.Services())...)
	}

	return services
}

func (g *documentGeneratorImpl) GetServices(serviceDescriptors protoreflect.ServiceDescriptors, nestedPrefixes ...string) []*registryv1alpha1.Service {
	services := make([]*registryv1alpha1.Service, 0, serviceDescriptors.Len())

	for i := 0; i < serviceDescriptors.Len(); i++ {
		serviceDescriptor := serviceDescriptors.Get(i)
		service := g.GetService(serviceDescriptor, nestedPrefixes...)
		services = append(services, service)
	}

	return services
}

func (g *documentGeneratorImpl) GetService(serviceDescriptor protoreflect.ServiceDescriptor, nestedPrefixes ...string) *registryv1alpha1.Service {
	// get location info
	serviceLocation := serviceDescriptor.ParentFile().SourceLocations().ByDescriptor(serviceDescriptor)
	// get options
	serviceOptions := protoutil.ProtoFromServiceDescriptor(serviceDescriptor).GetOptions()

	service := &registryv1alpha1.Service{
		Name:        string(serviceDescriptor.Name()),
		NestedName:  g.getNestedName(string(serviceDescriptor.Name()), nestedPrefixes...),
		FullName:    string(serviceDescriptor.FullName()),
		Description: serviceLocation.LeadingComments,
		FilePath:    serviceDescriptor.ParentFile().Path(),
		Location:    g.toProtoLocation(serviceLocation),
		ServiceOptions: &registryv1alpha1.ServiceOptions{
			Deprecated: serviceOptions.GetDeprecated(),
		},
		// TODO ImplicitlyDeprecated: false,
	}

	// methods
	// enum values
	methods := make([]*registryv1alpha1.Method, 0, serviceDescriptor.Methods().Len())
	for i := 0; i < serviceDescriptor.Methods().Len(); i++ {
		methodDescriptor := serviceDescriptor.Methods().Get(i)
		methodLocation := methodDescriptor.ParentFile().SourceLocations().ByDescriptor(methodDescriptor)
		methodOptions := protoutil.ProtoFromMethodDescriptor(methodDescriptor).GetOptions()
		requestDescriptor := methodDescriptor.Input()
		responseDescriptor := methodDescriptor.Output()

		method := &registryv1alpha1.Method{
			Name:        string(methodDescriptor.Name()),
			Description: methodLocation.LeadingComments,
			Request:     g.getMethodRequestResponse(methodDescriptor.IsStreamingClient(), requestDescriptor),
			Response:    g.getMethodRequestResponse(methodDescriptor.IsStreamingServer(), responseDescriptor),
			MethodOptions: &registryv1alpha1.MethodOptions{
				Deprecated:       methodOptions.GetDeprecated(),
				IdempotencyLevel: int32(methodOptions.GetIdempotencyLevel()),
			},
			// TODO ImplicitlyDeprecated: false,
		}

		methods = append(methods, method)
	}
	service.Methods = methods

	return service
}

func (g *documentGeneratorImpl) getMethodRequestResponse(streaming bool, descriptor protoreflect.MessageDescriptor, nestedPrefixes ...string) *registryv1alpha1.MethodRequestResponse {
	m := g.GetMessage(descriptor, nestedPrefixes...)

	r := &registryv1alpha1.MethodRequestResponse{
		NestedType: m.NestedName,
		FullType:   m.FullName,
		Streaming:  streaming,
		Message:    m,
	}
	if g.isDependency(descriptor.ParentFile()) {
		// 如果是外部依赖
		fileDescriptor := descriptor.ParentFile()
		identity := g.parserAccessorHandler.ModuleIdentity(fileDescriptor.Path())
		importModuleRef := &registryv1alpha1.ImportModuleRef{
			Remote:      identity.Remote(),
			Owner:       identity.Owner(),
			Repository:  identity.Repository(),
			Commit:      g.parserAccessorHandler.Commit(fileDescriptor.Path()),
			PackageName: string(fileDescriptor.Package()),
		}
		r.ImportModuleRef = importModuleRef
	}

	return r
}
