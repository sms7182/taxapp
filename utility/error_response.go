package utility

import "fmt"

type TaxErrorResponseManager struct {
	errors map[int]ErrorDetail
}
type ErrorDetail struct {
	key    string
	detail string
}

func (er TaxErrorResponseManager) populateErrors() {
	er.errors = make(map[int]ErrorDetail)
	er.errors[5000] = ErrorDetail{
		key:    "internal.server.error",
		detail: "خطای داخلی سرور",
	}
	er.errors[5001] = ErrorDetail{
		key:    "bad.request",
		detail: "زمانی که بسته مشکل داشته و در سمت سرور خطایی به صورت مشخص وجود نداشته باشد.",
	}
	er.errors[5002] = ErrorDetail{
		key:    "un.authorized",
		detail: "عدم دارا بودن دسترسی برای ارسال این درخواست",
	}

	er.errors[5003] = ErrorDetail{
		key:    "uid.format.is.not.valid",
		detail: "ارسال با فرمت اشتباه برای uid در packet",
	}

	er.errors[5004] = ErrorDetail{
		key:    "invalid.json.structure",
		detail: "ساختار جیسون درخواست اشتباه است.",
	}

	er.errors[5005] = ErrorDetail{
		key:    "duplicate.request.uid",
		detail: "ارسال درخواست با ایدی تکراری در بسته",
	}

	er.errors[5006] = ErrorDetail{
		key:    "packet.size.is.too.large",
		detail: "ارسال تعداد بسته ها بیش از حد مجاز است.",
	}

	er.errors[5007] = ErrorDetail{
		key:    "not.supported.packet-type",
		detail: "عدم پشتیبانی از نوع بسته ارسالی",
	}

	er.errors[5008] = ErrorDetail{
		key:    "encryption.key.id.not.valid",
		detail: "مقدار فیلد در ارسال رمز شده صورتحساب صحیح نمی باشد.encryptionKeyId",
	}

	er.errors[5009] = ErrorDetail{
		key:    "not.match.packet-type.with.request",
		detail: "عدم تطابق بسته ارسالی با نوع درخواست sync",
	}

	er.errors[5010] = ErrorDetail{
		key:    "request.time.has.passed",
		detail: "گذشت زمان مشخصی از زمان ارسال درخواست و دریافت ان توسط سرور",
	}

	er.errors[5011] = ErrorDetail{
		key:    "duplicate.request.trace.id",
		detail: "تکراری بودن در سرایند درخواست requestTraceId",
	}

	er.errors[5012] = ErrorDetail{
		key:    "fiscal.id.not.found",
		detail: "پیدا نشدن اطلاعات حافظه مالیاتی ارسالی",
	}

	er.errors[5013] = ErrorDetail{
		key:    "invalid.packet.signature",
		detail: "عدم اعتبار امضای درخواست",
	}

	er.errors[5014] = ErrorDetail{
		key:    "invalid.format",
		detail: "فرمت ساختار نادرست است.",
	}

	er.errors[5015] = ErrorDetail{
		key:    "invalid.token",
		detail: "توکن نامعتبر است",
	}

}
func (er TaxErrorResponseManager) GetDetail(code int) ErrorDetail {
	if len(er.errors) == 0 {
		er.populateErrors()
	}
	if val, ok := er.errors[code]; ok {
		return val
	}
	detail := fmt.Sprintf("Unknown error with %v code", code)
	err := ErrorDetail{
		key:    "Unknown",
		detail: detail,
	}
	er.errors[code] = err
	return err

}
