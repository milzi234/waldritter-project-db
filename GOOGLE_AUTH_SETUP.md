# Google OAuth Setup Instructions

This guide will help you set up Google OAuth authentication for the Waldritter admin interface with domain restriction to @waldritter.de emails.

## Prerequisites

- Google Cloud Console account
- Access to create OAuth credentials
- A Google Workspace domain for waldritter.de (for domain restriction)

## 1. Google Cloud Console Setup

### Create OAuth 2.0 Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to **APIs & Services** > **Credentials**
4. Click **+ CREATE CREDENTIALS** > **OAuth client ID**
5. If prompted, configure the OAuth consent screen first:
   - Choose **Internal** (if using Google Workspace) or **External**
   - Fill in the required fields
   - Add `waldritter.de` as an authorized domain
6. For Application type, select **Web application**
7. Give it a name (e.g., "Waldritter Admin")
8. Add Authorized JavaScript origins:
   - `http://localhost:5173`
   - `https://localhost:5173`
   - Your production URL when ready
9. Add Authorized redirect URIs:
   - `http://localhost:5173/login/callback`
   - `https://localhost:5173/login/callback`
   - Your production callback URL when ready
10. Click **Create** and save your Client ID

## 2. Local Development Setup

### Frontend Configuration

1. Copy the example environment file:
   ```bash
   cd website-project-db-admin-ui2
   cp .env.development.example .env.development
   ```

2. Update `.env.development` with your actual Client ID:
   ```bash
   VITE_GOOGLE_CLIENT_ID=your-actual-client-id.apps.googleusercontent.com
   VITE_API_URL=http://localhost:3000
   ```

### Backend Configuration

1. Copy the example environment file:
   ```bash
   cd website-project-db-api
   cp .env.development.example .env.development
   ```

2. Update `.env.development` with your actual Client ID:
   ```bash
   GOOGLE_CLIENT_ID=your-actual-client-id.apps.googleusercontent.com
   ```

**Note:** The `.env.development` files are gitignored and should never be committed to version control.

## 3. Running the Application

### Start the Rails Backend

```bash
cd website-project-db-api
rails server
```

The API will run on http://localhost:3000

### Start the Vue Frontend

```bash
cd website-project-db-admin-ui2
npm run dev
```

The frontend will run on https://localhost:5173 (with SSL)

## 4. Testing the Authentication

1. Visit https://localhost:5173
2. You should see the login page with a "Mit Google anmelden" button
3. Click the button to sign in
4. Use a @waldritter.de Google account
5. Non-waldritter.de accounts will be rejected with an error message

## 5. How It Works

### Authentication Flow

1. User clicks "Sign in with Google" on the Vue frontend
2. Google Identity Services handles the OAuth flow
3. After successful authentication, Google returns an ID token
4. Frontend validates the domain (must be @waldritter.de)
5. Frontend stores the token and includes it in API requests
6. Rails backend validates the Google token on each request
7. Backend creates/updates user record based on Google profile

### Security Features

- **Domain Restriction**: Only @waldritter.de emails can sign in
- **Token Validation**: Tokens are validated on both frontend and backend
- **HTTPS in Development**: Using mkcert for secure local development
- **Backward Compatibility**: Falls back to Dex/OIDC during migration period

## 6. Migration from OIDC/Dex

The system currently supports both authentication methods:
- Google OAuth (primary)
- Dex/OIDC (fallback for migration period)

To complete the migration:
1. Ensure all users have Google Workspace accounts
2. Test thoroughly with Google OAuth
3. Remove Dex/OIDC code when ready

## 7. Production Deployment

For production deployment, you'll need to:

1. Update Google Cloud Console with production URLs
2. Set environment variables on your production servers:
   - `GOOGLE_CLIENT_ID` for Rails
   - `VITE_GOOGLE_CLIENT_ID` for Vue build
3. Ensure HTTPS is configured (required for Google OAuth)
4. Update CORS settings if needed

## 8. Troubleshooting

### Common Issues

**"Google Sign-In ist nicht konfiguriert"**
- Check that VITE_GOOGLE_CLIENT_ID is set in .env.development

**"Invalid token" errors from API**
- Verify GOOGLE_CLIENT_ID matches in both frontend and backend
- Check that the googleauth gem is installed: `bundle install`

**"Nur waldritter.de E-Mail-Adressen sind erlaubt"**
- This is expected for non-waldritter.de accounts
- Ensure you're using a Google Workspace account with @waldritter.de domain

**SSL Certificate errors**
- Regenerate certificates: `mkcert localhost 127.0.0.1 ::1`
- Ensure mkcert CA is installed: `mkcert -install`

## 9. Additional Notes

- The Google token is stored in localStorage and included in API requests
- Tokens expire after 1 hour (handled automatically by Google)
- User information (email, name) is synced from Google profile
- The system creates User records automatically on first login