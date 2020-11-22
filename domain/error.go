package domain

import "errors"

var (
    // ErrInternalServerError will throw if any the Internal Server Error happen
    ErrInternalServerError = errors.New("Internal Server Error")
    // ErrNotFound will throw if the request item is not exists
    ErrNotFound = errors.New("Your requested Item is not found")
    // ErrConflict will throw if the current action already exists
    ErrConflict = errors.New("Your Item already exist")
    // ErrBadParamInput will throw if the given request-body or params is not valid
    ErrBadParamInput = errors.New("Given Param is not valid")
    // InvalidPassword will throw if the given requet is invalid password
    InvalidPassword = errors.New("password invalid")
    // InvalidUser will throw if the given user not found in database
    InvalidUser = errors.New("invalid user")
)