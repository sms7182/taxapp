package notify

import "context"

type NotificationClient interface {
	FailedNotify(ctx context.Context, failedBodys []FailedBody, failedTemplate string) (*NotifyResult, error)
}
