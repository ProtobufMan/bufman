package handlers

import (
	"context"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/validity"
	"github.com/bufbuild/connect-go"
)

type ResolveServiceHandler struct {
	resolveService services.ResolveService
	validator      validity.Validator
}

func NewResolveServiceHandler() *ResolveServiceHandler {
	return &ResolveServiceHandler{
		resolveService: services.NewResolveService(),
		validator:      validity.NewValidator(),
	}
}

func (handler *ResolveServiceHandler) GetModulePins(ctx context.Context, req *connect.Request[registryv1alpha1.GetModulePinsRequest]) (*connect.Response[registryv1alpha1.GetModulePinsResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 首先检查用户权限，是否对repo有访问权限
	var checkErr e.ResponseError
	repositoryMap := map[string]*model.Repository{}
	for _, moduleReference := range req.Msg.GetModuleReferences() {
		fullName := moduleReference.GetOwner() + "/" + moduleReference.GetRepository()
		repo, ok := repositoryMap[fullName]
		if !ok {
			repo, checkErr = handler.validator.CheckRepositoryCanAccess(userID, moduleReference.GetOwner(), moduleReference.GetRepository(), registryv1alpha1connect.ResolveServiceGetModulePinsProcedure)
			if checkErr != nil {
				return nil, connect.NewError(checkErr.Code(), checkErr.Err())
			}
			repositoryMap[fullName] = repo
		}
		repositoryMap[fullName] = repo
	}

	// 获取对应的commits
	commits, err := handler.resolveService.GetModulePins(repositoryMap, req.Msg.GetModuleReferences())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	retPins := commits.ToProtoModulePins()
	// 处理CurrentModulePins
	for _, currentModulePin := range req.Msg.GetCurrentModulePins() {
		//
		ownerName := currentModulePin.GetOwner()
		repositoryName := currentModulePin.GetRepository()
		for _, commit := range commits {
			// 如果current module pin在reference的查询出的commits内，则有breaking的可能
			if commit.UserName == ownerName && commit.RepositoryName == repositoryName {
				commitName := currentModulePin.GetCommit()
				if commit.CommitName != commitName {
					// 版本号不一样，存在breaking
					return nil, e.NewInvalidArgumentError(fmt.Sprintf("%s/%s (possible to cause breaking)", currentModulePin.GetOwner(), currentModulePin.GetRepository()))
				}
			}
		}

		// 当前pin没有breaking的可能性，加入到返回结果中
		retPins = append(retPins, currentModulePin)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetModulePinsResponse{
		ModulePins: retPins,
	})
	return resp, nil
}

func (handler *ResolveServiceHandler) GetGoVersion(ctx context.Context, req *connect.Request[registryv1alpha1.GetGoVersionRequest]) (*connect.Response[registryv1alpha1.GetGoVersionResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *ResolveServiceHandler) GetSwiftVersion(ctx context.Context, req *connect.Request[registryv1alpha1.GetSwiftVersionRequest]) (*connect.Response[registryv1alpha1.GetSwiftVersionResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *ResolveServiceHandler) GetMavenVersion(ctx context.Context, req *connect.Request[registryv1alpha1.GetMavenVersionRequest]) (*connect.Response[registryv1alpha1.GetMavenVersionResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *ResolveServiceHandler) GetNPMVersion(ctx context.Context, req *connect.Request[registryv1alpha1.GetNPMVersionRequest]) (*connect.Response[registryv1alpha1.GetNPMVersionResponse], error) {
	//TODO implement me
	panic("implement me")
}
