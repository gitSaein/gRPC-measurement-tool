package error

import (
	"log"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type ErrorStatus struct {
	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// A developer-facing error message, which should be in English. Any
	// user-facing error message should be localized and sent in the
	// [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// A list of messages that carry the error details.  There is a common set of
	// message types for APIs to use.
	Details []*any.Any `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
}

func CheckErrorState(err error, pid uint64) {
	if err != nil {

		st, _ := status.FromError(err)
		log.Printf("Error: [client-pid: %d] %v", pid, st)

		errorStatus := &ErrorStatus{}
		errorStatus.Code = st.Proto().Code
		errorStatus.Message = st.Proto().Message
		errorStatus.Details = st.Proto().Details

		errDetails := st.Details()

		for _, detail := range errDetails {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				log.Println("Oops! Your request was rejected by the server.")
				for _, violation := range t.GetFieldViolations() {
					log.Printf("The %q field was wrong:\n", violation.GetField())
					log.Printf("\t%s\n", violation.GetDescription())
				}
			default:
				log.Printf("did not connect: %v", st)

			}
		}
		log.Fatal(errorStatus)
	}
}
