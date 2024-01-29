package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var statusCodes = map[string]int{
	"100": http.StatusContinue,
	"101": http.StatusSwitchingProtocols,
	"102": http.StatusProcessing,
	"200": http.StatusOK,
	"201": http.StatusCreated,
	"202": http.StatusAccepted,
	"203": http.StatusNonAuthoritativeInfo,
	"204": http.StatusNoContent,
	"205": http.StatusResetContent,
	"206": http.StatusPartialContent,
	"300": http.StatusMultipleChoices,
	"301": http.StatusMovedPermanently,
	"302": http.StatusFound,
	"303": http.StatusSeeOther,
	"304": http.StatusNotModified,
	"305": http.StatusUseProxy,
	"307": http.StatusTemporaryRedirect,
	"308": http.StatusPermanentRedirect,
	"400": http.StatusBadRequest,
	"401": http.StatusUnauthorized,
	"402": http.StatusPaymentRequired,
	"403": http.StatusForbidden,
	"404": http.StatusNotFound,
	"405": http.StatusMethodNotAllowed,
	"406": http.StatusNotAcceptable,
	"407": http.StatusProxyAuthRequired,
	"408": http.StatusRequestTimeout,
	"409": http.StatusConflict,
	"410": http.StatusGone,
	"411": http.StatusLengthRequired,
	"412": http.StatusPreconditionFailed,
	"413": http.StatusRequestEntityTooLarge,
	"414": http.StatusRequestURITooLong,
	"415": http.StatusUnsupportedMediaType,
	"416": http.StatusRequestedRangeNotSatisfiable,
	"417": http.StatusExpectationFailed,
	"418": http.StatusTeapot,
	"422": http.StatusUnprocessableEntity,
	"423": http.StatusLocked,
	"424": http.StatusFailedDependency,
	"426": http.StatusUpgradeRequired,
	"428": http.StatusPreconditionRequired,
	"429": http.StatusTooManyRequests,
	"431": http.StatusRequestHeaderFieldsTooLarge,
	"451": http.StatusUnavailableForLegalReasons,
	"500": http.StatusInternalServerError,
	"501": http.StatusNotImplemented,
	"502": http.StatusBadGateway,
	"503": http.StatusServiceUnavailable,
	"504": http.StatusGatewayTimeout,
	"505": http.StatusHTTPVersionNotSupported,
	"506": http.StatusVariantAlsoNegotiates,
	"507": http.StatusInsufficientStorage,
	"508": http.StatusLoopDetected,
	"510": http.StatusNotExtended,
	"511": http.StatusNetworkAuthenticationRequired,
}

func statusCodeExists(key string, m map[string]int) bool {
	_, exists := m[key]
	return exists
}

func AppError(message string, statusCode int, w http.ResponseWriter) {
	response := make(map[string]interface{})
	response["message"] = message

	statusCodeStr := strconv.Itoa(statusCode)
	if len(statusCodeStr) > 0 && statusCodeStr[0] == '5' {
		response["status"] = "fail"
	} else {
		response["status"] = "error"
	}

	exists := statusCodeExists(statusCodeStr, statusCodes)

	if !exists {
		fmt.Println("Key does not exist")
		//    err:= errors.New(fmt.Sprintf("Invalid status code: %s", statusCodeStr))
		//    log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodes[statusCodeStr])

	json.NewEncoder(w).Encode(response)
}
