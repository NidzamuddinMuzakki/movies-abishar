package response

import (
	"context"
	"fmt"
)

type Err struct {
	Response
}

func (r *Err) Error() string {
	return fmt.Sprintf("%d %s", r.Code, r.Message)
}

func ParseErrorToHTTPCode(ctx context.Context, err error) (int, interface{}) {
	if resp, ok := err.(*Err); ok {
		return resp.ToHTTPCodeAndMap()
	}

	return InternalServerError(ctx).ToHTTPCodeAndMap()
}
