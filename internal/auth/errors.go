package auth

import "errors"

var AuthErrorsEmailOrPasswordAreEmpty = errors.New("Email or password are empty. ")
var AuthErrorsEmailOrOtpAreEmpty = errors.New("Email or OTP are empty. ")
var AuthErrorsEmailNotValid = errors.New("Email is not valid")
var AuthErrorsUserExisitsWithThisEmail = errors.New("User with this emails exists.")
var AuthErrorsPasswordCharLimit = errors.New("Password is less than 6 Chars")
var AuthErrorsUsernotFound = errors.New("User not found with this email.")
var AuthErrorsUserPasswordNotMatch = errors.New("Password not maching ")
