package studio

import "errors"

var StudioErrorsHandleInvalid = errors.New("handle is invalid")
var StudioErrorsNameInvalid = errors.New("name is invalid")
var StudioErrorsDescriptionInvalid = errors.New("description is invalid")

var StudioErrorsHandleLenMax = errors.New("handle is too big")
var StudioErrorsNameLenMax = errors.New("name is too big")
var StudioErrorsDescriptionLenMax = errors.New("description is too big")

var ErrHandleUnavailable = errors.New("handle unavailable")

var StudioErrorsTopicLenMax = errors.New("topic name too big")
