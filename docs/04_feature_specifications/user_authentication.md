# User Authentication Feature Specification

> **✅ INCLUDED IN FOUNDATION EDITION**
> This feature is fully available in the open-source Foundation edition.

**Document Version**: 1.0
**Last Updated**: January 2025
**Feature Owner**: Security Team
**Status**: Implemented
**Edition**: Foundation ✅ | Premium ✅

---

## Overview

The User Authentication system provides secure, scalable user registration, login, and session management for the Great Nigeria Library platform. It implements modern security practices including JWT tokens, multi-factor authentication, and OAuth integration.

## Feature Purpose

### Goals
1. **Secure Access**: Protect user accounts and platform resources
2. **User Experience**: Provide smooth, frictionless authentication flows
3. **Scalability**: Support millions of concurrent authenticated users
4. **Compliance**: Meet international security and privacy standards
5. **Flexibility**: Support multiple authentication methods and providers

### Success Metrics
- **Registration Conversion**: 85%+ of started registrations completed
- **Login Success Rate**: 98%+ successful login attempts
- **Security Incidents**: Zero successful unauthorized access attempts
- **User Satisfaction**: 90%+ satisfaction with authentication experience

## Technical Architecture

### Authentication System Architecture
Comprehensive user authentication and security infrastructure:
-- Main users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    bio TEXT,
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'creator', 'moderator', 'admin')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'deleted')),
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    email_verification_expires_at TIMESTAMP WITH TIME ZONE,
    password_reset_token VARCHAR(255),
    password_reset_expires_at TIMESTAMP WITH TIME ZONE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    preferences JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User sessions for tracking active logins
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL,
    device_info JSONB,
    ip_address INET,
    user_agent TEXT,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Multi-factor authentication
CREATE TABLE user_mfa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    method VARCHAR(20) NOT NULL CHECK (method IN ('totp', 'sms', 'email', 'backup_codes')),
    secret VARCHAR(255),
    phone_number VARCHAR(20),
    backup_codes TEXT[],
    enabled BOOLEAN DEFAULT FALSE,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, method)
);

-- OAuth integrations
CREATE TABLE oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255),
    provider_data JSONB,
    access_token VARCHAR(1000),
    refresh_token VARCHAR(1000),
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(provider, provider_user_id)
);

-- Role-based permissions
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    permissions TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User role assignments
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, role_id)
);

-- Security audit log
CREATE TABLE auth_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### API Endpoints

#### Public Authentication Endpoints

```yaml
# User Registration
POST /api/v1/auth/register:
  body:
    type: object
    required: [email, password, name]
    properties:
      email:
        type: string
        format: email
        maxLength: 255
      password:
        type: string
        minLength: 8
        pattern: "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]"
      name:
        type: string
        minLength: 2
        maxLength: 255
      username:
        type: string
        minLength: 3
        maxLength: 50
        pattern: "^[a-zA-Z0-9_-]+$"
  responses:
    201:
      description: User registered successfully
      schema:
        type: object
        properties:
          user:
            $ref: '#/components/schemas/User'
          tokens:
            $ref: '#/components/schemas/AuthTokens'
    400:
      description: Validation error
    409:
      description: Email or username already exists

# User Login
POST /api/v1/auth/login:
  body:
    type: object
    required: [email, password]
    properties:
      email:
        type: string
        format: email
      password:
        type: string
      remember_me:
        type: boolean
        default: false
      mfa_code:
        type: string
        description: Required if MFA is enabled
  responses:
    200:
      description: Login successful
    401:
      description: Invalid credentials
    423:
      description: Account locked due to too many failed attempts
    428:
      description: MFA code required

# Refresh Access Token
POST /api/v1/auth/refresh:
  body:
    type: object
    required: [refresh_token]
    properties:
      refresh_token:
        type: string
  responses:
    200:
      description: New access token generated
    401:
      description: Invalid or expired refresh token

# Password Reset Request
POST /api/v1/auth/password/reset:
  body:
    type: object
    required: [email]
    properties:
      email:
        type: string
        format: email
  responses:
    200:
      description: Reset email sent (always returns 200 for security)

# Password Reset Confirmation
POST /api/v1/auth/password/reset/confirm:
  body:
    type: object
    required: [token, password]
    properties:
      token:
        type: string
      password:
        type: string
        minLength: 8
  responses:
    200:
      description: Password reset successful
    400:
      description: Invalid or expired token

# Email Verification
POST /api/v1/auth/verify-email:
  body:
    type: object
    required: [token]
    properties:
      token:
        type: string
  responses:
    200:
      description: Email verified successfully
    400:
      description: Invalid or expired token
```

#### Protected Endpoints

```yaml
# Get Current User
GET /api/v1/auth/me:
  authentication: required
  responses:
    200:
      description: Current user information
      schema:
        $ref: '#/components/schemas/User'

# Update User Profile
PUT /api/v1/auth/me:
  authentication: required
  body:
    type: object
    properties:
      name:
        type: string
        maxLength: 255
      username:
        type: string
        maxLength: 50
      bio:
        type: string
        maxLength: 500
      preferences:
        type: object
  responses:
    200:
      description: Profile updated successfully

# Change Password
POST /api/v1/auth/password/change:
  authentication: required
  body:
    type: object
    required: [current_password, new_password]
    properties:
      current_password:
        type: string
      new_password:
        type: string
        minLength: 8
  responses:
    200:
      description: Password changed successfully
    400:
      description: Invalid current password

# Logout
POST /api/v1/auth/logout:
  authentication: required
  responses:
    200:
      description: Logged out successfully

# Logout All Sessions
POST /api/v1/auth/logout-all:
  authentication: required
  responses:
    200:
      description: All sessions terminated

# Enable MFA
POST /api/v1/auth/mfa/enable:
  authentication: required
  body:
    type: object
    required: [method]
    properties:
      method:
        type: string
        enum: [totp, sms]
      phone_number:
        type: string
        description: Required for SMS method
  responses:
    200:
      description: MFA setup initiated
      schema:
        type: object
        properties:
          qr_code:
            type: string
            description: Base64 QR code for TOTP
          backup_codes:
            type: array
            items:
              type: string

# Verify MFA Setup
POST /api/v1/auth/mfa/verify:
  authentication: required
  body:
    type: object
    required: [method, code]
    properties:
      method:
        type: string
        enum: [totp, sms]
      code:
        type: string
  responses:
    200:
      description: MFA verified and enabled

# Disable MFA
POST /api/v1/auth/mfa/disable:
  authentication: required
  body:
    type: object
    required: [password]
    properties:
      password:
        type: string
  responses:
    200:
      description: MFA disabled
```

#### OAuth Endpoints

```yaml
# OAuth Login (Google, Facebook, Twitter)
GET /api/v1/auth/oauth/{provider}:
  parameters:
    - provider: string (google|facebook|twitter|github)
    - redirect_uri: string
  responses:
    302:
      description: Redirect to OAuth provider

# OAuth Callback
GET /api/v1/auth/oauth/{provider}/callback:
  parameters:
    - code: string
    - state: string
  responses:
    302:
      description: Redirect to frontend with tokens

# Link OAuth Account
POST /api/v1/auth/oauth/{provider}/link:
  authentication: required
  body:
    type: object
    required: [access_token]
    properties:
      access_token:
        type: string
  responses:
    200:
      description: OAuth account linked

# Unlink OAuth Account
DELETE /api/v1/auth/oauth/{provider}/unlink:
  authentication: required
  responses:
    200:
      description: OAuth account unlinked
```

### Frontend Components

#### Authentication Components

```typescript
// Login Form Component
interface LoginFormProps {
  onSuccess?: (user: User) => void;
  redirectTo?: string;
}

export const LoginForm: React.FC<LoginFormProps> = ({ onSuccess, redirectTo }) => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    remember_me: false,
  });
  const [mfaRequired, setMfaRequired] = useState(false);
  const [mfaCode, setMfaCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await authService.login({
        ...formData,
        mfa_code: mfaCode || undefined,
      });

      if (response.data.mfa_required) {
        setMfaRequired(true);
        return;
      }

      // Store tokens and user data
      authService.setTokens(response.data.tokens);
      onSuccess?.(response.data.user);
      
      if (redirectTo) {
        window.location.href = redirectTo;
      }
    } catch (error) {
      if (error.code === 'MFA_REQUIRED') {
        setMfaRequired(true);
      } else {
        setError(error.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="login-form">
      <div className="form-group">
        <label htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          required
          disabled={loading}
        />
      </div>

      <div className="form-group">
        <label htmlFor="password">Password</label>
        <input
          id="password"
          type="password"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          required
          disabled={loading}
        />
      </div>

      {mfaRequired && (
        <div className="form-group">
          <label htmlFor="mfa-code">Authentication Code</label>
          <input
            id="mfa-code"
            type="text"
            value={mfaCode}
            onChange={(e) => setMfaCode(e.target.value)}
            placeholder="Enter 6-digit code"
            maxLength={6}
            required
          />
        </div>
      )}

      <div className="form-group">
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formData.remember_me}
            onChange={(e) => setFormData({ ...formData, remember_me: e.target.checked })}
          />
          Remember me
        </label>
      </div>

      {error && (
        <div className="error-message" role="alert">
          {error}
        </div>
      )}

      <button type="submit" disabled={loading} className="submit-button">
        {loading ? 'Signing in...' : 'Sign In'}
      </button>

      <div className="form-links">
        <Link to="/auth/forgot-password">Forgot your password?</Link>
        <Link to="/auth/register">Don't have an account? Sign up</Link>
      </div>

      <div className="oauth-providers">
        <button type="button" onClick={() => authService.loginWithGoogle()}>
          Continue with Google
        </button>
        <button type="button" onClick={() => authService.loginWithFacebook()}>
          Continue with Facebook
        </button>
      </div>
    </form>
  );
};

// Registration Form Component
interface RegistrationFormProps {
  onSuccess?: (user: User) => void;
}

export const RegistrationForm: React.FC<RegistrationFormProps> = ({ onSuccess }) => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
    username: '',
    acceptTerms: false,
  });
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.email) {
      newErrors.email = 'Email is required';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Email is invalid';
    }

    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
    } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])/.test(formData.password)) {
      newErrors.password = 'Password must contain uppercase, lowercase, number, and special character';
    }

    if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
    }

    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }

    if (!formData.acceptTerms) {
      newErrors.acceptTerms = 'You must accept the terms and conditions';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setLoading(true);

    try {
      const response = await authService.register({
        email: formData.email,
        password: formData.password,
        name: formData.name,
        username: formData.username || undefined,
      });

      // Store tokens and user data
      authService.setTokens(response.data.tokens);
      onSuccess?.(response.data.user);
    } catch (error) {
      setErrors({ general: error.message });
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="registration-form">
      {/* Form fields implementation */}
    </form>
  );
};

// MFA Setup Component
export const MFASetupModal: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  const [method, setMethod] = useState<'totp' | 'sms'>('totp');
  const [qrCode, setQrCode] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [backupCodes, setBackupCodes] = useState<string[]>([]);
  const [step, setStep] = useState<'method' | 'setup' | 'verify' | 'complete'>('method');

  // Implementation for MFA setup flow
  return (
    <Modal isOpen onClose={onClose}>
      <div className="mfa-setup">
        {/* MFA setup steps implementation */}
      </div>
    </Modal>
  );
};
```

#### Redux State Management

```typescript
// Auth slice
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  tokens: {
    access_token: string;
    refresh_token: string;
    expires_at: number;
  } | null;
}

const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  loading: false,
  error: null,
  tokens: null,
};

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setLoading: (state, action) => {
      state.loading = action.payload;
    },
    setError: (state, action) => {
      state.error = action.payload;
      state.loading = false;
    },
    setAuth: (state, action) => {
      state.user = action.payload.user;
      state.tokens = action.payload.tokens;
      state.isAuthenticated = true;
      state.loading = false;
      state.error = null;
    },
    clearAuth: (state) => {
      state.user = null;
      state.tokens = null;
      state.isAuthenticated = false;
      state.loading = false;
      state.error = null;
    },
    updateUser: (state, action) => {
      if (state.user) {
        state.user = { ...state.user, ...action.payload };
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginUser.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.user = action.payload.user;
        state.tokens = action.payload.tokens;
        state.isAuthenticated = true;
        state.loading = false;
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.error = action.error.message || 'Login failed';
        state.loading = false;
      });
  },
});

// Async thunks
export const loginUser = createAsyncThunk(
  'auth/login',
  async (credentials: LoginCredentials) => {
    const response = await authService.login(credentials);
    return response.data;
  }
);

export const registerUser = createAsyncThunk(
  'auth/register',
  async (userData: RegistrationData) => {
    const response = await authService.register(userData);
    return response.data;
  }
);
```

### Security Implementation

#### Password Security

```typescript
// Password hashing (backend)
import bcrypt from 'bcrypt';
import crypto from 'crypto';

export class PasswordService {
  private static readonly SALT_ROUNDS = 12;
  private static readonly PEPPER = process.env.PASSWORD_PEPPER || '';

  static async hashPassword(password: string): Promise<string> {
    const pepperedPassword = password + this.PEPPER;
    return bcrypt.hash(pepperedPassword, this.SALT_ROUNDS);
  }

  static async verifyPassword(password: string, hash: string): Promise<boolean> {
    const pepperedPassword = password + this.PEPPER;
    return bcrypt.compare(pepperedPassword, hash);
  }

  static generateSecureToken(): string {
    return crypto.randomBytes(32).toString('hex');
  }

  static validatePasswordStrength(password: string): {
    isValid: boolean;
    errors: string[];
  } {
    const errors: string[] = [];

    if (password.length < 8) {
      errors.push('Password must be at least 8 characters long');
    }

    if (!/[a-z]/.test(password)) {
      errors.push('Password must contain at least one lowercase letter');
    }

    if (!/[A-Z]/.test(password)) {
      errors.push('Password must contain at least one uppercase letter');
    }

    if (!/\d/.test(password)) {
      errors.push('Password must contain at least one number');
    }

    if (!/[@$!%*?&]/.test(password)) {
      errors.push('Password must contain at least one special character');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }
}
```

#### JWT Token Management

```typescript
// JWT service (backend)
import jwt from 'jsonwebtoken';

export class JWTService {
  private static readonly ACCESS_TOKEN_SECRET = process.env.JWT_SECRET!;
  private static readonly REFRESH_TOKEN_SECRET = process.env.JWT_REFRESH_SECRET!;
  private static readonly ACCESS_TOKEN_EXPIRY = '15m';
  private static readonly REFRESH_TOKEN_EXPIRY = '7d';

  static generateTokenPair(payload: TokenPayload): {
    access_token: string;
    refresh_token: string;
    expires_in: number;
  } {
    const access_token = jwt.sign(payload, this.ACCESS_TOKEN_SECRET, {
      expiresIn: this.ACCESS_TOKEN_EXPIRY,
    });

    const refresh_token = jwt.sign(
      { user_id: payload.sub },
      this.REFRESH_TOKEN_SECRET,
      { expiresIn: this.REFRESH_TOKEN_EXPIRY }
    );

    return {
      access_token,
      refresh_token,
      expires_in: 15 * 60, // 15 minutes in seconds
    };
  }

  static verifyAccessToken(token: string): TokenPayload {
    return jwt.verify(token, this.ACCESS_TOKEN_SECRET) as TokenPayload;
  }

  static verifyRefreshToken(token: string): { user_id: string } {
    return jwt.verify(token, this.REFRESH_TOKEN_SECRET) as { user_id: string };
  }

  static decodeToken(token: string): any {
    return jwt.decode(token);
  }
}
```

#### Rate Limiting and Brute Force Protection

```typescript
// Rate limiting for authentication endpoints
import rateLimit from 'express-rate-limit';
import RedisStore from 'rate-limit-redis';
import Redis from 'ioredis';

const redis = new Redis(process.env.REDIS_URL);

export const authRateLimiter = rateLimit({
  store: new RedisStore({
    sendCommand: (...args: string[]) => redis.call(...args),
  }),
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // 5 attempts per window
  message: {
    error: 'Too many authentication attempts, please try again later',
  },
  standardHeaders: true,
  legacyHeaders: false,
  skip: (req) => {
    // Skip rate limiting for successful requests
    return req.authSuccess === true;
  },
});

export const registrationRateLimiter = rateLimit({
  store: new RedisStore({
    sendCommand: (...args: string[]) => redis.call(...args),
  }),
  windowMs: 60 * 60 * 1000, // 1 hour
  max: 3, // 3 registrations per hour per IP
  message: {
    error: 'Too many registration attempts, please try again later',
  },
});
```

### Integration Points

#### OAuth Providers Configuration

```typescript
// OAuth provider configurations
export const oauthConfig = {
  google: {
    clientId: process.env.GOOGLE_CLIENT_ID!,
    clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    redirectUri: `${process.env.API_URL}/auth/oauth/google/callback`,
    scope: ['openid', 'email', 'profile'],
  },
  facebook: {
    clientId: process.env.FACEBOOK_CLIENT_ID!,
    clientSecret: process.env.FACEBOOK_CLIENT_SECRET!,
    redirectUri: `${process.env.API_URL}/auth/oauth/facebook/callback`,
    scope: ['email', 'public_profile'],
  },
  github: {
    clientId: process.env.GITHUB_CLIENT_ID!,
    clientSecret: process.env.GITHUB_CLIENT_SECRET!,
    redirectUri: `${process.env.API_URL}/auth/oauth/github/callback`,
    scope: ['user:email'],
  },
};
```

### Monitoring and Analytics

#### Security Event Logging

```typescript
// Security audit logging
export class SecurityAuditLogger {
  static async logAuthEvent(event: {
    user_id?: string;
    action: string;
    success: boolean;
    ip_address: string;
    user_agent: string;
    details?: any;
  }) {
    await db.auth_audit_log.create({
      data: {
        user_id: event.user_id,
        action: event.action,
        success: event.success,
        ip_address: event.ip_address,
        user_agent: event.user_agent,
        details: event.details,
      },
    });

    // Alert on suspicious activities
    if (!event.success && this.isSuspiciousActivity(event)) {
      await this.triggerSecurityAlert(event);
    }
  }

  private static isSuspiciousActivity(event: any): boolean {
    // Implement suspicious activity detection logic
    return false;
  }

  private static async triggerSecurityAlert(event: any) {
    // Implement security alerting logic
  }
}
```

---

*This feature specification serves as the complete guide for implementing and maintaining the user authentication system for the Great Nigeria Library platform.*
