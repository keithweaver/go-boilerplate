package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"go-boilerplate/common"
	"go-boilerplate/emails/accountlockemail"
	"go-boilerplate/emails/forgotpasswordemail"
	"go-boilerplate/emails/sessionunlockemail"
	"go-boilerplate/emails/signinemail"
	"go-boilerplate/emails/verifyemail"
	"go-boilerplate/logging"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode"

	// "github.com/google/uuid"
)

type Services struct {
	logger         logging.Logger
	userRepository Repository
	forgotPasswordRepository ForgotPasswordRepository
}

type ServiceContract interface {
	SignUp(ctx context.Context, userAgent *user_agent.UserAgent, currentIP string, body SignUpBody) (string, bool, *common.Error)
	SignIn(ctx context.Context, userAgent *user_agent.UserAgent, currentIP string, body SignInBody) (string, bool, *common.Error)
	LogOut(authToken string) error
}

func NewInstanceOfUserServices(logger logging.Logger, userRepository Repository, forgotPasswordRepository ForgotPasswordRepository) Services {
	return Services{logger, userRepository, forgotPasswordRepository}
}

// SignUp signs up the new account (or signs in the user).
func (s *Services) SignUp(ctx context.Context, userAgent *user_agent.UserAgent, currentIP string, body SignUpBody) (string, bool, *common.Error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "SignUp")

	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")

	// Verify password meets sign up requirements
	if !s.isValidPassword(ctx, body.Password) {
		return "", false, &common.Error{
			StatusCode: 400,
			Message: "error: Your password does not meet requirements.",
		}
	}

	// Check for user
	userExists, err := s.userRepository.DoesUserExist(emailTrimmed)
	if err != nil {
		return "", false, &common.Error{
			StatusCode: 500,
		}
	}
	if userExists {
		// Just try signing them in
		return s.signIn(ctx, false, userAgent, currentIP, emailTrimmed, body.Password)
	}
	encryptedPassword, err := s.getEncryptedPassword(body.Password)

	// Save the device they are signing up as a known device
	engine, engineVersion := userAgent.Engine()
	browserName, browserVersion := userAgent.Browser()
	knownDevices := []Device{{
		Name: "Sign Up Device",
		Mobile: userAgent.Mobile(),
		Bot: userAgent.Bot(),
		Mozilla: userAgent.Mozilla(),
		Platform: userAgent.Platform(),
		OperatingSystem: userAgent.OS(),
		Engine: engine,
		EngineVersion: engineVersion,
		Browser: browserName,
		BrowserVersion: browserVersion,
		ValidDevice: true,
	}}

	// Create verification code
	verificationCode := uuid.New().String()
	now := time.Now()
	verificationExpiry := now.Add(time.Hour * time.Duration(1)) // Expires in 1 hour

	// Sign up user
	newUser := User{
		Email:    emailTrimmed,
		Password: encryptedPassword,
		Name:     body.Name,
		Created:  time.Now(),
		VerifiedEmail: false,
		VerificationCode: verificationCode,
		VerificationExpiryTime: verificationExpiry,
		TrustedIPs: []IP{},
		InvalidIPs: []IP{},
		KnownDevices: knownDevices,
	}
	err = s.userRepository.SaveUser(newUser)
	if err != nil {
		s.logger.Warning(ctx, "failed to save user", err)
		return "", false, &common.Error{
			StatusCode: 500,
		}
	}

	// Send verification email
	err = verifyemail.SendVerifyEmail(newUser.Greeting(), newUser.Email, verificationCode)
	if err != nil {
		s.logger.Warning(ctx, "failed to send verify email", err)
		return "", false, &common.Error{
			StatusCode: 500,
		}
	}

	// Sign in user
	return s.signIn(ctx, true, userAgent, currentIP, emailTrimmed, body.Password)
}

func (s *Services) isValidPassword(ctx context.Context, password string) bool {
	// TODO - Change to return an error instead with the problem
	if len(password) < 8 {
		s.logger.Warning(ctx, "password is less than 8 characters", errors.New("error: invalid password"))
		// Length is less than 8
		return false
	}
	specialChar := 0
	for _, char := range password {
		if !s.isLetter(string(char)) {
			specialChar++
		}

		if specialChar >= 5 {
			break
		}
	}

	if specialChar < 5 {
		s.logger.Warning(ctx, "password is has less than 5 special characters", errors.New("error: invalid password"))
		// Requires at least 5 non-letters
		return false
	}

	return true
}

func (s *Services) isLetter(password string) bool {
	for _, c := range password {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}

func (s *Services) getEncryptedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func (s *Services) SignIn(ctx context.Context, userAgent *user_agent.UserAgent, currentIP string, body SignInBody) (string, bool, *common.Error) {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "SignIn")

	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")
	return s.signIn(ctx, false, userAgent, currentIP, emailTrimmed, body.Password)
}

func (s *Services) signIn(ctx context.Context, isSignUp bool, userAgent *user_agent.UserAgent, currentIP string, email string, password string) (string, bool, *common.Error) {
	ctx = context.WithValue(ctx, logging.CtxHelpMethods, logging.AddToHelperMethods(ctx, "signIn"))

	// Grab user
	found, user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		s.logger.Warning(ctx, "failed to get user by email", err)
		return "", false, &common.Error{
			StatusCode: 500,
		}
	}

	if !found {
		s.logger.Warning(ctx, "failed to find user", errors.New("error: unauthorized"))
		return "", false, &common.Error{
			StatusCode: 403,
		}
	}

	if !s.isUsersPassword(user.Password, password) {
		s.logger.Warning(ctx, "invalid password", errors.New("error: unauthorized"))
		return "", false, &common.Error{
			StatusCode: 403,
		}
	}

	// They now have a valid signed in
	if user.AccountLocked {
		return "", false, &common.Error{
			StatusCode: 403,
			Message: "error: Account has been locked. Please reset password.",
		}
	}

	// Check if trusted ip, level of legitamacy of the sign in
	lockSession, invalidSession, err := s.validateSignIn(ctx, isSignUp, user, userAgent, currentIP)
	if err != nil {
		return "", false, &common.Error{
			StatusCode: 500,
		}
	}

	if invalidSession {
		err := s.lockUserAccount(user)
		if err != nil {
			// Ignore the failure but worth notifying your dev team for
			s.logger.Error(ctx, "failed to lock the user account", err)
		}
		return "", false, &common.Error{
			StatusCode: 403,
		}
	}

	// Create session
	now := time.Now()
	expiryDate := now.AddDate(0, 0, 1)
	newSession := Session{
		Email:   email,
		Created: now,
		Expiry:  expiryDate,
		Locked: lockSession,
		UnlockCode: uuid.New().String(),
	}

	// Save the session
	token, err := s.userRepository.SaveSession(newSession)
	if err != nil {
		s.logger.Warning(ctx, "failed to save session", err)
		return "", false, &common.Error{
			StatusCode: 500,
		}
	}

	if lockSession {
		// Session has been locked. Send the user an email with a code to unlock it.
		s.logger.Info(ctx, "Session has been marked as locked, sending an email with the unlock code")
		err = sessionunlockemail.SendSessionUnLockEmail(user.Greeting(), user.Email, newSession.UnlockCode)
		if err != nil {
			s.logger.Warning(ctx, "failed to send session unlock email", err)
			return "", false, &common.Error{
				StatusCode: 500,
			}
		}
	} else {
		// Send sign in email (General sign in and not locked accounts)
		err = signinemail.SendSignInEmail("there", email, currentIP, "<browser>", "<os>") // TODO - Add in the actual values
		if err != nil {
			// Ignore the failure. This is my decision since it doesnt stop the user from signing
			// into their account. However, sending a sign in email is another layer of security.
			// You can return a 500 here if you want to make sure the email goes through.
			s.logger.Warning(ctx, "failed to send sign in email", err)
		}
		s.logger.Info(ctx, "Sent a sign in email")

		// Session is both valid and not locked. This is a trusted session. We should add the IP
		// to the list of trust IPs. At a minimum, the session they signed up with will have its
		// IP address stored as trusted.
		s.logger.Info(ctx, "Session is valid and trusted, adding to list of trusted IPs")
		err = s.userRepository.UpdateOrAddTrustedIPToUser(user.Email, IP{Address: currentIP, LocationFound: false})
		if err != nil {
			s.logger.Warning(ctx, "failed to save new trusted IP", err)
		}
	}

	return token, lockSession, nil
}

func (s *Services) isUsersPassword(storedPasswordHash string, plainTextInputtedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(plainTextInputtedPassword)) == nil
}

func (s *Services) validateSignIn(ctx context.Context, isSignUp bool, user User, userAgent *user_agent.UserAgent, currentIP string) (bool, bool, error) {
	ctx = context.WithValue(ctx, logging.CtxHelpMethods, logging.AddToHelperMethods(ctx, "validateSignIn"))

	if isSignUp {
		// User has just signed up so nothing to compare against.
		s.logger.Info(ctx, "User has just signed up, account is unlock and valid session")
		return false, false, nil
	}

	SESSION_LOCKED_LIMIT := 2
	warnings := 0
	// Check if IP is in trusted IPs
	if user.HasTrustedIP(currentIP) {
		// Log valid IP
		s.logger.Info(ctx, "User has a trusted IP")
	} else if user.HasInvalidIPs(currentIP) {
		// Current sign in is using an invalid IP
		warnings += 1
		s.logger.Info(ctx, "User has an invalid IP")
	} else {
		// Lock the session
		warnings += SESSION_LOCKED_LIMIT // Brand new IP
		s.logger.Info(ctx, "User has brand new IP")
	}

	// Check if bot
	if userAgent.Bot() {
		// Lock the session
		warnings += SESSION_LOCKED_LIMIT
		s.logger.Info(ctx, "User is a bot based on User Agent Header")
	}

	// Gather information about past devices compared to this one
	browserName, browserVersion := userAgent.Browser()
	s.logger.Info(ctx, fmt.Sprintf("User is currently on browser : %s %s", browserName, browserVersion))
	s.logger.Info(ctx, fmt.Sprintf("User is currently on OS : %s", userAgent.OS()))

	browserTypeFound := false
	browserVersionFound := false
	osFound := false
	for _, device := range user.KnownDevices {
		if device.ValidDevice && device.Browser == browserName && device.BrowserVersion == browserVersion {
			browserTypeFound = true
			browserVersionFound = true
		} else if device.ValidDevice &&  device.Browser == browserName {
			browserTypeFound = true
		}
		if device.ValidDevice && userAgent.OS() == device.OperatingSystem {
			osFound = true
		}
	}

	// Check if same browser
	if !browserTypeFound {
		// Lock the session
		warnings += SESSION_LOCKED_LIMIT
		s.logger.Info(ctx, "User has changed browsers as a previous session on User Agent Header")
	} else if browserTypeFound && !browserVersionFound {
		warnings += 1
		s.logger.Info(ctx, "User is on same browser (not version) as a previous session on User Agent Header")
	} else {
		s.logger.Info(ctx, "User is on same browser and version as a previous session on User Agent Header")
	}

	// Check if OS changed
	if !osFound {
		// Lock the session
		warnings += SESSION_LOCKED_LIMIT
		s.logger.Info(ctx, "User has changed OS's as a previous session on User Agent Header")
	} else {
		s.logger.Info(ctx, "User is on OS as a previous session on User Agent Header")
	}

	// You can add more details here such as location (Ex. country changed etc.)

	if warnings >= SESSION_LOCKED_LIMIT {
		// Lock the session and ask for secondary validation
		// One warning like old trusted IP is fine.
		s.logger.Info(ctx, "Too many warnings, locking session")
		return true, false, nil
	}
	// Do nothing
	return false, false, nil
}

func (s *Services) lockUserAccount(user User) error {
	err := accountlockemail.SendAccountLockedEmail("there", user.Email)
	if err != nil {
		return err
	}
	return s.userRepository.UpdateAccountLocked(user.Email, true)
}

func (s *Services) LogOut(ctx context.Context, authToken string) *common.Error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "LogOut")

	// Grab session
	found, _, err := s.userRepository.GetSessionById(authToken)
	if err != nil {
		s.logger.Warning(ctx, "failed to get session", err)
		return &common.Error{
			StatusCode: 500,
		}
	}
	if !found {
		s.logger.Warning(ctx, "session not found", errors.New("error: not found"))
		return &common.Error{
			StatusCode: 403,
		}
	}

	// Mark as expired
	err = s.userRepository.MarkSessionAsExpired(authToken)
	if err != nil {
		s.logger.Warning(ctx, "failed to mark session as expired", err)
		return &common.Error{
			StatusCode: 500,
		}
	}

	return nil
}

// UnlockSession  removes the lock flag from the session to allow to continue to make requests.
func (s *Services) UnlockSession(ctx context.Context, currentIP string, authToken string, body UnlockSessionBody) *common.Error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "UnlockSession")
	found, session, err := s.userRepository.GetSessionById(authToken)
	if err != nil {
		// TODO - Test not found
		s.logger.Warning(ctx, "failed to get the session by ID", err)
		return &common.Error{
			StatusCode: 500,
		}
	}
	if !found {
		s.logger.Warning(ctx, "failed to find session", errors.New("not found"))
		return &common.Error{
			StatusCode: 403,
		}
	}

	if !session.Locked {
		s.logger.Info(ctx, "Session was not locked")
		return nil
	}

	// Compare the code provided against code on session
	if session.UnlockCode == body.Code {
		// Valid code, update session
		err := s.userRepository.UnlockSession(authToken)
		if err != nil {
			s.logger.Warning(ctx, "failed to update the session to be unlocked", err)
			return &common.Error{
				StatusCode: 500,
			}
		}

		return nil
	}
	return &common.Error{
		StatusCode: 403,
	}
}

// SendForgotPassword sends a forgot password code via email to the user to reset their link. Since this is an
// unprotected endpoint (not auth) and we are calling a third party service (Sendgrid). We throttle the number of
// forgot password requests by a particular IP. 25 in 24 hours. There can be exceptions but this is on a case by
// case basis. I'd imagine places like universities would break this.
func (s *Services) SendForgotPassword(ctx context.Context, currentIP string, body SendForgotPasswordBody) *common.Error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "SendForgotPassword")
	if body.Email == "" {
		s.logger.Warning(ctx, "Email is missing", errors.New("error: email required"))
		return &common.Error{
			StatusCode: 400,
			Message: "Email is required",
		}
	}

	// Verify email exists
	userExists, user, err := s.userRepository.GetUserByEmail(body.GetFormattedEmail())
	if err != nil {
		s.logger.Warning(ctx, "failed to look up user", err)
		return &common.Error{
			StatusCode: 500,
		}
	}
	if !userExists {
		s.logger.Warning(ctx, "user does not exists", errors.New("not found"))
		return &common.Error{
			StatusCode: 200,
			Message: "Forgot password email sent",
		}
	}

	// Create instance
	code := uuid.New().String()
	now := time.Now()
	expiry := now.AddDate(0, 0, 1)
	err = s.forgotPasswordRepository.Save(ForgotPasswordCode{
		Email: body.GetFormattedEmail(),
		Code: code,
		Created: now,
		Expiry: expiry,
	})
	if err != nil {
		s.logger.Warning(ctx, "failed to save forgot password", err)
		return &common.Error{
			StatusCode: 500,
		}
	}

	// Send email
	err = forgotpasswordemail.SendForgotPasswordEmail(user.Greeting(), body.GetFormattedEmail(), code)
	if err != nil {
		s.logger.Warning(ctx, "failed to send the forgot password email", err)
		return &common.Error{
			StatusCode: 500,
		}
	}

	return nil
}

func (s *Services) ForgotPassword(ctx context.Context, currentIP string, body ResetForgotPasswordBody) *common.Error {
	ctx = context.WithValue(ctx, logging.CtxServiceMethod, "ForgotPassword")

	// Validate all fields
	if body.GetFormattedEmail() == "" {
		s.logger.Warning(ctx, "email is required", errors.New("error: invalid request"))
		return &common.Error{
			StatusCode: 400,
			Message: "Email is required",
		}
	}
	if body.Code == "" {
		s.logger.Warning(ctx, "email is required", errors.New("error: invalid request"))
		return &common.Error{
			StatusCode: 400,
			Message: "Code is required",
		}
	}

	// TODO - Add some throttling

	// Check if forgot password code exists
	exists, err := s.forgotPasswordRepository.Exists(body.GetFormattedEmail(), body.Code)
	if err != nil {
		s.logger.Warning(ctx, "failed to check if forgot password exists", err)
		return &common.Error{
			StatusCode: 500,
		}
	}
	if !exists {
		s.logger.Warning(ctx, "invalid email + code combination for forgot password", errors.New("unauthorized"))
		return &common.Error{
			StatusCode: 403,
		}
	}

	// Validate password strength
	if !s.isValidPassword(ctx, body.NewPassword) {
		s.logger.Error(ctx, "password does not meet requirements", errors.New("invalid password"))
		return &common.Error{
			StatusCode: 400,
			Message: "Password does not meet requirements",
		}
	}

	// Update password
	hash, err := s.getEncryptedPassword(body.NewPassword)
	if err != nil {
		s.logger.Warning(ctx, "failed to hash password", err)
		return &common.Error{
			StatusCode: 500,
		}
	}
	err = s.userRepository.UpdatePassword(body.GetFormattedEmail(), hash)
	if err != nil {
		s.logger.Warning(ctx, "failed to update password", err)
		return &common.Error{
			StatusCode: 500,
		}
	}

	// Update forgot password code
	err = s.forgotPasswordRepository.MarkCodeAsComplete(body.GetFormattedEmail(), body.Code)
	if err != nil {
		s.logger.Error(ctx, "failed to update mark code as completed", err)
		// Not returning error for UX
	}

	return nil
}