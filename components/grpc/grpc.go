package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"

	trackerApi "github.com/somatech1/mikros/apis/tracker"
	mcontext "github.com/somatech1/mikros/components/context"
	"github.com/somatech1/mikros/components/service"
	merrors "github.com/somatech1/mikros/internal/components/errors"
)

// ClientConnectionOptions gathers custom options to establish a connection with
// a gRPC client.
type ClientConnectionOptions struct {
	ServiceName           service.Name
	ClientName            service.Name
	Context               *mcontext.ServiceContext
	Connection            ConnectionOptions
	AlternativeConnection *ConnectionOptions
	Tracker               trackerApi.Tracker
}

type ConnectionOptions struct {
	Host      string
	Namespace string
	Port      int32
}

// ClientConnection establishes a connection with a gRPC service and returns its
// connection.
//
// This method provides a mechanism to a service to connect into several other
// gRPC services to access their APIs.
func ClientConnection(options *ClientConnectionOptions) (*grpc.ClientConn, error) {
	address := getClientConnectionAddress(options)

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(gRPCClientUnaryInterceptor(options.Context, options.Tracker, options.ServiceName, options.ClientName)),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getClientConnectionAddress(options *ClientConnectionOptions) string {
	getAddress := func(prefix string, c *ConnectionOptions) string {
		if c.Host != "" {
			return fmt.Sprintf("%s:%d", c.Host, c.Port)
		}

		return fmt.Sprintf("%s.%v:%d", prefix, c.Namespace, c.Port)
	}

	addr := getAddress(options.ClientName.String(), &options.Connection)
	if options.AlternativeConnection != nil {
		addr = getAddress(options.ClientName.String(), options.AlternativeConnection)
	}

	return addr
}

func gRPCClientUnaryInterceptor(svcCtx *mcontext.ServiceContext, tracker trackerApi.Tracker, from, to service.Name) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Add gzip compression to the request.
		opts = append(opts, grpc.UseCompressor(gzip.Name))

		if tracker != nil {
			trackId := tracker.Generate()

			// If we already have a tracker ID, we need to use for subsequent calls.
			if trk, ok := tracker.Retrieve(ctx); ok {
				trackId = trk
			}

			// Adds the track ID on the context.
			ctx = tracker.Add(ctx, trackId)
		}

		// Calls invoker with a new context.
		if err := invoker(mcontext.AppendServiceContext(ctx, svcCtx), method, req, reply, cc, opts...); err != nil {
			// Return the proper inner service error for the caller.
			if st, ok := status.FromError(err); ok {
				return merrors.FromGRPCStatus(st, from, to)
			}

			// Should not fall here
			return err
		}

		return nil
	}
}
