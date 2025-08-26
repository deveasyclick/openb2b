# Authentication Documentation

This project uses [Clerk](https://clerk.com) for authentication and session management. Clerk handles **sign-in, sign-up, and other authentication flows**, while our application manages user and organization data through webhooks, onboarding, and custom session claims.

---

## Setup

### 1. Create a Clerk Account & Project

1. Go to [Clerk](https://clerk.com) and create an account.
2. Create a new project from the dashboard.
3. Configure your frontend to use Clerk for authentication.

For more details, see the [Clerk Documentation](https://clerk.com/docs).

---

### 2. Add Clerk Secrets to Your Backend

Your backend requires **two secrets** from Clerk:

1. **Clerk API Key (Secret)** – for verifying sessions and interacting with the Clerk API.
2. **Clerk Webhook Secret** – for verifying incoming webhook requests from Clerk.

**Steps:**

1. In the Clerk dashboard, navigate to **API Keys** and copy your **API Key (Secret)**.
2. Navigate to **Webhooks**, create your webhook, and copy the **Signing Secret**.
3. Add both secrets to your `.env` file:

```bash
# Clerk API key for backend server requests
CLERK_API_KEY=your-clerk-api-secret

# Clerk Webhook signing secret to verify incoming webhook events
CLERK_WEBHOOK_SECRET=your-webhook-signing-secret
```

---

### 3. Configure Custom Session Claims

Under **Session Management** in Clerk settings, add the following custom session claims:

```json
{
  "user_id": "{{user.external_id}}",
  "org_id": "{{user.public_metadata.org_id}}"
}
```

* `user_id` → Maps to the `external_id` of the user in Clerk.
* `org_id` → Pulled from the user's `public_metadata` once their organization is created through onboarding.

These claims allow our backend to easily identify users and their associated organizations.

---

### 4. Setup Webhooks

We use Clerk webhooks to sync user data into our database.

1. Go to **Webhooks** in the Clerk dashboard.
2. Add a new webhook pointing to:

```bash
https://f69b-185-107-56-125.ngrok-free.app/api/v1/webhooks/handleEvents
```

3. After creating the webhook, click it and copy the **Signing Secret**.
4. Add the secret to your `.env` file (see step 2 above).

When a user signs up, Clerk sends an event to our webhook. The webhook handler then creates a corresponding user in our database.

---

## Organization Flow

* **User Signup**: The user signs up via Clerk (sign in / sign up handled entirely by Clerk).
* **Onboarding Form**: After signup, the user is taken through an **onboarding form** where they create an **Organization** in our system.
* **Org Creation**: Once the org is created, we set the `org_id` in the user’s `public_metadata` on Clerk.
* **Custom Claims Update**: The `org_id` is included as part of the user’s session claims, ensuring it’s available on every authenticated request.
* **Homepage Redirect**: After onboarding is complete, the user is redirected to the homepage with their session enriched by the `org_id`.

---

## Summary

* Clerk manages **authentication (sign in / sign up / session handling)**.
* Our app manages:

  * **User persistence** via Clerk webhooks.
  * **Organization creation** via the onboarding form.
  * **Session enrichment** with `org_id` once the org is created.
* Custom session claims (`user_id`, `org_id`) ensure we can associate requests with users and organizations.
* **Backend securely communicates with Clerk** using the `CLERK_API_KEY` and verifies webhooks using `CLERK_WEBHOOK_SECRET`.
