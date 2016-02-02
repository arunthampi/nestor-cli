package errors

import "fmt"

var UnexpectedServerError error = fmt.Errorf("Unexpected response from the Nestor API. Try again after a while or contact help@asknestor.me")
