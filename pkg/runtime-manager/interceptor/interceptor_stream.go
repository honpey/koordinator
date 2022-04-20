package interceptor

import (
	"context"
	"fmt"
	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/config"
	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/klog/v2"
	"net"
	"os"

	"github.com/koordinator-sh/koordinator/pkg/runtime-manager/dispatcher"
)

const (
	defaultRuntimeSocketPath = "/tmp/socket.sock"
)

type CriInterceptor struct {
	dispatcher    *dispatcher.RuntimeDispatcher
	director      StreamDirector
	backendConn   *grpc.ClientConn
	runtimeClient runtimeapi.RuntimeServiceClient
}

func NewCriInterceptor(dispatcher *dispatcher.RuntimeDispatcher) *CriInterceptor {
	criInterceptor := &CriInterceptor{
		dispatcher: dispatcher,
	}
	return criInterceptor
}

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

type StreamDirector func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error)

func (ci *CriInterceptor) generateInterceptServer(opts ...grpc.ServerOption) *grpc.Server {
	opts = append(opts, func() grpc.ServerOption {
		return grpc.UnknownServiceHandler(ci.Handle)
	}())
	return grpc.NewServer(opts...)
}

func (ci *CriInterceptor) Setup() error {
	os.Remove(defaultRuntimeSocketPath)
	lis, err := net.Listen("unix", defaultRuntimeSocketPath)
	if err != nil {
		fmt.Printf("fail to create the lis %v", err)
		return err
	}
	ci.Init("/run/containerd/containerd.sock")
	grpcServer := grpc.NewServer()
	runtimeapi.RegisterRuntimeServiceServer(grpcServer, ci)
	fmt.Printf("pre to run the HHH")
	err = grpcServer.Serve(lis)
	fmt.Printf("fail to create the client %v", err)
	return nil

}

func (ci *CriInterceptor) Setup3() error {
	os.Remove(defaultRuntimeSocketPath)
	lis, err := net.Listen("unix", defaultRuntimeSocketPath)

	if err != nil {
		fmt.Printf("fail to create the lis %v", err)
		return err
	}

	// grpc.ClientConn
	// 这里完全是把containerd当成了后端的代理
	backendConn, err := grpc.Dial("unix:///run/containerd/containerd.sock",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		/*grpc.WithDefaultCallOptions(grpc.ForceCodec(encoding.{}))*/)
	ci.backendConn = backendConn

	proxySvc := proxy.NewProxy(backendConn)

	klog.Infof("enter the sock %v", defaultRuntimeSocketPath)

	go func() {
		if err := proxySvc.Serve(lis); err != nil {
			klog.Infof("fail to server backend.sock")
		}

		os.Exit(1)
	}()

	return nil
}

func (ci *CriInterceptor) Setup2() error {
	os.Remove(defaultRuntimeSocketPath)
	lis, err := net.Listen("unix", defaultRuntimeSocketPath)
	if err != nil {

		fmt.Printf("fail to create the lis %v", err)
		return err
	}

	// grpc.ClientConn
	backendConn, err := grpc.Dial("unix:///run/containerd/containerd.sock", grpc.WithTransportCredentials(insecure.NewCredentials()))
	ci.backendConn = backendConn

	interceptServer := ci.generateInterceptServer()

	go func() {
		interceptServer.Serve(lis)
		klog.Infof("fail to server backend.sock")
		os.Exit(1)
	}()

	return nil
}

func (ci *CriInterceptor) Name() string {
	return "CRI"
}

func (ci *CriInterceptor) generateHookPath(method string) config.RuntimeRequestPath {
	return config.RunPodSandbox
}

func (ci *CriInterceptor) Handle(srv interface{}, proxyStream grpc.ServerStream) error {
	return ci.HandleInternal2(srv, proxyStream)
}

func (ci *CriInterceptor) HandleInternal2(srv interface{}, serverStream grpc.ServerStream) error {
	// little bit of gRPC internals never hurt anyone
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}
	md, _ := metadata.FromIncomingContext(serverStream.Context())
	outgoingCtx := metadata.NewOutgoingContext(serverStream.Context(), md.Copy())
	clientCtx, clientCancel := context.WithCancel(outgoingCtx)
	defer clientCancel()
	/*
		outgoingCtx, backendConn, err := s.director(serverStream.Context(), fullMethodName)
		if err != nil {
			return err
		}

		clientCtx, clientCancel := context.WithCancel(outgoingCtx)

		defer clientCancel()
	*/
	// TODO(mwitkow): Add a `forwarded` header to metadata, https://en.wikipedia.org/wiki/X-Forwarded-For.
	clientStream, err := grpc.NewClientStream(clientCtx, clientStreamDescForProxying, ci.backendConn, fullMethodName)
	if err != nil {
		return err
	}
	// Explicitly *do not close* s2cErrChan and c2sErrChan, otherwise the select below will not terminate.
	// Channels do not have to be closed, it is just a control flow mechanism, see
	// https://groups.google.com/forum/#!msg/golang-nuts/pZwdYRGxCIk/qpbHxRRPJdUJ

	s2cErrChan := ci.forwardClientToBackend(serverStream, clientStream)
	c2sErrChan := ci.forwardBackendToClient(clientStream, serverStream)
	// We don't know which side is going to stop sending first, so we need a select between the two.
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
				// the clientStream>serverStream may continue pumping though.
				fmt.Println("正常释放tcp连接..")
				clientStream.CloseSend()
			} else {
				// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case we need
				// to cancel the clientStream to the backend, let all of its goroutines be freed up by the CancelFunc and
				// exit with an error to the stack
				clientCancel()
				return status.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
			// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
			// will be nil.
			serverStream.SetTrailer(clientStream.Trailer())
			// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}

	}
	return status.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
}

func (ci *CriInterceptor) HandleInternal(srv interface{}, proxyStream grpc.ServerStream) error {
	//func (s *handler) handler(srv interface{}, proxyStream grpc.ServerStream) error {

	fullMethodName, ok := grpc.MethodFromServerStream(proxyStream)
	if !ok {
		return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}
	/*
		outgoingCtx, backendConn, err := ci.DefaultDirector(proxyStream.Context(), fullMethodName)
		if err != nil {
			return err
		}
		clientCtx, clientCancel := context.WithCancel(outgoingCtx)
		defer clientCancel()
	*/
	hookPath := ci.generateHookPath(fullMethodName)
	// TODO(mwitkow): Add a `forwarded` header to metadata, https://en.wikipedia.org/wiki/X-Forwarded-For.
	backendStream, err := grpc.NewClientStream(proxyStream.Context(), clientStreamDescForProxying, ci.backendConn, fullMethodName)
	if err != nil {
		return err
	}
	// =====================================
	// receive message from client(kubelet)
	// =====================================
	frame := &anypb.Any{}
	for {
		err := proxyStream.RecvMsg(frame)
		if err == io.EOF {
			break
		}
	}
	//  fullMethodName ||  frame
	// dockerd
	// frame ===> cri
	preHookType := hookPath.PreHookType()
	if preHookType != config.NoneRuntimeHookType {
		// slo-agent & runtime-manager
		ci.dispatcher.Dispatch(proxyStream.Context(), hookPath, frame)
	}
	// ===> serve
	// fram
	// =====================================
	// send message to backend(containerd)
	// =====================================
	if err := backendStream.SendMsg(frame); err == nil {
		fmt.Printf("success: 发送数据到后端结束\n")
	} else {
		fmt.Printf("fail: 发送数据到后端: %v\n", err)
	}

	// =====================================
	// receive message from backend(containerd)
	// =====================================
	result := &anypb.Any{}
	for {
		err = backendStream.RecvMsg(result)
		if err == io.EOF {
			fmt.Printf("fail to recvMsg: %v\n", err)
			break
		}
		if err != nil {
			fmt.Printf("fail to create: %v\n", err)
			break
		}
		fmt.Printf("success: get Msg success (%v)\n", result)
	}

	postHookType := hookPath.PostHookType()
	if postHookType != config.NoneRuntimeHookType {
		ci.dispatcher.Dispatch(proxyStream.Context(), hookPath, frame)
	}
	// =====================================
	// send message to client(containerd)
	// =====================================
	err = proxyStream.SendMsg(result)
	if err == nil {
		fmt.Printf("success: 发送数据到客户端\n")
	} else {
		fmt.Printf("fail: 发送数据到客户端: %v\n", err)
	}
	return nil
}

func (_ *CriInterceptor) forwardBackendToClient(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &anypb.Any{}
		for i := 0; ; i++ {
			fmt.Printf("forward client to client %v\n", i)
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if i == 0 {
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (_ *CriInterceptor) forwardClientToBackend(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &anypb.Any{}
		for i := 0; ; i++ {
			fmt.Printf("forward server to client %v\n", i)
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}
